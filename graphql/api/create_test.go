package api

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/storage"
	"github.com/stretchr/testify/assert"
)

type fakeCreateStorage struct {
	storage.Storage
	Input *storage.CreateInput
}

const (
	id     = "1234567890"
	url    = "https://foo.bar/do"
	method = "POST"
	accept = "application/json"
	body   = "{ \"foo\": \"bar\" }"
)

//goland:noinspection GoUnusedParameter
func (srv *fakeCreateStorage) Create(
	ctx context.Context,
	input *storage.CreateInput) (string, error) {

	srv.Input = input

	return id, nil
}

func Test_Create_Resolve_Success(t *testing.T) {
	dueAt := time.Now().Add(time.Minute * 1)

	db := fakeCreateStorage{}
	factory := NewFactory(&db)
	field := factory.Create()

	res, err := field.Resolve(graphql.ResolveParams{
		Args: map[string]interface{}{
			"dueAt":  dueAt,
			"url":    url,
			"method": method,
			"headers": map[string]string{
				"accept": accept,
			},
			"body": body,
		},
	})

	assert.Equal(t, id, res.(string))
	assert.Nil(t, err)
	assert.Equal(t, dueAt.Unix(), db.Input.DueAt.Unix())
	assert.Equal(t, url, db.Input.URL)
	assert.Equal(t, method, db.Input.Method)
	assert.Equal(t, accept, db.Input.Headers["accept"])
	assert.Equal(t, body, db.Input.Body)
}

func Test_Create_Resolve_Fail_Load_Struct_Error(t *testing.T) {
	realLoadStruct := loadStruct
	loadStruct = func(i interface{}, i2 interface{}) error {
		return fmt.Errorf("load struct error")
	}

	db := fakeCreateStorage{}
	factory := NewFactory(&db)
	field := factory.Create()

	res, err := field.Resolve(graphql.ResolveParams{
		Args: map[string]interface{}{
			"dueAt":  time.Now().Add(-time.Minute * 1),
			"url":    url,
			"method": method,
			"headers": map[string]string{
				"accept": accept,
			},
			"body": body,
		},
	})

	assert.Nil(t, res)
	assert.NotNil(t, err)

	loadStruct = realLoadStruct
}

func Test_Create_Resolve_Fail_Not_Future_Date(t *testing.T) {
	db := fakeCreateStorage{}
	factory := NewFactory(&db)
	field := factory.Create()

	res, err := field.Resolve(graphql.ResolveParams{
		Args: map[string]interface{}{
			"dueAt":  time.Now().Add(-time.Minute * 1),
			"url":    url,
			"method": method,
			"headers": map[string]string{
				"accept": accept,
			},
			"body": body,
		},
	})

	assert.Nil(t, res)
	assert.NotNil(t, err)
}
