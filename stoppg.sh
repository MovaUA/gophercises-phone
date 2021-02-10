#!/bin/bash

set -ex

docker stop postgres
docker stop pgadmin

docker network rm postgresql
