#!/bin/bash

# Start all services
echo "Starting all services..."
if ! docker-compose up -d; then
    echo "Error starting services. Check the docker-compose.yml file."
    exit 1
fi

# Wait a moment for services to start
echo "Waiting for services to initialize..."
sleep 10

# Run the Grafana setup script
echo "Setting up Grafana dashboards..."
if ! ./setup-grafana.sh; then
    echo "Error setting up Grafana dashboards. Check the setup-grafana.sh script."
    exit 1
fi

echo "Setup complete! Your monitoring stack is now ready."
echo "API available at: http://localhost:8082"
echo "Prometheus available at: http://localhost:9090"
echo "Grafana available at: http://localhost:3000" 