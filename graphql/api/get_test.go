package api

import (
	"context"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/storage"
	"github.com/stretchr/testify/assert"
)

type fakeGetStorage struct {
	storage.Storage
	ID string
}

//goland:noinspection GoUnusedParameter
func (srv *fakeGetStorage) Get(ctx context.Context, id string) (*storage.Schedule, error) {
	srv.ID = id

	return &storage.Schedule{}, nil
}

func Test_Get_Resolve_Success(t *testing.T) {
	const id = "1234567890"

	db := fakeGetStorage{}
	factory := NewFactory(&db)
	field := factory.Get()

	res, err := field.Resolve(graphql.ResolveParams{
		Args: map[string]interface{}{
			"id": id,
		},
	})

	assert.NotNil(t, res)
	assert.Nil(t, err)
	assert.Equal(t, id, db.ID)
}

func Test_Get_Resolve_Fail_Missing_ID(t *testing.T) {
	db := fakeGetStorage{}
	factory := NewFactory(&db)
	field := factory.Get()

	res, err := field.Resolve(graphql.ResolveParams{
		Args: map[string]interface{}{},
	})

	assert.Nil(t, res)
	assert.NotNil(t, err)
}
