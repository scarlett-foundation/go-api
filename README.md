# Go Groq API Wrapper

A lean Go API wrapper for the Groq chat completions endpoint using the Echo framework. Fully compatible with OpenAI's API conventions.

## Features

- Wraps Groq's chat completions API
- 100% OpenAI-compatible request/response format
- Supports both streaming and non-streaming responses
- Environment-based configuration
- CORS enabled
- Error handling and logging
- API key authentication
- Swagger/OpenAPI documentation

## Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and set your Groq API key:
   ```
   GROQ_API_KEY=your_api_key_here
   PORT=8080  # Optional, defaults to 8080
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Run the server:
   ```bash
   go run main.go
   ```

## API Usage

### Authentication

The API requires an API key passed in the Authorization header with the Bearer token format:

```
Authorization: Bearer your-api-key
```

**Important**: The "Bearer " prefix is required. Requests without this prefix will be rejected with a 401 Unauthorized error.

### Chat Completions Endpoint

**Endpoint:** `POST /chat/completions`

**Request Body:**
```json
{
  "messages": [
    {
      "role": "system",
      "content": "you are a helpful agent"
    },
    {
      "role": "user",
      "content": "Hello!"
    }
  ],
  "model": "deepseek-r1-distill-llama-70b",
  "temperature": 0.6,
  "max_tokens": 4096,
  "top_p": 0.95,
  "frequency_penalty": 0.0,
  "presence_penalty": 0.0,
  "n": 1,
  "stream": true,
  "user": "user-123"
}
```

**Response Format (non-streaming):**
```json
{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677652288,
  "model": "deepseek-r1-distill-llama-70b",
  "system_fingerprint": "fp-1234",
  "choices": [{
    "index": 0,
    "message": {
      "role": "assistant",
      "content": "Hello! How can I help you today?"
    },
    "finish_reason": "stop"
  }],
  "usage": {
    "prompt_tokens": 9,
    "completion_tokens": 12,
    "total_tokens": 21
  }
}
```

**Example curl:**
```bash
curl -X POST "http://localhost:8080/chat/completions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-api-key" \
  -d '{
    "messages": [
      {
        "role": "system",
        "content": "you are a helpful agent"
      },
      {
        "role": "user",
        "content": "Hello!"
      }
    ],
    "model": "deepseek-r1-distill-llama-70b",
    "temperature": 0.6,
    "stream": true
  }'
```

## Error Handling

The API returns OpenAI-compatible error responses:

```json
{
  "error": {
    "message": "Error message here",
    "type": "invalid_request_error",
    "param": null,
    "code": null
  }
}
```

Error types include:
- `invalid_request_error`: Invalid request parameters
- `api_error`: Error communicating with Groq API
- `internal_error`: Server-side errors

## License

MIT License