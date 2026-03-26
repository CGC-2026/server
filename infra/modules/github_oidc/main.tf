
data "aws_caller_identity" "current" {}

# Github OIDC provider
resource "aws_iam_openid_connect_provider" "github" {
  url = "https://token.actions.githubusercontent.com"

  client_id_list = [
    "sts.amazonaws.com"
  ]

  # Thumbprint required for github OIDC
  thumbprint_list = ["6938fd4d98bab03faadb97b34396831e3780aea1"]

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-github-oidc"
  })
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    actions = ["sts:AssumeRoleWithWebIdentity"]
    effect  = "Allow"

    principals {
      type        = "Federated"
      identifiers = [aws_iam_openid_connect_provider.github.arn]
    }

    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:aud"
      values   = ["sts.amazonaws.com"]
    }

    # Restrict to this repo and branch
    condition {
      test     = "StringLike"
      variable = "token.actions.githubusercontent.com:sub"
      values = [
        "repo:${var.github_repo}:ref:refs/heads/${var.github_branch}"
      ]
    }
  }
}

resource "aws_iam_role" "deploy" {
  name               = "${var.name_prefix}-github-deploy"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-github-deploy"
  })
}

# Allow pushing to ECR and sending ssm commands to ec2 instance
data "aws_iam_policy_document" "deploy" {
  statement {
    sid       = "EcrAuth"
    effect    = "Allow"
    actions   = ["ecr:GetAuthorizationToken"]
    resources = ["*"]
  }

  statement {
    sid    = "EcrPushPullRepo"
    effect = "Allow"
    actions = [
      "ecr:BatchCheckLayerAvailability",
      "ecr:BatchGetImage",
      "ecr:CompleteLayerUpload",
      "ecr:DescribeRepositories",
      "ecr:InitiateLayerUpload",
      "ecr:PutImage",
      "ecr:UploadLayerPart",
      "ecr:GetDownloadUrlForLayer"
    ]
    resources = [
      "arn:aws:ecr:${var.aws_region}:${data.aws_caller_identity.current.account_id}:repository/${var.ecr_repository_name}"
    ]
  }

  statement {
    sid    = "SsmSendCommand"
    effect = "Allow"
    actions = [
      "ssm:SendCommand",
      "ssm:GetCommandInvocation",
      "ssm:ListCommandInvocations"
    ]
    resources = [
      "*"
    ]
  }

  statement {
    sid    = "DescribeForDeploy"
    effect = "Allow"
    actions = [
      "ec2:DescribeInstances"
    ]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "deploy" {
  name   = "${var.name_prefix}-github-deploy-policy"
  policy = data.aws_iam_policy_document.deploy.json
  tags   = var.tags
}

resource "aws_iam_role_policy_attachment" "attach" {
  role       = aws_iam_role.deploy.name
  policy_arn = aws_iam_policy.deploy.arn
}

# Allow manual control of resources through github actions
data "aws_iam_policy_document" "github_start_stop" {

  statement {
    sid    = "EC2Describe"
    effect = "Allow"
    actions = [
      "ec2:DescribeInstances",
      "ec2:DescribeInstanceStatus",
      "ec2:DescribeTags"
    ]
    resources = ["*"]
  }

  statement {
    sid    = "EC2StartStop"
    effect = "Allow"
    actions = [
      "ec2:StartInstances",
      "ec2:StopInstances"
    ]
    resources = ["arn:aws:ec2:*:*:instance/*"]

    condition {
      test     = "StringEquals"
      variable = "ec2:ResourceTag/Project"
      values   = ["cgc-2026"]
    }
  }

  statement {
    sid    = "RDSDescribe"
    effect = "Allow"
    actions = [
      "rds:DescribeDBInstances",
      "rds:ListTagsForResource"
    ]
    resources = ["*"]
  }

  statement {
    sid    = "RDSStartStop"
    effect = "Allow"
    actions = [
      "rds:StartDBInstance",
      "rds:StopDBInstance"
    ]
    resources = ["arn:aws:rds:*:*:db:*"]

    condition {
      test     = "StringEquals"
      variable = "aws:ResourceTag/Project"
      values   = ["cgc-2026"]
    }
  }
}

resource "aws_iam_policy" "github_start_stop" {
  name   = "${var.name_prefix}-start-stop-env"
  policy = data.aws_iam_policy_document.github_start_stop.json
  tags   = var.tags
}

resource "aws_iam_role_policy_attachment" "github_start_stop" {
  role       = aws_iam_role.deploy.name
  policy_arn = aws_iam_policy.github_start_stop.arn
}