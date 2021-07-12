#!/bin/sh
# proto    ， pb.go     types/   , dplatformos_path    dplatformos   proto
dplatformos_path=$(go list -f '{{.Dir}}' "github.com/D-PlatformOperatingSystem/dplatformos")
protoc --go_out=plugins=grpc:../types ./*.proto --proto_path=. --proto_path="${dplatformos_path}/types/proto/"
