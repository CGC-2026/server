#!/usr/bin/env bash

set -euo pipefail

: "${AWS_REGION:?}"
: "${ECR_REPO:?}"
: "${GITHUB_SHA:?}"
: "${EC2_INSTANCE_ID:?}"
: "${CONTAINER_NAME:?}"
: "${LOG_GROUP_NAME:?}"

ECR_REGISTRY="$(aws sts get-caller-identity --query Account --output text).dkr.ecr.${AWS_REGION}.amazonaws.com"
APP_IMAGE_URI="${ECR_REGISTRY}/${ECR_REPO}:${GITHUB_SHA}"

echo "Deploying image ${APP_IMAGE_URI} to instance"

COMMANDS_JSON=$(cat <<JSON
[
    "set -euo pipefail",
    "AWS_REGION=${AWS_REGION}",
    "APP_IMAGE_URI=${APP_IMAGE_URI}",
    "CONTAINER_NAME=${CONTAINER_NAME}",
    "LOG_GROUP_NAME=${LOG_GROUP_NAME}",
    "echo APP_IMAGE_URI=$APP_IMAGE_URI",
    "if ! command -v docker >/dev/null 2>&1; then sudo dnf -y install docker; sudo systemctl enable --now docker; fi",
    "sudo usermod -aG docker ec2-user || true",
    "ECR_REGISTRY=\${APP_IMAGE_URI%%/*}",
    "aws ecr get-login-password --region $AWS_REGION | sudo docker login --username AWS --password-stdin \$ECR_REGISTRY",
    "DBURL=\$(aws ssm get-parameter --region $AWS_REGION --with-decryption --name /cgc-2026-prod/api/database_url --query Parameter.Value --output text)",
    "CLERK_SECRET_KEY=\$(aws ssm get-parameter --region $AWS_REGION --with-decryption --name /cgc-2026-prod/api/clerk_secret_key --query Parameter.Value --output text)",
    "CLERK_WEBHOOK_SECRET=\$(aws ssm get-parameter --region $AWS_REGION --with-decryption --name /cgc-2026-prod/api/clerk_webhook_secret --query Parameter.Value --output text)",
    "sudo docker pull $APP_IMAGE_URI",
    "sudo docker rm -f $CONTAINER_NAME || true",
    "sudo docker run -d --restart unless-stopped --name $CONTAINER_NAME -p 8080:8080 --log-driver=awslogs --log-opt awslogs-region=${AWS_REGION} --log-opt awslogs-group=${LOG_GROUP_NAME} --log-opt awslogs-stream=${CONTAINER_NAME}  -e PORT=8080 -e DATABASE_URL=\"\$DBURL\" -e CLERK_SECRET_KEY=\"\$CLERK_SECRET_KEY\" -e CLERK_WEBHOOK_SECRET=\"\$CLERK_WEBHOOK_SECRET\" $APP_IMAGE_URI",
    "sleep 2",
    "curl -fsS http://localhost:8080/health"
]
JSON
)

bash ./scripts/ci/ssm_run_and_wait.sh "${COMMANDS_JSON}"