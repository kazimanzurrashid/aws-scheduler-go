#!/usr/bin/env bash

cd ../graphql || exit
mkdir -p dist
CGO_ENABLED=0 go build -o dist/main

cd ../collector || exit
mkdir -p dist
CGO_ENABLED=0 go build -o dist/main

cd ../worker || exit
mkdir -p dist
CGO_ENABLED=0 go build -o dist/main

cd ../stack || exit
cdk bootstrap
cdk deploy --require-approval never
