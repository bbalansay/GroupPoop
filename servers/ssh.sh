#!/bin/bash
docker rm -f apiGateway
docker rm -f database
docker rm -f bathroomServer1
docker network rm ServerNetwork

docker volume rm $(docker volume ls -qf dangling=true)

docker pull bowerw2/grouppoop_api_gateway
docker pull bowerw2/grouppoop_database
docker pull bower2/grouppoop_bathroom_server


export TLSCERT=/etc/letsencrypt/live/api.grouppoop.icu/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.grouppoop.icu/privkey.pem
export MYSQL_ROOT_PASSWORD="password123"
export DSN="root:$MYSQL_ROOT_PASSWORD@tcp(database:3306)/users"
export REDISADDR="redisServer:6379"
export SESSIONKEY=$(echo -n "Message" | openssl dgst -sha256 -hmac "secret" -binary | base64)
export SUMMARYADDR="http://summaryServer1:80,http://summaryServer2:81"
export MESSAGESADDR="http://messagesServer1:80,http://messagesServer2:81"


docker network create ServerNetwork

docker run -d \
--name database \
--network ServerNetwork \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=users \
bowerw2/grouppoop_database

docker run -d \
-p 443:443 \
--name apiGateway \
--network ServerNetwork \
--restart=unless-stopped \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e DSN=$DSN \
-e REDISADDR=$REDISADDR \
-e SESSIONKEY=$SESSIONKEY \
-e SUMMARYADDR=$SUMMARYADDR \
-e MESSAGESADDR=$MESSAGESADDR \
bowerw2/grouppoop_api_gateway

docker run -d \
--name bathroomServer1 \
--network ServerNetwork \
--restart=unless-stopped \
-e BATHROOMPORT=":80" \
-e DBHOST="database" \
-e DBPORT="3306" \
-e DBUSER="root" \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e DBNAME="users" \
bowerw2/grouppoop_bathroom_server
