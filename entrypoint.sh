#!/bin/bash

# Start all services
echo "Starting all services..."
docker-compose up -d

# Wait a moment for services to start
echo "Waiting for services to initialize..."
sleep 10

# Run the Grafana setup script
echo "Setting up Grafana dashboards..."
./setup-grafana.sh

echo "Setup complete! Your monitoring stack is now ready."
echo "API available at: http://localhost:8082"
echo "Prometheus available at: http://localhost:9090"
echo "Grafana available at: http://localhost:3000" 