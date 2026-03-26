variable "name_prefix" {
  type = string
}

variable "tags" {
  type    = map(string)
  default = {}
}


variable "api_log_retentions_in_days" {
  type    = number
  default = 60
}

variable "ec2_role_name" {
  type = string
}