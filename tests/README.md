# API Testing

This directory contains tests for the Go API.

## Test Structure

- `integration/`: Contains integration tests that run against a live instance of the API
  - Tests real behavior of the system without mocks
  - Verifies API authentication is working properly
  - Tests against valid and invalid API keys

## Running Tests

### Integration Tests

To run the integration tests:

```bash
make test-integration
```

This will:
1. Build a test version of the API server
2. Start the server
3. Run the tests against the live server
4. Shut down the server when tests are complete

### Test Coverage

The integration tests specifically cover:

1. Authentication with a valid API key
2. Rejecting requests with an invalid API key
3. Rejecting requests with no API key

## Adding New Tests

When adding new integration tests:

1. Place them in the appropriate directory (`integration/`)
2. Follow the Ginkgo BDD-style test patterns
3. Make sure they check real behavior of the system
4. Keep tests independent and idempotent 