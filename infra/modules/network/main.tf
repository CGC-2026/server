data "aws_availability_zones" "available" {
    state = "available"
}

locals {
    az_a = data.aws_availability_zones.available.names[0]
    az_b = data.aws_availability_zones.available.names[1]

    default_tags = merge(var.tags, {
        Name = var.name_prefix
    })
}

# VPC is responsible for private network
resource "aws_vpc" "this" {
  cidr_block = var.vpc_cidr
  enable_dns_support = true
  enable_dns_hostnames = true

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-vpc"
  })
}

# Gateway allows public subnets reach the internet
resource "aws_internet_gateway" "this" {
    vpc_id = aws_vpc.this.id

    tags = merge(var.tags, {
        Name = "${var.name_prefix}-igw"
    })
}

resource "aws_subnet" "private" {
  vpc_id = aws_vpc.this.id
  cidr_block = var.private_subnet_cidr
  availability_zone = local.az_a

    tags = merge(var.tags, {
    Name = "${var.name_prefix}-private-subnet"
    Tier = "private"
  })
}

# Second private subnet for RDS required by aws
resource "aws_subnet" "private_b" {
  vpc_id = aws_vpc.this.id
  cidr_block = var.private_subnet_cidr_b
  availability_zone = local.az_b

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-private-subnet-b"
    Tier = "private"
  })
}

# Public Subnet for EC2. Has public IP
resource "aws_subnet" "public" {
  vpc_id = aws_vpc.this.id
  cidr_block = var.public_subnet_cidr
  availability_zone = local.az_a
  map_public_ip_on_launch = true

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-public-subnet"
    Tier = "public"
  })
}

# Second public subnet for RDS
resource "aws_subnet" "public_b" {
  vpc_id = aws_vpc.this.id
  cidr_block = var.public_subnet_cidr_b
  availability_zone = local.az_b
  map_public_ip_on_launch = true

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-public-subnet-b"
    Tier = "public"
  })
}

# Public routing table to make subnet public
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.this.id

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-public-rt"
  })
}

resource "aws_route" "public_default" {
    route_table_id = aws_route_table.public.id
    destination_cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.this.id
}

resource "aws_route_table_association" "public_assoc" {
    subnet_id = aws_subnet.public.id
    route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "public_assoc_b" {
  subnet_id = aws_subnet.public_b.id
  route_table_id = aws_route_table.public.id
}

# Security group for EC2
resource "aws_security_group" "ec2" {
  name = "${var.name_prefix}-ec2-sg"
  description = "Allow inbound app traffic to EC2"
  vpc_id = aws_vpc.this.id

  # Set ports to 22 for debugging
  ingress {
    description = "App port"
    from_port = var.app_port
    to_port = var.app_port
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "All outbound"
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-ec2-sg"
  })
}

# Security group for RDS (Postgres)
resource "aws_security_group" "rds" {
  name = "${var.name_prefix}-rds-sg"
  description = "Allow Postgres only from EC2"
  vpc_id = aws_vpc.this.id

  ingress {
    description = "Postgres from EC2"
    from_port = 5432
    to_port = 5432
    protocol = "tcp"
    security_groups = [aws_security_group.ec2.id]
  }

  egress {
    description = "All outbound"
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-rds-sg"
  })
}