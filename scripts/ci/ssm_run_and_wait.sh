#!/usr/bin/env bash
set -euo pipefail

: "${AWS_REGION:?}"
: "${EC2_INSTANCE_ID:?}"


COMMANDS_JSON="${1:?usage: ssm_run_and_wait.sh '<json-array>'}"

COMMAND_ID=$(aws ssm send-command \
            --region "${AWS_REGION}" \
            --instance-ids "${EC2_INSTANCE_ID}" \
            --document-name "AWS-RunShellScript" \
            --parameters "commands=${COMMANDS_JSON}" \
            --query 'Command.CommandId' \
            --output text)

          echo "Waiting for SSM command ${COMMAND_ID} to complete..."

          for i in {1..60}; do
            STATUS=$(aws ssm get-command-invocation \
              --region "${AWS_REGION}" \
              --command-id "${COMMAND_ID}" \
              --instance-id "${EC2_INSTANCE_ID}" \
              --query 'Status' \
              --output text 2>/dev/null || echo "Pending")
            
            case "${STATUS}" in
              Success)
                echo "SSM command succeeded"
                exit 0
                ;;
              Failed|Cancelled|TimedOut)
                echo "SSM command failed with status: ${STATUS}"
                aws ssm get-command-invocation \
                  --region "${AWS_REGION}" \
                  --command-id "${COMMAND_ID}" \
                  --instance-id "${EC2_INSTANCE_ID}" \
                  --query 'StandardErrorContent' \
                  --output text
                exit 1
                ;;
              *)
                echo "Status: ${STATUS} (attempt $i/60)"
                sleep 5
                ;;
            esac
          done

          echo "Timeout: SSM command did not reach terminal status"
          exit 1
