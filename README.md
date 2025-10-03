# Todo Backend API

This is a robust and scalable backend API for a Todo application, built with Go (Golang) and the Fiber web framework. It provides a complete set of features for user management and todo list functionality, including JWT-based authentication, pagination, and filtering.

## Features

- **User Management:**
  - User registration
  - User login and logout
  - Secure password hashing using bcrypt
  - JWT-based authentication
  - User profile management
- **Todo Management:**
  - Create, read, update, and delete (CRUD) operations for todos
  - Mark todos as complete
  - Pagination for listing todos
  - Filtering todos by completion status
- **API:**
  - RESTful API
  - Rate limiting to prevent abuse
  - CORS (Cross-Origin Resource Sharing) support
  - Structured and consistent JSON responses
- **Database:**
  - PostgreSQL database
  - Automatic table creation on startup

## Technologies Used

- **Backend:**
  - [Go (Golang)](https://golang.org/)
  - [Fiber](https://gofiber.io/) - A fast and expressive web framework for Go
  - [PostgreSQL](https://www.postgresql.org/) - A powerful, open-source object-relational database system
  - [Docker](https://www.docker.com/) - For containerizing and running the PostgreSQL database
- **Libraries:**
  - [github.com/gofiber/fiber/v2](https://github.com/gofiber/fiber/v2) - The Fiber web framework
  - [github.com/lib/pq](https://github.com/lib/pq) - The PostgreSQL driver for Go
  - [github.com/golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) - For creating and signing JWTs
  - [github.com/google/uuid](https://github.com/google/uuid) - For generating and working with UUIDs
  - [github.com/joho/godotenv](https://github.com/joho/godotenv) - For loading environment variables from a `.env` file
  - [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - For hashing passwords

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.18 or higher)
- [Docker](https://docs.docker.com/get-docker/)
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

### Installation

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/todo-backend.git
    cd todo-backend
    ```

2.  **Install the dependencies:**

    ```bash
    go mod tidy
    ```

### Configuration

1.  **Create a `.env` file:**

    Create a file named `.env` in the root of the project and add the following environment variables. You can copy the `.env.example` file to get started.

    ```env
    # Server configuration
    PORT=8000
    HOST=localhost

    # Database configuration
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=postgres
    DB_NAME=postgres
    DB_SSLMODE=disable

    # JWT configuration
    JWT_SECRET_KEY=your-secret-key
    JWT_EXPIRY_HOURS=24

    # CORS configuration
    CORS_ORIGINS=http://localhost:3000
    ```

2.  **Start the PostgreSQL database:**

    You can use the provided `docker-compose.yml` file to start a PostgreSQL database in a Docker container.

    ```bash
    docker-compose -f postgres/docker-compose.yml up -d
    ```

### Running the Application

To run the application, use the following command:

```bash
go run main.go
```

The server will start on the port specified in the `.env` file (default is `8000`).

## API Endpoints

All endpoints are prefixed with `/api/v1`.

### Authentication

| Method | Endpoint         | Description              | Request Body                 | Response                       |
| ------ | ---------------- | ------------------------ | ---------------------------- | ------------------------------ |
| `POST` | `/auth/register` | Register a new user      | `registerUserRequest`        | `register_loginUserResponse`   |
| `POST` | `/auth/login`    | Login an existing user   | `loginUserRequest`           | `register_loginUserResponse`   |
| `GET`  | `/auth/logout`   | Logout the current user  | -                            | `200 OK`                       |
| `GET`  | `/auth/profile`  | Get the current user's profile | -                        | `register_loginUserResponse`   |

### Todos

| Method   | Endpoint            | Description                | Request Body                 | Response                  |
| -------- | ------------------- | -------------------------- | ---------------------------- | ------------------------- |
| `POST`   | `/todos/create`     | Create a new todo          | `Create_UpdateTodoRequest`   | `TodoResponse`            |
| `GET`    | `/todos/list`       | Get a list of todos        | -                            | `PaginatedTodoResponse`   |
| `PUT`    | `/todos/update/:id` | Update a todo's title     | `Create_UpdateTodoRequest`   | `TodoResponse`            |
| `PATCH`  | `/todos/complete/:id` | Mark a todo as complete    | `CompleteTodoRequest`        | `TodoResponse`            |
| `DELETE` | `/todos/delete/:id` | Delete a todo              | -                            | `200 OK`                  |

## Project Structure

```
.
├── apps
│   ├── todos
│   │   ├── controller.go
│   │   ├── models.go
│   │   ├── serializers.go
│   │   └── sql.go
│   └── users
│       ├── controllers.go
│       ├── models.go
│       ├── serializers.go
│       └── sql.go
├── backend
│   ├── config
│   │   └── config.go
│   ├── database
│   │   └── db.go
│   ├── middleware
│   │   ├── auth.go
│   │   ├── cors.go
│   │   ├── limiter.go
│   │   ├── logger.go
│   │   ├── recover.go
│   │   └── user.go
│   ├── response
│   │   └── response.go
│   ├── router
│   │   └── router.go
│   └── utils
│       ├── constraints.go
│       ├── encryption.go
│       ├── structure.go
│       ├── timeParser.go
│       └── token.go
├── postgres
│   └── docker-compose.yml
├── test
│   └── test.go
├── .dockerignore
├── .env.example
├── .gitignore
├── Dockerfile
├── go.mod
├── go.sum
└── main.go
```

## Database Schema

### `users`

| Column      | Type        | Description                  |
| ----------- | ----------- | ---------------------------- |
| `id`        | `UUID`      | Primary key                  |
| `name`      | `TEXT`      | The user's name             |
| `email`     | `TEXT`      | The user's email (unique)   |
| `image`     | `TEXT`      | The user's profile image    |
| `password`  | `TEXT`      | The user's hashed password  |
| `jwt`       | `UUID`      | Foreign key to `jwt_tokens`  |
| `created_at`| `TIMESTAMPTZ` | The time the user was created|
| `updated_at`| `TIMESTAMPTZ` | The time the user was last updated |

### `jwt_tokens`

| Column     | Type        | Description                  |
| ---------- | ----------- | ---------------------------- |
| `id`       | `UUID`      | Primary key                  |
| `token`    | `TEXT`      | The JWT                      |
| `expires_at`| `TIMESTAMPTZ` | The time the JWT expires     |
| `created_at`| `TIMESTAMPTZ` | The time the JWT was created |

### `todos`

| Column      | Type        | Description                  |
| ----------- | ----------- | ---------------------------- |
| `id`        | `UUID`      | Primary key                  |
| `title`     | `TEXT`      | The title of the todo        |
| `completed` | `BOOLEAN`   | The completion status of the todo |
| `owner`     | `UUID`      | Foreign key to `users`       |
| `created_at`| `TIMESTAMPTZ` | The time the todo was created|

## Contributing

Contributions are welcome! Please feel free to submit a pull request.
