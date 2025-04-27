# Wait-to-Go Demo Frontend

This is the demo version of Wait-to-Go's frontend, designed to work without a backend server. It simulates all backend functionality in JavaScript, making it perfect for static hosting and demonstrations.

## Features

- Identical UI to the production version
- Simulated backend functionality
- No server requirements
- Perfect for GitHub Pages or other static hosting
- Great for testing and demonstrations

## Usage

Simply open `index.html` in a web browser or host the contents of this directory on any static file hosting service.

## Development

The demo version uses:
- HTML5 for structure
- CSS3 for styling
- JavaScript for both UI and simulated backend functionality

### Key Differences from Production

- Uses `MockBackend` class instead of real API calls
- Stores data in browser's localStorage
- Simulates network latency for realistic feel
- Maintains queue state in memory

## Files

- `index.html` - Main entry point
- `styles/` - CSS stylesheets
- `scripts/` - JavaScript files, including backend simulation
- `assets/` - Images and other static assets 