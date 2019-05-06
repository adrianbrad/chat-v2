#!/bin/bash

export CODECOV_TOKEN="a05673cd-dc8a-4e07-8f59-45ddd5c7a9d9"

parent_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
cd $parent_path
cd ../..
#add -race flag in the future
go test -coverprofile=./test/codecov/coverage.txt -covermode=atomic {./test/...,./pkg/...,./cmd/...,./configs/...,./internal/...}
wait $!

bash <(curl -s https://codecov.io/bash)
