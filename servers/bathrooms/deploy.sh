#!/bin/bash
echo Deploying bathroom server now.
./build.sh

sudo docker push bowerw2/grouppoop_bathroom_server
