#!/bin/bash
set -e

BACKUP_DIR=/root/backups/prometheus
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_PATH=$BACKUP_DIR/data_$TIMESTAMP
LOG=/var/log/prometheus-backup.log

echo "$(date): Starting Prometheus backup" >> $LOG
cd /root/scarlett

echo "$(date): Stopping Prometheus" >> $LOG
docker-compose stop prometheus

echo "$(date): Creating backup" >> $LOG
docker cp prometheus:/prometheus $BACKUP_PATH

echo "$(date): Starting Prometheus" >> $LOG
docker-compose start prometheus

echo "$(date): Compressing backup" >> $LOG
tar -czf $BACKUP_PATH.tar.gz -C $BACKUP_DIR data_$TIMESTAMP
rm -rf $BACKUP_PATH

echo "$(date): Cleaning old backups (older than 7 days)" >> $LOG
find $BACKUP_DIR -name "data_*.tar.gz" -mtime +7 -delete

echo "$(date): Backup completed successfully: data_$TIMESTAMP.tar.gz" >> $LOG 