#!/bin/bash
# Deploy the MySQL database
cd ./database
./deploy.sh
cd ..

# Deploy the API gateway
cd ./gateway
./deploy.sh
cd ..

ssh root@api.grouppoop.icu < ./ssh.sh