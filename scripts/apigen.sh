#!/bin/bash

protoc -I api api/chatservice.proto --go_out=plugins=grpc:api
