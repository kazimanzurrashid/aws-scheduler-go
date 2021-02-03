package storage

import (
	"container/list"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database", func() {
	const (
		table = "scheduler_v1"
		id    = "1234567890"
	)

	var (
		dynamo fakeDynamoDB
		db     *Database
	)

	BeforeEach(func() {
		_ = os.Setenv("SCHEDULER_TABLE_NAME", table)

		dynamo = fakeDynamoDB{}
		db = NewDatabase(&dynamo)
	})

	Describe("Update", func() {
		Describe("success", func() {
			Context("matching schedules", func() {
				var (
					err              error
					queryInputs      []*dynamodb.QueryInput
					batchWriteInputs []*dynamodb.BatchWriteItemInput
				)

				BeforeEach(func() {
					dynamo.PushQueryOutput(&dynamodb.QueryOutput{
						Items: []map[string]*dynamodb.AttributeValue{
							{
								"id":     {S: aws.String(id)},
								"status": {S: aws.String(scheduleStatusIdle)},
							},
						},
						LastEvaluatedKey: map[string]*dynamodb.AttributeValue{
							"id":     {S: aws.String(id)},
							"dueAt":  {N: aws.String("77627362")},
							"status": {S: aws.String(scheduleStatusIdle)},
						},
					})

					dynamo.PushQueryOutput(&dynamodb.QueryOutput{
						Items: []map[string]*dynamodb.AttributeValue{
							{
								"id":     {S: aws.String(id)},
								"status": {S: aws.String(scheduleStatusIdle)},
							},
						},
					})

					dynamo.PushBatchWriteOutput(&dynamodb.BatchWriteItemOutput{})
					dynamo.PushBatchWriteOutput(&dynamodb.BatchWriteItemOutput{
						UnprocessedItems: map[string][]*dynamodb.WriteRequest{
							table: {
								&dynamodb.WriteRequest{
									PutRequest: &dynamodb.PutRequest{
										Item: map[string]*dynamodb.AttributeValue{
											"id": {S: aws.String(id)},
											"status": {S: aws.String(
												scheduleStatusQueued)},
										},
									},
								},
							},
						},
					})
					dynamo.PushBatchWriteOutput(&dynamodb.BatchWriteItemOutput{})

					err = db.Update(context.TODO())

					queryInputs = []*dynamodb.QueryInput{
						dynamo.PullQueryInput(),
						dynamo.PullQueryInput(),
					}

					batchWriteInputs = []*dynamodb.BatchWriteItemInput{
						dynamo.PullBatchWriteInput(),
						dynamo.PullBatchWriteInput(),
						dynamo.PullBatchWriteInput(),
					}
				})

				It("reads table name from env", func() {
					for _, queryInput := range queryInputs {
						Expect(*queryInput.TableName).To(Equal(table))
					}

					for _, batchWriteInput := range batchWriteInputs {
						Expect(batchWriteInput.RequestItems[table]).NotTo(BeNil())
					}
				})

				It("uses index to query", func() {
					for _, queryInput := range queryInputs {
						Expect(*queryInput.IndexName).To(Equal("ix_status_dueAt"))
					}
				})

				It("uses idle status to match", func() {
					for _, queryInput := range queryInputs {
						Expect(*queryInput.ExpressionAttributeValues[":s"].S).To(
							Equal(scheduleStatusIdle))
					}
				})

				It("uses dueAt to match", func() {
					for _, queryInput := range queryInputs {
						Expect(
							*queryInput.ExpressionAttributeValues[":da"].N).NotTo(
							Equal(""))
					}
				})

				It("updates status to queued", func() {
					for _, batchWriteInput := range batchWriteInputs {
						for _, r := range batchWriteInput.RequestItems[table] {
							Expect(*r.PutRequest.Item["status"].S).To(
								Equal(scheduleStatusQueued))
						}
					}
				})

				It("does not return error", func() {
					Expect(err).To(BeNil())
				})

				AfterEach(func() {
					dynamo.ClearState()
				})
			})

			Context("no matching schedule", func() {
				var err error

				BeforeEach(func() {
					dynamo.PushQueryOutput(&dynamodb.QueryOutput{})

					err = db.Update(context.TODO())
				})

				It("does not return error", func() {
					Expect(err).To(BeNil())
				})

				AfterEach(func() {
					dynamo.ClearState()
				})
			})
		})

		Describe("fail", func() {
			Context("query error", func() {
				var err error

				BeforeEach(func() {
					dynamo.QueryError = fmt.Errorf("query error")

					err = db.Update(context.TODO())
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})

				AfterEach(func() {
					dynamo.ClearState()
				})
			})

			Context("batch write error", func() {
				var err error

				BeforeEach(func() {
					dynamo.BatchWriteError = fmt.Errorf("batch write error")

					dynamo.PushQueryOutput(&dynamodb.QueryOutput{
						Items: []map[string]*dynamodb.AttributeValue{
							{
								"id": {S: aws.String(id)},
							},
						},
					})

					err = db.Update(context.TODO())
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})

				AfterEach(func() {
					dynamo.ClearState()
				})
			})
		})
	})

	Describe("chunkBy", func() {
		Context("100 by 25", func() {
			var items []map[string]*dynamodb.AttributeValue
			var ret [][]map[string]*dynamodb.AttributeValue

			BeforeEach(func() {
				items = make([]map[string]*dynamodb.AttributeValue, 100)
				ret = chunkBy(items, 25)
			})

			It("returns 4 chunk", func() {
				Expect(ret).To(HaveLen(4))
			})

			It("each chunk have 25 items", func() {
				for _, chunk := range ret {
					Expect(chunk).To(HaveLen(25))
				}
			})
		})

		Context("73 by 25", func() {
			var items []map[string]*dynamodb.AttributeValue
			var ret [][]map[string]*dynamodb.AttributeValue

			BeforeEach(func() {
				items = make([]map[string]*dynamodb.AttributeValue, 73)
				ret = chunkBy(items, 25)
			})

			It("returns 3 chunk", func() {
				Expect(ret).To(HaveLen(3))
			})

			It("first two chunk 25 and last one 23", func() {
				for i, chunk := range ret {
					if i < 2 {
						Expect(chunk).To(HaveLen(25))
					} else {
						Expect(chunk).To(HaveLen(23))
					}
				}
			})
		})

		Context("20 by 25", func() {
			var items []map[string]*dynamodb.AttributeValue
			var ret [][]map[string]*dynamodb.AttributeValue

			BeforeEach(func() {
				items = make([]map[string]*dynamodb.AttributeValue, 20)
				ret = chunkBy(items, 25)
			})

			It("returns 1 chunk", func() {
				Expect(ret).To(HaveLen(1))
			})

			It("chunk has 20", func() {
				Expect(ret[0]).To(HaveLen(20))
			})
		})
	})
})

type fakeDynamoDB struct {
	dynamodbiface.DynamoDBAPI

	QueryError      error
	BatchWriteError error

	queryInputs  list.List
	queryOutputs list.List

	batchWriteInputs  list.List
	batchWriteOutputs list.List
}

//goland:noinspection GoUnusedParameter
func (db *fakeDynamoDB) QueryWithContext(
	ctx aws.Context,
	input *dynamodb.QueryInput,
	options ...request.Option) (*dynamodb.QueryOutput, error) {

	db.queryInputs.PushBack(input)

	output := db.queryOutputs.Front()

	if output == nil {
		return nil, db.QueryError
	}

	db.queryOutputs.Remove(output)

	return output.Value.(*dynamodb.QueryOutput), db.QueryError
}

//goland:noinspection GoUnusedParameter
func (db *fakeDynamoDB) BatchWriteItemWithContext(
	ctx aws.Context,
	input *dynamodb.BatchWriteItemInput,
	options ...request.Option) (*dynamodb.BatchWriteItemOutput, error) {

	db.batchWriteInputs.PushBack(input)

	output := db.batchWriteOutputs.Front()

	if output == nil {
		return nil, db.BatchWriteError
	}

	db.batchWriteOutputs.Remove(output)

	return output.Value.(*dynamodb.BatchWriteItemOutput), db.BatchWriteError
}

func (db *fakeDynamoDB) PushQueryOutput(output *dynamodb.QueryOutput) {
	db.queryOutputs.PushBack(output)
}

func (db *fakeDynamoDB) PullQueryInput() *dynamodb.QueryInput {
	input := db.queryInputs.Front()

	if input == nil {
		return nil
	}

	db.queryInputs.Remove(input)

	return input.Value.(*dynamodb.QueryInput)
}

func (db *fakeDynamoDB) PushBatchWriteOutput(output *dynamodb.BatchWriteItemOutput) {
	db.batchWriteOutputs.PushBack(output)
}

func (db *fakeDynamoDB) PullBatchWriteInput() *dynamodb.BatchWriteItemInput {
	input := db.batchWriteInputs.Front()

	if input == nil {
		return nil
	}

	db.batchWriteInputs.Remove(input)

	return input.Value.(*dynamodb.BatchWriteItemInput)
}

func (db *fakeDynamoDB) ClearState() {
	db.QueryError = nil
	db.BatchWriteError = nil
	db.queryInputs.Init()
	db.queryOutputs.Init()
	db.batchWriteInputs.Init()
	db.batchWriteOutputs.Init()
}
