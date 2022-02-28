@echo off
cd ..
 .\proto\protoc.exe -I=.\proto --plugin=protoc-gen-go=.\proto\protoc-gen-go.exe --go_out=.\core\command\pb .\proto\online.proto
 .\proto\protoc.exe -I=.\proto --plugin=protoc-gen-go=.\proto\protoc-gen-go.exe --go_out=.\core\command\pb .\proto\command.proto
 .\proto\protoc.exe -I=.\proto --plugin=protoc-gen-go=.\proto\protoc-gen-go.exe --go_out=.\core\command\pb .\proto\pvp.proto
 .\proto\protoc.exe -I=.\proto --plugin=protoc-gen-go=.\proto\protoc-gen-go.exe --go_out=.\core\command\pb .\proto\gateway.proto

echo compile proto ok

pause
