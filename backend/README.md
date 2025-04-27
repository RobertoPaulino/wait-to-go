 

This is the backend service for the Wait-to-Go queue management system. It provides a RESTful API for managing a waiting queue system.

## Prerequisites

- Go 1.24 or later
- PostgreSQL database
- Docker (optional, for running PostgreSQL in container)

## Environment Variables

The service uses the following environment variables:

- `DB_HOST` (default: "localhost")
- `DB_PORT` (default: "5432")
- `DB_USER` (default: "postgres")
- `DB_PASSWORD` (default: "sicreto")
- `DB_NAME` (default: "gopgtest")
- `DB_SSL_MODE` (default: "disable")

## API Endpoints

- `POST /join` - Add a new entry to the queue
- `GET /queue` - Get all waiting entries
- `POST /next` - Notify the next person in queue
- `POST /serve` - Mark an entry as served
- `GET /status/{id}` - Get status of a specific entry
- `POST /clear` - Clear the queue

## Running the Service

1. Make sure PostgreSQL is running and accessible
2. Set environment variables if needed
3. Run the service:
   ```bash
   go run .
   ```

The service will start on port 8080.

## Example Usage

Add a new entry:
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"firstName":"John","lastName":"Doe","email":"john@example.com","phoneNumber":"1234567890"}' \
  http://localhost:8080/join
```

Get queue status:
```bash
curl http://localhost:8080/queue
```