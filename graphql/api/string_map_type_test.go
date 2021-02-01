package api

import (
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("stringMapType", func() {
	Describe("identity", func() {
		var (
			in  interface{}
			ret interface{}
		)

		BeforeEach(func() {
			in = "foo-bar"
			ret = identity(in)
		})

		It("returns the same", func() {
			Expect(ret).To(Equal(in))
		})
	})

	Describe("parseLiteral", func() {
		Context("string", func() {
			var ret interface{}

			BeforeEach(func() {
				ret = parseLiteral(&ast.StringValue{
					Value: "foo-bar",
					Kind:  kinds.StringValue,
				})
			})

			It("returns value", func() {
				Expect(ret).To(Equal("foo-bar"))
			})
		})

		Context("object", func() {
			var ret interface{}

			BeforeEach(func() {
				ret = parseLiteral(&ast.ObjectValue{
					Fields: []*ast.ObjectField{
						{
							Name: &ast.Name{
								Kind:  kinds.StringValue,
								Value: "foo",
							},
							Value: &ast.StringValue{
								Kind:  kinds.StringValue,
								Value: "bar",
							},
							Kind: kinds.StringValue,
						},
					},
					Kind: kinds.ObjectValue,
				})
			})

			It("returns key/value pairs", func() {
				pair := ret.(map[string]string)
				Expect(pair["foo"]).To(Equal("bar"))
			})
		})

		Context("other", func() {
			var ret interface{}

			BeforeEach(func() {
				ret = parseLiteral(&ast.IntValue{
					Value: "123",
					Kind:  kinds.IntValue,
				})
			})

			It("returns nothing", func() {
				Expect(ret).To(BeNil())
			})
		})
	})
})
