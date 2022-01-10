package services

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

	Describe("Update", func() {
		BeforeEach(func() {
			_ = os.Setenv("SCHEDULER_TABLE_NAME", table)

			dynamo = fakeDynamoDB{}
			db = NewDatabase(&dynamo)
		})

		Describe("success", func() {
			var (
				err              error
				batchWriteInputs []*dynamodb.BatchWriteItemInput
			)

			BeforeEach(func() {
				dynamo.PushBatchWriteOutput(&dynamodb.BatchWriteItemOutput{
					UnprocessedItems: map[string][]*dynamodb.WriteRequest{
						table: {
							&dynamodb.WriteRequest{
								PutRequest: &dynamodb.PutRequest{
									Item: map[string]*dynamodb.AttributeValue{
										"id": {S: aws.String(id)},
									},
								},
							},
						},
					},
				})
				dynamo.PushBatchWriteOutput(&dynamodb.BatchWriteItemOutput{})

				err = db.Update(context.TODO(), []*UpdateInput{
					{
						ID: id,
					},
				})

				batchWriteInputs = []*dynamodb.BatchWriteItemInput{
					dynamo.PullBatchWriteInput(),
					dynamo.PullBatchWriteInput(),
				}
			})

			It("reads table name from env", func() {
				for _, batchWriteInput := range batchWriteInputs {
					Expect(batchWriteInput.RequestItems[table]).NotTo(BeNil())
				}
			})

			It("sets put item form input", func() {
				for _, batchWriteInput := range batchWriteInputs {
					for _, r := range batchWriteInput.RequestItems[table] {
						Expect(*r.PutRequest.Item["id"].S).To(Equal(id))
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

		Describe("fail", func() {
			Context("marshal error", func() {
				var (
					realMarshal marshalStorage
					err         error
				)

				BeforeEach(func() {
					realMarshal = marshalStorageStruct

					marshalStorageStruct = func(
						in interface{}) (
						map[string]*dynamodb.AttributeValue, error) {

						return nil, fmt.Errorf("marshal error")
					}

					err = db.Update(context.TODO(), []*UpdateInput{
						{
							ID: id,
						},
					})
				})

				It("return error", func() {
					Expect(err).NotTo(BeNil())
				})

				AfterEach(func() {
					marshalStorageStruct = realMarshal
				})
			})

			Context("batch write error", func() {
				var err error

				BeforeEach(func() {
					dynamo.BatchWriteError = fmt.Errorf("batch write error")

					err = db.Update(context.TODO(), []*UpdateInput{
						{
							ID: id,
						},
					})
				})

				It("return error", func() {
					Expect(err).NotTo(BeNil())
				})

				AfterEach(func() {
					dynamo.BatchWriteError = nil
				})
			})
		})
	})
	Describe("chunkBy", func() {
		Context("100 by 25", func() {
			var items []*UpdateInput
			var ret [][]*UpdateInput

			BeforeEach(func() {
				items = make([]*UpdateInput, 100)
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
			var items []*UpdateInput
			var ret [][]*UpdateInput

			BeforeEach(func() {
				items = make([]*UpdateInput, 73)
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
			var items []*UpdateInput
			var ret [][]*UpdateInput

			BeforeEach(func() {
				items = make([]*UpdateInput, 20)
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

	BatchWriteError error

	batchWriteInputs  list.List
	batchWriteOutputs list.List
}

func (db *fakeDynamoDB) BatchWriteItemWithContext(
	_ aws.Context,
	input *dynamodb.BatchWriteItemInput,
	_ ...request.Option) (*dynamodb.BatchWriteItemOutput, error) {

	db.batchWriteInputs.PushBack(input)

	output := db.batchWriteOutputs.Front()

	if output == nil {
		return nil, db.BatchWriteError
	}

	db.batchWriteOutputs.Remove(output)

	return output.Value.(*dynamodb.BatchWriteItemOutput),
		db.BatchWriteError
}

func (db *fakeDynamoDB) PushBatchWriteOutput(
	output *dynamodb.BatchWriteItemOutput) {
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
	db.BatchWriteError = nil
	db.batchWriteInputs.Init()
	db.batchWriteOutputs.Init()
}
