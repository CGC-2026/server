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

  name_prefix = local.name_prefix
  tags        = local.tags

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

# EC2 runtime host for backend API
module "compute" {
  source = "../../modules/compute"

  name_prefix = local.name_prefix
  tags        = local.tags

  public_subnet_id = module.network.public_subnet_ids[0]
  ec2_sg_id        = module.network.ec2_sg_id

  ssm_dbname_param   = module.rds.ssm_dbname_param
  ssm_username_param = module.rds.ssm_username_param
  ssm_password_param = module.rds.ssm_password_param

  app_port = 8080

  ami_id = var.ami_id
}

output "ec2_instance_id" { value = module.compute.instance_id }
output "ec2_public_ip" { value = module.compute.public_ip }

module "ecr" {
  source = "../../modules/ecr"

  name_prefix = local.name_prefix
  tags        = local.tags
}

output "ecr_repository_url" { value = module.ecr.repository_url }
output "ecr_repository_name" { value = module.ecr.repository_name }

module "github_oidc" {
  source = "../../modules/github_oidc"

  name_prefix = local.name_prefix
  tags        = local.tags

  github_repo   = "CGC-2026/server"
  github_branch = "main"

  aws_region          = var.aws_region
  ecr_repository_name = module.ecr.repository_name
  ec2_instance_id     = module.compute.instance_id
}

output "github_deploy_role_arn" { value = module.github_oidc.deploy_role_arn }

module "logs" {
  source = "../../modules/logs"

  name_prefix = local.name_prefix
  tags        = local.tags

  ec2_role_name              = module.compute.role_name
  api_log_retentions_in_days = var.api_log_retention_in_days
}

output "api_log_group_name" { value = module.logs.api_log_group_name }
output "api_log_group_arn" { value = module.logs.api_log_group_arn }