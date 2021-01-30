package main

import (
	"context"
	"testing"

	"github.com/kazimanzurrashid/aws-scheduler-go/collector/storage"
	"github.com/stretchr/testify/assert"
)

type fakeStorage struct {
	storage.Storage

	Called bool
}

//goland:noinspection GoUnusedParameter
func (srv *fakeStorage) Update(ctx context.Context) error  {
	srv.Called = true
	return nil
}

func Test_handler_Success(t *testing.T) {
	fake := fakeStorage{}
	database = &fake

	_ = handler(context.TODO())

	assert.True(t, fake.Called)
}
