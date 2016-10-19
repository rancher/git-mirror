#!/bin/bash

sudo apt-get install -y git python-setuptools
git clone https://github.com/LLParse/aws-ec2-assign-elastic-ip /tmp/aws-ec2-assign-elastic-ip
cd /tmp/aws-ec2-assign-elastic-ip
sudo python setup.py install
cd /tmp
sudo rm -rf aws-ec2-assign-elastic-ip
