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
	Create(context.Context, CreateInput) (string, error)

	Cancel(context.Context, string) (bool, error)

	Get(context.Context, string) (*Schedule, error)

	List(context.Context, ListInput) (*List, error)
}

type Database struct {
	dynamodb dynamodbiface.DynamoDBAPI
}

func NewDatabase(dynamodb dynamodbiface.DynamoDBAPI) *Database {
	return &Database{dynamodb}
}

const dummyValue = "-"

type (
	idGenerate func() (string, error)
	marshal    func(interface{}) (map[string]*dynamodb.AttributeValue, error)
	unmarshal  func(
		map[string]*dynamodb.AttributeValue,
		interface{}) error
	unmarshalList func(
		[]map[string]*dynamodb.AttributeValue,
		interface{}) error
)

func nanoIDGenerate() (string, error) {
	return gonanoid.Generate(
		"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		16)
}

func dynamoDBMarshalStruct(
	in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(in)
}

func dynamoDBUnmarshalMap(
	m map[string]*dynamodb.AttributeValue, out interface{}) error {
	return dynamodbattribute.UnmarshalMap(m, out)
}

func dynamoDBUnmarshalListOfMap(
	l []map[string]*dynamodb.AttributeValue, out interface{}) error {
	return dynamodbattribute.UnmarshalListOfMaps(l, out)
}

var (
	generateID         idGenerate    = nanoIDGenerate
	marshalStruct      marshal       = dynamoDBMarshalStruct
	unmarshalMap       unmarshal     = dynamoDBUnmarshalMap
	unmarshalListOfMap unmarshalList = dynamoDBUnmarshalListOfMap
)

func (srv *Database) Create(
	ctx context.Context,
	input CreateInput) (string, error) {

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
	item["dummy"] = &dynamodb.AttributeValue{S: aws.String(dummyValue)}
	item["createdAt"] = &dynamodb.AttributeValue{
		N: aws.String(strconv.FormatInt(time.Now().Unix(), 10)),
	}

	params := &dynamodb.PutItemInput{
		TableName: aws.String(tableName()),
		Item:      item,
		ReturnConsumedCapacity: aws.String(
			dynamodb.ReturnConsumedCapacityNone),
		ReturnItemCollectionMetrics: aws.String(
			dynamodb.ReturnItemCollectionMetricsNone),
		ReturnValues: aws.String(dynamodb.ReturnValueNone),
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
		ReturnConsumedCapacity: aws.String(
			dynamodb.ReturnConsumedCapacityNone),
		ReturnValues: aws.String(
			dynamodb.ReturnValueNone),
	}

	if _, err := srv.dynamodb.UpdateItemWithContext(ctx, params); err != nil {
		if ccf, ok := err.(awserr.RequestFailure); ok &&
			ccf.Code() == "ConditionalCheckFailedException" {
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

	if err = unmarshalMap(res.Item, &s); err != nil {
		return nil, err
	}

	return &s, nil
}

func (srv *Database) List(ctx context.Context, input ListInput) (*List, error) {
	params := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName()),
		IndexName:                 aws.String("ix_dummy_dueAt"),
		Limit:                     aws.Int64(input.Limit),
		ReturnConsumedCapacity:    aws.String(dynamodb.ReturnConsumedCapacityNone),
		ExpressionAttributeNames:  map[string]*string{},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{},
		ScanIndexForward:          aws.Bool(true),
	}

	if input.Status != "" {
		params.IndexName = aws.String("ix_status_dueAt")
		params.ExpressionAttributeNames["#s"] = aws.String("status")
		params.ExpressionAttributeValues[":s"] = &dynamodb.AttributeValue{
			S: aws.String(input.Status),
		}

		if input.DueAt == nil {
			params.KeyConditionExpression = aws.String("#s = :s")
		} else {
			params.ExpressionAttributeNames["#da"] = aws.String("dueAt")
			params.ExpressionAttributeValues[":da1"] = &dynamodb.AttributeValue{
				N: aws.String(strconv.FormatInt(input.DueAt.From.Unix(), 10)),
			}
			params.ExpressionAttributeValues[":da2"] = &dynamodb.AttributeValue{
				N: aws.String(strconv.FormatInt(input.DueAt.To.Unix(), 10)),
			}
			params.KeyConditionExpression = aws.String(
				"#s = :s AND #da BETWEEN :da1 AND :da2")
		}
	} else if input.DueAt != nil {
		params.ExpressionAttributeNames["#d"] = aws.String("dummy")
		params.ExpressionAttributeValues[":d"] = &dynamodb.AttributeValue{
			S: aws.String(dummyValue),
		}
		params.ExpressionAttributeNames["#da"] = aws.String("dueAt")
		params.ExpressionAttributeValues[":da1"] = &dynamodb.AttributeValue{
			N: aws.String(strconv.FormatInt(input.DueAt.From.Unix(), 10)),
		}
		params.ExpressionAttributeValues[":da2"] = &dynamodb.AttributeValue{
			N: aws.String(strconv.FormatInt(input.DueAt.To.Unix(), 10)),
		}
		params.KeyConditionExpression = aws.String(
			"#d = :d AND #da BETWEEN :da1 AND :da2")
	} else {
		params.KeyConditionExpression = aws.String("#d = :d")
		params.ExpressionAttributeNames["#d"] = aws.String("dummy")
		params.ExpressionAttributeValues[":d"] = &dynamodb.AttributeValue{
			S: aws.String(dummyValue),
		}
	}

	if input.StartKey != nil {
		startKey, err := marshalStruct(input.StartKey)

		if err != nil {
			return nil, err
		}

		if *params.IndexName == "ix_dummy_dueAt" {
			startKey["dummy"] = &dynamodb.AttributeValue{
				S: aws.String(dummyValue),
			}
		}

		params.ExclusiveStartKey = startKey
	}

	res, err := srv.dynamodb.QueryWithContext(ctx, params)

	if err != nil {
		return nil, err
	}

	var nextKey *ListKey

	if len(res.LastEvaluatedKey) > 0 {
		var nk ListKey

		err = unmarshalMap(res.LastEvaluatedKey, &nk)

		if err != nil {
			return nil, err
		}

		nextKey = &nk
	}

	var schedules []*Schedule

	if err = unmarshalListOfMap(res.Items, &schedules); err != nil {
		return nil, err
	}

	return &List{Schedules: schedules, NextKey: nextKey}, nil
}

func tableName() string {
	return os.Getenv("SCHEDULER_TABLE_NAME")
}
