
# **Full-Stack Image Upload Application**

This project is a full-stack image upload application built using Golang, gRPC, and PostgreSQL. The app allows users to upload images and interact with the system in various ways, including signing up, signing in, and accessing secure endpoints via authentication tokens.


## Features

- User Authentication: Users can sign up, sign in, and manage their sessions.
- Image Upload: Users can upload images to the server and store them securely.
- Token-Based Authentication: Access secure endpoints with JWT authentication tokens.
- gRPC API: The application uses gRPC for communication between services, ensuring fast and efficient interactions.

## Tech Stack

- Golang: The backend is built using Golang for its performance and scalability.
- gRPC: Communication between services is handled via gRPC.
- PostgreSQL: The database used for storing user data and image metadata.
- JWT: JSON Web Tokens are used for authenticating and securing endpoints.
- Docker: The application uses Docker to containerize the PostgreSQL database and ensure a consistent environment.
  
## Installation 

Before running the application, make sure you have the following installed:

- Golang
- Docker
- PostgreSQL (Dockerized in this case)


**1.**
```bash 
  git clone https://github.com/EmreZURNACI/PostUploadApp.git
  cd PostUploadApp
```
    
**2.**

You can use Docker to run PostgreSQL and connect the application to it. 
Use the following command to start a PostgreSQL container:
```bash 
  docker-compose up -d
```

**3.**

Run the following command to install the required dependencies:
```bash 
  go mod tidy
```

**4.**

Run the following command to install the required dependencies:
```bash 
  cd Server
  go run Server
```

**5.**

Run the following command to install the required dependencies:
```bash 
  cd Client
  go run Client
```