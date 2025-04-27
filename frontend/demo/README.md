# Wait-to-Go Demo Frontend

This is a demo version of the Wait-to-Go queue management system that runs entirely in the browser without requiring a backend server. It uses browser's local storage to simulate a backend, making it perfect for demonstrations and testing.

## Features

- Join the queue with your information
- Check your position in the queue using your queue ID
- Admin interface to manage the queue
- Simulated backend using browser's local storage
- Real-time queue updates in the admin view

## Running the Demo

1. Simply serve this directory using any HTTP server. For example:
   ```bash
   python -m http.server 3000
   ```
   or
   ```bash
   npx serve .
   ```

2. Open your browser and navigate to:
   - If using Python: `http://localhost:3000`
   - If using npx serve: `http://localhost:3000`

## How It Works

- All data is stored in your browser's local storage
- Queue operations are simulated with a small delay to mimic real network requests
- The queue state persists until you clear your browser data or use the "Clear Queue" function
- No actual backend server is required

## Testing the Demo

1. Join Queue:
   - Fill in the form with your name and optional contact details
   - Submit to receive a queue ID

2. Check Status:
   - Use your queue ID to check your position
   - See who is currently being served

3. Admin Functions:
   - View all entries in the queue
   - Call the next person
   - Mark entries as served
   - Clear the entire queue

## Note

This is a demo version for testing and demonstration purposes. For production use, please use the main Wait-to-Go application with a proper backend server. 