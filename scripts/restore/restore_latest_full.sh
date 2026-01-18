#!/bin/bash
set -euo pipefail

S3_BUCKET="${S3_BUCKET:?S3_BUCKET is required}"
S3_PREFIX="${S3_PREFIX:?S3_PREFIX is required}"
S3_FULL_PREFIX="s3://${S3_BUCKET}/${S3_PREFIX}/full"

WORKDIR="/opt/fit-profi/restore"
KEY_FILE="/opt/fit-profi/keys/xb.key"

MYSQL_CONTAINER="fitprofi-mysql"
MYSQL_DATADIR="/var/lib/mysql"

echo "[RESTORE] S3: ${S3_FULL_PREFIX}"
echo "[RESTORE] WORKDIR: ${WORKDIR}"

mkdir -p "${WORKDIR}"
cd "${WORKDIR}"

echo "[RESTORE] Finding latest full backup in S3..."
LATEST_KEY="$(aws s3 ls "${S3_FULL_PREFIX}/" | awk '{print $4}' | sort | tail -n 1)"
if [[ -z "${LATEST_KEY}" ]]; then
  echo "[RESTORE][ERROR] No full backups found in ${S3_FULL_PREFIX}/"
  exit 1
fi

echo "[RESTORE] Latest: ${LATEST_KEY}"
echo "[RESTORE] Downloading..."
aws s3 cp "${S3_FULL_PREFIX}/${LATEST_KEY}" ./full.tar.gz

echo "[RESTORE] Extracting..."
rm -rf ./full_*
tar -xzf full.tar.gz

FULL_DIR="$(find . -maxdepth 1 -type d -name 'full_*' | head -n 1)"
if [[ -z "${FULL_DIR}" ]]; then
  echo "[RESTORE][ERROR] Cannot find extracted full_* directory"
  exit 1
fi

echo "[RESTORE] Extracted dir: ${FULL_DIR}"

echo "[RESTORE] Decrypting *.xbcrypt..."
find "${FULL_DIR}" -type f -name "*.xbcrypt" -print0 | while IFS= read -r -d '' f; do
  docker run --rm \
    -v "${WORKDIR}/${FULL_DIR}:/backup" \
    -v "${KEY_FILE}:/keys/xb.key:ro" \
    percona/percona-xtrabackup:8.0 \
    sh -lc "xtrabackup --decrypt=AES256 --encrypt-key-file=/keys/xb.key --target-dir=/backup/$(dirname "${f#${FULL_DIR}/}")"
done

find "${FULL_DIR}" -type f -name "*.xbcrypt" -delete

echo "[RESTORE] Preparing backup..."
docker run --rm \
  -v "${WORKDIR}/${FULL_DIR}:/backup" \
  percona/percona-xtrabackup:8.0 \
  sh -lc "xtrabackup --prepare --target-dir=/backup"

echo "[RESTORE] Stopping MySQL..."
docker stop "${MYSQL_CONTAINER}"

echo "[RESTORE] Cleaning datadir in volume..."
docker run --rm \
  --volumes-from "${MYSQL_CONTAINER}" \
  alpine:3.19 \
  sh -lc "rm -rf ${MYSQL_DATADIR}/*"

echo "[RESTORE] Copy-back..."
docker run --rm \
  --volumes-from "${MYSQL_CONTAINER}" \
  -v "${WORKDIR}/${FULL_DIR}:/backup" \
  percona/percona-xtrabackup:8.0 \
  sh -lc "xtrabackup --copy-back --target-dir=/backup"

echo "[RESTORE] Fix ownership..."
docker run --rm \
  --volumes-from "${MYSQL_CONTAINER}" \
  alpine:3.19 \
  sh -lc "chown -R 999:999 ${MYSQL_DATADIR}"

echo "[RESTORE] Starting MySQL..."
docker start "${MYSQL_CONTAINER}"

echo "[RESTORE] Done. Logs:"
echo "docker logs -f ${MYSQL_CONTAINER}"
