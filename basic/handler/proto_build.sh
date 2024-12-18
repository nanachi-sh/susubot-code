#!/bin/bash
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. basic/handler/protos/handler/request/*.proto
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. basic/handler/protos/handler/response/*.proto
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. basic/handler/protos/connector/*.proto
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. basic/handler/protos/fileweb/*.proto