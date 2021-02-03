package api

import "github.com/graphql-go/graphql"

var httpMethodType = graphql.NewEnum(graphql.EnumConfig{
	Name: "HTTPMethod",
	Values: map[string]*graphql.EnumValueConfig{
		"GET":    {Value: "GET"},
		"POST":   {Value: "POST"},
		"PUT":    {Value: "PUT"},
		"PATCH":  {Value: "PATCH"},
		"DELETE": {Value: "DELETE"},
	},
})
