package api

import (
	"github.com/graphql-go/graphql"

	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Factory", func() {
	Describe("NewFactory", func() {
		var factory *Factory

		BeforeEach(func() {
			db := fakeStorage{}
			factory = NewFactory(&db)
		})

		It("returns new factory", func() {
			Expect(factory).NotTo(BeNil())
		})
	})

	Describe("Schema", func() {
		var (
			schema graphql.Schema
			err    error
		)

		BeforeEach(func() {
			db := fakeStorage{}
			factory := NewFactory(&db)

			schema, err = factory.Schema()
		})

		It("return schema", func() {
			Expect(schema).NotTo(BeNil())
		})

		It("does not return error", func() {
			Expect(err).To(BeNil())
		})

		It("has get in query", func() {
			Expect(schema.QueryType().Fields()["get"]).NotTo(BeNil())
		})

		It("has list in query", func() {
			Expect(schema.QueryType().Fields()["list"]).NotTo(BeNil())
		})

		It("has create in mutation", func() {
			Expect(schema.MutationType().Fields()["create"]).NotTo(BeNil())
		})

		It("has cancel in mutation", func() {
			Expect(schema.MutationType().Fields()["cancel"]).NotTo(BeNil())
		})
	})
})

type fakeStorage struct {
	storage.Storage
}
