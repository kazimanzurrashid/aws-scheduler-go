package main

import (
	"context"
	"github.com/kazimanzurrashid/aws-scheduler-go/collector/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("handler", func() {
	var (
		storage fakeStorage
		err error
	)

	BeforeEach(func() {
		storage = fakeStorage{}
		database = &storage

		err = handler(context.TODO())
	})

	It("calls database update", func() {
		Expect(storage.Called).To(BeTrue())
	})

	It("does not return error", func() {
		Expect(err).To(BeNil())
	})
})

type fakeStorage struct {
	storage.Storage

	Called bool
}

//goland:noinspection GoUnusedParameter
func (srv *fakeStorage) Update(ctx context.Context) error  {
	srv.Called = true
	return nil
}
