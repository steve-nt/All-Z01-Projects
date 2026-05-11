#!/bin/bash

echo "Building TCPChat..."
go build -o TCPChat .
echo "Build complete. Run with './TCPChat' or './TCPChat <port>'" 