package storage

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database", func() {
	const (
		table  = "scheduler_v1"
		id     = "1234567890"
		url    = "https://foo.bar/do"
		method = "POST"
		accept = "application/json"
		body   = "{ \"foo\": \"bar\" }"
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

	Describe("NewDatabase", func() {
		It("returns new database", func() {
			Expect(db).NotTo(BeNil())
		})
	})

	Describe("Create", func() {
		Describe("success", func() {
			var (
				res string
				err error

				dueAt   time.Time
				headers map[string]string
			)

			BeforeEach(func() {
				dueAt = time.Now().Add(time.Minute * 1)
				headers = map[string]string{
					"accept": accept,
				}

				res, err = db.Create(context.TODO(), CreateInput{
					DueAt:   dueAt,
					URL:     url,
					Method:  method,
					Headers: headers,
					Body:    body,
				})
			})

			It("reads table name from env", func() {
				Expect(*dynamo.PutInput.TableName).To(Equal(table))
			})

			It("includes new id in put", func() {
				Expect(*dynamo.PutInput.Item["id"].S).NotTo(Equal(""))
			})

			It("includes status idle in put", func() {
				Expect(*dynamo.PutInput.Item["status"].S).To(
					Equal(ScheduleStatusIdle))
			})

			It("includes dummy in put", func() {
				Expect(*dynamo.PutInput.Item["dummy"].S).To(Equal(dummyValue))
			})

			It("includes createdAt in put", func() {
				Expect(dynamo.PutInput.Item["createdAt"].N).NotTo(Equal(""))
			})

			It("sets put from input", func() {
				Expect(*dynamo.PutInput.Item["dueAt"].N).To(
					Equal(strconv.FormatInt(dueAt.Unix(), 10)))
				Expect(*dynamo.PutInput.Item["url"].S).To(Equal(url))
				Expect(*dynamo.PutInput.Item["method"].S).To(Equal(method))
				Expect(*dynamo.PutInput.Item["headers"].M["accept"].S).To(
					Equal(accept))
				Expect(*dynamo.PutInput.Item["body"].S).To(Equal(body))
			})

			It("returns new id", func() {
				Expect(res).NotTo(Equal(""))
			})

			It("does not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("fail", func() {
			Context("id generate error", func() {
				var (
					realIdGen idGenerate
					res       string
					err       error
				)

				BeforeEach(func() {
					realIdGen = generateID
					generateID = func() (string, error) {
						return "", fmt.Errorf("id generate error")
					}

					res, err = db.Create(context.TODO(), CreateInput{})
				})

				It("does not return id", func() {
					Expect(res).To(Equal(""))
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})

				AfterEach(func() {
					generateID = realIdGen
				})
			})

			Context("input marshal error", func() {
				var (
					realMarshal marshal
					res         string
					err         error
				)

				BeforeEach(func() {
					realMarshal = marshalStruct
					marshalStruct = func(in interface{}) (
						map[string]*dynamodb.AttributeValue,
						error) {

						return nil, fmt.Errorf("marshal error")
					}

					res, err = db.Create(context.TODO(), CreateInput{})
				})

				It("does not return id", func() {
					Expect(res).To(Equal(""))
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})

				AfterEach(func() {
					marshalStruct = realMarshal
				})
			})

			Context("put error", func() {
				var (
					res string
					err error
				)

				BeforeEach(func() {
					dynamo.Error = awserr.New(
						"InternalError",
						"InternalError",
						nil)

					res, err = db.Create(context.TODO(), CreateInput{
						DueAt:  time.Now().Add(time.Minute * 5),
						URL:    url,
						Method: method,
						Body:   body,
					})
				})

				It("does not return id", func() {
					Expect(res).To(Equal(""))
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})

				AfterEach(func() {
					dynamo.Error = nil
				})
			})
		})
	})

	Describe("Cancel", func() {
		Describe("success", func() {
			var (
				res bool
				err error
			)

			BeforeEach(func() {
				res, err = db.Cancel(context.TODO(), id)
			})

			It("reads table name from env", func() {
				Expect(*dynamo.UpdateInput.TableName).To(Equal(table))
			})

			It("sets given id to match schedule", func() {
				Expect(*dynamo.UpdateInput.Key["id"].S).To(Equal(id))
			})

			It("only cancels idle schedule", func() {
				Expect(*dynamo.UpdateInput.ConditionExpression).NotTo(Equal(""))
				Expect(
					*dynamo.UpdateInput.ExpressionAttributeValues[":s2"].S).To(
					Equal(ScheduleStatusIdle))
			})

			It("updates status to canceled", func() {
				Expect(*dynamo.UpdateInput.UpdateExpression).NotTo(Equal(""))
				Expect(
					*dynamo.UpdateInput.ExpressionAttributeValues[":s1"].S).To(
					Equal(ScheduleStatusCanceled))
			})

			It("updates sets canceledAt", func() {
				Expect(dynamo.UpdateInput.UpdateExpression).NotTo(Equal(""))
				Expect(
					dynamo.UpdateInput.ExpressionAttributeValues[":ca"].N).NotTo(
					Equal(""))
			})

			It("returns success", func() {
				Expect(res).To(BeTrue())
			})

			It("does not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("fail", func() {
			Context("status is not idle", func() {
				var (
					res bool
					err error
				)

				BeforeEach(func() {
					dynamo.Error = awserr.NewRequestFailure(
						awserr.New(
							"ConditionalCheckFailedException",
							"NotFound",
							nil),
						400,
						"")

					res, err = db.Cancel(context.TODO(), id)
				})

				It("returns fail", func() {
					Expect(res).To(BeFalse())
				})

				It("does not return error", func() {
					Expect(err).To(BeNil())
				})

				AfterEach(func() {
					dynamo.Error = nil
				})
			})

			Context("update error", func() {
				var (
					res bool
					err error
				)

				BeforeEach(func() {
					dynamo.Error = awserr.New(
						"InternalError",
						"InternalError",
						nil)

					res, err = db.Cancel(context.TODO(), id)
				})

				It("returns fail", func() {
					Expect(res).To(BeFalse())
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})

				AfterEach(func() {
					dynamo.Error = nil
				})
			})
		})
	})

	Describe("Get", func() {
		Describe("success", func() {
			var (
				res *Schedule
				err error

				dueAt time.Time
				item  map[string]*dynamodb.AttributeValue
			)

			BeforeEach(func() {
				dueAt = time.Now().Add(-time.Hour * 24 * 7)

				item, _ = dynamodbattribute.MarshalMap(Schedule{
					ID:     id,
					DueAt:  dueAt,
					URL:    url,
					Method: method,
					Headers: map[string]string{
						"accept": accept,
					},
					Body:   body,
					Status: ScheduleStatusSucceeded,
				})

				dynamo.GetOutput = &dynamodb.GetItemOutput{Item: item}

				res, err = db.Get(context.TODO(), id)
			})

			It("reads table name from env", func() {
				Expect(*dynamo.GetInput.TableName).To(Equal(table))
			})

			It("sets given id to match schedule", func() {
				Expect(*dynamo.GetInput.Key["id"].S).To(Equal(id))
			})

			It("returns matching schedule", func() {
				Expect(res.ID).To(Equal(id))
				Expect(res.DueAt.Unix()).To(Equal(dueAt.Unix()))
				Expect(res.URL).To(Equal(url))
				Expect(res.Method).To(Equal(method))
				Expect(res.Headers["accept"]).To(Equal(accept))
				Expect(res.Body).To(Equal(body))
				Expect(res.Status).To(Equal(ScheduleStatusSucceeded))
			})

			It("does not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("fail", func() {
			Context("non existent schedule", func() {
				var (
					res *Schedule
					err error
				)

				BeforeEach(func() {
					dynamo.GetOutput = &dynamodb.GetItemOutput{}

					res, err = db.Get(context.TODO(), id)
				})

				It("does not return any schedule", func() {
					Expect(res).To(BeNil())
				})

				It("does not return error", func() {
					Expect(err).To(BeNil())
				})
			})

			Context("get error", func() {
				var (
					res *Schedule
					err error
				)

				BeforeEach(func() {
					dynamo.Error = awserr.New(
						"InternalError",
						"InternalError",
						nil)

					res, err = db.Get(context.TODO(), id)
				})

				It("does not return any schedule", func() {
					Expect(res).To(BeNil())
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})

				AfterEach(func() {
					dynamo.Error = nil
				})
			})

			Context("unmarshal error", func() {
				var (
					res *Schedule
					err error

					realUnmarshal unmarshal
				)

				BeforeEach(func() {
					realUnmarshal = unmarshalMap

					unmarshalMap = func(
						m map[string]*dynamodb.AttributeValue,
						out interface{}) error {

						return fmt.Errorf("unmarshal error")
					}

					dynamo.GetOutput = &dynamodb.GetItemOutput{
						Item: map[string]*dynamodb.AttributeValue{
							"id": {S: aws.String(id)},
						},
					}

					res, err = db.Get(context.TODO(), id)
				})

				It("does not return any schedule", func() {
					Expect(res).To(BeNil())
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})

				AfterEach(func() {
					unmarshalMap = realUnmarshal
				})
			})
		})
	})

	Describe("List", func() {
		Describe("empty conditions", func() {
			Describe("success", func() {
				Context("empty list", func() {
					var (
						res *List
						err error
					)

					BeforeEach(func() {
						dynamo.ScanOutput = &dynamodb.ScanOutput{}

						res, err = db.List(context.TODO(), ListInput{})
					})

					It("returns empty list", func() {
						Expect(res.Schedules).To(HaveLen(0))
					})

					It("does not return next key", func() {
						Expect(res.NextKey).To(BeNil())
					})

					It("does not return error", func() {
						Expect(err).To(BeNil())
					})
				})

				Context("non-empty list with large result", func() {
					var (
						res *List
						err error
					)

					BeforeEach(func() {
						items := make([]map[string]*dynamodb.AttributeValue, 5)

						for i := 0; i < len(items); i++ {
							item, _ := dynamodbattribute.MarshalMap(Schedule{
								ID: "12345",
							})

							items[i] = item
						}

						lastKey, _ := dynamodbattribute.MarshalMap(ListKey{
							ID:     "67890",
							DueAt:  time.Now().Add(-time.Hour * 100),
							Status: ScheduleStatusIdle,
						})

						dynamo.ScanOutput = &dynamodb.ScanOutput{
							Items:            items,
							LastEvaluatedKey: lastKey,
						}

						res, err = db.List(context.TODO(), ListInput{})
					})

					It("returns non-empty list", func() {
						Expect(res.Schedules).To(HaveLen(5))
					})

					It("returns next key", func() {
						Expect(res.NextKey).NotTo(BeNil())
					})

					It("does not return error", func() {
						Expect(err).To(BeNil())
					})
				})
			})

			Describe("fail", func() {
				Context("start key marshal error", func() {
					var (
						realMarshal marshal
						res         *List
						err         error
					)

					BeforeEach(func() {
						realMarshal = marshalStruct
						marshalStruct = func(in interface{}) (
							map[string]*dynamodb.AttributeValue,
							error) {

							return nil, fmt.Errorf("marshal error")
						}

						res, err = db.List(context.TODO(), ListInput{
							StartKey: &ListKey{
								ID:    "6568",
								DueAt: time.Now().Add(-time.Hour * 54),
							},
						})
					})

					It("does not return any result", func() {
						Expect(res).To(BeNil())
					})

					It("returns error", func() {
						Expect(err).NotTo(BeNil())
					})

					AfterEach(func() {
						marshalStruct = realMarshal
					})
				})

				Context("scan error", func() {
					var (
						res *List
						err error
					)

					BeforeEach(func() {
						dynamo.Error = awserr.New(
							"InternalError",
							"InternalError",
							nil)

						res, err = db.List(context.TODO(), ListInput{})
					})

					It("does not return any result", func() {
						Expect(res).To(BeNil())
					})

					It("returns error", func() {
						Expect(err).NotTo(BeNil())
					})

					AfterEach(func() {
						dynamo.Error = nil
					})
				})

				Context("next key unmarshal error", func() {
					var (
						realUnmarshal unmarshal
						res           *List
						err           error
					)

					BeforeEach(func() {
						realUnmarshal = unmarshalMap

						unmarshalMap = func(
							m map[string]*dynamodb.AttributeValue,
							out interface{}) error {

							return fmt.Errorf("unmarshal error")
						}
						dynamo.ScanOutput = &dynamodb.ScanOutput{
							Items: make([]map[string]*dynamodb.AttributeValue, 0),
							LastEvaluatedKey: map[string]*dynamodb.AttributeValue{
								"id": {S: aws.String("123")},
							},
						}

						res, err = db.List(context.TODO(), ListInput{})
					})

					It("does not return any result", func() {
						Expect(res).To(BeNil())
					})

					It("returns error", func() {
						Expect(err).NotTo(BeNil())
					})

					AfterEach(func() {
						unmarshalMap = realUnmarshal
					})
				})

				Context("result unmarshal error", func() {
					var (
						realUnmarshalList unmarshalList
						res               *List
						err               error
					)

					BeforeEach(func() {
						realUnmarshalList = unmarshalListOfMap
						unmarshalListOfMap = func(
							l []map[string]*dynamodb.AttributeValue,
							out interface{}) error {
							return fmt.Errorf("unmarshal list error")
						}

						dynamo.ScanOutput = &dynamodb.ScanOutput{
							Items: make([]map[string]*dynamodb.AttributeValue, 0),
						}

						res, err = db.List(context.TODO(), ListInput{})
					})

					It("does not return any result", func() {
						Expect(res).To(BeNil())
					})

					It("returns error", func() {
						Expect(err).NotTo(BeNil())
					})

					AfterEach(func() {
						unmarshalListOfMap = realUnmarshalList
					})
				})
			})
		})

		Describe("non empty conditions", func() {
			Describe("success", func() {
				Describe("with both status and dueAt", func() {
					var (
						res *List
						err error
					)

					BeforeEach(func() {
						items := make([]map[string]*dynamodb.AttributeValue, 5)

						for i := 0; i < len(items); i++ {
							item, _ := dynamodbattribute.MarshalMap(Schedule{
								ID: "12345",
							})

							items[i] = item
						}

						lastKey, _ := dynamodbattribute.MarshalMap(ListKey{
							ID:     "67890",
							DueAt:  time.Now().Add(-time.Hour * 100),
							Status: ScheduleStatusIdle,
						})

						dynamo.QueryOutput = &dynamodb.QueryOutput{
							Items:            items,
							LastEvaluatedKey: lastKey,
						}

						res, err = db.List(context.TODO(), ListInput{
							Status: ScheduleStatusIdle,
							DueAt: &DateRange{
								From: time.Now().Add(-time.Hour * 100),
								To:   time.Now().Add(-time.Hour * 50),
							},
						})
					})

					It("uses status and dueAt as condition", func() {
						Expect(dynamo.QueryInput.KeyConditionExpression).NotTo(
							BeNil())
						Expect(
							*dynamo.QueryInput.ExpressionAttributeValues[":s"].S).ToNot(
							Equal(""))
						Expect(dynamo.QueryInput.ExpressionAttributeValues[":d"]).To(
							BeNil())
						Expect(
							*dynamo.QueryInput.ExpressionAttributeValues[":da1"].N).ToNot(
							Equal(""))
						Expect(
							*dynamo.QueryInput.ExpressionAttributeValues[":da2"].N).ToNot(
							Equal(""))
					})

					It("uses ix_status_dueAt index", func() {
						Expect(*dynamo.QueryInput.IndexName).To(
							Equal("ix_status_dueAt"))
					})

					It("returns non-empty list", func() {
						Expect(res.Schedules).To(HaveLen(5))
					})

					It("returns next key", func() {
						Expect(res.NextKey).NotTo(BeNil())
					})

					It("does not return error", func() {
						Expect(err).To(BeNil())
					})
				})

				Describe("with only status", func() {
					var (
						res *List
						err error
					)

					BeforeEach(func() {
						items := make([]map[string]*dynamodb.AttributeValue, 5)

						for i := 0; i < len(items); i++ {
							item, _ := dynamodbattribute.MarshalMap(Schedule{
								ID: "12345",
							})

							items[i] = item
						}

						lastKey, _ := dynamodbattribute.MarshalMap(ListKey{
							ID:     "67890",
							DueAt:  time.Now().Add(-time.Hour * 100),
							Status: ScheduleStatusIdle,
						})

						dynamo.QueryOutput = &dynamodb.QueryOutput{
							Items:            items,
							LastEvaluatedKey: lastKey,
						}

						res, err = db.List(context.TODO(), ListInput{
							Status: ScheduleStatusQueued,
						})
					})

					It("uses only status as condition", func() {
						Expect(dynamo.QueryInput.KeyConditionExpression).NotTo(
							BeNil())
						Expect(
							*dynamo.QueryInput.ExpressionAttributeValues[":s"].S).ToNot(
							Equal(""))
						Expect(dynamo.QueryInput.ExpressionAttributeValues[":d"]).To(
							BeNil())
						Expect(dynamo.QueryInput.ExpressionAttributeValues[":da1"]).To(
							BeNil())
						Expect(dynamo.QueryInput.ExpressionAttributeValues[":da2"]).To(
							BeNil())
					})

					It("uses ix_status_dueAt index", func() {
						Expect(*dynamo.QueryInput.IndexName).To(
							Equal("ix_status_dueAt"))
					})

					It("returns non-empty list", func() {
						Expect(res.Schedules).To(HaveLen(5))
					})

					It("returns next key", func() {
						Expect(res.NextKey).NotTo(BeNil())
					})

					It("does not return error", func() {
						Expect(err).To(BeNil())
					})
				})

				Describe("with only dueAt", func() {
					var (
						res *List
						err error
					)

					BeforeEach(func() {
						items := make([]map[string]*dynamodb.AttributeValue, 5)

						for i := 0; i < len(items); i++ {
							item, _ := dynamodbattribute.MarshalMap(Schedule{
								ID: "12345",
							})

							items[i] = item
						}

						lastKey, _ := dynamodbattribute.MarshalMap(ListKey{
							ID:     "67890",
							DueAt:  time.Now().Add(-time.Hour * 100),
							Status: ScheduleStatusIdle,
						})

						dynamo.QueryOutput = &dynamodb.QueryOutput{
							Items:            items,
							LastEvaluatedKey: lastKey,
						}

						res, err = db.List(context.TODO(), ListInput{
							DueAt: &DateRange{
								From: time.Now().Add(-time.Hour * 20),
								To:   time.Now().Add(-time.Hour * 10),
							},
							StartKey: &ListKey{
								ID:    "6568",
								DueAt: time.Now().Add(-time.Hour * 5),
							},
						})
					})

					It("uses dummy and dueAt as condition", func() {
						Expect(dynamo.QueryInput.KeyConditionExpression).NotTo(
							BeNil())
						Expect(
							*dynamo.QueryInput.ExpressionAttributeValues[":d"].S).To(
							Equal(dummyValue))
						Expect(
							*dynamo.QueryInput.ExpressionAttributeValues[":da1"].N).ToNot(
							Equal(""))
						Expect(
							*dynamo.QueryInput.ExpressionAttributeValues[":da2"].N).ToNot(
							Equal(""))
						Expect(dynamo.QueryInput.ExpressionAttributeValues[":s"]).To(
							BeNil())
					})

					It("uses ix_dummy_dueAt index", func() {
						Expect(*dynamo.QueryInput.IndexName).To(
							Equal("ix_dummy_dueAt"))
					})

					It("auto included dummy value in start key", func() {
						Expect(
							*dynamo.QueryInput.ExclusiveStartKey["dummy"].S).To(
							Equal(dummyValue))
					})

					It("returns non-empty list", func() {
						Expect(res.Schedules).To(HaveLen(5))
					})

					It("returns next key", func() {
						Expect(res.NextKey).NotTo(BeNil())
					})

					It("does not return error", func() {
						Expect(err).To(BeNil())
					})
				})
			})
		})
	})
})

type fakeDynamoDB struct {
	dynamodbiface.DynamoDBAPI

	Error error

	PutInput    *dynamodb.PutItemInput
	UpdateInput *dynamodb.UpdateItemInput
	GetInput    *dynamodb.GetItemInput
	GetOutput   *dynamodb.GetItemOutput
	QueryInput  *dynamodb.QueryInput
	QueryOutput *dynamodb.QueryOutput
	ScanInput   *dynamodb.ScanInput
	ScanOutput  *dynamodb.ScanOutput
}

//goland:noinspection GoUnusedParameter
func (db *fakeDynamoDB) PutItemWithContext(
	ctx aws.Context,
	input *dynamodb.PutItemInput,
	options ...request.Option) (*dynamodb.PutItemOutput, error) {

	db.PutInput = input
	return nil, db.Error
}

//goland:noinspection GoUnusedParameter
func (db *fakeDynamoDB) UpdateItemWithContext(
	ctx aws.Context,
	input *dynamodb.UpdateItemInput,
	options ...request.Option) (*dynamodb.UpdateItemOutput, error) {

	db.UpdateInput = input
	return nil, db.Error
}

//goland:noinspection GoUnusedParameter
func (db *fakeDynamoDB) GetItemWithContext(
	ctx aws.Context,
	input *dynamodb.GetItemInput,
	options ...request.Option) (*dynamodb.GetItemOutput, error) {

	db.GetInput = input
	return db.GetOutput, db.Error
}

//goland:noinspection GoUnusedParameter
func (db *fakeDynamoDB) QueryWithContext(
	ctx aws.Context,
	input *dynamodb.QueryInput,
	options ...request.Option) (*dynamodb.QueryOutput, error) {

	db.QueryInput = input

	return db.QueryOutput, db.Error
}

//goland:noinspection GoUnusedParameter
func (db *fakeDynamoDB) ScanWithContext(
	ctx aws.Context,
	input *dynamodb.ScanInput,
	options ...request.Option) (*dynamodb.ScanOutput, error) {

	db.ScanInput = input

	return db.ScanOutput, db.Error
}
