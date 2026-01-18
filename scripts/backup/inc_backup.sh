#!/bin/bash
set -euo pipefail

DATE=$(date +%Y-%m-%d_%H-%M)
BASE_DIR="/data"
INC_DIR="${BASE_DIR}/inc_${DATE}"
KEY_FILE="/keys/xb.key"

mkdir -p "${INC_DIR}"

LAST_FULL_DIR=$(ls -dt ${BASE_DIR}/full_* 2>/dev/null | head -n1 || true)
if [ -z "${LAST_FULL_DIR}" ]; then
  echo "[INC] No full backup found locally. Skipping incremental."
  exit 0
fi

echo "[INC] Using base-dir=${LAST_FULL_DIR}"
echo "[INC] Starting incremental backup to ${INC_DIR}..."

xtrabackup \
  --backup \
  --host="${MYSQL_HOST}" \
  --port="${MYSQL_PORT}" \
  --user="${MYSQL_USER}" \
  --password="${MYSQL_PASSWORD}" \
  --target-dir="${INC_DIR}" \
  --incremental-basedir="${LAST_FULL_DIR}" \
  --encrypt=AES256 \
  --encrypt-key-file="${KEY_FILE}"

tar -czf "${INC_DIR}.tar.gz" -C "${BASE_DIR}" "$(basename "${INC_DIR}")"
aws s3 cp "${INC_DIR}.tar.gz" "s3://${S3_BUCKET}/${S3_PREFIX}/inc/${DATE}.tar.gz"

rm -rf "${INC_DIR}" "${INC_DIR}.tar.gz"

echo "[INC] Incremental backup completed and uploaded to S3."
