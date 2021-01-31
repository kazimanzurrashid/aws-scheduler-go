package services

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

	BatchWriteReturnError error

	batchWriteInputs  list.List
	batchWriteOutputs list.List
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

	fake.PushBatchWriteOutput(&dynamodb.BatchWriteItemOutput{
		UnprocessedItems: map[string][]*dynamodb.WriteRequest{
			table: {
				&dynamodb.WriteRequest{
					PutRequest: &dynamodb.PutRequest{
						Item: map[string]*dynamodb.AttributeValue{
							"id":     {S: aws.String(id)},
						},
					},
				},
			},
		},
	})
	fake.PushBatchWriteOutput(&dynamodb.BatchWriteItemOutput{})

	db := NewDatabase(&fake)

	err := db.Update(ctx, []*UpdateInput{
		{
			ID: id,
		},
	})

	batchWriteInputs := []*dynamodb.BatchWriteItemInput{
		fake.PullBatchWriteInput(),
		fake.PullBatchWriteInput(),
	}

	assert.Nil(t, err)

	for _, batchWriteInput := range batchWriteInputs {
		assert.NotNil(t, batchWriteInput.RequestItems[table])

		for _, r := range batchWriteInput.RequestItems[table] {
			assert.Equal(t, id, *r.PutRequest.Item["id"].S)
		}
	}
}

func Test_Database_Update_Fail_Marshal_Error(t *testing.T) {
	setTableName(t)
	realMarshal := marshalStorageStruct

	marshalStorageStruct = func(
		in interface{}) (map[string]*dynamodb.AttributeValue, error) {

		return nil, fmt.Errorf("marshal error")
	}

	fake := fakeDynamoDB{}

	db := NewDatabase(&fake)

	err := db.Update(ctx, []*UpdateInput{
		{
			ID: id,
		},
	})

	assert.NotNil(t, err)

	marshalStorageStruct = realMarshal
}

func Test_Database_Update_Fail_Batch_Write_Error(t *testing.T) {
	setTableName(t)

	fake := fakeDynamoDB{
		BatchWriteReturnError: fmt.Errorf("batch write error"),
	}

	db := NewDatabase(&fake)

	err := db.Update(ctx, []*UpdateInput{
		{
			ID: id,
		},
	})

	assert.NotNil(t, err)
}
