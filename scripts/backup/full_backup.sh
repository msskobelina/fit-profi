#!/bin/bash
set -euo pipefail

DATE=$(date +%Y-%m-%d_%H-%M)
BASE_DIR="/data"
FULL_DIR="${BASE_DIR}/full_${DATE}"
KEY_FILE="/keys/xb.key"

mkdir -p "${FULL_DIR}"

echo "[FULL] Starting full backup to ${FULL_DIR}..."

xtrabackup \
  --backup \
  --host="${MYSQL_HOST}" \
  --port="${MYSQL_PORT}" \
  --user="${MYSQL_USER}" \
  --password="${MYSQL_PASSWORD}" \
  --target-dir="${FULL_DIR}" \
  --encrypt=AES256 \
  --encrypt-key-file="${KEY_FILE}"

echo "[FULL] Decrypting backup files..."
find "${FULL_DIR}" -type f -name "*.xbcrypt" -print0 | while IFS= read -r -d '' f; do
  xtrabackup --decrypt=AES256 --encrypt-key-file="${KEY_FILE}" --target-dir="$(dirname "$f")"
done

find "${FULL_DIR}" -type f -name "*.xbcrypt" -delete

echo "[FULL] Preparing backup..."
xtrabackup --prepare --target-dir="${FULL_DIR}"

echo "[FULL] Archiving..."
tar -czf "${FULL_DIR}.tar.gz" -C "${BASE_DIR}" "$(basename "${FULL_DIR}")"

echo "[FULL] Uploading to S3..."
aws s3 cp "${FULL_DIR}.tar.gz" "s3://${S3_BUCKET}/${S3_PREFIX}/full/${DATE}.tar.gz"

echo "[FULL] Keeping only the latest full locally..."
ls -dt ${BASE_DIR}/full_* 2>/dev/null | tail -n +2 | xargs -r rm -rf

rm -f "${FULL_DIR}.tar.gz"

echo "[FULL] Backup completed and uploaded to S3."
