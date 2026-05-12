#!/bin/bash
nohup go run cmd/api/main.go > backend.log 2>&1 </dev/null &
