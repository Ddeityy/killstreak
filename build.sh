#!/usr/bin/bash

cd rust/
echo "building .so"
cargo build --release
echo "copying dynamic library"
cp target/release/librust.so ../lib/rust.so
cd ..
echo "building go"
go build -o bin/killstreak -ldflags="-r ./lib" .