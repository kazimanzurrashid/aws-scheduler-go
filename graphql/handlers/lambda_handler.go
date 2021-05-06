package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

func lambdaStatus(code int, err error) (events.APIGatewayV2HTTPResponse, error) {
	buff, _ := marshalStruct(struct {
		Result    string `json:"result"`
		Timestamp string `json:"timestamp"`
	}{
		Result:    http.StatusText(code),
		Timestamp: time.Now().Format(time.RFC3339),
	})

	res := events.APIGatewayV2HTTPResponse{
		StatusCode: code,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"Content-Type":                "application/json;charset=utf-8",
		},
		Body: string(buff),
	}

	return res, err
}

func Lambda(
	ctx context.Context,
	req events.APIGatewayV2HTTPRequest) (
	events.APIGatewayV2HTTPResponse, error) {

	httpMethod := strings.ToUpper(req.RequestContext.HTTP.Method)

	if httpMethod == http.MethodGet &&
		req.RawPath == fmt.Sprintf("/%s/", req.RequestContext.Stage) {

		var buffer strings.Builder

		if err := playgroundTemplate.Execute(&buffer, struct {
			Endpoint string
		}{
			Endpoint: fmt.Sprintf(
				"https://%s/%s/graphql",
				req.RequestContext.DomainName,
				req.RequestContext.Stage),
		}); err != nil {
			return lambdaStatus(http.StatusInternalServerError, err)
		}

		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"Cache-Control": "private,max-age=31536000,immutable",
				"Content-Type":  "text/html;charset=utf-8",
			},
			Body: buffer.String(),
		}, nil
	}

	if httpMethod != http.MethodPost ||
		req.RawPath != fmt.Sprintf("/%s/graphql", req.RequestContext.Stage) {
		return lambdaStatus(http.StatusNotFound, nil)
	}

	body := strings.TrimSpace(req.Body)
	ret, statusCode := executeGraphQL(ctx, body)

	if statusCode != http.StatusOK {
		return lambdaStatus(http.StatusBadRequest, nil)
	}

	buff, err := marshalStruct(ret)

	if err != nil {
		return lambdaStatus(http.StatusInternalServerError, err)
	}

	res := events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"Content-Type":                "application/json;charset=utf-8",
		},
		Body: string(buff),
	}

	return res, nil
}
