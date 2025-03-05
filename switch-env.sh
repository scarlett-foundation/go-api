#!/bin/bash

# Script to switch between development and production environments

function usage {
  echo "Usage: $0 [dev|prod]"
  echo "  dev  - Switch to development environment"
  echo "  prod - Switch to production environment"
  exit 1
}

# Check if argument is provided
if [ "$#" -ne 1 ]; then
  usage
fi

case "$1" in
  dev)
    echo "Switching to development environment..."
    cp .env.development .env
    echo "Environment set to development. Use 'docker-compose -f docker-compose.dev.yml up' to start services."
    ;;
  prod)
    echo "Switching to production environment..."
    cp .env.production .env
    echo "Environment set to production. Use 'docker-compose up' to start services."
    ;;
  *)
    usage
    ;;
esac 