#!/bin/sh
protoc --go_out=plugins=grpc:../types ./*.proto --proto_path=. --proto_path="../../../../vendor/github.com/D-PlatformOperatingSystem/dplatformos/types/proto/"
