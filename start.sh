#!/usr/bin/env bash

go build ./...

echo "starting redis..."
docker run -p 6379:6379 --name redis_cabapi -d redis:4.0.10

mkdir -p ~/tmp_mysql_data_cabapi

echo "starting mysql..."
docker run -p 3306:3306 --name mysql_cabapi -v ~/tmp_mysql_data_cabapi:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=root -d mysql:5.7.22

echo "waiting mysql and redis server before continuing (about a minute)..."
sleep 60s

mysql -hlocalhost -P3306 --protocol=tcp -uroot -proot -e "create database cabapi";

echo "importing database..."
mysql -hlocalhost -P3306 --protocol=tcp -uroot -proot -Dcabapi < ./ny_cab_data_cab_trip_data_full.sql

go run api/main.go

echo "stoping redis"
docker stop redis_cabapi

echo "stoping mysql"
docker stop mysql_cabapi

rm -rf ~/tmp_mysql_data_cabapi

docker rm redis_cabapi
docker rm mysql_cabapi