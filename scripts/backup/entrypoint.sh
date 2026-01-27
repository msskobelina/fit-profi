#!/bin/bash
set -e

cat <<EOF > /etc/crontab
SHELL=/bin/bash
PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin

0 3 * * * root /scripts/full_backup.sh >> /var/log/full_backup.log 2>&1
0 * * * * root /scripts/inc_backup.sh >> /var/log/inc_backup.log 2>&1

EOF

touch /var/log/full_backup.log /var/log/inc_backup.log

echo "[ENTRYPOINT] Starting crond..."
crond -n
