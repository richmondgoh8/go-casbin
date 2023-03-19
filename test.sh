#!/bin/bash

CVPKG=$(go list ./internal/... | grep -v mocks | grep -v cmd | tr '\n', ',')

go test -short -count 1 -coverpkg $CVPKG -coverprofile coverage.out ./...
go tool cover -func coverage.out

# https://blog.seriesci.com/how-to-measure-code-coverage-in-go/
# https://github.com/AlexBeauchemin/gobadge
# Generate Local Badge
gobadge -filename='coverage.out' -value=$(go tool cover -func coverage.out | grep total | awk '{print $3}')