#!/bin/sh
set -eu

spawn-fcgi -s /run/fcgi.sock /usr/bin/fcgiwrap && nginx

