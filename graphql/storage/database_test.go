package storage

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/stretchr/testify/assert"
)

type fakeDynamoDB struct {
	dynamodbiface.DynamoDBAPI

	ReturnError error

	PutInput    *dynamodb.PutItemInput
	UpdateInput *dynamodb.UpdateItemInput
	GetInput    *dynamodb.GetItemInput
	GetOutput   *dynamodb.GetItemOutput
}

//goland:noinspection GoUnusedParameter
func (db *fakeDynamoDB) PutItemWithContext(
	ctx aws.Context,
	input *dynamodb.PutItemInput,
	options ...request.Option) (*dynamodb.PutItemOutput, error) {

	db.PutInput = input
	return nil, db.ReturnError
}

//goland:noinspection GoUnusedParameter
func (db *fakeDynamoDB) UpdateItemWithContext(
	ctx aws.Context,
	input *dynamodb.UpdateItemInput,
	options ...request.Option) (*dynamodb.UpdateItemOutput, error) {

	db.UpdateInput = input
	return nil, db.ReturnError
}

//goland:noinspection GoUnusedParameter
func (db *fakeDynamoDB) GetItemWithContext(
	ctx aws.Context,
	input *dynamodb.GetItemInput,
	options ...request.Option) (*dynamodb.GetItemOutput, error) {

	db.GetInput = input
	return db.GetOutput, db.ReturnError
}

const (
	table  = "scheduler_v1"
	id     = "1234567890"
	url    = "https://foo.bar/do"
	method = "POST"
	accept = "application/json"
	body   = "{ \"foo\": \"bar\" }"
)

var ctx = context.TODO()

func setTableName(t *testing.T)  {
	err := os.Setenv("SCHEDULER_TABLE_NAME", table)
	if err != nil {
		t.Error(err)
	}
}

func Test_Database_NewDatabase_Success(t *testing.T) {
	db := NewDatabase(&fakeDynamoDB{})

	assert.NotNil(t, db)
}

func Test_Database_Create_Success(t *testing.T) {
	dueAt := time.Now().Add(time.Minute * 1)
	headers := map[string]string{
		"accept": accept,
	}

	setTableName(t)

	fake := fakeDynamoDB{}
	db := NewDatabase(&fake)

	res, err := db.Create(ctx, &CreateInput{
		DueAt:   dueAt,
		URL:     url,
		Method:  method,
		Headers: headers,
		Body:    body,
	})

	assert.NotEqual(t, res, "")
	assert.Nil(t, err)
	assert.Equal(t, table, *fake.PutInput.TableName)
	assert.Equal(t, res, *fake.PutInput.Item["id"].S)
	assert.Equal(
		t,
		strconv.FormatInt(dueAt.Unix(), 10),
		*fake.PutInput.Item["dueAt"].N)
	assert.Equal(t, ScheduleStatusIdle, *fake.PutInput.Item["status"].S)
	assert.Equal(t, url, *fake.PutInput.Item["url"].S)
	assert.Equal(t, method, *fake.PutInput.Item["method"].S)
	assert.Equal(t, accept, *fake.PutInput.Item["headers"].M["accept"].S)
	assert.Equal(t, body, *fake.PutInput.Item["body"].S)
	assert.NotEqual(t, "", *fake.PutInput.Item["createdAt"].N)
}

func Test_Database_Create_Fail_ID_Generate_Error(t *testing.T) {
	realIdGen := generateID
	generateID = func() (string, error) {
		return "", fmt.Errorf("id generate error")
	}

	fake := fakeDynamoDB{}
	db := NewDatabase(&fake)

	res, err := db.Create(ctx, &CreateInput{})

	assert.Equal(t, res, "")
	assert.NotNil(t, err)

	generateID = realIdGen
}

func Test_Database_Create_Fail_Input_Marshal_Error(t *testing.T) {
	realMarshal := marshalStruct

	marshalStruct = func(in interface{}) (
		map[string]*dynamodb.AttributeValue,
		error) {

		return nil, fmt.Errorf("marshal error")
	}

	fake := fakeDynamoDB{}
	db := NewDatabase(&fake)

	res, err := db.Create(ctx, &CreateInput{})

	assert.Equal(t, res, "")
	assert.NotNil(t, err)

	marshalStruct = realMarshal
}

func Test_Database_Create_Fail_Internal_Error(t *testing.T) {
	dueAt := time.Now().Add(time.Minute * 1)
	headers := map[string]string{
		"accept": accept,
	}

	fake := fakeDynamoDB{
		ReturnError: awserr.New("InternalError", "InternalError", nil),
	}
	db := NewDatabase(&fake)

	res, err := db.Create(ctx, &CreateInput{
		DueAt:   dueAt,
		URL:     url,
		Method:  method,
		Headers: headers,
		Body:    body,
	})

	assert.Equal(t, res, "")
	assert.NotNil(t, err)
}

func Test_Database_Cancel_Success(t *testing.T) {
	setTableName(t)

	fake := fakeDynamoDB{}
	db := NewDatabase(&fake)

	res, err := db.Cancel(ctx, id)

	assert.True(t, res)
	assert.Nil(t, err)
	assert.Equal(t, table, *fake.UpdateInput.TableName)
	assert.Equal(t, id, *fake.UpdateInput.Key["id"].S)
	assert.NotEqual(t, "", *fake.UpdateInput.UpdateExpression)
	assert.NotEqual(t, "", *fake.UpdateInput.ConditionExpression)
	assert.Equal(
		t,
		ScheduleStatusCanceled,
		*fake.UpdateInput.ExpressionAttributeValues[":s1"].S)
	assert.Equal(
		t,
		ScheduleStatusIdle,
		*fake.UpdateInput.ExpressionAttributeValues[":s2"].S)
	assert.NotEqual(
		t,
		"",
		*fake.UpdateInput.ExpressionAttributeValues[":ca"].N)
}

func Test_Database_Cancel_Success_Conditional_Check_Fail(t *testing.T) {
	fake := fakeDynamoDB{
		ReturnError: awserr.NewRequestFailure(
			awserr.New(
				"ConditionalCheckFailedException",
				"NotFound",
				nil),
			400,
			""),
	}
	db := NewDatabase(&fake)

	res, err := db.Cancel(ctx, id)

	assert.False(t, res)
	assert.Nil(t, err)
}

func Test_Database_Cancel_Fail_Internal_Error(t *testing.T) {
	fake := fakeDynamoDB{
		ReturnError: awserr.New("InternalError", "InternalError", nil),
	}
	db := NewDatabase(&fake)

	res, err := db.Cancel(ctx, id)

	assert.False(t, res)
	assert.NotNil(t, err)
}

func Test_Database_Get_Success_Existent_Record(t *testing.T) {
	dueAt := time.Now().Add(-time.Hour * 24 * 7)

	setTableName(t)

	item, err := dynamodbattribute.MarshalMap(Schedule{
		ID:     id,
		DueAt:  dueAt,
		URL:    url,
		Method: method,
		Headers: map[string]string{
			"accept": accept,
		},
		Body: body,
		Status: ScheduleStatusSucceeded,
		StartedAt: dueAt.Add(time.Minute * 2),
		CompletedAt: time.Now().Add(time.Minute * 3),
		CreatedAt: dueAt.Add(-time.Hour * 24 * 3),
	})

	if err != nil {
		t.Error(err)
	}

	fake := fakeDynamoDB{
		GetOutput: &dynamodb.GetItemOutput{Item: item},
	}
	db := NewDatabase(&fake)

	res, err := db.Get(ctx, id)

	assert.NotNil(t, res)
	assert.Nil(t, err)

	assert.Equal(t, id, res.ID)
	assert.Equal(t, dueAt.Unix(), res.DueAt.Unix())
	assert.Equal(t, url, res.URL)
	assert.Equal(t, method, res.Method)
	assert.Equal(t, accept, res.Headers["accept"])
	assert.Equal(t, body, res.Body)
	assert.Equal(t, table, *fake.GetInput.TableName)
	assert.Equal(t, id, *fake.GetInput.Key["id"].S)
}

func Test_Database_Get_Success_NonExistent_Record(t *testing.T) {
	fake := fakeDynamoDB{
		GetOutput: &dynamodb.GetItemOutput{},
	}
	db := NewDatabase(&fake)

	res, err := db.Get(ctx, id)

	assert.Nil(t, res)
	assert.Nil(t, err)
}

func Test_Database_Get_Fail_Internal_Error(t *testing.T) {
	fake := fakeDynamoDB{
		ReturnError: awserr.New("InternalError", "InternalError", nil),
	}
	db := NewDatabase(&fake)

	res, err := db.Get(ctx, id)

	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func Test_Database_Get_Fail_Output_Unmarshal_Error(t *testing.T) {
	realUnmarshal := unmarshalMap
	unmarshalMap = func(
		m map[string]*dynamodb.AttributeValue,
		out interface{}) error {

		return fmt.Errorf("unmarshal error")
	}

	fake := fakeDynamoDB{
		GetOutput: &dynamodb.GetItemOutput{
			Item: map[string]*dynamodb.AttributeValue{
				"id": {S: aws.String(id)},
		}},
	}
	db := NewDatabase(&fake)

	res, err := db.Get(ctx, id)

	assert.Nil(t, res)
	assert.NotNil(t, err)

	unmarshalMap = realUnmarshal
}
