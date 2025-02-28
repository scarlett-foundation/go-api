# API Documentation

This directory contains the Swagger documentation for the Go API.

## Overview

The API provides a proxy to the Groq API, offering endpoints for chat completions. Documentation is generated using Swagger/OpenAPI.

## Accessing Documentation

When the server is running, you can access the Swagger UI at:

```
http://localhost:8082/swagger/index.html
```

This UI provides interactive documentation where you can:
- Explore available endpoints
- See request and response formats
- Test API calls directly from the browser

## Authentication

API calls require authentication using an API key in the Authorization header:

```
Authorization: Bearer your-api-key
```

Valid API keys are configured in `api-keys.yaml` in the project root.

## Example Usage

```bash
curl -X POST http://localhost:8082/chat/completions \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "deepseek-r1-distill-llama-70b",
    "messages": [
      {
        "role": "user",
        "content": "Hello, how are you?"
      }
    ],
    "temperature": 0.7,
    "max_tokens": 50
  }'
```

## Regenerating Documentation

To update the Swagger documentation after making changes to the API, run:

```bash
make swagger
```

This will scan the codebase for Swagger annotations and regenerate the documentation files. 