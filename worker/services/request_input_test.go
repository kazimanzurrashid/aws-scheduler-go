package services

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func Test_CreateRequestInput_Return_Input(t *testing.T) {
	attrs := map[string]events.DynamoDBAttributeValue {
		"url": events.NewStringAttribute("https://foo.bar/do"),
		"method": events.NewStringAttribute("POST"),
		"headers": events.NewMapAttribute(map[string]events.DynamoDBAttributeValue{
			"authorization": events.NewStringAttribute("token 123"),
		}),
		"body": events.NewStringAttribute("{ \"foo\": \"bar\" }"),
	}

	ri := CreateRequestInput(attrs)

	assert.Equal(t, "https://foo.bar/do", ri.URL)
	assert.Equal(t, "POST", ri.Method)
	assert.Equal(t, "token 123", ri.Headers["authorization"])
	assert.Equal(t, "{ \"foo\": \"bar\" }", ri.Body)
}

func Test_CreateRequestInput_With_Empty_Header_Adds_Default_Header(t *testing.T) {
	attrs := map[string]events.DynamoDBAttributeValue {
		"url": events.NewStringAttribute("https://foo.bar/do"),
	}

	ri := CreateRequestInput(attrs)

	assert.NotEqual(t, "", ri.Headers["Accept"])
	assert.NotEqual(t, "", ri.Headers["Content-Type"])
}
