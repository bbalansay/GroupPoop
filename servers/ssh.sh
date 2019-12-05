#!/bin/bash
docker rm -f apiGateway
docker rm -f redisServer
docker rm -f database
docker rm -f bathroomsServer
docker rm -f authServer
docker rm -f usersServer
docker network rm ServerNetwork

docker volume rm $(docker volume ls -qf dangling=true)

docker pull bowerw2/grouppoop_api_gateway
docker pull bowerw2/grouppoop_database
docker pull bowerw2/grouppoop_bathrooms_server
docker pull bowerw2/grouppoop_auth_server
docker pull bowerw2/grouppoop_users_server


export TLSCERT=/etc/letsencrypt/live/api.grouppoop.icu/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.grouppoop.icu/privkey.pem
export MYSQL_ROOT_PASSWORD="password123"
export DSN="root:$MYSQL_ROOT_PASSWORD@tcp(database:3306)/users"
export REDISADDR="redisServer:6379"
export SESSIONKEY=$(echo -n "Message" | openssl dgst -sha256 -hmac "secret" -binary | base64)
export BATHROOMADDR="http://bathroomsServer:80"
export AUTHADDR="http://authServer:80"
export USERSADDR="http://usersServer:80"


docker network create ServerNetwork

docker run -d --name redisServer --network ServerNetwork redis

docker run -d \
--name database \
--network ServerNetwork \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=users \
bowerw2/grouppoop_database

docker run -d \
--name bathroomsServer \
--network ServerNetwork \
--restart=unless-stopped \
-e BATHROOMSPORT=":80" \
-e DBHOST="database" \
-e DBPORT="3306" \
-e DBUSER="root" \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e DBNAME="users" \
bowerw2/grouppoop_bathrooms_server

docker run -d \
--name authServer \
--network ServerNetwork \
--restart=unless-stopped \
-e AUTHPORT=":80" \
-e DSN=$DSN \
-e REDISADDR=$REDISADDR \
-e SESSIONKEY=$SESSIONKEY \
bowerw2/grouppoop_auth_server

docker run -d \
--name usersServer \
--network ServerNetwork \
--restart=unless-stopped \
-e USERSPORT=":80" \
-e DSN=$DSN \
-e REDISADDR=$REDISADDR \
-e SESSIONKEY=$SESSIONKEY \
bowerw2/grouppoop_users_server

docker run -d \
-p 443:443 \
--name apiGateway \
--network ServerNetwork \
--restart=unless-stopped \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e REDISADDR=$REDISADDR \
-e SESSIONKEY=$SESSIONKEY \
-e BATHROOMADDR=$BATHROOMADDR \
-e AUTHADDR=$AUTHADDR \
-e USERSADDR=$USERSADDR \
bowerw2/grouppoop_api_gateway
