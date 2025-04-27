# Wait-to-Go Production Frontend

This is the production frontend for Wait-to-Go, designed to work with the Go backend service.

## Setup

1. Make sure the backend server is running (see [../../backend/README.md](../../backend/README.md))

2. Serve this directory using any static file server. For example:
   ```bash
   # Using Python's built-in server
   python -m http.server 3000
   
   # Or using Node's http-server
   npx http-server -p 3000
   ```

3. Visit `http://localhost:3000` in your browser

## Development

The frontend is built with vanilla web technologies:
- HTML5 for structure
- CSS3 for styling
- JavaScript for functionality

No build process is required - edit the files directly and refresh your browser to see changes.

## Files

- `index.html` - Main entry point
- `styles/` - CSS stylesheets
- `scripts/` - JavaScript files
- `assets/` - Images and other static assets 