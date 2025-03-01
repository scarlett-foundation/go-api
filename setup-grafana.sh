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
    
    # Function to import a dashboard
    import_dashboard() {
        local dashboard_file=$1
        local dashboard_name=$2
        
        echo "Importing $dashboard_name dashboard..."
        TEMP_DASHBOARD=$(mktemp)
        cat "$dashboard_file" | sed "s/\${DS_PROMETHEUS}/$PROM_UID/g" > "$TEMP_DASHBOARD"
        
        RESULT=$(curl -s -X POST -H "Content-Type: application/json" \
          -d "{\"dashboard\": $(cat $TEMP_DASHBOARD), \"overwrite\": true}" \
          -u admin:admin http://localhost:3000/api/dashboards/db)
        
        rm "$TEMP_DASHBOARD"
        
        if echo "$RESULT" | grep -q '"status":"success"'; then
            echo "✓ Successfully imported $dashboard_name"
        else
            echo "✗ Failed to import $dashboard_name"
            echo "$RESULT"
        fi
    }
    
    # Import all dashboards
    import_dashboard "./grafana/dashboards/api_key_metrics.json" "API Key Metrics"
    import_dashboard "./grafana/dashboards/token_usage_metrics.json" "Token Usage Overview"
    
    # Disable public dashboards feature globally in Grafana
    echo "Disabling public dashboards feature globally..."
    curl -s -X PUT -H "Content-Type: application/json" \
      -d '{"enabled": false}' \
      -u admin:admin http://localhost:3000/api/admin/settings/public-dashboards > /dev/null
    
    echo "Setup complete!"
    echo "Grafana is ready to use at: http://localhost:3000"
    echo "Username: admin"
    echo "Password: admin"
else
    echo "Error: Could not find Prometheus datasource UID"
    exit 1
fi 