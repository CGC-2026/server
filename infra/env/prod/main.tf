locals {
  name_prefix = "${var.project}-${var.environment}"

  tags = {
    Project     = var.project
    Environment = var.environment
    ManagedBy   = "opentofu"
  }
}

output "name_prefix" {
  value = local.name_prefix
}

module "network" {
  source = "../../modules/network"

  name_prefix = local.name_prefix
  tags        = local.tags

  app_port = 8080

}

output "vpc_id" { value = module.network.vpc_id }
output "ec2_sg_id" { value = module.network.ec2_sg_id }
output "rds_sg_id" { value = module.network.rds_sg_id }
output "private_subnet_ids" { value = module.network.private_subnet_ids }
output "public_subnet_ids" { value = module.network.public_subnet_ids }

module "rds" {
  source = "../../modules/rds"

  name_prefix        = local.name_prefix
  tags               = local.tags

  # DB lives in private subnets, accessed only by the app server (EC2)
  private_subnet_ids = module.network.private_subnet_ids
  rds_sg_id          = module.network.rds_sg_id

  db_name     = "campusbuddy"
  db_username = "appuser"
}

output "db_endpoint" { value = module.rds.db_endpoint }
output "db_port" { value = module.rds.db_port }
output "ssm_password_param" { value = module.rds.ssm_password_param }
output "ssm_username_param" { value = module.rds.ssm_username_param }
output "ssm_dbname_param" { value = module.rds.ssm_dbname_param }