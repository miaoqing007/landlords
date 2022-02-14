@echo off

.\proto\protoc.exe -I=.\proto --plugin=protoc-gen-go=.\proto\protoc-gen-go.exe --go_out=.\command\pb .\proto\game.proto

echo compile proto ok

pause