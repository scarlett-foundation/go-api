#!/bin/bash

# Make the script executable
chmod +x setup-grafana.sh

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
    
    # Replace the variable in the dashboard file
    sed -i '' "s/\${DS_PROMETHEUS}/$PROM_UID/g" ./grafana/dashboards/api_key_metrics.json
    
    echo "Updated dashboard with correct datasource UID"
    echo "Grafana is ready to use at: http://localhost:3000"
    echo "Username: admin"
    echo "Password: admin"
else
    echo "Error: Could not find Prometheus datasource UID"
fi 