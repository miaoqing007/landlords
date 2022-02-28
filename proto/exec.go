//go:generate ./protoc.exe -I=../proto --go_out=plugins=grpc:../src/command/pb ../proto/grpc.proto

package command
