package api

import (
	"net/http"

	"github.com/graphql-go/graphql"
)

var httpMethodType = graphql.NewEnum(graphql.EnumConfig{
	Name: "HTTPMethod",
	Values: map[string]*graphql.EnumValueConfig{
		"GET":    {Value: http.MethodGet},
		"POST":   {Value: http.MethodPost},
		"PUT":    {Value: http.MethodPut},
		"PATCH":  {Value: http.MethodPatch},
		"DELETE": {Value: http.MethodDelete},
	},
})
