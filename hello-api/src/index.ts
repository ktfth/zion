// Import the Express library
import express, { Request, Response } from 'express';

// Create a new Express application
const app = express();
const port = 3000;

// Define a route handler for the '/' endpoint
app.get('/', (req: Request, res: Response) => {
  res.send('Hello World!');
});

// Start the server and listen on the specified port
app.listen(port, () => {
  console.log("Server listening on port " + port);
});
