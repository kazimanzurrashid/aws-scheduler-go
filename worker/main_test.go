package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"

	"github.com/kazimanzurrashid/aws-scheduler-go/worker/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("handler", func() {
	var (
		fc fakeClient
		fs fakeStorage

		ro services.ResponseOutput

		err error
	)

	BeforeEach(func() {
		ro = services.ResponseOutput{
			Status: services.ScheduleStatusSucceeded,
			Result: "dummy result",
		}

		fc = fakeClient{
			Output: &ro,
		}

		fs = fakeStorage{}

		httpClient = &fc
		database = &fs

		err = handler(context.TODO(), events.DynamoDBEvent{
			Records: []events.DynamoDBEventRecord{
				{
					EventName: "INSERT",
					Change: events.DynamoDBStreamRecord{
						NewImage: map[string]events.DynamoDBAttributeValue{
							"status": events.NewStringAttribute(
								services.ScheduleStatusQueued),
						},
					},
				},
				{
					EventName: "MODIFY",
					Change: events.DynamoDBStreamRecord{
						NewImage: map[string]events.DynamoDBAttributeValue{
							"id":    events.NewStringAttribute("1234"),
							"dueAt": events.NewNumberAttribute("9876543"),
							"url": events.NewStringAttribute(
								"https://foo.bar/do"),
							"method":    events.NewStringAttribute("POST"),
							"createdAt": events.NewNumberAttribute("343334232"),
							"status": events.NewStringAttribute(
								services.ScheduleStatusQueued),
						},
					},
				},
				{
					EventName: "MODIFY",
					Change: events.DynamoDBStreamRecord{
						NewImage: map[string]events.DynamoDBAttributeValue{
							"status": events.NewStringAttribute(
								services.ScheduleStatusSucceeded),
						},
					},
				},
				{
					EventName: "REMOVE",
					Change: events.DynamoDBStreamRecord{
						NewImage: map[string]events.DynamoDBAttributeValue{
							"status": events.NewStringAttribute(
								services.ScheduleStatusFailed),
						},
					},
				},
			},
		})
	})

	It("only sends completed schedules to database to persist", func() {
		Expect(fs.Inputs).To(HaveLen(1))
	})

	It("updates status of completed schedules", func() {
		for _, input := range fs.Inputs {
			Expect(input.Status).To(Equal(ro.Status))
		}
	})

	It("updates result of completed schedules", func() {
		for _, input := range fs.Inputs {
			Expect(*input.Result).To(Equal(ro.Result))
		}
	})

	It("updates startedAt of completed schedules", func() {
		for _, input := range fs.Inputs {
			Expect(*input.StartedAt).NotTo(Equal(0))
		}
	})

	It("updates completedAt of completed schedules", func() {
		for _, input := range fs.Inputs {
			Expect(*input.CompletedAt).NotTo(Equal(0))
		}
	})

	It("does not return error", func() {
		Expect(err).To(BeNil())
	})
})

type fakeClient struct {
	services.Client

	Output *services.ResponseOutput
}

//goland:noinspection GoUnusedParameter
func (fc *fakeClient) Request(
	ctx context.Context,
	ri *services.RequestInput) *services.ResponseOutput {
	return fc.Output
}

type fakeStorage struct {
	services.Storage

	Inputs []*services.UpdateInput
}

//goland:noinspection GoUnusedParameter
func (fs *fakeStorage) Update(
	ctx context.Context,
	inputs []*services.UpdateInput) error {

	fs.Inputs = inputs

	return nil
}
