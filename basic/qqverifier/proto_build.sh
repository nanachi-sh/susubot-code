#!/bin/bash
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --grpc-gateway_out=paths=source_relative:. basic/qqverifier/protos/qqverifier/*.proto