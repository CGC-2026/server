
# Machine Image
data "aws_ami" "al2023" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-2023.*-kernel-6.1-x86_64"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }
}

# IAM role for EC2
data "aws_iam_policy_document" "ec2_assume" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "this" {
  name               = "${var.name_prefix}-ec2-role"
  assume_role_policy = data.aws_iam_policy_document.ec2_assume.json
  tags               = merge(var.tags, { Name = "${var.name_prefix}-ec2-role" })
}

# Allow SSM
resource "aws_iam_role_policy_attachment" "ssm" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

# Allow ECR pull for docker
resource "aws_iam_role_policy_attachment" "ecr_read" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

# Allow reading SSM parameters for DB credentials
data "aws_iam_policy_document" "ssm_read" {
  statement {
    actions = [
      "ssm:GetParameter",
      "ssm:GetParameters",
      "ssm:GetParametersByPath"
    ]
    resources = ["*"]
  }

}

resource "aws_iam_policy" "ssm_read" {
  name   = "${var.name_prefix}-ssm-read"
  policy = data.aws_iam_policy_document.ssm_read.json
}

resource "aws_iam_role_policy_attachment" "ssm_read" {
  role       = aws_iam_role.this.name
  policy_arn = aws_iam_policy.ssm_read.arn
}

# Attaches IAM role to EC2
resource "aws_iam_instance_profile" "this" {
  name = "${var.name_prefix}-ec2-profile"
  role = aws_iam_role.this.name
}


# bootstraps the instance, installing docker and pulling image
locals {
  user_data = <<-EOF
    #!/bin/bash
    set -euo pipefail

    # Install ssm agent
    dnf install -y amazon-ssm-agent
    systemctl enable amazon-ssm-agent
    systemctl start amazon-ssm-agent

    # Install docker + aws cli tools
    dnf update -y
    dnf install -y docker
    systemctl enable docker
    systemctl start docker
    usermod -aG docker ec2-user

    echo "DB_NAME_PARAM=${var.ssm_dbname_param}" >> /etc/environment
    echo "DB_USER_PARAM=${var.ssm_username_param}" >> /etc/environment
    echo "DB_PASS_PARAM=${var.ssm_password_param}" >> /etc/environment

    echo "EC2 bootstrap complete" > /var/log/app-bootstrap.txt
  EOF
}

# EC2 instance hosts backend container
resource "aws_instance" "this" {
  ami                         = data.aws_ami.al2023.id
  instance_type               = var.instance_type
  subnet_id                   = var.public_subnet_id
  vpc_security_group_ids      = [var.ec2_sg_id]
  iam_instance_profile        = aws_iam_instance_profile.this.name
  associate_public_ip_address = true

  user_data = local.user_data

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-ec2"
  })
}