#!/bin/sh
set -eu

tini -s -- fcgiwrap
tini -s -- nginx
