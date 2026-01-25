locals {
    name_prefix = "${var.project}-${var.environment}"
}

output "name_prefix" {
    value = local.name_prefix
}