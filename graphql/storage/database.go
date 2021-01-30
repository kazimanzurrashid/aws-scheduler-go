package storage

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/matoous/go-nanoid"
)

type Storage interface {
	Create(context.Context, *CreateInput) (string, error)

	Cancel(context.Context, string) (bool, error)

	Get(context.Context, string) (*Schedule, error)
}

type Database struct {
	dynamodb dynamodbiface.DynamoDBAPI
}

func NewDatabase(dynamodb dynamodbiface.DynamoDBAPI) *Database {
	return &Database{dynamodb}
}

type (
	idGenerate func() (string, error)
	marshal    func(in interface{}) (map[string]*dynamodb.AttributeValue, error)
	unmarshal  func(
		m map[string]*dynamodb.AttributeValue,
		out interface{}) error
)

func nanoIDGenerate() (string, error) {
	return gonanoid.Generate(
		"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		16)
}

func dynamoDBMarshal(
	in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(in)
}

func dynamoDBUnmarshal(
	m map[string]*dynamodb.AttributeValue, out interface{}) error {
	return dynamodbattribute.UnmarshalMap(m, out)
}

var (
	generateID    idGenerate = nanoIDGenerate
	marshalStruct marshal    = dynamoDBMarshal
	unmarshalMap  unmarshal  = dynamoDBUnmarshal
)

func (srv *Database) Create(
	ctx context.Context,
	input *CreateInput) (string, error) {

	id, err := generateID()

	if err != nil {
		return "", err
	}

	item, err := marshalStruct(input)

	if err != nil {
		return "", err
	}

	item["id"] = &dynamodb.AttributeValue{S: aws.String(id)}
	item["status"] = &dynamodb.AttributeValue{
		S: aws.String(ScheduleStatusIdle),
	}
	item["createdAt"] = &dynamodb.AttributeValue{
		N: aws.String(strconv.FormatInt(time.Now().Unix(), 10)),
	}

	params := &dynamodb.PutItemInput{
		TableName:                   aws.String(tableName()),
		Item:                        item,
		ReturnConsumedCapacity:      aws.String(
			dynamodb.ReturnConsumedCapacityNone),
		ReturnItemCollectionMetrics: aws.String(
			dynamodb.ReturnItemCollectionMetricsNone),
		ReturnValues:                aws.String(dynamodb.ReturnValueNone),
	}

	if _, err = srv.dynamodb.PutItemWithContext(ctx, params); err != nil {
		return "", err
	}

	return id, nil
}

func (srv *Database) Cancel(ctx context.Context, id string) (bool, error) {
	params := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName()),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
		UpdateExpression:    aws.String("SET #s = :s1, #ca = :ca"),
		ConditionExpression: aws.String("#s = :s2"),
		ExpressionAttributeNames: map[string]*string{
			"#s":  aws.String("status"),
			"#ca": aws.String("canceledAt"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s1": {S: aws.String(ScheduleStatusCanceled)},
			":s2": {S: aws.String(ScheduleStatusIdle)},
			":ca": {N: aws.String(strconv.FormatInt(time.Now().Unix(), 10))},
		},
		ReturnItemCollectionMetrics: aws.String(
			dynamodb.ReturnItemCollectionMetricsNone),
		ReturnConsumedCapacity:      aws.String(
			dynamodb.ReturnConsumedCapacityNone),
		ReturnValues:                aws.String(
			dynamodb.ReturnValueNone),
	}

	if _, err := srv.dynamodb.UpdateItemWithContext(ctx, params); err != nil {
		if ccf, ok := err.(awserr.RequestFailure);
		ok && ccf.Code() == "ConditionalCheckFailedException" {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (srv *Database) Get(ctx context.Context, id string) (*Schedule, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(tableName()),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
		ReturnConsumedCapacity: aws.String(dynamodb.ReturnConsumedCapacityNone),
	}

	res, err := srv.dynamodb.GetItemWithContext(ctx, params)

	if err != nil {
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, nil
	}

	var s Schedule

	if err := unmarshalMap(res.Item, &s); err != nil {
		return nil, err
	}

	return &s, nil
}

func tableName() string {
	return os.Getenv("SCHEDULER_TABLE_NAME")
}
