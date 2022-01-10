package api

import (
	"context"

	"github.com/graphql-go/graphql"

	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cancel", func() {
	var (
		field *graphql.Field
		db    fakeCancelStorage
	)

	BeforeEach(func() {
		db = fakeCancelStorage{}
		factory := NewFactory(&db)

		field = factory.Cancel()
	})

	Describe("Args", func() {
		It("has id as non-nullable ID", func() {
			t := field.Args["id"].Type

			Expect(t).To(BeAssignableToTypeOf(&graphql.NonNull{}))
			Expect(t.(*graphql.NonNull).OfType).To(Equal(graphql.ID))
		})
	})

	Describe("Resolve", func() {
		Context("valid id", func() {
			const id = "1234567890"

			var (
				res interface{}
				err error
			)

			BeforeEach(func() {
				res, err = field.Resolve(graphql.ResolveParams{
					Args: map[string]interface{}{
						"id": id,
					},
				})
			})

			It("sends id to db", func() {
				Expect(db.ID).To(Equal(id))
			})

			It("returns success", func() {
				Expect(res).To(Equal(true))
			})

			It("does not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("missing id", func() {
			var (
				res interface{}
				err error
			)

			BeforeEach(func() {
				res, err = field.Resolve(graphql.ResolveParams{
					Args: map[string]interface{}{
						"id": "",
					},
				})
			})

			It("returns fail", func() {
				Expect(res).To(Equal(false))
			})

			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})

	Describe("Type", func() {
		It("returns boolean", func() {
			Expect(field.Type).To(Equal(graphql.Boolean))
		})
	})
})

type fakeCancelStorage struct {
	storage.Storage
	ID string
}

func (srv *fakeCancelStorage) Cancel(
	_ context.Context,
	id string) (bool, error) {

	srv.ID = id

	return true, nil
}
