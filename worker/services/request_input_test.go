package services

import (
	"github.com/aws/aws-lambda-go/events"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateRequestInput", func() {

	Context("with header", func() {
		var (
			attrs map[string]events.DynamoDBAttributeValue
			ri    *RequestInput
		)

		BeforeEach(func() {
			attrs = map[string]events.DynamoDBAttributeValue{
				"url":    events.NewStringAttribute("https://foo.bar/do"),
				"method": events.NewStringAttribute("POST"),
				"headers": events.NewMapAttribute(
					map[string]events.DynamoDBAttributeValue{
						"authorization": events.NewStringAttribute("token 123"),
					}),
				"body": events.NewStringAttribute("{ \"foo\": \"bar\" }"),
			}

			ri = CreateRequestInput(attrs)
		})

		It("sets url", func() {
			Expect(ri.URL).To(Equal("https://foo.bar/do"))
		})

		It("sets method", func() {
			Expect(ri.Method).To(Equal("POST"))
		})

		It("sets header", func() {
			Expect(ri.Headers["authorization"]).To(Equal("token 123"))
		})

		It("sets body", func() {
			Expect(ri.Body).To(Equal("{ \"foo\": \"bar\" }"))
		})
	})

	Context("without header", func() {
		var (
			attrs map[string]events.DynamoDBAttributeValue
			ri    *RequestInput
		)

		BeforeEach(func() {
			attrs = map[string]events.DynamoDBAttributeValue{
				"url": events.NewStringAttribute("https://foo.bar/do"),
				"method": events.NewStringAttribute("GET"),
			}
			ri = CreateRequestInput(attrs)
		})

		It("sets url", func() {
			Expect(ri.URL).To(Equal("https://foo.bar/do"))
		})

		It("sets method", func() {
			Expect(ri.Method).To(Equal("GET"))
		})

		It("sets default header", func() {
			Expect(ri.Headers["accept"]).NotTo(Equal(""))
			Expect(ri.Headers["content-type"]).NotTo(Equal(""))
		})
	})
})
