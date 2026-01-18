#!/bin/bash
set -euo pipefail

echo "[ENTRYPOINT] Starting backup cron..."

cat >/etc/cron.d/fitprofi-mysql-backup <<'CRON'
SHELL=/bin/bash
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

0 3 * * * root /scripts/full_backup.sh >> /var/log/full_backup.log 2>&1
0 * * * * root /scripts/inc_backup.sh  >> /var/log/inc_backup.log 2>&1
CRON

chmod 0644 /etc/cron.d/fitprofi-mysql-backup
touch /var/log/full_backup.log /var/log/inc_backup.log

cron
tail -F /var/log/full_backup.log /var/log/inc_backup.log
