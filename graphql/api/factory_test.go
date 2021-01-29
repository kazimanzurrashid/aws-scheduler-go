package api

import (
	"testing"

	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/storage"
	"github.com/stretchr/testify/assert"
)

type fakeStorage struct {
	storage.Storage
}

func Test_Factory_NewFactory_Success(t *testing.T) {
	db := fakeStorage{}
	factory := NewFactory(&db)

	assert.NotNil(t, factory)
}

