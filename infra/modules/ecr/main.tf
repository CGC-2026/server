resource "aws_ecr_repository" "this" {
    name = "${var.name_prefix}-backend"
    image_tag_mutability = "IMMUTABLE"

    image_scanning_configuration {
      scan_on_push = true
    }

    tags = merge(var.tags, {
        Name = "${var.name_prefix}-backend-ecr"
    })
}

resource "aws_ecr_lifecycle_policy" "this" {
  repository = aws_ecr_repository.this.name

  policy = core::jsonencode({
    rules = [{
        rulePriority = 1
        description = "Keep last ${var.keep_last_images} images"
        selection = {
            tagStatus = "any"
            countType = "imageCountMoreThan"
            countNumber = var.keep_last_images
        }
        action = { type = "expire" }
    }]
  })
}