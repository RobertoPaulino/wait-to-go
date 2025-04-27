# Wait-to-Go

<div align="center">

[![License](https://img.shields.io/badge/license-GPL%20v2-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.24-blue.svg)](https://golang.org/dl/)
[![PostgreSQL](https://img.shields.io/badge/postgresql-latest-blue.svg)](https://www.postgresql.org/)

A modern, secure queue management system built with Go and JavaScript.

[Getting Started](#getting-started) •
[Features](#features) •
[Documentation](#documentation) •
[Contributing](#contributing)

</div>

## Overview

Wait-to-Go is a robust queue management system designed for businesses that need to manage customer waiting lists efficiently. It features a secure backend built with Go, a responsive frontend interface, and comprehensive authentication systems for both customers and administrators.

### Key Features

- Secure authentication for customers and admins
- Responsive web interface
- Real-time queue updates
- Admin dashboard
- Easy deployment
- Comprehensive test coverage

## Project Structure

```
wait-to-go/
├── backend/         # Go backend service
│   ├── auth/       # Authentication system
│   ├── tests/      # Test suites
│   └── README.md   # Backend documentation
├── frontend/       # Frontend implementations
│   ├── production/ # Production frontend
│   └── demo/       # Demo version (static)
```

## Getting Started

### Prerequisites

- Go 1.24 or later
- PostgreSQL
- Node.js (optional, for development)
- Docker (optional)

### Quick Start

1. **Clone the Repository**
   ```bash
   git clone https://github.com/yourusername/wait-to-go.git
   cd wait-to-go
   ```

2. **Set Up the Backend**
   ```bash
   cd backend
   cp .env.example .env  # Configure your environment variables
   go mod download
   go run .
   ```

3. **Start the Frontend**
   ```bash
   cd frontend/production
   python -m http.server 3000  # Or any static file server
   ```

4. Visit `http://localhost:3000` in your browser

### Demo Version

For a quick demo without backend setup:
```bash
cd frontend/demo
python -m http.server 3000
```

## Features

### Customer Features
- Join queue with name and phone number
- Receive unique access token
- Track position in queue
- Real-time status updates
- Mobile-friendly interface

### Admin Features
- Secure admin dashboard
- View and manage entire queue
- Notify next customer
- Mark customers as served
- Queue analytics
- Multiple admin support

### Security Features
- JWT authentication for customers
- Bcrypt-hashed admin keys
- Rate limiting protection
- Input validation
- Secure headers
- CORS protection

## Documentation

### API Endpoints

#### Public Endpoints
- `POST /join` - Add new entry to queue
- `GET /status/{id}` - Check entry status (requires JWT)

#### Admin Endpoints
- `GET /queue` - View current queue
- `POST /next` - Notify next in line
- `POST /serve` - Mark as served
- `POST /clear` - Clear queue

For detailed API documentation, see [backend/README.md](backend/README.md).

## Development

### Running Tests
```bash
cd backend
go test ./... -v
```

### Environment Variables

Backend configuration (`.env`):
```bash
JWT_SECRET=your-256-bit-secret
ADMIN_API_KEY=your-admin-key
DB_HOST=localhost
DB_PORT=5432
DB_NAME=waitdb
DB_USER=postgres
DB_PASSWORD=secret
```

## Deployment

### Docker Deployment
```bash
docker-compose up -d
```

### Manual Deployment
1. Set up PostgreSQL database
2. Configure environment variables
3. Build and run the Go backend
4. Serve the frontend using Nginx/Apache

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Development Process
1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to your fork
5. Submit a pull request

## License

This project is licensed under the GNU General Public License v2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Go](https://golang.org/) - Backend language
- [PostgreSQL](https://www.postgresql.org/) - Database
- [JWT](https://jwt.io/) - Authentication

---

<div align="center">
An open source project
</div>
