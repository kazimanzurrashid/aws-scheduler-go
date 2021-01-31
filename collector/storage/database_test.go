package storage

import (
	"container/list"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
)

type fakeDynamoDB struct {
	dynamodbiface.DynamoDBAPI

	QueryReturnError      error
	BatchWriteReturnError error

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
		return nil, db.QueryReturnError
	}

	db.queryOutputs.Remove(output)

	return output.Value.(*dynamodb.QueryOutput), db.QueryReturnError
}

//goland:noinspection GoUnusedParameter
func (db *fakeDynamoDB) BatchWriteItemWithContext(
	ctx aws.Context,
	input *dynamodb.BatchWriteItemInput,
	options ...request.Option) (*dynamodb.BatchWriteItemOutput, error) {

	db.batchWriteInputs.PushBack(input)

	output := db.batchWriteOutputs.Front()

	if output == nil {
		return nil, db.BatchWriteReturnError
	}

	db.batchWriteOutputs.Remove(output)

	return output.Value.(*dynamodb.BatchWriteItemOutput), db.BatchWriteReturnError
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

const (
	table = "scheduler_v1"
	id    = "1234567890"
)

var ctx = context.TODO()

func setTableName(t *testing.T) {
	err := os.Setenv("SCHEDULER_TABLE_NAME", table)
	if err != nil {
		t.Error(err)
	}
}

func Test_Database_Update_Success(t *testing.T) {
	setTableName(t)

	fake := fakeDynamoDB{}

	fake.PushQueryOutput(&dynamodb.QueryOutput{
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

	fake.PushQueryOutput(&dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			{
				"id":     {S: aws.String(id)},
				"status": {S: aws.String(scheduleStatusIdle)},
			},
		},
	})

	fake.PushBatchWriteOutput(&dynamodb.BatchWriteItemOutput{})
	fake.PushBatchWriteOutput(&dynamodb.BatchWriteItemOutput{
		UnprocessedItems: map[string][]*dynamodb.WriteRequest{
			table: {
				&dynamodb.WriteRequest{
					PutRequest: &dynamodb.PutRequest{
						Item: map[string]*dynamodb.AttributeValue{
							"id":     {S: aws.String(id)},
							"status": {S: aws.String(scheduleStatusQueued)},
						},
					},
				},
			},
		},
	})
	fake.PushBatchWriteOutput(&dynamodb.BatchWriteItemOutput{})

	db := NewDatabase(&fake)

	err := db.Update(ctx)

	queryInputs := []*dynamodb.QueryInput{
		fake.PullQueryInput(),
		fake.PullQueryInput(),
	}

	batchWriteInputs := []*dynamodb.BatchWriteItemInput{
		fake.PullBatchWriteInput(),
		fake.PullBatchWriteInput(),
		fake.PullBatchWriteInput(),
	}

	assert.Nil(t, err)

	for _, queryInput := range queryInputs {
		assert.Equal(t, table, *queryInput.TableName)
		assert.NotEqual(t, "", *queryInput.IndexName)
		assert.NotNil(t, queryInput.ExpressionAttributeValues[":s"])
		assert.NotNil(t, queryInput.ExpressionAttributeValues[":da"])
	}

	for _, batchWriteInput := range batchWriteInputs {
		assert.NotNil(t, batchWriteInput.RequestItems[table])

		for _, r := range batchWriteInput.RequestItems[table] {
			assert.Equal(t, scheduleStatusQueued, *r.PutRequest.Item["status"].S)
		}
	}
}

func Test_Database_Update_Success_Empty_Records(t *testing.T) {
	setTableName(t)

	fake := fakeDynamoDB{}
	fake.PushQueryOutput(&dynamodb.QueryOutput{})

	db := NewDatabase(&fake)

	err := db.Update(ctx)

	assert.Nil(t, err)
}

func Test_Database_Update_Fail_Query_Error(t *testing.T) {
	setTableName(t)

	fake := fakeDynamoDB{
		QueryReturnError: fmt.Errorf("query error"),
	}

	db := NewDatabase(&fake)

	err := db.Update(ctx)

	assert.NotNil(t, err)
}

func Test_Database_Update_Fail_Batch_Write_Error(t *testing.T) {
	setTableName(t)

	fake := fakeDynamoDB{
		BatchWriteReturnError: fmt.Errorf("batch write error"),
	}

	fake.PushQueryOutput(&dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			{
				"id": {S: aws.String(id)},
			},
		},
	})

	db := NewDatabase(&fake)

	err := db.Update(ctx)

	assert.NotNil(t, err)
}
