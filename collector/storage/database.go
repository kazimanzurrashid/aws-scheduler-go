package storage

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"golang.org/x/sync/errgroup"
)

type Storage interface {
	Delete(context.Context) error
}

type Database struct {
	dynamodb dynamodbiface.DynamoDBAPI
}

func NewDatabase(dynamodb dynamodbiface.DynamoDBAPI) *Database {
	return &Database{dynamodb}
}

func (srv *Database) Delete(ctx context.Context) error {
	startKey := make(map[string]*dynamodb.AttributeValue)
	now := strconv.FormatInt(time.Now().Unix(), 10)
	g, _ := errgroup.WithContext(ctx)

	for {
		params := &dynamodb.QueryInput{
			TableName:              aws.String(tableName()),
			IndexName:              aws.String("ix_due_at"),
			KeyConditionExpression: aws.String("#d = :d AND #da < :da"),
			ProjectionExpression:   aws.String("#i"),
			ExpressionAttributeNames: map[string]*string{
				"#i":  aws.String("id"),
				"#d":  aws.String("dummy"),
				"#da": aws.String("dueAt"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":d":  {S: aws.String("s")},
				":da": {N: aws.String(now)},
			},
			ReturnConsumedCapacity: aws.String(dynamodb.ReturnConsumedCapacityNone),
		}

		if len(startKey) > 0 {
			params.ExclusiveStartKey = startKey
		}

		res, err := srv.dynamodb.QueryWithContext(ctx, params)

		if err != nil {
			return err
		}

		length := len(res.Items)

		if length == 0 {
			return nil
		}

		ids := make([]string, length)

		for index, row := range res.Items {
			ids[index] = *row["id"].S
		}

		for _, chunk := range chunkBy(ids, 25) {
			g.Go(func() error {
				writes := make([]*dynamodb.WriteRequest, len(chunk))

				for index, id := range chunk {
					write := &dynamodb.WriteRequest{
						DeleteRequest: &dynamodb.DeleteRequest{
							Key: map[string]*dynamodb.AttributeValue{
								"id": {S: aws.String(id)},
							},
						},
					}
					writes[index] = write
				}

				return srv.delete(ctx, writes)
			})
		}

		if len(res.LastEvaluatedKey) > 0 {
			startKey = res.LastEvaluatedKey
		} else {
			break
		}
	}

	return g.Wait()
}

func (srv *Database) delete(ctx context.Context, writes []*dynamodb.WriteRequest) error {
	table := tableName()
	params := &dynamodb.BatchWriteItemInput{
		RequestItems:                map[string][]*dynamodb.WriteRequest{},
		ReturnItemCollectionMetrics: aws.String(dynamodb.ReturnItemCollectionMetricsNone),
		ReturnConsumedCapacity:      aws.String(dynamodb.ReturnConsumedCapacityNone),
	}

	params.RequestItems[table] = writes

	res, err := srv.dynamodb.BatchWriteItemWithContext(ctx, params)

	if err != nil {
		return err
	}

	if len(res.UnprocessedItems) > 0 {
		ui, ok := res.UnprocessedItems[table]

		if ok && len(ui) > 0 {
			return srv.delete(ctx, ui)
		}
	}

	return nil
}

func chunkBy(items []string, size int) (chunks [][]string) {
	for size < len(items) {
		items, chunks = items[size:], append(chunks, items[0:size:size])
	}

	return append(chunks, items)
}

func tableName() string {
	return os.Getenv("SCHEDULER_TABLE_NAME")
}
