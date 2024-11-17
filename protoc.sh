#!/usr/bin/env bash

protoc --go_out=./bop --go_opt=paths=source_relative --go-grpc_out=./bop --go-grpc_opt=paths=source_relative bop.proto
