terraform {
  backend "s3" {
    bucket         = "cgc-2026-tfstate-nick"
    key            = "prod/infra.tfstate"
    region         = "ca-central-1"
    dynamodb_table = "cgc-2026-tf-locks"
    encrypt        = true
  }
}