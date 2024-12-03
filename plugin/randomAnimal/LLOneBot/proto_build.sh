#!/bin/bash
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. plugin/randomAnimal/LLOneBot/protos/randomAnimal/*.proto
