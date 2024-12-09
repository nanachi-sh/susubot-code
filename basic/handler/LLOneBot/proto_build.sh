#!/bin/bash
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. basic/handler/LLOneBot/protos/handler/request/*.proto
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. basic/handler/LLOneBot/protos/handler/response/*.proto
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. basic/handler/LLOneBot/protos/connector/*.proto
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. basic/handler/LLOneBot/protos/fileweb/*.proto