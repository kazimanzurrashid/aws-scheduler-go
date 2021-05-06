package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/api"
	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lambda", func() {
	BeforeEach(func() {
		f := api.NewFactory(&fakeStorage{})
		s, _ := f.Schema()
		schema = s
	})

	Context("single request", func() {
		var gatewayRequest events.APIGatewayV2HTTPRequest

		BeforeEach(func() {
			bodyStruct := request{
				Query: "query Get($id: ID!) { get(id: $id) { id url method } }",
				Variables: map[string]interface{}{
					"id": "01234567890",
				},
			}

			bodyBuff, _ := json.Marshal(bodyStruct)

			gatewayRequest = events.APIGatewayV2HTTPRequest{
				RawPath: "/v1/graphql",
				Body:    string(bodyBuff),
				RequestContext: events.APIGatewayV2HTTPRequestContext{
					Stage: "v1",
					HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
						Method: "POST",
					},
				},
			}
		})

		Context("valid", func() {
			var gatewayResponse events.APIGatewayV2HTTPResponse

			BeforeEach(func() {
				gatewayResponse, _ = Lambda(context.TODO(), gatewayRequest)
			})

			It("returns status code OK", func() {
				Expect(gatewayResponse.StatusCode).To(Equal(http.StatusOK))
			})

			It("returns graphql response", func() {
				Expect(gatewayResponse.Body).NotTo(Equal(""))
			})
		})

		Context("invalid", func() {
			var (
				gatewayResponse events.APIGatewayV2HTTPResponse
				realUnmarshal   unmarshal
			)

			BeforeEach(func() {
				realUnmarshal = unmarshalStruct

				unmarshalStruct = func(bytes []byte, i interface{}) error {
					return fmt.Errorf("unmarshal error")
				}

				gatewayResponse, _ = Lambda(context.TODO(), gatewayRequest)
			})

			It("returns status code Bad Request", func() {
				Expect(gatewayResponse.StatusCode).To(
					Equal(http.StatusBadRequest))
			})

			AfterEach(func() {
				unmarshalStruct = realUnmarshal
			})
		})
	})

	Context("multiple request", func() {
		var gatewayRequest events.APIGatewayV2HTTPRequest

		BeforeEach(func() {
			bodyStruct := []request{
				{
					Query: "query Get($id: ID!) { get(id: $id) { id url method } }",
					Variables: map[string]interface{}{
						"id": "01234567890",
					},
				},
				{
					Query: "query Get($id: ID!) { get(id: $id) { id url method } }",
					Variables: map[string]interface{}{
						"id": "01234567890",
					},
				},
			}

			bodyBuff, _ := json.Marshal(bodyStruct)

			gatewayRequest = events.APIGatewayV2HTTPRequest{
				RawPath: "/v1/graphql",
				Body:    string(bodyBuff),
				RequestContext: events.APIGatewayV2HTTPRequestContext{
					Stage: "v1",
					HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
						Method: "POST",
					},
				},
			}
		})

		Context("valid", func() {
			var gatewayResponse events.APIGatewayV2HTTPResponse

			BeforeEach(func() {
				gatewayResponse, _ = Lambda(context.TODO(), gatewayRequest)
			})

			It("returns status code OK", func() {
				Expect(gatewayResponse.StatusCode).To(Equal(http.StatusOK))
			})

			It("returns graphql response", func() {
				Expect(gatewayResponse.Body).NotTo(Equal(""))
			})
		})

		Context("invalid", func() {
			var (
				gatewayResponse events.APIGatewayV2HTTPResponse
				realUnmarshal   unmarshal
			)

			BeforeEach(func() {
				realUnmarshal = unmarshalStruct

				unmarshalStruct = func(bytes []byte, i interface{}) error {
					return fmt.Errorf("unmarshal error")
				}

				gatewayResponse, _ = Lambda(context.TODO(), gatewayRequest)
			})

			It("returns status code Bad Request", func() {
				Expect(gatewayResponse.StatusCode).To(
					Equal(http.StatusBadRequest))
			})

			AfterEach(func() {
				unmarshalStruct = realUnmarshal
			})
		})
	})

	Context("any request", func() {
		Context("empty body", func() {
			var gatewayResponse events.APIGatewayV2HTTPResponse

			BeforeEach(func() {
				gatewayRequest := events.APIGatewayV2HTTPRequest{
					RawPath: "/v1/graphql",
					RequestContext: events.APIGatewayV2HTTPRequestContext{
						Stage: "v1",
						HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
							Method: "POST",
						},
					},
				}
				gatewayResponse, _ = Lambda(context.TODO(), gatewayRequest)
			})

			It("returns status code Bad Request", func() {
				Expect(gatewayResponse.StatusCode).To(
					Equal(http.StatusBadRequest))
			})
		})

		Context("unrecognized body", func() {
			var gatewayResponse events.APIGatewayV2HTTPResponse

			BeforeEach(func() {
				gatewayRequest := events.APIGatewayV2HTTPRequest{
					RawPath: "/v1/graphql",
					Body:    "foo-bar",
					RequestContext: events.APIGatewayV2HTTPRequestContext{
						Stage: "v1",
						HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
							Method: "POST",
						},
					},
				}

				gatewayResponse, _ = Lambda(context.TODO(), gatewayRequest)
			})

			It("returns status code Bad Request", func() {
				Expect(gatewayResponse.StatusCode).To(Equal(
					http.StatusBadRequest))
			})
		})

		Context("error serializing graphql response", func() {
			var (
				gatewayResponse events.APIGatewayV2HTTPResponse
				gatewayError    error
				realMarshal     marshal
			)

			BeforeEach(func() {
				realMarshal = marshalStruct

				marshalStruct = func(i interface{}) ([]byte, error) {
					return nil, fmt.Errorf("marshal error")
				}

				bodyStruct := request{
					Query: "query Get($id: ID!) { get(id: $id) { id url method } }",
					Variables: map[string]interface{}{
						"id": "01234567890",
					},
				}

				bodyBuff, _ := json.Marshal(bodyStruct)

				gatewayRequest := events.APIGatewayV2HTTPRequest{
					RawPath: "/v1/graphql",
					Body:    string(bodyBuff),
					RequestContext: events.APIGatewayV2HTTPRequestContext{
						Stage: "v1",
						HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
							Method: "POST",
						},
					},
				}

				gatewayResponse, gatewayError = Lambda(
					context.TODO(),
					gatewayRequest)
			})

			It("returns status code Internal Server Error", func() {
				Expect(gatewayResponse.StatusCode).To(
					Equal(http.StatusInternalServerError))
			})

			It("returns error", func() {
				Expect(gatewayError).NotTo(BeNil())
			})

			AfterEach(func() {
				marshalStruct = realMarshal
			})
		})
	})
})

type fakeStorage struct {
	storage.Storage
}

//goland:noinspection GoUnusedParameter
func (srv *fakeStorage) Get(
	ctx context.Context,
	id string) (*storage.Schedule, error) {
	return &storage.Schedule{
		ID:     "1234567890",
		DueAt:  time.Now(),
		URL:    "https://foo.bar/do",
		Method: "POST",
		Headers: map[string]string{
			"accept": "application/json",
		},
		Body: aws.String("{ \"foo\": \"bar\" }"),
	}, nil
}
