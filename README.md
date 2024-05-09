# Basic Chat App with Angular and Golang

This is a basic chat application built with Angular for the frontend, Golang for the backend, and WebSocket library to facilitate real-time communication between the frontend and backend.

## Features

- Real-time chat functionality allowing users to send and receive messages instantly.
- Simple and intuitive user interface built with Angular.
- Secure WebSocket-based communication between the frontend and backend.
- Scalable Golang backend capable of handling multiple concurrent connections.

## Technologies Used

- **Angular**: A popular frontend framework for building dynamic web applications.
- **Golang**: A powerful programming language known for its simplicity, concurrency support, and performance.
- **WebSocket Library**: Used to establish a full-duplex communication channel between the client (Angular) and server (Golang) over a single, long-lived connection.

## Getting Started

To get started with this application, follow these steps:

1. Clone this repository to your local machine:

    ```bash
    git clone https://github.com/samehel/Golang-ChatApp
    ```

2. Navigate to the `frontend` directory and install dependencies:

    ```bash
    cd frontend/socket
    npm install
    ```

3. Run the Angular development server:

    ```bash
    ng serve
    ```

4. Open a new terminal window, navigate to the `backend` directory, and run the Golang server:

    ```bash
    cd server
    go run server.go
    ```

5. Once both the frontend and backend servers are running, open your web browser and navigate to `http://localhost:4200` to access the chat application.

## Contributing

This project was made for practice purposes and was made by following a [guide](https://www.thepolyglotdeveloper.com/2016/12/create-real-time-chat-app-golang-angular-2-websockets/) and was updated since the old code contained depreciated libraries and outdated methods.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
