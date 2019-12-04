#!/bin/bash
echo Building API server now
CGO_ENABLED=0 GOOS=linux go build
sudo docker build -t bowerw2/grouppoop_users_server .
go clean