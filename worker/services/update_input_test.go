package services

import (
	"github.com/aws/aws-lambda-go/events"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateUpdateInput", func() {
	var (
		attrs map[string]events.DynamoDBAttributeValue
		ui *UpdateInput
	)

	BeforeEach(func() {
		attrs = map[string]events.DynamoDBAttributeValue{
			"id":     events.NewStringAttribute("1234"),
			"dueAt":  events.NewNumberAttribute("9876543"),
			"url":    events.NewStringAttribute("https://foo.bar/do"),
			"method": events.NewStringAttribute("PATCH"),
			"headers": events.NewMapAttribute(
				map[string]events.DynamoDBAttributeValue{
					"authorization": events.NewStringAttribute("token 123"),
				}),
			"body":      events.NewStringAttribute("{ \"foo\": \"bar\" }"),
			"createdAt": events.NewNumberAttribute("343334232"),

			"startedAt":   events.NewNumberAttribute("53454344"),
			"completedAt": events.NewNumberAttribute("2256r5454"),
			"status":      events.NewStringAttribute(ScheduleStatusQueued),
			"result":      events.NewStringAttribute("dummy result"),
		}

		ui = CreateUpdateInput(attrs)
	})

	It("sets id", func() {
		Expect(ui.ID).To(Equal("1234"))
	})

	It("sets dueAt", func() {
		Expect(ui.DueAt).To(BeEquivalentTo(9876543))
	})

	It("sets url", func() {
		Expect(ui.URL).To(Equal("https://foo.bar/do"))
	})

	It("sets method", func() {
		Expect(ui.Method).To(Equal("PATCH"))
	})

	It("sets header", func() {
		Expect(ui.Headers["authorization"]).To(Equal("token 123"))
	})

	It("sets body", func() {
		Expect(ui.Body).To(Equal("{ \"foo\": \"bar\" }"))
	})

	It("sets createdAt", func() {
		Expect(ui.CreatedAt).To(BeEquivalentTo(343334232))
	})

	It("never sets status", func() {
		Expect(ui.Status).To(Equal(""))
	})

	It("never sets result", func() {
		Expect(ui.Result).To(Equal(""))
	})

	It("never sets startedAt", func() {
		Expect(ui.StartedAt).To(BeEquivalentTo(0))
	})

	It("never sets completedAt", func() {
		Expect(ui.CompletedAt).To(BeEquivalentTo(0))
	})
})
