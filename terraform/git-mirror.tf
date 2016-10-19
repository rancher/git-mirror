variable "region" {
  default = "sa-east-1"
}

variable "key_name" {
  default = "git"
}

// EFS is only supported in eu-west-1, us-east-1, us-east-2, us-west-2

variable "region_count" {
  default = {
    us-east-1      = 4
    us-east-2      = 3
    us-west-1      = 2
    us-west-2      = 3
    eu-west-1      = 3
    eu-central-1   = 2
    ap-northeast-1 = 2
    ap-northeast-2 = 2
    ap-southeast-1 = 2
    ap-southeast-2 = 3
    ap-south-1     = 2
    sa-east-1      = 2
  }
}

variable "region_az" {
  default = {
    us-east-1      = "us-east-1a,us-east-1b,us-east-1c,us-east-1e"
    us-east-2      = "us-east-2a,us-east-2b,us-east-2c"
    us-west-1      = "us-west-1a,us-west-1b"
    us-west-2      = "us-west-2a,us-west-2b,us-west-2c"
    eu-west-1      = "eu-west-1a,eu-west-1b,eu-west-1c"
    eu-central-1   = "eu-central-1a,eu-central-1b"
    ap-northeast-1 = "ap-northeast-1a,ap-northeast-1c"
    ap-northeast-2 = "ap-northeast-2a,ap-northeast-2c"
    ap-southeast-1 = "ap-southeast-1a,ap-southeast-1b"
    ap-southeast-2 = "ap-southeast-2a,ap-southeast-2b,ap-southeast-2c"
    ap-south-1     = "ap-south-1a,ap-south-1b"
    sa-east-1      = "sa-east-1a,sa-east-1c"                           # t2.micro not available in sa-east-1b
  }
}

// must be HVM
variable "instance_type" {
  default = "t2.micro"
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
    us-east-1      = "ami-ee6f3ef9"
    us-east-2      = "ami-d5045eb0"
    us-west-1      = "ami-43db9323"
    us-west-2      = "ami-6e67c20e"
    eu-west-1      = "ami-6f6c231c"
    eu-central-1   = "ami-a2e51ccd"
    ap-northeast-1 = "ami-6462c505"
    ap-northeast-2 = "ami-5df32733"
    ap-southeast-1 = "ami-cc6fc9af"
    ap-southeast-2 = "ami-d67845b5"
    ap-south-1     = "ami-1560147a"
    sa-east-1      = "ami-2173ee4d"
  }
}

variable "aws_access_key" {
  default = ""
}

variable "aws_secret_key" {
  default = ""
}

provider "aws" {
  access_key = "${var.aws_access_key}"
  secret_key = "${var.aws_secret_key}"
  region = "${var.region}"
}

resource "aws_launch_configuration" "git_mirror" {
  name          = "git-mirror-lc"
  image_id      = "${lookup(var.region_ami, var.region)}"
  instance_type = "${var.instance_type}"
  key_name      = "${var.key_name}"
  root_block_device {
    volume_type           = "gp2"
    volume_size           = 20
    delete_on_termination = true
  }
  security_groups   = ["${aws_security_group.git_mirror.name}"]
  user_data       = "#!/bin/bash -ex
exec > >(tee /var/log/user-data.log|logger -t user-data -s 2>/dev/console) 2>&1
sleep 30
sudo aws-ec2-assign-elastic-ip --region ${var.region} --access-key ${var.aws_access_key} --secret-key ${var.aws_secret_key} --valid-ips ${join(",", aws_eip.git_mirror.*.public_ip)}
sleep 15
sudo docker run -d -v git-mirror:/var/git --net=host --restart=always --name=git-mirror llparse/git-mirror
sudo docker run -d -v git-mirror:/var/git --net=host --restart=always --name=git-serve -v /var/log/nginx:/var/log/nginx llparse/git-serve"
}

resource "aws_autoscaling_group" "git_mirror" {
  name                 = "git-mirror-asg"
  availability_zones   = ["${split(",", lookup(var.region_az, var.region))}"]
  max_size             = "${lookup(var.region_count, var.region)}"
  min_size             = "${lookup(var.region_count, var.region)}"
  desired_capacity     = "${lookup(var.region_count, var.region)}"
  launch_configuration = "${aws_launch_configuration.git_mirror.name}"
  lifecycle {
    create_before_destroy = true
  }
  tag {
    key = "Name"
    value = "git-mirror"
    propagate_at_launch = true
  }
}

# decoupled from the instance
resource "aws_eip" "git_mirror" {
  count    = "${lookup(var.region_count, var.region)}"
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

resource "aws_route53_record" "git_mirror" {
  zone_id = "${var.zone_id}"
  name    = "git.${var.region}.${lookup(var.zones, var.zone_id)}"
  type    = "A"
  ttl     = 300
  records = ["${aws_eip.git_mirror.*.public_ip}"]
}
