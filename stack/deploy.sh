#!/usr/bin/env bash

cd ../graphql || exit
mkdir -p dist
CGO_ENABLED=0 go build -o dist/main
cd dist && zip -r -9 scheduler-graphql-v1.zip ./* && rm main

cd ../../collector || exit
mkdir -p dist
CGO_ENABLED=0 go build -o dist/main
cd dist && zip -r -9 scheduler-collector-v1.zip ./* && rm main

cd ../../worker || exit
mkdir -p dist
CGO_ENABLED=0 go build -o dist/main
cd dist && zip -r -9 scheduler-worker-v1.zip ./* && rm main

cd ../../stack || exit
cdk deploy --require-approval never
