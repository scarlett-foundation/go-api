# Prometheus Configuration

This directory contains the Prometheus configuration and management scripts for the Scarlett API monitoring setup.

## Directory Structure

```
prometheus/
├── config/
│   └── prometheus.yml    # Main Prometheus configuration
├── scripts/
│   ├── backup.sh        # Automated backup script
│   └── restore.sh       # Data restoration script
└── README.md            # This file
```

## Backup System

The backup system automatically creates daily backups of Prometheus data at 2 AM. Backups are stored in `/root/backups/prometheus/` on the server.

### Features

- Daily automated backups
- 7-day retention policy
- Compressed storage
- Logging to `/var/log/prometheus-backup.log`
- Safe shutdown/startup during backup

### Manual Backup

To manually trigger a backup:

```bash
/root/backup-prometheus.sh
```

### Restore from Backup

To restore from a backup:

```bash
/root/restore-prometheus.sh /root/backups/prometheus/data_YYYYMMDD_HHMMSS.tar.gz
```

## Configuration

The `prometheus.yml` file contains the main Prometheus configuration. Key settings:

- Scrape interval: 15s
- Evaluation interval: 15s
- Target: API service on port 8082

## Deployment

The configuration is automatically deployed via docker-compose. The Prometheus service:

- Uses the official Prometheus v2.45.0 image
- Mounts the configuration file read-only
- Uses a persistent volume for data
- Includes health checks
- Exposes port 9090 internally

## Maintenance

- Backups are automatically pruned after 7 days
- All backup/restore operations are logged
- The backup system preserves data during container updates 