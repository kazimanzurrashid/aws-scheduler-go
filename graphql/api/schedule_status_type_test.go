package api

import (
	"github.com/graphql-go/graphql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ScheduleStatus", func() {
	Describe("Name", func() {
		It("is ScheduleStatus", func() {
			Expect(scheduleStatusType.Name()).To(Equal("ScheduleStatus"))
		})
	})

	Describe("Values", func() {
		var (
			create = func(s string) *graphql.EnumValueDefinition {
				return &graphql.EnumValueDefinition{
					Name:  s,
					Value: s,
				}
			}
			values []*graphql.EnumValueDefinition
		)

		BeforeEach(func() {
			values = scheduleStatusType.Values()
		})

		It("has IDLE", func() {
			Expect(values).To(ContainElements(create("IDLE")))
		})

		It("has QUEUED", func() {
			Expect(values).To(ContainElements(create("QUEUED")))
		})

		It("has SUCCEEDED", func() {
			Expect(values).To(ContainElements(create("SUCCEEDED")))
		})

		It("has CANCELED", func() {
			Expect(values).To(ContainElements(create("CANCELED")))
		})

		It("has FAILED", func() {
			Expect(values).To(ContainElements(create("FAILED")))
		})
	})
})
