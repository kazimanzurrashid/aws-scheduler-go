package services

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type marshalClient func(interface{}) string

func jsonMarshal(v interface{}) string {
	buff, err := json.Marshal(v)

	if err != nil {
		return ""
	}

	return string(buff)
}

var marshalClientStruct marshalClient = jsonMarshal

type result struct {
	Error      string            `json:"error,omitempty"`
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}

type Client interface {
	Request(context.Context, *RequestInput) *ResponseOutput
}

type HttpClient struct {
	http *http.Client
}

func NewHttpClient(c *http.Client) *HttpClient {
	return &HttpClient{c}
}

func (hc *HttpClient) Request(
	ctx context.Context,
	ri *RequestInput) *ResponseOutput {

	req, err := http.NewRequestWithContext(
		ctx,
		ri.Method,
		ri.URL,
		bytes.NewBufferString(ri.Body))

	if err != nil {
		return &ResponseOutput{
			Status: ScheduleStatusFailed,
			Result: marshalClientStruct(result{Error: err.Error()}),
		}
	}

	for k, v := range ri.Headers {
		req.Header.Set(k, v)
	}

	res, err := hc.http.Do(req)

	if err != nil {
		return &ResponseOutput{
			Status: ScheduleStatusFailed,
			Result: marshalClientStruct(result{Error: err.Error()}),
		}
	}

	defer func() {
		_ = res.Body.Close()
	}()

	headers := make(map[string]string)

	for k, v := range res.Header {
		headers[k] = strings.Join(v, ";")
	}

	body, _ := ioutil.ReadAll(res.Body)

	var status string

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		status = ScheduleStatusSucceeded
	} else {
		status = ScheduleStatusFailed
	}

	return &ResponseOutput{
		Status: status,
		Result: marshalClientStruct(result{
			StatusCode: res.StatusCode,
			Headers:    headers,
			Body:       string(body),
		}),
	}
}
