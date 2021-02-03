package services

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HttpClient", func() {
	const (
		url      = "https://foo.bar/do"
		method   = "POST"
		mimeType = "application/json"
		reqBody  = "{ \"foo\": \"bar\" }"
	)

	var (
		ft fakeTransport
		hc *HttpClient
	)

	BeforeEach(func() {
		ft = fakeTransport{}

		hc = NewHttpClient(&http.Client{
			Transport: &ft,
		})
	})

	Describe("Request", func() {
		Describe("success", func() {
			var ro *ResponseOutput

			BeforeEach(func() {
				res := http.Response{
					StatusCode: http.StatusOK,
					Header: http.Header{
						"content-type": []string{mimeType},
					},
					Body: ioutil.NopCloser(
						bytes.NewBufferString("{ \"baz\": \"qux\" }")),
				}

				ft.Response = &res

				ro = hc.Request(context.TODO(), &RequestInput{
					URL:    url,
					Method: method,
					Headers: map[string]string{
						"accept": mimeType,
					},
					Body: reqBody,
				})
			})

			It("sets url", func() {
				Expect(ft.Request.URL.String()).To(Equal(url))
			})

			It("sets method", func() {
				Expect(ft.Request.Method).To(Equal(method))
			})

			It("sets header", func() {
				Expect(ft.Request.Header.Get("accept")).To(Equal(mimeType))
			})

			It("returns succeeded status", func() {
				Expect(ro.Status, ScheduleStatusSucceeded)
			})

			It("returns result with http status code", func() {
				Expect(ro.Result).To(
					ContainSubstring(
						fmt.Sprintf(
							"\"statusCode\":%v",
							http.StatusOK)))
			})

			It("returns result with header", func() {
				Expect(ro.Result).To(
					ContainSubstring(
						fmt.Sprintf(
							"\"headers\":{\"content-type\":\"%v\"}",
							mimeType)))
			})

			It("returns result with body", func() {
				Expect(ro.Result).To(
					ContainSubstring("\"body\":\"{ \\\"baz\\\": \\\"qux\\\" }\""))
			})

			AfterEach(func() {
				ft.Response = nil
			})
		})

		Describe("fail", func() {
			Context("invalid request", func() {
				var ro *ResponseOutput

				BeforeEach(func() {
					ro = hc.Request(context.TODO(), &RequestInput{
						URL: "~!@#$%",
					})
				})

				It("returns failed status", func() {
					Expect(ro.Status, ScheduleStatusFailed)
				})

				It("returns result with with", func() {
					Expect(ro.Result).To(ContainSubstring("\"error\":"))
				})
			})

			Context("internal error", func() {
				var ro *ResponseOutput

				BeforeEach(func() {
					ft.Error = fmt.Errorf("internal error")

					ro = hc.Request(context.TODO(), &RequestInput{
						URL:    url,
						Method: method,
					})
				})

				It("returns failed status", func() {
					Expect(ro.Status, ScheduleStatusFailed)
				})

				It("returns result with with", func() {
					Expect(ro.Result).To(ContainSubstring("\"error\":"))
				})

				AfterEach(func() {
					ft.Error = nil
				})
			})

			Context("unsuccessful http status code", func() {
				var ro *ResponseOutput

				BeforeEach(func() {
					res := http.Response{
						StatusCode: http.StatusInternalServerError,
						Header: http.Header{
							"content-type": []string{mimeType},
						},
						Body: ioutil.NopCloser(
							bytes.NewBufferString("{ \"baz\": \"qux\" }")),
					}

					ft.Response = &res

					ro = hc.Request(context.TODO(), &RequestInput{
						URL:    url,
						Method: method,
						Headers: map[string]string{
							"accept": mimeType,
						},
						Body: reqBody,
					})
				})

				It("returns failed status", func() {
					Expect(ro.Status, ScheduleStatusFailed)
				})

				It("returns result with http status code", func() {
					Expect(ro.Result).To(
						ContainSubstring(
							fmt.Sprintf(
								"\"statusCode\":%v",
								http.StatusInternalServerError)))
				})

				It("returns result with header", func() {
					Expect(ro.Result).To(
						ContainSubstring(
							fmt.Sprintf(
								"\"headers\":{\"content-type\":\"%v\"}",
								mimeType)))
				})

				It("returns result with body", func() {
					Expect(ro.Result).To(
						ContainSubstring(
							"\"body\":\"{ \\\"baz\\\": \\\"qux\\\" }\""))
				})

				AfterEach(func() {
					ft.Response = nil
				})
			})
		})
	})
})

type fakeTransport struct {
	Request  *http.Request
	Response *http.Response
	Error    error
}

func (ft *fakeTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	ft.Request = r

	return ft.Response, ft.Error
}
