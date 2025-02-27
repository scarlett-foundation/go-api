# Go API Developer Documentation

## Overview
This is a Go-based API service that acts as a proxy to the Groq API, providing chat completion functionality. The service is built using the Echo web framework and follows idiomatic Go practices.

## Project Structure
```
.
├── cmd/            # Command line tools
├── internal/       # Private application code
│   ├── api/       # API-specific implementations
│   └── types/     # Internal type definitions
├── pkg/           # Public libraries
│   └── groq/      # Groq API integration
├── main.go        # Application entry point
├── go.mod         # Go module definition
├── go.sum         # Go module checksums
└── .env           # Environment configuration
```

## Prerequisites
- Go 1.x (version specified in go.mod)
- Environment variables configured (see Configuration section)

## Configuration
The application uses environment variables for configuration. Create a `.env` file in the root directory with the following variables:

```env
PORT=8080              # Optional: Default is 8080
GROQ_API_KEY=your_key  # Required: Your Groq API key
```

An example configuration is provided in `.env.example`.

## API Endpoints

### POST /chat/completions
Proxies requests to Groq's chat completions API.

#### Request Body
```json
{
  "messages": [
    {
      "role": "string",
      "content": "string"
    }
  ],
  "model": "string",
  "temperature": float,
  "max_tokens": integer,
  "top_p": float,
  "frequency_penalty": float,
  "presence_penalty": float,
  "stream": boolean,
  "stop": ["string"],
  "n": integer,
  "user": "string"
}
```

#### Response
Returns the Groq API response directly, including:
```json
{
  "id": "string",
  "object": "string",
  "created": integer,
  "model": "string",
  "system_fingerprint": "string",
  "choices": [
    {
      "index": integer,
      "message": {
        "role": "string",
        "content": "string"
      },
      "finish_reason": "string"
    }
  ],
  "usage": {
    "prompt_tokens": integer,
    "completion_tokens": integer,
    "total_tokens": integer
  }
}
```

## Examples

### Basic Chat Completion Request
```bash
curl -X POST http://localhost:8080/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "deepseek-coder-33b-instruct",
    "messages": [
      {
        "role": "user",
        "content": "What is the capital of France?"
      }
    ],
    "temperature": 0.7,
    "max_tokens": 150
  }'
```

### Streaming Response Request
```bash
curl -X POST http://localhost:8082/chat/completions \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -d '{
    "model": "deepseek-r1-distill-llama-70b",
    "messages": [
      {
        "role": "user",
        "content": "Suggest a dogecoin trading strategy for the next 30 days"
      }
    ],
    "stream": true,
    "temperature": 0.7,
    "max_tokens": 150
  }'
```

## Error Handling
The API returns standard HTTP status codes and JSON error responses:

```json
{
  "error": {
    "message": "string",
    "type": "string",
    "param": "string",
    "code": "string"
  }
}
```

Common error types:
- `invalid_request_error`: Invalid request body or parameters
- `api_error`: Error communicating with Groq API
- `internal_error`: Server-side errors

## Development

### Running Locally
1. Clone the repository
2. Copy `.env.example` to `.env` and configure your environment variables
3. Run the application:
   ```bash
   go run main.go
   ```

### Building
```bash
go build -o api-server
```

### Testing
```bash
go test ./...
```

## Streaming Support
The API supports both streaming and non-streaming responses. To use streaming:
1. Set `stream: true` in your request
2. Handle Server-Sent Events (SSE) in your client
3. The response will be streamed as chunks of data

## Security Considerations
- API keys are read from environment variables
- CORS middleware is enabled
- Request validation is performed
- Error messages are sanitized

## Dependencies
Key dependencies include:
- `github.com/labstack/echo/v4`: Web framework
- `github.com/joho/godotenv`: Environment variable management

## Contributing
1. Fork the repository
2. Create a feature branch
3. Commit your changes using conventional commit messages
4. Push to your branch
5. Create a Pull Request

## License
See LICENSE file in the repository root. 