
resource "aws_cloudwatch_log_group" "api" {
  name              = "/${var.name_prefix}/api"
  retention_in_days = var.api_log_retentions_in_days
  tags              = var.tags
}

data "aws_iam_policy_document" "ec2_cloudwatch_logs" {
  statement {
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
      "logs:DescribeLogStreams"
    ]

    resources = [
      aws_cloudwatch_log_group.api.arn,
      "${aws_cloudwatch_log_group.api.arn}:*"
    ]
  }
}

resource "aws_iam_policy" "ec2_cloudwatch_logs" {
  name   = "${var.name_prefix}-ec2-cloudwatch-logs"
  policy = data.aws_iam_policy_document.ec2_cloudwatch_logs.json
}

resource "aws_iam_role_policy_attachment" "ec2_cloudwatch_logs" {
  role       = var.ec2_role_name
  policy_arn = aws_iam_policy.ec2_cloudwatch_logs.arn
}