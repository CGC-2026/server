#!/usr/bin/env bash
set -euo pipefail
          
: "${AWS_REGION:?}"
: "${ECR_REPO:?}"
: "${GITHUB_SHA:?}"
: "${EC2_INSTANCE_ID:?}"

ECR_REGISTRY="$(aws sts get-caller-identity --query Account --output text).dkr.ecr.${AWS_REGION}.amazonaws.com"
MIGRATOR_IMAGE_URI="${ECR_REGISTRY}/${ECR_REPO}:${GITHUB_SHA}-migrate"

COMMANDS_JSON=$(cat <<JSON
[
    "set -euo pipefail",
    "AWS_REGION=${AWS_REGION}",
    "MIGRATOR_IMAGE_URI=${MIGRATOR_IMAGE_URI}",
    "if ! command -v docker >/dev/null 2>&1; then sudo dnf -y install docker; sudo systemctl enable --now docker; fi",
    "ECR_REGISTRY=\${MIGRATOR_IMAGE_URI%%/*}",
    "aws ecr get-login-password --region $AWS_REGION | sudo docker login --username AWS --password-stdin \$ECR_REGISTRY",
    "DBURL=\$(aws ssm get-parameter --region $AWS_REGION --with-decryption --name /cgc-2026-prod/api/database_url --query Parameter.Value --output text)",
    "sudo docker pull $MIGRATOR_IMAGE_URI",
    "sudo docker run --rm -e GOOSE_DRIVER=postgres -e GOOSE_DBSTRING=\"\$DBURL\" $MIGRATOR_IMAGE_URI -dir ./db/migrations up"
]
JSON
)

./scripts/ci/ssm_run_and_wait.sh "${COMMANDS_JSON}"