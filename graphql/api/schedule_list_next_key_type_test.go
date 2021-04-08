package api

import (
	"github.com/graphql-go/graphql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ScheduleListNextKey", func() {
	Describe("Name", func() {
		It("is ScheduleListNextKey", func() {
			Expect(scheduleListNextKeyType.Name()).To(
				Equal("ScheduleListNextKey"))
		})
	})

	Describe("Fields", func() {
		It("has id as non-nullable ID", func() {
			t := scheduleListNextKeyType.Fields()["id"].Type

			Expect(t).To(BeAssignableToTypeOf(&graphql.NonNull{}))
			Expect(t.(*graphql.NonNull).OfType).To(Equal(graphql.ID))
		})

		It("has dueAt as nullable Int", func() {
			Expect(scheduleListNextKeyType.Fields()["dueAt"].Type).To(
				Equal(graphql.Int))
		})

		It("has status as nullable ScheduleStatus", func() {
			Expect(scheduleListNextKeyType.Fields()["status"].Type).To(
				Equal(scheduleStatusType))
		})
	})
})
