#!/usr/bin/bash

cd rust/
echo "building .so for linux"
cargo build --release
cargo build --release --target=x86_64-pc-windows-gnu
echo "copying dynamic libraries and executables"
cp target/x86_64-pc-windows-gnu/release/cut_demo.exe ../bin/
cp target/x86_64-pc-windows-gnu/release/parse_demo.exe ../bin/
cp target/release/librust.so ../lib/rust.so
cd ..
echo "building go for linux"
go build -o bin/killstreak -ldflags="-r ./lib" .
echo "building go for windows"
env CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o bin/killstreak.exe .