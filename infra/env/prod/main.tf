locals {
  name_prefix = "${var.project}-${var.environment}"
}

output "name_prefix" {
  value = local.name_prefix
}

module "network" {
  source = "../../modules/network"

  name_prefix = local.name_prefix

  app_port = 8080

}

output "vpc_id" {
  value = module.network.vpc_id
}

output "public_subnet_id" {
  value = module.network.public_subnet_id
}

output "private_subnet_id" {
  value = module.network.private_subnet_id
}

output "ec2_sg_id" {
  value = module.network.ec2_sg_id
}

output "rds_sg_id" {
  value = module.network.rds_sg_id
}