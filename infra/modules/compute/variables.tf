variable "name_prefix" {
  type = string
}

variable "tags" {
  type    = map(string)
  default = {}
}

variable "public_subnet_id" { type = string }
variable "ec2_sg_id" { type = string }

variable "app_port" {
  type    = number
  default = 8080
}

# Parameters for DB connection

variable "ssm_dbname_param" { type = string }
variable "ssm_username_param" { type = string }
variable "ssm_password_param" { type = string }

variable "ecr_image" {
  type    = string
  default = ""
}

variable "instance_type" {
  type    = string
  default = "t3.micro"
}