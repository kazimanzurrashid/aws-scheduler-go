package api

import (
	"github.com/graphql-go/graphql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HTTPMethod", func() {
	Describe("Name", func() {
		It("is HTTPMethod", func() {
			Expect(httpMethodType.Name()).To(Equal("HTTPMethod"))
		})
	})

	Describe("Values", func() {
		var (
			create = func(m string) *graphql.EnumValueDefinition {
				return &graphql.EnumValueDefinition{
					Name:  m,
					Value: m,
				}
			}
			values []*graphql.EnumValueDefinition
		)

		BeforeEach(func() {
			values = httpMethodType.Values()
		})

		It("has GET", func() {
			Expect(values).To(ContainElements(create("GET")))
		})

		It("has POST", func() {
			Expect(values).To(ContainElements(create("POST")))
		})

		It("has PUT", func() {
			Expect(values).To(ContainElements(create("PUT")))
		})

		It("has PATCH", func() {
			Expect(values).To(ContainElements(create("PATCH")))
		})

		It("has DELETE", func() {
			Expect(values).To(ContainElements(create("DELETE")))
		})
	})
})
