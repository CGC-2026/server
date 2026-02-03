#!/usr/bin/env bash
set -euo pipefail

: "${AWS_REGION:?}"
: "${EC2_INSTANCE_ID:?}"


          PUBLIC_IP=$(aws ec2 describe-instances \
            --region "${AWS_REGION}" \
            --instance-ids "${EC2_INSTANCE_ID}" \
            --query 'Reservations[0].Instances[0].PublicIpAddress' \
            --output text)

            test "${PUBLIC_IP}" != "None"

          URL="http://${PUBLIC_IP}:8080/health"
          echo "Checking ${URL}"

          for i in {1..20}; do
            if curl -fsS "${URL}" >/dev/null; then
              echo "Health check passed"
              exit 0
            fi
            echo "Health not ready yet (attempt $i/20) — retrying..."
            sleep 3
          done

          echo "Health check failed"
          exit 1