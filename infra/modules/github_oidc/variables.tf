variable "name_prefix" {
  type = string
}

variable "tags" {
  type    = map(string)
  default = {}
}

variable "github_repo" {
  type = string
}

variable "github_branch" {
  type    = string
  default = "main"
}

variable "aws_region" {
  type    = string
  default = "ca-central-1"
}

variable "ecr_repository_name" {
  type = string
}

variable "ec2_instance_id" {
  type = string
}