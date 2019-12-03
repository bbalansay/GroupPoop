#!/bin/bash
echo Deploying API server now.
./build.sh

sudo docker push bowerw2/grouppoop_api_gateway
