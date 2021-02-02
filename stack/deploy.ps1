$env:CGO_ENABLED = "0"
$env:GOOS = "linux"

Set-Location -Path ../graphql
New-Item -Name dist -ItemType directory -Force
go build -o dist/main

Set-Location -Path ../collector
New-Item -Name dist -ItemType directory -Force
go build -o dist/main

Set-Location -Path ../worker
New-Item -Name dist -ItemType directory  -Force
go build -o dist/main

Set-Location -Path ../stack
cdk bootstrap
cdk deploy --require-approval never
