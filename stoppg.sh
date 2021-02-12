#!/bin/bash

set -x

docker stop postgres
docker stop pgadmin

docker network rm postgresql
