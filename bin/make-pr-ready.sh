#!/bin/bash -x

# run housekeeping tasks

go test ./...
go vet ./...
go fmt ./...
go generate ./...
go mod tidy
prettier --w pkg/js/helpers.js
./bin/fmtjson pkg/js/parse_tests/*json
for i in $(ls pkg/js/parse_tests/*js); do
  dnscontrol fmt -i $i -o $i
done
