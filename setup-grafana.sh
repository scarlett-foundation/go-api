#!/bin/bash

# Check for required commands
if ! command -v jq &> /dev/null; then
    echo "Error: jq is required but not installed. Please install jq."
    exit 1
fi

# Ensure we're in the project root
if [ ! -d "./grafana" ]; then
    echo "Error: This script must be run from the project root directory"
    exit 1
fi

# Wait for Grafana to be up
echo "Waiting for Grafana to start..."
until $(curl --output /dev/null --silent --head --fail http://localhost:3000); do
    printf '.'
    sleep 5
done

echo "Grafana is up and running!"

# Get the Prometheus datasource UID
PROM_UID=$(curl -s -u admin:admin http://localhost:3000/api/datasources | jq -r '.[0].uid')

if [ -n "$PROM_UID" ] && [ "$PROM_UID" != "null" ]; then
    echo "Found Prometheus datasource with UID: $PROM_UID"
    
    # Replace the variable in the dashboard file with cross-platform compatibility
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        sed -i '' "s/\${DS_PROMETHEUS}/$PROM_UID/g" ./grafana/dashboards/api_key_metrics.json
        sed -i '' "s/PBFA97CFB590B2093/$PROM_UID/g" ./grafana/dashboards/api_key_metrics.json
    else
        # Linux
        sed -i "s/\${DS_PROMETHEUS}/$PROM_UID/g" ./grafana/dashboards/api_key_metrics.json
        sed -i "s/PBFA97CFB590B2093/$PROM_UID/g" ./grafana/dashboards/api_key_metrics.json
    fi
    
    echo "Updated dashboard with correct datasource UID"
    echo "Grafana is ready to use at: http://localhost:3000"
    echo "Username: admin"
    echo "Password: admin"
else
    echo "Error: Could not find Prometheus datasource UID"
    exit 1
fi 