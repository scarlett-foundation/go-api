# Prometheus Monitoring for Scarlett API

This document explains how to use Prometheus to monitor the Scarlett API, including tracking API key usage and other important metrics.

## Overview

The Scarlett API includes built-in Prometheus monitoring to track:
- Total HTTP requests by status code, method, and path
- HTTP request duration
- API key usage (with masked keys for privacy)
- Token usage from LLM inference requests (prompt, completion, and total tokens)

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

#### API Key Masking

The API uses the following rules for masking API keys in metrics:

- Keys longer than 8 characters: First 4 characters + "..." + last 4 characters (e.g., "test...y123")
- Keys 8 characters or shorter: Labeled as "short_key" to prevent identification
- Missing or malformed keys: Labeled as "unknown"

The "short_key" label in metrics represents any API key that was too short to apply the standard masking pattern. This could include both valid and invalid keys that are 8 characters or fewer in length.

### Token Usage Metrics

The API tracks token usage for each request, broken down by API key:

- Prompt tokens (input tokens):
  ```
  token_usage_prompt_total
  ```

- Completion tokens (output tokens):
  ```
  token_usage_completion_total
  ```

- Total tokens:
  ```
  token_usage_total
  ```

You can use these metrics to track token consumption by different API keys, monitor costs, and plan capacity.

### HTTP Request Metrics

- Total requests by status code, method, and path:
  ```
  http_requests_total
  ```

- Request duration:
  ```
  http_request_duration_seconds
  ```

- Average request duration:
  ```
  http_request_duration_seconds_sum / http_request_duration_seconds_count
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
          "w": 12,
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
      },
      {
        "id": 2,
        "title": "Token Usage by API Key",
        "type": "bargauge",
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 0
        },
        "targets": [
          {
            "expr": "sum by (api_key) (token_usage_total)",
            "legendFormat": "{{api_key}}",
            "refId": "A"
          }
        ]
      },
      {
        "id": 3,
        "title": "Prompt vs Completion Tokens",
        "type": "timeseries",
        "gridPos": {
          "h": 8,
          "w": 24,
          "x": 0,
          "y": 8
        },
        "targets": [
          {
            "expr": "sum(token_usage_prompt_total)",
            "legendFormat": "Prompt Tokens",
            "refId": "A"
          },
          {
            "expr": "sum(token_usage_completion_total)",
            "legendFormat": "Completion Tokens",
            "refId": "B"
          }
        ]
      }
    ],
    "schemaVersion": 16,
    "version": 0
  }
}
```

### Token Usage Dashboard

To create a dedicated dashboard for token usage metrics, use the following configuration:

```bash
curl -X POST -H "Content-Type: application/json" -u admin:admin -d '{
  "dashboard": {
    "id": null,
    "title": "Token Usage Metrics",
    "tags": ["api", "prometheus", "tokens"],
    "timezone": "browser",
    "panels": [
      {
        "id": 1,
        "title": "Total Tokens Used by API Key",
        "type": "piechart",
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 0,
          "y": 0
        },
        "targets": [
          {
            "expr": "sum by (api_key) (token_usage_total)",
            "legendFormat": "{{api_key}}",
            "refId": "A"
          }
        ]
      },
      {
        "id": 2,
        "title": "Prompt vs Completion Tokens",
        "type": "piechart",
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 0
        },
        "targets": [
          {
            "expr": "sum(token_usage_prompt_total)",
            "legendFormat": "Prompt",
            "refId": "A"
          },
          {
            "expr": "sum(token_usage_completion_total)",
            "legendFormat": "Completion",
            "refId": "B"
          }
        ]
      },
      {
        "id": 3,
        "title": "Token Usage Over Time",
        "type": "timeseries",
        "gridPos": {
          "h": 8,
          "w": 24,
          "x": 0,
          "y": 8
        },
        "targets": [
          {
            "expr": "sum(rate(token_usage_total[5m]))",
            "legendFormat": "Tokens per second (5m avg)",
            "refId": "A"
          }
        ]
      },
      {
        "id": 4,
        "title": "Input/Output Token Ratio by API Key",
        "type": "bargauge",
        "gridPos": {
          "h": 8,
          "w": 24,
          "x": 0,
          "y": 16
        },
        "options": {
          "orientation": "horizontal"
        },
        "targets": [
          {
            "expr": "sum by (api_key) (token_usage_completion_total) / sum by (api_key) (token_usage_prompt_total)",
            "legendFormat": "{{api_key}}",
            "refId": "A"
          }
        ]
      }
    ],
    "schemaVersion": 16,
    "version": 0
  },
  "folderId": 0,
  "overwrite": false
}' http://localhost:3000/api/dashboards/db
```

## Custom Prometheus Configuration

The default Prometheus configuration is in `prometheus.yml`. You can modify this file to change scrape intervals or add additional targets.

## Security Considerations

- The `/metrics` endpoint is publicly accessible by default. In a production environment, you should secure this endpoint.
- API keys are masked in the metrics to protect sensitive information.
- The masking mechanism ensures that full API keys are never exposed in metrics.

## Troubleshooting

If metrics are not appearing in Prometheus:

1. Check that the API is running and accessible
2. Verify Prometheus is scraping the API by checking the Targets page in Prometheus UI
3. Ensure the network configuration allows Prometheus to reach the API

### Common Issues

- **No API key metrics:** Make sure requests include a properly formatted Authorization header (`Bearer your-api-key`)
- **Missing metrics endpoint:** Check that `RegisterPrometheusHandler` is being called in main.go
- **"short_key" metrics high:** Could indicate potential brute force attempts with short keys
- **No token metrics:** Only successful chat completion requests (HTTP 200) will record token metrics

## Additional Resources

- [Prometheus Documentation](https://prometheus.io/docs/introduction/overview/)
- [Grafana Documentation](https://grafana.com/docs/) 