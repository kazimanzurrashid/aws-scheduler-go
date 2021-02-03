package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/graphql-go/graphql"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/xray"

	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/api"
	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/storage"
)

type request struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

var schema graphql.Schema

var headers = map[string]string{
	"Content-Type": "application/json;charset=utf-8",
}

type (
	marshal   func(interface{}) ([]byte, error)
	unmarshal func([]byte, interface{}) error
)

func jsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func jsonUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

var (
	marshalStruct   marshal   = jsonMarshal
	unmarshalStruct unmarshal = jsonUnmarshal
)

func status(code int, err error) (events.APIGatewayV2HTTPResponse, error) {
	buff, _ := marshalStruct(struct {
		Result    string `json:"result"`
		Timestamp string `json:"timestamp"`
	}{
		Result:    http.StatusText(code),
		Timestamp: time.Now().Format(time.RFC3339),
	})

	res := events.APIGatewayV2HTTPResponse{
		StatusCode: code,
		Headers:    headers,
		Body:       string(buff),
	}

	return res, err
}

func handler(
	ctx context.Context,
	req events.APIGatewayV2HTTPRequest) (
	events.APIGatewayV2HTTPResponse, error) {

	body := strings.TrimSpace(req.Body)

	if body == "" {
		return status(http.StatusBadRequest, nil)
	}

	var ret interface{}

	if strings.HasPrefix(body, "{") && strings.HasSuffix(body, "}") {
		var payload request

		if err := unmarshalStruct([]byte(body), &payload); err != nil {
			return status(http.StatusInternalServerError, err)
		}

		ret = graphql.Do(graphql.Params{
			Context:        ctx,
			Schema:         schema,
			RequestString:  payload.Query,
			OperationName:  payload.OperationName,
			VariableValues: payload.Variables,
		})
	} else if strings.HasPrefix(body, "[") && strings.HasSuffix(body, "]") {
		var payloads []request

		if err := unmarshalStruct([]byte(body), &payloads); err != nil {
			return status(http.StatusInternalServerError, err)
		}

		rets := make([]*graphql.Result, len(payloads))
		var wg sync.WaitGroup

		for i, p := range payloads {
			go func(payload request, index int) {
				defer wg.Done()

				out := graphql.Do(graphql.Params{
					Context:        ctx,
					Schema:         schema,
					RequestString:  payload.Query,
					OperationName:  payload.OperationName,
					VariableValues: payload.Variables,
				})

				rets[index] = out
			}(p, i)
			wg.Add(1)
		}

		wg.Wait()

		ret = rets
	} else {
		return status(http.StatusBadRequest, nil)
	}

	buff, err := marshalStruct(ret)

	if err != nil {
		return status(http.StatusInternalServerError, err)
	}

	res := events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       string(buff),
	}

	return res, nil
}

func init() {
	ses := session.Must(session.NewSession())

	ddbc := dynamodb.New(ses)
	xray.AWS(ddbc.Client)

	database := storage.NewDatabase(ddbc)

	f := api.NewFactory(database)
	s, err := f.Schema()

	if err != nil {
		log.Fatalf("schema create error: %v", err)
		return
	}

	schema = s
}

func main() {
	lambda.Start(handler)
}
