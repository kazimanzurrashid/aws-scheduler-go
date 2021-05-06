package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/handlers"
)

func main() {
	if os.Getenv("LAMBDA_TASK_ROOT") != "" {
		lambda.Start(handlers.Lambda)
	} else {
		handlers.Http()
	}
}
