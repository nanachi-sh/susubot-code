#!/bin/bash
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. handler/LLOneBot/protos/handler/*.proto handler/LLOneBot/protos/handler/*/*.proto
