package api

import (
	"github.com/graphql-go/graphql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ScheduleList", func() {
	Describe("Name", func() {
		It("is ScheduleList", func() {
			Expect(scheduleListType.Name()).To(Equal("ScheduleList"))
		})
	})

	Describe("Fields", func() {
		It("has schedules as list of Schedule", func() {
			t := scheduleListType.Fields()["schedules"].Type

			Expect(t).To(BeAssignableToTypeOf(&graphql.List{}))
			Expect(t.(*graphql.List).OfType).To(Equal(scheduleType))
		})

		It("has nextKey as nullable ScheduleListNextKey", func() {
			Expect(scheduleListType.Fields()["nextKey"].Type).To(
				Equal(scheduleListNextKeyType))
		})
	})
})
