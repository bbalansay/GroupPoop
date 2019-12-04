#!/bin/bash
echo Deploying messages server now.
./build.sh

sudo docker push bowerw2/grouppoop_bathrooms_server
