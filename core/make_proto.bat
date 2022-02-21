@echo off

.\proto\protoc.exe -I=.\proto --plugin=protoc-gen-go=.\proto\protoc-gen-go.exe --go_out=.\src\command\pb .\proto\online.proto
.\proto\protoc.exe -I=.\proto --plugin=protoc-gen-go=.\proto\protoc-gen-go.exe --go_out=.\src\command\pb .\proto\command.proto
.\proto\protoc.exe --go_out=plugins=.\proto\grpc:.\src\command\pb .\proto\grpc.proto
echo compile proto ok

pause
