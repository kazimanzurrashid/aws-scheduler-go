package handlers

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/xray"

	"github.com/graphql-go/graphql"
	"github.com/joho/godotenv"

	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/api"
	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/storage"
)

type request struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

type (
	marshal   func(interface{}) ([]byte, error)
	unmarshal func([]byte, interface{}) error
)

func jsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func jsonUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

var (
	marshalStruct   marshal   = jsonMarshal
	unmarshalStruct unmarshal = jsonUnmarshal
)

var schema graphql.Schema
var playgroundTemplate *template.Template

func executeGraphQL(ctx context.Context, statement string) (interface{}, int) {
	if statement == "" {
		return nil, http.StatusBadRequest
	}

	var ret interface{}

	if strings.HasPrefix(statement, "{") && strings.HasSuffix(statement, "}") {
		var payload request

		if err := unmarshalStruct([]byte(statement), &payload); err != nil {
			return nil, http.StatusInternalServerError
		}

		ret = graphql.Do(graphql.Params{
			Context:        ctx,
			Schema:         schema,
			RequestString:  payload.Query,
			OperationName:  payload.OperationName,
			VariableValues: payload.Variables,
		})
	} else if strings.HasPrefix(statement, "[") &&
		strings.HasSuffix(statement, "]") {
		var payloads []request

		if err := unmarshalStruct([]byte(statement), &payloads); err != nil {
			return nil, http.StatusInternalServerError
		}

		rets := make([]*graphql.Result, len(payloads))
		var wg sync.WaitGroup

		for i, p := range payloads {
			wg.Add(1)
			go func(payload request, index int) {
				defer wg.Done()

				out := graphql.Do(graphql.Params{
					Context:        ctx,
					Schema:         schema,
					RequestString:  payload.Query,
					OperationName:  payload.OperationName,
					VariableValues: payload.Variables,
				})

				rets[index] = out
			}(p, i)
		}

		wg.Wait()

		ret = rets
	} else {
		return nil, http.StatusBadRequest
	}

	return ret, http.StatusOK
}

func init() {
	inLambda := os.Getenv("LAMBDA_TASK_ROOT") != ""

	if inLambda {
		currentDir, _ := os.Getwd()
		playgroundTemplate = template.Must(
			template.ParseFiles(
				filepath.Join(currentDir, "/pages/playground.html")))
	} else {
		_, currentFile, _, _ := runtime.Caller(0)
		currentDir := path.Dir(currentFile)

		envFile := filepath.Join(currentDir, "./../.env")

		if _, err := os.Stat(envFile); err == nil {
			if err := godotenv.Load(envFile); err != nil {
				log.Fatalf("env file load error: %v", err)
				return
			}
		}

		templatePath := filepath.Join(currentDir, "./../pages/playground.html")
		playgroundTemplate = template.Must(template.ParseFiles(templatePath))
	}

	ses := session.Must(session.NewSession())

	ddbc := dynamodb.New(ses)

	if inLambda {
		xray.AWS(ddbc.Client)
	}

	database := storage.NewDatabase(ddbc)

	f := api.NewFactory(database)
	s, err := f.Schema()

	if err != nil {
		log.Fatalf("schema create error: %v", err)
		return
	}

	schema = s
}
