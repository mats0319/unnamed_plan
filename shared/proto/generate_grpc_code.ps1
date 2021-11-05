# Quick Start: https://grpc.io/docs/languages/go/quickstart/
# require:
# 1. protoc, download: https://github.com/protocolbuffers/protobuf/releases, remember set PATH
# 2. protoc-gen-go, install script: go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
# 3. protoc-gen-go-grpc, install script: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

$path = Get-Location

Set-Location $PSScriptRoot

    if (!(Test-Path "./impl")) {
        mkdir "./impl"
    }

    # out path is relative on proto file
    protoc --go_out=./impl --go_opt=paths=source_relative `
    --go-grpc_out=./impl --go-grpc_opt=paths=source_relative `
    .\common.proto `
    .\user.proto `
    .\cloud_file.proto `
    .\thinking_note.proto

Set-Location $path
