#!/bin/bash -x

load_service_index() {
  local meta_url="http://169.254.169.250/2016-07-29"
  local service_index=$(curl -s $meta_url/self/container/service_index)

  # wait for metadata to wake up
  while [ "$service_index" == "" ]; do
    sleep 1
    service_index=$(curl -s $meta_url/self/container/service_index)
  done

  export RANCHER_SERVICE_INDEX="$service_index"
}

load_service_index

confd -onetime -backend env

exec "$@"
