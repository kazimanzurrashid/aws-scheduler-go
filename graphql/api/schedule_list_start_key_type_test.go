package api

import (
	"github.com/graphql-go/graphql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ScheduleListStartKey", func() {
	Describe("Name", func() {
		It("is ScheduleListStartKey", func() {
			Expect(scheduleListStartKeyType.Name()).To(
				Equal("ScheduleListStartKey"))
		})
	})

	Describe("Fields", func() {
		It("has id as non-nullable ID", func() {
			t := scheduleListStartKeyType.Fields()["id"].Type

			Expect(t).To(BeAssignableToTypeOf(&graphql.NonNull{}))
			Expect(t.(*graphql.NonNull).OfType).To(Equal(graphql.ID))
		})

		It("has dueAt as non-nullable DateTime", func() {
			t := scheduleListStartKeyType.Fields()["dueAt"].Type

			Expect(t).To(BeAssignableToTypeOf(&graphql.NonNull{}))
			Expect(t.(*graphql.NonNull).OfType).To(Equal(graphql.DateTime))
		})

		It("has status as nullable ScheduleStatus", func() {
			Expect(scheduleListStartKeyType.Fields()["status"].Type).To(
				Equal(scheduleStatusType))
		})
	})
})
