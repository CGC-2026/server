variable "name_prefix" {
    type = string
}

variable "tags" { 
    type = map(string) 
    default = {}
}

variable "private_subnet_ids" {
    type = list(string)
}

variable "rds_sg_id" {
    type = string
}

variable "db_name" {
  type = string
  default = "app"
}

variable "db_username" {
  type = string
  default = "appuser"
}

variable "instance_class" {
    type = string
    default = "db.t3.micro"
}

variable "allocated_storage" {
    type = number
    default = 20
}