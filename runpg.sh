#!/bin/bash

set -ex

docker network create --driver bridge postgresql

docker run \
  -p 5432:5432 \
  --name postgres \
  --rm \
  --hostname postgres \
  --network postgresql \
  -e POSTGRESQL_POSTGRES_PASSWORD=test \
  -e POSTGRESQL_DATABASE=testdb \
  -e POSTGRESQL_USERNAME=user \
  -e POSTGRESQL_PASSWORD=test \
  -d bitnami/postgresql:13.1.0-debian-10-r80

docker run \
  -p 5050:5050 \
  --name pgadmin \
  --rm \
  --network postgresql \
  -e PGADMIN_DEFAULT_EMAIL=vmo@ciklum.com \
  -e PGADMIN_DEFAULT_PASSWORD=test \
  -e PGADMIN_LISTEN_PORT=5050 \
  -d dpage/pgadmin4:4.30
