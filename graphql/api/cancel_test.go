package api

import (
	"context"

	"github.com/graphql-go/graphql"

	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cancel", func() {
	Describe("Resolve", func() {
		var (
			field *graphql.Field
			db    fakeCancelStorage
		)

		BeforeEach(func() {
			db = fakeCancelStorage{}
			factory := NewFactory(&db)

			field = factory.Cancel()
		})

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
})

type fakeCancelStorage struct {
	storage.Storage
	ID string
}

//goland:noinspection GoUnusedParameter
func (srv *fakeCancelStorage) Cancel(
	ctx context.Context,
	id string) (bool, error) {

	srv.ID = id

	return true, nil
}
