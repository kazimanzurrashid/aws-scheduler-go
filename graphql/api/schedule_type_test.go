package api

import (
	"github.com/graphql-go/graphql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Schedule", func() {
	Describe("Name", func() {
		It("is Schedule", func() {
			Expect(scheduleType.Name()).To(Equal("Schedule"))
		})
	})

	Describe("Fields", func() {
		It("has id as non-nullable ID", func() {
			t := scheduleType.Fields()["id"].Type

			Expect(t).To(BeAssignableToTypeOf(&graphql.NonNull{}))
			Expect(t.(*graphql.NonNull).OfType).To(Equal(graphql.ID))
		})

		It("has dueAt as non-nullable DateTime", func() {
			t := scheduleType.Fields()["dueAt"].Type

			Expect(t).To(BeAssignableToTypeOf(&graphql.NonNull{}))
			Expect(t.(*graphql.NonNull).OfType).To(Equal(graphql.DateTime))
		})

		It("has url as non-nullable String", func() {
			t := scheduleType.Fields()["url"].Type

			Expect(t).To(BeAssignableToTypeOf(&graphql.NonNull{}))
			Expect(t.(*graphql.NonNull).OfType).To(Equal(graphql.String))
		})

		It("has method as non-nullable HTTPMethod", func() {
			t := scheduleType.Fields()["method"].Type

			Expect(t).To(BeAssignableToTypeOf(&graphql.NonNull{}))
			Expect(t.(*graphql.NonNull).OfType).To(Equal(httpMethodType))
		})

		It("has headers as nullable StringMap", func() {
			Expect(scheduleType.Fields()["headers"].Type).To(
				Equal(stringMapType))
		})

		It("has body as nullable String", func() {
			Expect(scheduleType.Fields()["body"].Type).To(Equal(graphql.String))
		})

		It("has status as non-nullable ScheduleStatus", func() {
			t := scheduleType.Fields()["status"].Type

			Expect(t).To(BeAssignableToTypeOf(&graphql.NonNull{}))
			Expect(t.(*graphql.NonNull).OfType).To(Equal(scheduleStatusType))
		})

		It("has startedAt as nullable DateTime", func() {
			Expect(scheduleType.Fields()["startedAt"].Type).To(
				Equal(graphql.DateTime))
		})

		It("has completedAt as nullable DateTime", func() {
			Expect(scheduleType.Fields()["completedAt"].Type).To(
				Equal(graphql.DateTime))
		})

		It("has canceledAt as nullable DateTime", func() {
			Expect(scheduleType.Fields()["canceledAt"].Type).To(
				Equal(graphql.DateTime))
		})

		It("has result as nullable String", func() {
			Expect(scheduleType.Fields()["result"].Type).To(
				Equal(graphql.String))
		})

		It("has createdAt as non-nullable DateTime", func() {
			t := scheduleType.Fields()["createdAt"].Type

			Expect(t).To(BeAssignableToTypeOf(&graphql.NonNull{}))
			Expect(t.(*graphql.NonNull).OfType).To(Equal(graphql.DateTime))
		})
	})
})
