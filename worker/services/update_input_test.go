package services

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func Test_CreateUpdateInput_Returns_Input(t *testing.T) {
	attrs := map[string]events.DynamoDBAttributeValue{
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

	ui := CreateUpdateInput(attrs)

	assert.Equal(t, "1234", ui.ID)
	assert.EqualValues(t, 9876543, ui.DueAt)
	assert.Equal(t, "PATCH", ui.Method)
	assert.Equal(t, "token 123", ui.Headers["authorization"])
	assert.Equal(t, "{ \"foo\": \"bar\" }", ui.Body)
	assert.EqualValues(t, 343334232, ui.CreatedAt)

	assert.EqualValues(t, 0, ui.StartedAt)
	assert.EqualValues(t, 0, ui.CompletedAt)
	assert.Equal(t, "", ui.Status)
	assert.Equal(t, "", ui.Result)
}
