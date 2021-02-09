#!/bin/bash

cp ../../authentication/authsvc .
cp ../../api/apisvc .

docker build -t microservices:v1 .
docker inspect microservices:v1 
