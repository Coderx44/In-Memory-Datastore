# In-Memory-Datastore

In memory key value database
A simple in memory key value datastore that performs operations on it based on certain
commands and uses REST API for communication.
Transport
Read commands via HTTP REST API.
Use JSON encoding for requests and responses.
Use appropriate HTTP status codes for responses.
Query Validation
Ensure that the input command is valid before processing the request.
Storage
Use appropriate data structures.
Must support concurrent operations.

The project structure is divided into different layers, for efficient understanding and debugging of code.

API: 
Post: '/commands'
INPUT:
{
"command" :"actual command"
}

render link : https://in-memory-datastore.onrender.com/commands

I have hosted the api on render, you must use postman (or similar software) to send a post request to the above mentioned url, and in the body (application/json) send the input commands. 
