
variable "aws_region" {
  type    = string
  default = "ca-central-1"
}

variable "project" {
  type    = string
  default = "cgc-2026"
}

variable "environment" {
  type    = string
  default = "prod"
}

variable "api_log_retention_in_days" {
  type    = number
  default = 60
}

variable "ami_id" {
  type = string
}

variable "instance_type" {
  type    = string
  default = "t3.xlarge"
}