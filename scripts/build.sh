#!/usr/bin/bash

cd ../rust/
cargo build --release
cargo build --release --target=x86_64-pc-windows-gnu
cp target/release/librust.so ../lib/rust.so
cp target/x86_64-pc-windows-gnu/release/rust.dll ../lib/rust.dll
cd ..
go build -o bin/killstreak -ldflags="-r ./lib" main_unix.go
env CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o bin/killstreak.exe -ldflags="-r ./lib" main_windows.go