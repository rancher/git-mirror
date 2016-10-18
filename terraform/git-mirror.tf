variable "region" {
  default = "eu-west-1"
}

variable "key_name" {
  default = "git"
}

variable "region_count" {
  default = {
    us-west-1      = 2
    us-east-1      = 2
    us-west-2      = 0
    us-west-1      = 0
    eu-west-1      = 2
    eu-central-1   = 0
    ap-northeast-1 = 0
    ap-northeast-2 = 0
    ap-southeast-1 = 0
    ap-southeast-2 = 0
    sa-east-1      = 0
  }
}

variable "zone_id" {
  default = "Z2Y4C5EBXP7TX6"
}

variable "zones" {
  default = {
    Z2Y4C5EBXP7TX6 = "rancher-test.com"
    Z3EMIF7NU6YP0B = "rancher.space"
  }
}

variable "region_ami" {
  default = {
    us-west-1      = "ami-f47c3494"
    us-east-1      = "ami-94510583"
    us-west-2      = "ami-6715cf07"
    us-west-1      = "ami-f47c3494"
    eu-west-1      = "ami-9de6abee"
    eu-central-1   = "ami-ea5ba585"
    ap-northeast-1 = "ami-58359039"
    ap-northeast-2 = "ami-beda0ed0"
    ap-southeast-1 = "ami-0dfb5d6e"
    ap-southeast-2 = "ami-f1ddef92"
    sa-east-1      = "ami-f8801d94"
  }
}

provider "aws" {
  region = "${var.region}"
}

resource "aws_instance" "git_mirror" {
  ami             = "${lookup(var.region_ami, var.region)}"
  count           = "${lookup(var.region_count, var.region)}"
  instance_type   = "t2.micro"
  key_name        = "${var.key_name}"
  monitoring      = true
  security_groups = ["${aws_security_group.git_mirror.name}"]
  tags {
    Name = "${format("git-mirror-%d", count.index + 1)}"
  }
  user_data       = "<<EOF
  #!/bin/bash
  docker run -d -v git-mirror:/var/git --net=host --restart=always --name=git-mirror llparse/git-mirror
  docker run -d -v git-mirror:/var/git -v /var/log/nginx:/var/log/nginx --net=host --restart=always --name=git-serve llparse/git-serve
  EOF"
}

resource "aws_eip" "git_mirror" {
  count    = "${lookup(var.region_count, var.region)}"
  instance = "${element(aws_instance.git_mirror.*.id, count.index)}"
}

resource "aws_security_group" "git_mirror" {
  name = "git-mirror"
  description = "Allow traffic for git-mirror instances"
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 4141
    to_port     = 4141
    protocol    = "tcp"
    cidr_blocks = ["192.30.252.0/22"] # Github
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }
}

/*resource "aws_route53_record" "www" {
  count = 1
  zone_id = "${var.zone_id}"
  name = "git.${var.region}.${lookup(var.zones, var.zone_id)}"
  type = "A"
  ttl = 300
  records = ["${aws_eip.git_mirror.*.public_ip}"]
}*/