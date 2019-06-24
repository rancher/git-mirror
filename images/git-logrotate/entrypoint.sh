#!/bin/bash -x

confd -onetime -backend env

/sbin/tini -s -- "$@"
