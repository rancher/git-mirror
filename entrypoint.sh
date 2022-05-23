#!/bin/bash

echo "fs.file-max = 12000500" >> /etc/sysctl.conf
echo "fs.nr_open = 20000500" >> /etc/sysctl.conf
echo "net.core.somaxconn = 512" >> /etc/sysctl.conf
echo "# <domain> <type> <item>  <value>" >> /etc/security/limits.d/limits.conf
echo "    *       soft  nofile  20000" >> /etc/security/limits.d/limits.conf
echo "    *       hard  nofile  20000" >> /etc/security/limits.d/limits.conf

sysctl -p

ulimit -n 200000

# Turn up services
echo 'FCGI_CHILDREN="5"' >> /etc/default/fcgiwrap
service fcgiwrap start
nginx

