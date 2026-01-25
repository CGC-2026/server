output "db_endpoint" {
  value = aws_db_instance.this.address
}

output "db_port" {
  value = aws_db_instance.this.port
}

output "ssm_password_param" {
  value = aws_ssm_parameter.db_password.name
}

output "ssm_username_param" {
  value = aws_ssm_parameter.db_username.name
}

output "ssm_dbname_param" {
  value = aws_ssm_parameter.db_name.name
}

output "ssm_database_url_param" {
  value = aws_ssm_parameter.database_url.name
}