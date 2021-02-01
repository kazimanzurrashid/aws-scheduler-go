package api

import (
	"context"
	"fmt"
	"time"

	"github.com/graphql-go/graphql"

	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", func() {
	Describe("Resolve", func() {
		const (
			id     = "1234567890"
			url    = "https://foo.bar/do"
			method = "POST"
			accept = "application/json"
			body   = "{ \"foo\": \"bar\" }"
		)

		var (
			field *graphql.Field
			db    fakeCreateStorage
		)

		BeforeEach(func() {
			db = fakeCreateStorage{}
			factory := NewFactory(&db)

			field = factory.Create()
		})

		Describe("valid input", func() {
			var (
				res interface{}
				err error

				dueAt time.Time
			)

			BeforeEach(func() {
				db.ReturnID = id

				dueAt = time.Now().Add(time.Minute * 1)

				res, err = field.Resolve(graphql.ResolveParams{
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
			})

			It("sends input to db", func() {
				Expect(db.Input.DueAt.Unix()).To(Equal(dueAt.Unix()))
				Expect(db.Input.URL).To(Equal(url))
				Expect(db.Input.Method).To(Equal(method))
				Expect(db.Input.Headers["accept"]).To(Equal(accept))
				Expect(db.Input.Body).To(Equal(body))
			})

			It("returns newly created id", func() {
				Expect(res).To(Equal(id))
			})

			It("does not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("invalid input", func() {
			Context("not future dua at", func() {
				var (
					res interface{}
					err error
				)

				BeforeEach(func() {
					res, err = field.Resolve(graphql.ResolveParams{
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
				})

				It("does not return id", func() {
					Expect(res).To(BeNil())
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})
			})

			Context("not future dua at", func() {
				var (
					res interface{}
					err error
				)

				BeforeEach(func() {
					res, err = field.Resolve(graphql.ResolveParams{
						Args: map[string]interface{}{
							"dueAt":  time.Now().Add(time.Minute * 1),
							"url":    "~!@#$%",
							"method": method,
							"headers": map[string]string{
								"accept": accept,
							},
							"body": body,
						},
					})
				})

				It("does not return id", func() {
					Expect(res).To(BeNil())
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})
			})
		})

		Describe("any input", func() {
			Context("deserializing input error", func() {

				var (
					res            interface{}
					err            error
					realLoadStruct load
				)

				BeforeEach(func() {
					realLoadStruct = loadStruct

					loadStruct = func(i interface{}, i2 interface{}) error {
						return fmt.Errorf("load struct error")
					}

					res, err = field.Resolve(graphql.ResolveParams{
						Args: map[string]interface{}{
							"dueAt":  time.Now().Add(time.Hour * 1),
							"url":    url,
							"method": method,
							"headers": map[string]string{
								"accept": accept,
							},
							"body": body,
						},
					})
				})

				It("does not return id", func() {
					Expect(res).To(BeNil())
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})

				AfterEach(func() {
					loadStruct = realLoadStruct
				})
			})
		})
	})
})

type fakeCreateStorage struct {
	storage.Storage
	Input storage.CreateInput

	ReturnID string
}

//goland:noinspection GoUnusedParameter
func (srv *fakeCreateStorage) Create(
	ctx context.Context,
	input storage.CreateInput) (string, error) {

	srv.Input = input

	return srv.ReturnID, nil
}
