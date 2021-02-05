package api

import (
	"github.com/graphql-go/graphql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DateRange", func() {
	Describe("Name", func() {
		It("is DateRange", func() {
			Expect(dataRangeType.Name()).To(Equal("DateRange"))
		})
	})

	Describe("Fields", func() {
		It("has from as non-nullable DateTime", func() {
			t := dataRangeType.Fields()["from"].Type

			Expect(t).To(BeAssignableToTypeOf(&graphql.NonNull{}))
			Expect(t.(*graphql.NonNull).OfType).To(Equal(graphql.DateTime))
		})

		It("has to as non-nullable DateTime", func() {
			t := dataRangeType.Fields()["to"].Type

			Expect(t).To(BeAssignableToTypeOf(&graphql.NonNull{}))
			Expect(t.(*graphql.NonNull).OfType).To(Equal(graphql.DateTime))
		})
	})
})
