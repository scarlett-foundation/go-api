# Grafana Dashboard Setup

This directory contains the provisioning configuration for Grafana. The dashboards and data sources are automatically loaded when the container starts up.

## Directory Structure

```
grafana/
├── dashboards/             # Contains the dashboard JSON definitions
│   └── api_key_metrics.json
├── provisioning/
│   ├── dashboards/         # Dashboard provisioning configuration
│   │   └── default.yaml
│   └── datasources/        # Data source provisioning configuration
│       └── prometheus.yaml
└── README.md               # This file
```

## Automatic Setup

When you start the Docker Compose environment, Grafana will automatically:

1. Create the Prometheus data source
2. Load the API Key Metrics dashboard

## Manual Adjustments

If you need to adjust the setup or fix any issues with variable replacement, run the included setup script:

```bash
./setup-grafana.sh
```

This script will:
1. Wait for Grafana to be up and running
2. Get the Prometheus data source UID
3. Update the dashboard JSON with the correct UID

## Accessing Grafana

Once the containers are up and running, you can access Grafana at:

http://localhost:3000

Default credentials:
- Username: admin
- Password: admin

## Adding New Dashboards

To add a new dashboard:

1. Export the dashboard JSON from Grafana UI
2. Save it in the `grafana/dashboards/` directory
3. Make sure it uses the `${DS_PROMETHEUS}` variable for data source UID
4. Restart the containers or run `./setup-grafana.sh`

## Troubleshooting

If your dashboard doesn't show data:
- Check that the Prometheus data source is properly configured
- Verify that the dashboard uses the correct data source UID
- Ensure your metrics are being collected by Prometheus 