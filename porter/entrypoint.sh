#!/bin/sh
set -eu

spawn-fcgi -s /run/fcgi.sock \
	   -u nginx \
	   -g www-data \
	   -M 0700 \
           /usr/bin/fcgiwrap && nginx

