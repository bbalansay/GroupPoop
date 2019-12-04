#!/bin/bash
echo Deploying summary server now.
./build.sh

sudo docker push bowerw2/grouppoop_users_server
