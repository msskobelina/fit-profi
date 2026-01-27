#!/bin/sh
set -e

echo "[ENTRYPOINT] Setting up cron jobs..."

mkdir -p /etc/crontabs

cat > /etc/crontabs/root <<'EOF'
0 3 * * * /scripts/full_backup.sh >> /var/log/full_backup.log 2>&1
0 * * * * /scripts/inc_backup.sh >> /var/log/inc_backup.log 2>&1
EOF

chmod 600 /etc/crontabs/root
touch /var/log/full_backup.log /var/log/inc_backup.log

echo "[ENTRYPOINT] Cron jobs installed:"
cat /etc/crontabs/root

echo "[ENTRYPOINT] Starting crond..."
exec crond -n
