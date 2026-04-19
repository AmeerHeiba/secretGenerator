# JWT Secret Generator

A secure, anonymous Golang web service for generating random JWT secrets.

## Features

- Generates cryptographically secure random 256-bit secrets
- Base64 encoded for compatibility
- Rate limiting (10 requests per minute per IP)
- CORS enabled
- No logging of generated secrets
- HTTPS recommended for production

## API

### POST /generate-secret

Generates a new random JWT secret.

**Request:** No body required

**Response:** JSON
```json
{
  "secret": "base64encodedrandomkey"
}
```

## Running

### Locally

```bash
go run main.go
```

Server starts on http://localhost:8080

### Docker

```bash
docker build -t jwt-secret-generator .
docker run -p 8080:8080 jwt-secret-generator
```

## Security Notes

- Use HTTPS in production
- The service does not store or log generated secrets
- Rate limiting prevents abuse
- Secrets are generated using crypto/rand for maximum security