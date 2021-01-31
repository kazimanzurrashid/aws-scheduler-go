package services

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeTransport struct {
	Request  *http.Request
	Response *http.Response
	Error    error
}

func (ft *fakeTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	ft.Request = r

	return ft.Response, ft.Error
}

const (
	url      = "https://foo.bar/do"
	method   = "POST"
	mimeType = "application/json"
	reqBody  = "{ \"foo\": \"bar\" }"
)

func Test_HttpClient_Request_Success(t *testing.T) {
	res := http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"content-type": []string{mimeType},
		},
		Body: ioutil.NopCloser(bytes.NewBufferString("{ \"baz\": \"qux\" }")),
	}

	ft := fakeTransport{
		Response: &res,
	}

	hc := NewHttpClient(&http.Client{
		Transport: &ft,
	})

	ro := hc.Request(context.TODO(), &RequestInput{
		URL:    url,
		Method: method,
		Headers: map[string]string{
			"accept": mimeType,
		},
		Body: reqBody,
	})

	assert.Equal(t, url, ft.Request.URL.String())
	assert.Equal(t, method, ft.Request.Method)
	assert.Equal(t, mimeType, ft.Request.Header.Get("accept"))

	assert.Equal(t, ScheduleStatusSucceeded, ro.Status)
	assert.Contains(
		t,
		ro.Result,
		fmt.Sprintf("\"statusCode\":%v", http.StatusOK))

	assert.Contains(
		t,
		ro.Result,
		fmt.Sprintf("\"headers\":{\"content-type\":\"%v\"}", mimeType))

	assert.Contains(
		t,
		ro.Result,
		"\"body\":\"{ \\\"baz\\\": \\\"qux\\\" }\"")
}

func Test_HttpClient_Request_Fail_Invalid_Request(t *testing.T) {
	hc := NewHttpClient(&http.Client{
		Transport: &fakeTransport{},
	})

	ro := hc.Request(context.TODO(), &RequestInput{
		URL:    "~!@#$%",
	})

	assert.Equal(t, ScheduleStatusFailed, ro.Status)
	assert.Contains(
		t,
		ro.Result,
		"\"error\":")
}

func Test_HttpClient_Request_Fail_Internal_Error(t *testing.T) {
	hc := NewHttpClient(&http.Client{
		Transport: &fakeTransport{
			Error: fmt.Errorf("internal error"),
		},
	})

	ro := hc.Request(context.TODO(), &RequestInput{
		URL:    url,
		Method: method,
		Headers: map[string]string{
			"accept": mimeType,
		},
		Body: reqBody,
	})

	assert.Equal(t, ScheduleStatusFailed, ro.Status)
	assert.Contains(
		t,
		ro.Result,
		"\"error\":")
}

func Test_HttpClient_Request_Fail_Unsuccessful_Status_Code(t *testing.T) {
	res := http.Response{
		StatusCode: http.StatusInternalServerError,
		Header: http.Header{
			"content-type": []string{mimeType},
		},
		Body: ioutil.NopCloser(bytes.NewBufferString("{ \"baz\": \"qux\" }")),
	}

	ft := fakeTransport{
		Response: &res,
	}

	hc := NewHttpClient(&http.Client{
		Transport: &ft,
	})

	ro := hc.Request(context.TODO(), &RequestInput{
		URL:    url,
		Method: method,
		Headers: map[string]string{
			"accept": mimeType,
		},
		Body: reqBody,
	})

	assert.Equal(t, url, ft.Request.URL.String())
	assert.Equal(t, method, ft.Request.Method)
	assert.Equal(t, mimeType, ft.Request.Header.Get("accept"))

	assert.Equal(t, ScheduleStatusFailed, ro.Status)
	assert.Contains(
		t,
		ro.Result,
		fmt.Sprintf("\"statusCode\":%v", http.StatusInternalServerError))

	assert.Contains(
		t,
		ro.Result,
		fmt.Sprintf("\"headers\":{\"content-type\":\"%v\"}", mimeType))

	assert.Contains(
		t,
		ro.Result,
		"\"body\":\"{ \\\"baz\\\": \\\"qux\\\" }\"")
}
