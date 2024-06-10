# Quiz App Backend

## Description

The Quiz App Backend is built using Go and connects with MongoDB to manage quiz data. It features JWT authentication for secure API access and robust error handling with appropriate status codes. The application allows admins to create and manage quizzes, and students to search for and take quizzes. The backend can be tested using Postman for API requests.

MongoDB text-based indexing is used for efficient searching by `username` and `topic`. Text indexes in MongoDB support searching for words and phrases within string content, enhancing the performance of search operations.

## Installation

### Install Go

1. **Download and install Go** from the official website: [https://golang.org/dl/](https://golang.org/dl/)
2. Follow the instructions for your operating system.

### Install MongoDB

1. **Download and install MongoDB** from the official website: [https://www.mongodb.com/try/download/community](https://www.mongodb.com/try/download/community)
2. Follow the instructions for your operating system.

### Set Up the Project

1. **Clone the repository:**
    ```sh
    git clone git@github.com:zeekhoks/quiz-backend.git
    ```

2. **Navigate to the project directory:**
    ```sh
    cd quiz-backend
    ```

3. **Install dependencies:**
    ```sh
    go mod tidy
    ```

4. **Create an environment file (`.env`):**
    ```
    PORT=your_port
    MONGODB_URI=your_mongodb_uri
    SIGNING_KEY=your_signing_key
    ```

5. **Start the server:**
    ```sh
    go run main.go
    ```

## MongoDB Indexing

To enable efficient text-based searching, create text indexes on the `username` and `topic` fields in your MongoDB collections:

```sh
db.users.createIndex({ username: "text" })
db.quizzes.createIndex({ topic: "text" })
```
Text indexes in MongoDB allow for efficient searching of string content within fields. This can improve query performance significantly for search operations.

## API Testing
Use Postman to test the API endpoints. Ensure you include the JWT token in the Authorization header for secure endpoints.

## Environment Variables
The application requires the following environment variables:

PORT: The port on which the server will run.
MONGODB_URI: The URI for connecting to MongoDB.
SIGNING_KEY: The key used for signing JWT tokens.

## Features
JWT Authentication: Secure access to API endpoints.
Admin Functionality: Create, edit, and delete quizzes.
Student Functionality: Search for quizzes by topic and take quizzes.
Error Handling: Robust error handling with appropriate status codes.
Roadmap
Planned future updates and improvements for the project include:

## Adding more detailed logging.
- Introducing more complex query capabilities for quizzes.
- Implementing user roles and permissions.

## Contact Information

For any questions or suggestions, feel free to reach out:

- **Email:** [khokawala.z@northeastern.edu](mailto:khokawala.z@northeastern.edu)
- **LinkedIn:** [Zainab Khokawala](https://www.linkedin.com/in/zainabkhokawala/)

## License

This project is licensed under the terms of the MIT License.