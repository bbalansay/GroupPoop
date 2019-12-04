#!/bin/bash
# Deploy the MySQL database
cd ./database
./deploy.sh
cd ..

# Deploy the auth server
cd ./auth
./deploy.sh
cd ..

# Deploy the users server
cd ./users
./deploy.sh
cd ..

# Deploy the bathrooms server
cd ./bathrooms
./deploy.sh
cd ..

# Deploy the API gateway
cd ./gateway
./deploy.sh
cd ..

ssh root@api.grouppoop.icu < ./ssh.sh