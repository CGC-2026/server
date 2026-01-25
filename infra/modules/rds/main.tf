

# Specify subnet group for rds
resource "aws_db_subnet_group" "this" {
  name       = "${var.name_prefix}-db-subnets"
  subnet_ids = var.private_subnet_ids

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-db-subnets"
  })
}

# Generate random password
resource "random_password" "db" {
  length  = 24
  special = true
}

# Store credentials in SSM parameter store
resource "aws_ssm_parameter" "db_password" {
  name  = "/${var.name_prefix}/db/password"
  type  = "SecureString"
  value = random_password.db.result
}

resource "aws_ssm_parameter" "db_username" {
  name  = "/${var.name_prefix}/db/username"
  type  = "String"
  value = var.db_name
}

resource "aws_ssm_parameter" "db_name" {
  name  = "/${var.name_prefix}/db/name"
  type  = "String"
  value = var.db_name
}

resource "aws_ssm_parameter" "database_url" {
  name  = "/${var.name_prefix}/api/database_url"
  type  = "SecureString"
  value = "postgres://${var.db_username}:${random_password.db.result}@${aws_db_instance.this.address}:${aws_db_instance.this.port}/${var.db_name}?sslmode=require"

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-database-url"
  })
}

# Define Actual Postgres instance inside private subnets
resource "aws_db_instance" "this" {
  identifier = "${var.name_prefix}-postgres"

  engine         = "postgres"
  engine_version = "16"
  instance_class = var.instance_class

  allocated_storage = var.allocated_storage
  storage_type      = "gp3"

  db_name  = var.db_name
  username = var.db_username
  password = random_password.db.result

  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = [var.rds_sg_id]

  publicly_accessible = false

  backup_retention_period = 7
  skip_final_snapshot     = true
  deletion_protection     = false

  apply_immediately = true

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-postgres"
  })
}