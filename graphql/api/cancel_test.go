package api

import (
	"context"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/storage"
	"github.com/stretchr/testify/assert"
)

type fakeCancelStorage struct {
	storage.Storage
	ID string
}

//goland:noinspection GoUnusedParameter
func (srv *fakeCancelStorage) Cancel(ctx context.Context, id string) (bool, error) {
	srv.ID = id

	return true, nil
}

func Test_Cancel_Resolve_Success(t *testing.T) {
	const id = "1234567890"

	db := fakeCancelStorage{}
	factory := NewFactory(&db)
	field := factory.Cancel()

	res, err := field.Resolve(graphql.ResolveParams{
		Args: map[string]interface{}{
			"id": id,
		},
	})

	assert.True(t, res.(bool))
	assert.Nil(t, err)
	assert.Equal(t, id, db.ID)
}

func Test_Cancel_Resolve_Fail_Missing_ID(t *testing.T) {
	db := fakeCancelStorage{}
	factory := NewFactory(&db)
	field := factory.Cancel()

	res, err := field.Resolve(graphql.ResolveParams{
		Args: map[string]interface{}{},
	})

	assert.False(t, res.(bool))
	assert.NotNil(t, err)
}
