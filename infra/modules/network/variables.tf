variable "name_prefix" {
    type = string
}

variable "tags" {
    type = map(string)
    default = {}
}

variable "vpc_cidr" {
  type = string
  default = "10.0.0.0/16"
}

variable "public_subnet_cidr" {
    type = string
    default = "10.0.1.0/24"
}

variable "private_subnet_cidr" {
    type = string
    default = "10.0.2.0/24"
}

variable "app_port" {
  type = number
  default = 8080
}