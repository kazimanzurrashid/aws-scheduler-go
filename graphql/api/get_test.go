package api

import (
	"context"

	"github.com/graphql-go/graphql"

	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Get", func() {
	Describe("Resolve", func() {
		var (
			field *graphql.Field
			db    fakeGetStorage
		)

		BeforeEach(func() {
			db = fakeGetStorage{}
			factory := NewFactory(&db)

			field = factory.Get()
		})

		Context("valid id", func() {
			const id = "1234567890"

			var (
				res interface{}
				err error
			)

			BeforeEach(func() {
				db.Schedule = &storage.Schedule{}

				res, err = field.Resolve(graphql.ResolveParams{
					Args: map[string]interface{}{
						"id": id,
					},
				})
			})

			It("sends id to db", func() {
				Expect(db.ID).To(Equal(id))
			})

			It("returns matching schedule", func() {
				Expect(res).To(Equal(db.Schedule))
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

			It("does not return any schedule", func() {
				Expect(res).To(BeNil())
			})

			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})
})

type fakeGetStorage struct {
	storage.Storage
	ID       string
	Schedule *storage.Schedule
}

//goland:noinspection GoUnusedParameter
func (srv *fakeGetStorage) Get(
	ctx context.Context,
	id string) (*storage.Schedule, error) {

	srv.ID = id

	return srv.Schedule, nil
}
