#!/usr/bin/env bash
set -euo pipefail

: "${AWS_REGION:?}"
: "${ECR_REPO:?}"
: "${GITHUB_SHA:?}"

ECR_REGISTRY="$(aws sts get-caller-identity --query Account --output text).dkr.ecr.${AWS_REGION}.amazonaws.com"
          APP_IMAGE_URI="${ECR_REGISTRY}/${ECR_REPO}:${GITHUB_SHA}"
          MIGRATOR_IMAGE_URI="${ECR_REGISTRY}/${ECR_REPO}:${GITHUB_SHA}-migrate"

          docker build -f Dockerfile.prod --target app -t "${APP_IMAGE_URI}" .
          docker push "${APP_IMAGE_URI}"

          docker build -f Dockerfile.prod --target migrator -t "${MIGRATOR_IMAGE_URI}" .
          docker push "${MIGRATOR_IMAGE_URI}"