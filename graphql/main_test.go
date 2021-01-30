package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/api"
	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/storage"
	"github.com/stretchr/testify/assert"
)

type fakeStorage struct {
	storage.Storage
}

//goland:noinspection GoUnusedParameter
func (srv *fakeStorage) Get(
	ctx context.Context,
	id string) (*storage.Schedule, error) {
	return &storage.Schedule{
		ID:     "1234567890",
		DueAt:  time.Now(),
		URL:    "https://foo.bar/do",
		Method: "POST",
		Headers: map[string]string{
			"accept": "application/json",
		},
		Body: "{ \"foo\": \"bar\" }",
	}, nil
}

func Test_handler_Single_Request_Success(t *testing.T) {
	realSchema := schema
	f := api.NewFactory(&fakeStorage{})
	s, err := f.Schema()
	if err != nil {
		t.Error(err)
	}
	schema = s

	bodyStruct := request{
		Query: "query Get($id: ID!) { get(id: $id) { id url method } }",
		Variables: map[string]interface{}{
			"id": "01234567890",
		},
	}

	bodyBuff, err := json.Marshal(bodyStruct)
	if err != nil {
		t.Error(err)
	}

	req := events.APIGatewayV2HTTPRequest{
		Body: string(bodyBuff),
	}

	res, err := handler(context.TODO(), req)

	assert.NotNil(t, res)
	//goland:noinspection GoNilness
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Nil(t, err)

	schema = realSchema
}

func Test_handler_Single_Request_Fail_Unmarshal_Error(t *testing.T) {
	realSchema := schema
	f := api.NewFactory(&fakeStorage{})
	s, err := f.Schema()
	if err != nil {
		t.Error(err)
	}
	schema = s

	realUnmarshal := unmarshalStruct

	unmarshalStruct = func(bytes []byte, i interface{}) error {
		return fmt.Errorf("unmarshal error")
	}

	bodyStruct := request{
		Query: "query Get($id: ID!) { get(id: $id) { id url method } }",
		Variables: map[string]interface{}{
			"id": "01234567890",
		},
	}

	bodyBuff, err := json.Marshal(bodyStruct)
	if err != nil {
		t.Error(err)
	}

	req := events.APIGatewayV2HTTPRequest{
		Body: string(bodyBuff),
	}

	res, err := handler(context.TODO(), req)

	assert.NotNil(t, res)
	//goland:noinspection GoNilness
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.NotNil(t, err)

	unmarshalStruct = realUnmarshal
	schema = realSchema
}

func Test_handler_Multi_Request_Success(t *testing.T) {
	realSchema := schema
	f := api.NewFactory(&fakeStorage{})
	s, err := f.Schema()
	if err != nil {
		t.Error(err)
	}
	schema = s

	bodyStruct := []request{
		{
			Query: "query Get($id: ID!) { get(id: $id) { id url method } }",
			Variables: map[string]interface{}{
				"id": "01234567890",
			},
		},
		{
			Query: "query Get($id: ID!) { get(id: $id) { id url method } }",
			Variables: map[string]interface{}{
				"id": "01234567890",
			},
		},
	}

	bodyBuff, err := json.Marshal(bodyStruct)
	if err != nil {
		t.Error(err)
	}

	req := events.APIGatewayV2HTTPRequest{
		Body: string(bodyBuff),
	}

	res, err := handler(context.TODO(), req)

	assert.NotNil(t, res)
	//goland:noinspection GoNilness
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Nil(t, err)

	schema = realSchema
}

func Test_handler_Multi_Request_Fail_Unmarshal_Error(t *testing.T) {
	realSchema := schema
	f := api.NewFactory(&fakeStorage{})
	s, err := f.Schema()
	if err != nil {
		t.Error(err)
	}
	schema = s

	realUnmarshal := unmarshalStruct

	unmarshalStruct = func(bytes []byte, i interface{}) error {
		return fmt.Errorf("unmarshal error")
	}

	bodyStruct := []request{
		{
			Query: "query Get($id: ID!) { get(id: $id) { id url method } }",
			Variables: map[string]interface{}{
				"id": "01234567890",
			},
		},
		{
			Query: "query Get($id: ID!) { get(id: $id) { id url method } }",
			Variables: map[string]interface{}{
				"id": "01234567890",
			},
		},
	}

	bodyBuff, err := json.Marshal(bodyStruct)
	if err != nil {
		t.Error(err)
	}

	req := events.APIGatewayV2HTTPRequest{
		Body: string(bodyBuff),
	}

	res, err := handler(context.TODO(), req)

	assert.NotNil(t, res)
	//goland:noinspection GoNilness
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.NotNil(t, err)

	unmarshalStruct = realUnmarshal
	schema = realSchema
}

func Test_handler_Fail_Empty_Body(t *testing.T) {
	req := events.APIGatewayV2HTTPRequest{}

	res, err := handler(context.TODO(), req)

	assert.NotNil(t, res)
	//goland:noinspection GoNilness
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Nil(t, err)

}

func Test_handler_Fail_Invalid_Body(t *testing.T) {
	req := events.APIGatewayV2HTTPRequest{
		Body: "foo=bar",
	}

	res, err := handler(context.TODO(), req)

	assert.NotNil(t, res)
	//goland:noinspection GoNilness
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Nil(t, err)
}

func Test_handler_Fail_Marshal_Error(t *testing.T) {
	realSchema := schema
	f := api.NewFactory(&fakeStorage{})
	s, err := f.Schema()
	if err != nil {
		t.Error(err)
	}
	schema = s

	realMarshal := marshalStruct

	marshalStruct = func(i interface{}) ([]byte, error) {
		return nil, fmt.Errorf("marshal error")
	}

	bodyStruct := request{
		Query: "query Get($id: ID!) { get(id: $id) { id url method } }",
		Variables: map[string]interface{}{
			"id": "01234567890",
		},
	}

	bodyBuff, err := json.Marshal(bodyStruct)
	if err != nil {
		t.Error(err)
	}

	req := events.APIGatewayV2HTTPRequest{
		Body: string(bodyBuff),
	}

	res, err := handler(context.TODO(), req)

	assert.NotNil(t, res)
	//goland:noinspection GoNilness
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.NotNil(t, err)

	marshalStruct = realMarshal
	schema = realSchema
}
