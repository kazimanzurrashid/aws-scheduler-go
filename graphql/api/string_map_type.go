package api

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
)

func identity(value interface{}) interface{} {
	return value
}

func parseLiteral(astValue ast.Value) interface{} {
	kind := astValue.GetKind()

	switch kind {
	case kinds.StringValue:
		return astValue.GetValue()
	case kinds.ObjectValue:
		obj := make(map[string]string)
		for _, v := range astValue.GetValue().([]*ast.ObjectField) {
			obj[v.Name.Value] = parseLiteral(v.Value).(string)
		}
		return obj
	default:
		return nil
	}
}

var stringMapType = graphql.NewScalar(
	graphql.ScalarConfig{
		Name:         "StringMap",
		Description:  "String Key/Value pair list",
		Serialize:    identity,
		ParseValue:   identity,
		ParseLiteral: parseLiteral,
	},
)
