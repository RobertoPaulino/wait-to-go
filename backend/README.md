# Wait-to-Go Backend Service

This is the backend service for the Wait-to-Go queue management system. It provides a RESTful API for managing a waiting queue system with authentication.

## Prerequisites

- Go 1.24 or later
- PostgreSQL database
- Docker (optional, for running PostgreSQL in container)

## Environment Variables

The service uses the following environment variables:

### Database Configuration
- `DB_HOST` (default: "localhost")
- `DB_PORT` (default: "5432")
- `DB_USER` (default: "postgres")
- `DB_PASSWORD` (default: "sicreto")
- `DB_NAME` (default: "gopgtest")
- `DB_SSL_MODE` (default: "disable")

### Security Configuration
- `JWT_SECRET` (default: "your-256-bit-secret") - Change this in production!
- `ADMIN_API_KEY` (default: none) - Initial admin API key. Additional keys can be added programmatically.

## API Endpoints

### Public Endpoints
- `POST /join` - Add a new entry to the queue
  - Returns a JWT token for authentication

### Protected Customer Endpoints (requires JWT)
- `GET /status/{id}` - Get status of a specific entry
  - Requires Bearer token authentication
  - Only accessible by the entry owner

### Protected Admin Endpoints (requires API Key)
- `GET /queue` - Get all waiting entries
- `POST /next` - Notify the next person in queue
- `POST /serve` - Mark an entry as served
- `POST /clear` - Clear the queue

## Authentication

### Customer Authentication
When a customer joins the queue, they receive a JWT token. This token should be included in subsequent requests to check their status:
```bash
Authorization: Bearer <token>
```

Features:
- JWT tokens expire after 24 hours
- Contains customer ID and phone number
- Rate limited to 30 requests per minute per IP

### Admin Authentication
Admin endpoints require an API key to be included in the request header:
```bash
X-API-Key: <your-admin-key>
```

Features:
- Support for multiple admin API keys
- Keys are securely hashed using bcrypt
- Each key stores creation and last used timestamps
- Rate limited to 100 requests per minute per IP
- Keys can be added and removed programmatically

## Security Features

1. Rate Limiting
   - Customer endpoints: 30 requests per minute per IP
   - Admin endpoints: 100 requests per minute per IP
   - Prevents brute force attacks and DoS attempts

2. JWT Security
   - Tokens expire after 24 hours
   - Signed using HMAC-SHA256
   - Contains minimal required claims
   - Validates signing method

3. Admin API Key Security
   - Keys are hashed using bcrypt
   - Support for multiple active keys
   - Key rotation capability
   - Usage tracking
   - No plaintext storage

## Running the Service

1. Copy `.env.example` to `.env` and configure your environment variables:
   ```bash
   cp .env.example .env
   ```

2. Make sure PostgreSQL is running and accessible

3. Run the service:
   ```bash
   go run .
   ```

The service will start on port 8080 by default.

## Example Usage

Join the queue:
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"firstName":"John","lastName":"Doe","email":"john@example.com","phoneNumber":"1234567890"}' \
  http://localhost:8080/join
```

Check status (with token):
```bash
curl -H "Authorization: Bearer <your-token>" \
  http://localhost:8080/status/<id>
```

Get queue (admin):
```bash
curl -H "X-API-Key: <your-admin-key>" \
  http://localhost:8080/queue
```

## Security Considerations

1. Always change the default JWT secret in production
2. Use HTTPS in production
3. Implement proper key rotation procedures
4. Monitor rate limit violations
5. Store sensitive environment variables securely
6. Regularly audit admin key usage
7. Consider implementing IP allowlisting for admin endpoints
8. Use secure headers (HSTS, CSP, etc.) in production