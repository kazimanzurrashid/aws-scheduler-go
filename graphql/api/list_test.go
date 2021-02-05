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

var _ = Describe("List", func() {
	var (
		field *graphql.Field
		db    fakeListStorage
	)

	BeforeEach(func() {
		db = fakeListStorage{}
		factory := NewFactory(&db)

		field = factory.List()
	})

	Describe("Args", func() {
		It("has status as nullable ScheduleStatus", func() {
			Expect(field.Args["status"].Type).To(Equal(scheduleStatusType))
		})

		It("has dueAt as nullable DateRange", func() {
			Expect(field.Args["dueAt"].Type).To(Equal(dataRangeType))
		})

		It("has startKey as nullable ScheduleListStartKey", func() {
			Expect(field.Args["startKey"].Type).To(
				Equal(scheduleListStartKeyType))
		})

		It("has limit as nullable Int with default value", func() {
			a := field.Args["limit"]

			Expect(a.Type).To(Equal(graphql.Int))
			Expect(a.DefaultValue).To(Equal(25))
		})
	})

	Describe("Resolve", func() {
		const (
			id     = "1234567890"
			status = storage.ScheduleStatusIdle
			limit  = int64(10)
		)

		Describe("valid input", func() {
			var (
				res interface{}
				err error

				from time.Time
				to   time.Time
			)

			BeforeEach(func() {
				db.ReturnList = &storage.List{}

				from = time.Now().Add(-time.Minute * 3)
				to = time.Now().Add(-time.Minute * 2)

				res, err = field.Resolve(graphql.ResolveParams{
					Args: map[string]interface{}{
						"status": status,
						"dueAt": map[string]interface{}{
							"from": from,
							"to":   to,
						},
						"startKey": map[string]interface{}{
							"id": id,
						},
						"limit": limit,
					},
				})
			})

			It("sends input to db", func() {
				Expect(db.Input.Status).To(Equal(status))
				Expect(db.Input.DueAt.From.Unix()).To(Equal(from.Unix()))
				Expect(db.Input.DueAt.To.Unix()).To(Equal(to.Unix()))
				Expect(db.Input.StartKey.ID).To(Equal(id))
				Expect(db.Input.Limit).To(Equal(limit))
			})

			It("returns result", func() {
				Expect(res).ToNot(BeNil())
			})

			It("does not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("invalid input", func() {
			Context("from is less than to of dueAt", func() {
				var (
					res interface{}
					err error
				)

				BeforeEach(func() {
					res, err = field.Resolve(graphql.ResolveParams{
						Args: map[string]interface{}{
							"dueAt": map[string]interface{}{
								"from": time.Now().Add(-time.Minute * 2),
								"to":   time.Now().Add(-time.Minute * 3),
							},
						},
					})
				})

				It("does not return result", func() {
					Expect(res).To(BeNil())
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})
			})

			Context("limit is less than 1", func() {
				var (
					res interface{}
					err error
				)

				BeforeEach(func() {
					res, err = field.Resolve(graphql.ResolveParams{
						Args: map[string]interface{}{
							"limit": 0,
						},
					})
				})

				It("does not return result", func() {
					Expect(res).To(BeNil())
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})
			})

			Context("limit is greater than 100", func() {
				var (
					res interface{}
					err error
				)

				BeforeEach(func() {
					res, err = field.Resolve(graphql.ResolveParams{
						Args: map[string]interface{}{
							"limit": 101,
						},
					})
				})

				It("does not return result", func() {
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
						Args: make(map[string]interface{}),
					})
				})

				It("does not return result", func() {
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

	Describe("Type", func() {
		It("returns ScheduleList", func() {
			Expect(field.Type).To(Equal(scheduleListType))
		})
	})
})

type fakeListStorage struct {
	storage.Storage
	Input storage.ListInput

	ReturnList *storage.List
}

//goland:noinspection GoUnusedParameter
func (srv *fakeListStorage) List(
	ctx context.Context,
	input storage.ListInput) (*storage.List, error) {

	srv.Input = input

	return srv.ReturnList, nil
}
