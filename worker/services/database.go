package services

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"golang.org/x/sync/errgroup"
)

type marshalStorage func(in interface{}) (
	map[string]*dynamodb.AttributeValue, error)

func dynamoDBMarshal(
	in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(in)
}

var marshalStorageStruct marshalStorage = dynamoDBMarshal

type Storage interface {
	Update(context.Context, []*UpdateInput) error
}

type Database struct {
	dynamodb dynamodbiface.DynamoDBAPI
}

func NewDatabase(dynamodb dynamodbiface.DynamoDBAPI) *Database {
	return &Database{dynamodb}
}

func (srv *Database) Update(ctx context.Context, inputs []*UpdateInput) error {
	table := tableName()
	g, _ := errgroup.WithContext(ctx)

	for _, chunk := range chunkBy(inputs, 25) {
		localInputs := chunk

		g.Go(func() error {
			writes := make([]*dynamodb.WriteRequest, len(localInputs))

			for index, input := range localInputs {
				item, err := marshalStorageStruct(input)

				if err != nil {
					return err
				}

				write := &dynamodb.WriteRequest{
					PutRequest: &dynamodb.PutRequest{
						Item: item,
					},
				}

				writes[index] = write
			}

			return srv.update(ctx, table, writes)
		})
	}

	return g.Wait()
}

func (srv *Database) update(
	ctx context.Context,
	table string,
	writes []*dynamodb.WriteRequest) error {

	params := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{},
		ReturnItemCollectionMetrics: aws.String(
			dynamodb.ReturnItemCollectionMetricsNone),
		ReturnConsumedCapacity: aws.String(
			dynamodb.ReturnConsumedCapacityNone),
	}

	params.RequestItems[table] = writes

	res, err := srv.dynamodb.BatchWriteItemWithContext(ctx, params)

	if err != nil {
		return err
	}

	if len(res.UnprocessedItems) > 0 {
		ui, ok := res.UnprocessedItems[table]

		if ok && len(ui) > 0 {
			return srv.update(ctx, table, ui)
		}
	}

	return nil
}

func chunkBy(items []*UpdateInput, size int) (chunks [][]*UpdateInput) {
	for size < len(items) {
		items, chunks = items[size:], append(chunks, items[0:size:size])
	}

	return append(chunks, items)
}

func tableName() string {
	return os.Getenv("SCHEDULER_TABLE_NAME")
}
