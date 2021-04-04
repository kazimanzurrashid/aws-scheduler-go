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

const (
	scheduleStatusIdle   = "IDLE"
	scheduleStatusQueued = "QUEUED"
)

type Storage interface {
	Update(context.Context) error
}

type Database struct {
	dynamodb dynamodbiface.DynamoDBAPI
}

func NewDatabase(dynamodb dynamodbiface.DynamoDBAPI) *Database {
	return &Database{dynamodb}
}

func (srv *Database) Update(ctx context.Context) error {
	table := tableName()
	startKey := make(map[string]*dynamodb.AttributeValue)
	now := strconv.FormatInt(time.Now().Unix(), 10)

	g, _ := errgroup.WithContext(ctx)

	for {
		params := &dynamodb.QueryInput{
			TableName:              aws.String(table),
			IndexName:              aws.String("ix_status_dueAt"),
			KeyConditionExpression: aws.String("#s = :s AND #da <= :da"),
			ExpressionAttributeNames: map[string]*string{
				"#s":  aws.String("status"),
				"#da": aws.String("dueAt"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":s":  {S: aws.String(scheduleStatusIdle)},
				":da": {N: aws.String(now)},
			},
			ReturnConsumedCapacity: aws.String(
				dynamodb.ReturnConsumedCapacityNone),
		}

		if len(startKey) > 0 {
			params.ExclusiveStartKey = startKey
		}

		res, err := srv.dynamodb.QueryWithContext(ctx, params)

		if err != nil {
			return err
		}

		if len(res.Items) == 0 {
			return nil
		}

		for _, chunk := range chunkBy(res.Items, 25) {
			localItems := chunk

			g.Go(func() error {
				writes := make([]*dynamodb.WriteRequest, len(localItems))

				for index, item := range localItems {
					item["status"] = &dynamodb.AttributeValue{
						S: aws.String(scheduleStatusQueued),
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

		if len(res.LastEvaluatedKey) > 0 {
			startKey = res.LastEvaluatedKey
		} else {
			break
		}
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

func chunkBy(
	items []map[string]*dynamodb.AttributeValue,
	size int) (chunks [][]map[string]*dynamodb.AttributeValue) {

	for size < len(items) {
		items, chunks = items[size:], append(chunks, items[0:size:size])
	}

	return append(chunks, items)
}

func tableName() string {
	return os.Getenv("SCHEDULER_TABLE_NAME")
}
