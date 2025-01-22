#!/bin/bash
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. plugin/randomanimal/protos/randomanimal/*.proto
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. plugin/randomanimal/protos/fileweb/*.proto
