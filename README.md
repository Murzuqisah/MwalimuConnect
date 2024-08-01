# MwalimuConnect

MwalimuConnect is a peer-to-peer tutoring platform designed to facilitate real-time communication between tutors and students in Kenya. The platform leverages Go APIs for seamless interaction and blockchain technology for secure payment transactions. This README provides an overview of the project's scope, technologies used, and instructions for setup and usage.

## Table of Contents

1. [Project Overview](#project-overview)
2. [Technology Stack](#technology-stack)
3. [Features](#features)
4. [Setup and Installation](#setup-and-installation)
5. [Configuration](#configuration)
6. [Running the Application](#running-the-application)
7. [API Endpoints](#api-endpoints)
8. [Blockchain Integration](#blockchain-integration)
9. [Contributing](#contributing)
10. [License](#license)

## Project Overview

MwalimuConnect aims to bridge the gap between students and tutors by providing a platform where they can connect, communicate, and transact securely. The platform supports real-time communication via video and chat, allows for scheduling and managing tutoring sessions, and ensures that payments are processed securely through blockchain technology.

## Technology Stack

- **Programming Language**: Go (Golang)
- **Frameworks & Libraries**:
  - [Gin](https://github.com/gin-gonic/gin) - HTTP web framework for Go
  - [Gorm](https://gorm.io/) - ORM library for Go
  - [WebRTC](https://webrtc.org/) - Real-time communication technology
  - [Blockchain](https://ethereum.org/) - Ethereum-based smart contracts for payment transactions
- **Database**:
  - PostgreSQL - Relational database management system
- **Front-end**:
  - HTML, CSS, JavaScript (for web-based UI)
- **Testing**:
  - Go testing framework

## Features

- **User Registration and Authentication**: Secure user registration and login mechanisms.
- **Real-time Communication**: Video and chat functionality for tutoring sessions.
- **Session Scheduling**: Ability to schedule and manage tutoring sessions.
- **Payment Integration**: Secure payments via blockchain technology.
- **User Profiles**: Customizable user profiles for tutors and students.
- **Session History**: Track past tutoring sessions and payments.

## Setup and Installation

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/Murzuqisah/MwalimuConnect.git
   cd MwalimuConnect
   ```

2. **Install Dependencies**:

   Ensure you have Go installed. Run the following command to install required Go modules:

   ```bash
   go mod download
   ```

3. **Set Up the Database**:

   Create a PostgreSQL database and update the database configuration in the `.env` file.

4. **Run Migrations**:

   Apply the database migrations to set up the initial schema:

   ```bash
   go run migrate.go
   ```

5. **Start the Application**:

   Run the application on port 8080:

   ```bash
   go run main.go
   ```



## Running the Application

To start the application, execute:

```bash
go run main.go
```

The application will be available at `http://localhost:8080`.

## API Endpoints

- **POST /api/register**: Register a new user
- **POST /api/login**: Authenticate a user
- **POST /api/schedule**: Schedule a new tutoring session
- **POST /api/complete**: Mark a session as completed and process payment
- **GET /api/profile**: Retrieve user profile details
- **GET /api/sessions**: Retrieve session history

Refer to the [API documentation](https://github.com/Murzuqisah/MwalimuConnect/wiki/API-Documentation) for detailed information on endpoints and request/response formats.

## Blockchain Integration

Payments are processed through Ethereum-based smart contracts. 

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes. For more details, refer to the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

