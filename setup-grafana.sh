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
    
    # Fix corrupted queries in the dashboard
    echo "Fixing corrupted Prometheus queries in the dashboard..."
    
    # Get the dashboard from Grafana
    DASHBOARD_JSON=$(curl -s -u admin:admin http://localhost:3000/api/dashboards/uid/api-key-metrics)
    DASHBOARD_DATA=$(echo $DASHBOARD_JSON | jq -r '.dashboard')
    
    # Create a temporary file with the dashboard JSON
    TMP_DASHBOARD_FILE=$(mktemp)
    echo $DASHBOARD_DATA > $TMP_DASHBOARD_FILE
    
    # Fix the corrupted queries by setting them to the correct expressions
    # Panel 2 (API Calls per Hour)
    jq '.panels[2].targets[0].expr = "sum(rate(api_key_requests_total[1h]) * 3600)"' $TMP_DASHBOARD_FILE > tmp.json && mv tmp.json $TMP_DASHBOARD_FILE
    jq '.panels[2].targets[1].expr = "sum by(api_key) (rate(api_key_requests_total[1h]) * 3600)"' $TMP_DASHBOARD_FILE > tmp.json && mv tmp.json $TMP_DASHBOARD_FILE
    
    # Panel 3 (Tokens per Hour)
    jq '.panels[3].targets[0].expr = "sum(rate(token_usage_total[1h]) * 3600)"' $TMP_DASHBOARD_FILE > tmp.json && mv tmp.json $TMP_DASHBOARD_FILE
    jq '.panels[3].targets[1].expr = "sum by(api_key) (rate(token_usage_total[1h]) * 3600)"' $TMP_DASHBOARD_FILE > tmp.json && mv tmp.json $TMP_DASHBOARD_FILE
    
    # Disable public dashboards explicitly
    echo "Disabling public dashboards in the dashboard..."
    jq '.publicDashboardEnabled = false' $TMP_DASHBOARD_FILE > tmp.json && mv tmp.json $TMP_DASHBOARD_FILE
    
    # Update the dashboard in Grafana
    curl -s -X POST -H "Content-Type: application/json" \
      -d "{\"dashboard\": $(cat $TMP_DASHBOARD_FILE), \"overwrite\": true}" \
      -u admin:admin http://localhost:3000/api/dashboards/db > /dev/null
    
    rm $TMP_DASHBOARD_FILE
    
    # Disable public dashboards feature globally in Grafana
    echo "Disabling public dashboards feature globally..."
    
    # In Grafana 9.x+, we can disable the feature via API
    curl -s -X PUT -H "Content-Type: application/json" \
      -d '{"enabled": false}' \
      -u admin:admin http://localhost:3000/api/admin/settings/public-dashboards > /dev/null
    
    echo "Public dashboards have been explicitly disabled."
    echo "Corrupted Prometheus queries have been fixed."
    echo "Sharing dashboards via the Share button will still work (link sharing, snapshot, etc.)"
    
    echo "Grafana is ready to use at: http://localhost:3000"
    echo "Username: admin"
    echo "Password: admin"
else
    echo "Error: Could not find Prometheus datasource UID"
    exit 1
fi 