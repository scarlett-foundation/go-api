# Prometheus Monitoring for Scarlett API

This document explains how to use Prometheus to monitor the Scarlett API, including tracking API key usage and other important metrics.

## Overview

The Scarlett API includes built-in Prometheus monitoring to track:
- Total HTTP requests by status code, method, and path
- HTTP request duration
- API key usage (with masked keys for privacy)

## Running with Docker Compose

The easiest way to run the API with Prometheus monitoring is using Docker Compose:

```bash
docker-compose up -d
```

This will start:
- The Scarlett API on port 8082
- Prometheus on port 9090
- Grafana on port 3000 (for visualization)

## Accessing Metrics

### Raw Metrics

You can access the raw Prometheus metrics directly from the API at:

```
http://localhost:8082/metrics
```

### Prometheus UI

Prometheus provides a UI for querying metrics at:

```
http://localhost:9090
```

### Grafana (Visualization)

For better visualization, you can use Grafana at:

```
http://localhost:3000
```

Default login:
- Username: admin
- Password: admin

## Important Metrics

### API Key Usage

Track API key usage with the following Prometheus query:

```
api_key_requests_total
```

This will show the number of requests per (masked) API key. For privacy, API keys are masked to show only the first and last 4 characters.

### HTTP Request Metrics

- Total requests by status code, method, and path:
  ```
  http_requests_total
  ```

- Request duration:
  ```
  http_request_duration_seconds
  ```

## Setting Up Grafana Dashboards

1. Log in to Grafana at http://localhost:3000
2. Add Prometheus as a data source:
   - Go to Configuration > Data Sources
   - Add Prometheus with URL: http://prometheus:9090
3. Import or create dashboards to visualize the metrics

### Automated Dashboard Setup

The API Key Usage dashboard has been automatically configured through the Docker Compose setup. You can access it at:

```
http://localhost:3000/d/feejooeqieio0c/api-key-usage
```

If you need to recreate this dashboard or create it manually, use the following configuration:

```json
{
  "dashboard": {
    "id": null,
    "title": "API Key Usage",
    "tags": ["api", "prometheus"],
    "timezone": "browser",
    "panels": [
      {
        "id": 1,
        "title": "API Key Usage",
        "type": "bargauge",
        "gridPos": {
          "h": 8,
          "w": 24,
          "x": 0,
          "y": 0
        },
        "targets": [
          {
            "expr": "sum by (api_key) (api_key_requests_total)",
            "legendFormat": "{{api_key}}",
            "refId": "A"
          }
        ]
      }
    ],
    "schemaVersion": 16,
    "version": 0
  }
}
```

This can be created via the Grafana API with:

```bash
curl -X POST -H "Content-Type: application/json" -u admin:admin -d '{"dashboard":{"id":null,"title":"API Key Usage","tags":["api","prometheus"],"timezone":"browser","panels":[{"id":1,"title":"API Key Usage","type":"bargauge","gridPos":{"h":8,"w":24,"x":0,"y":0},"targets":[{"expr":"sum by (api_key) (api_key_requests_total)","legendFormat":"{{api_key}}","refId":"A"}]}],"schemaVersion":16,"version":0},"folderId":0,"overwrite":false}' http://localhost:3000/api/dashboards/db
```

## Custom Prometheus Configuration

The default Prometheus configuration is in `prometheus.yml`. You can modify this file to change scrape intervals or add additional targets.

## Security Considerations

- The `/metrics` endpoint is publicly accessible by default. In a production environment, you should secure this endpoint.
- API keys are masked in the metrics to protect sensitive information.

## Troubleshooting

If metrics are not appearing in Prometheus:

1. Check that the API is running and accessible
2. Verify Prometheus is scraping the API by checking the Targets page in Prometheus UI
3. Ensure the network configuration allows Prometheus to reach the API

## Additional Resources

- [Prometheus Documentation](https://prometheus.io/docs/introduction/overview/)
- [Grafana Documentation](https://grafana.com/docs/) 