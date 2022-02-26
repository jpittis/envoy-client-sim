package proto

//go:generate protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative backend.proto
