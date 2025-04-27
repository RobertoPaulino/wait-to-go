# Wait-to-Go

Wait-to-Go is an open-source queue management system that helps businesses manage their waiting lists efficiently. It provides both a real-time queue management solution and a demo version for showcase purposes.

## Project Structure

```
wait-to-go/
├── backend/         # Go backend service
├── frontend/        # Frontend implementations
│   ├── production/ # Production frontend (connects to real backend)
│   └── demo/       # Demo frontend (simulated backend for static hosting)
```

## Components

### Backend

The backend is built with Go and PostgreSQL, providing a RESTful API for queue management. See [backend/README.md](backend/README.md) for detailed setup and API documentation.

### Production Frontend

Located in `frontend/production/`, this is the main frontend implementation that connects to the Go backend. It's built with:
- HTML5
- CSS3
- Vanilla JavaScript
- Real-time backend communication

### Demo Frontend

Located in `frontend/demo/`, this is a standalone demo version that can be hosted on static websites. It:
- Uses the same UI as the production version
- Simulates backend functionality in JavaScript
- Requires no server setup
- Perfect for showcasing the project

## Getting Started

### Full Production Setup

1. Set up the backend:
   ```bash
   cd backend
   go run .
   ```

2. Serve the production frontend:
   ```bash
   cd frontend/production
   # Use any static file server
   python -m http.server 3000
   ```

3. Visit `http://localhost:3000` in your browser

### Demo Version

Simply open `frontend/demo/index.html` in a web browser or host it on any static hosting service.

## Features

- Real-time queue management
- Simple and intuitive interface
- Customer status tracking
- Queue position notifications
- Easy deployment options
- Demo mode for testing and showcase

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the GNU General Public License v2.0 - see the [LICENSE](LICENSE) file for details.
