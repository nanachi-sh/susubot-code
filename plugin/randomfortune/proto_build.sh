#!/bin/bash
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. plugin/randomfortune/protos/randomfortune/*.proto
