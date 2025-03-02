#!/bin/bash
set -e

# Check if backup file is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <backup_file>"
    echo "Example: $0 /root/backups/prometheus/data_20250302_220044.tar.gz"
    exit 1
fi

BACKUP_FILE=$1
TEMP_DIR=/tmp/prometheus_restore_$(date +%s)
LOG=/var/log/prometheus-backup.log

echo "$(date): Starting Prometheus restore from $BACKUP_FILE" >> $LOG

# Create temp directory
mkdir -p $TEMP_DIR

# Extract backup
echo "$(date): Extracting backup" >> $LOG
tar -xzf $BACKUP_FILE -C $TEMP_DIR

# Get the extracted directory name
EXTRACTED_DIR=$(ls $TEMP_DIR)

echo "$(date): Stopping Prometheus" >> $LOG
cd /root/scarlett
docker-compose stop prometheus

echo "$(date): Restoring data" >> $LOG
docker cp $TEMP_DIR/$EXTRACTED_DIR/. prometheus:/prometheus/

echo "$(date): Starting Prometheus" >> $LOG
docker-compose start prometheus

# Cleanup
echo "$(date): Cleaning up temporary files" >> $LOG
rm -rf $TEMP_DIR

echo "$(date): Restore completed successfully" >> $LOG 