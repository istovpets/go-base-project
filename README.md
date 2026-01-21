# Go Base Project

This is a template project for creating web services in Go, based on the principles of Clean Architecture.

## Features

- **Clean Architecture**: The code is divided into layers for better maintainability and testability.
- **REST API**: A ready-to-use RESTful API server.
- **Configuration**: Configuration management using environment files (`.env`).
- **PostgreSQL**: Integration with a PostgreSQL database.
- **Migrations**: Database migration management.
- **Makefile**: Simplified project management using `make`.

## Project Structure

```
/
├───cmd/                    # Application entry point
├───internal/               # Internal application logic
│   ├───app/                # Application
│   ├───config/             # Configuration
│   ├───delivery/           # Delivery layer (REST, gRPC, etc.)
│   ├───domain/             # Domain entities and business logic
│   ├───infrastructure/     # Infrastructure layer (repositories, etc.)
│   └───usecase/            # Use cases layer
├───migrations/             # Database migrations
├───.env.example            # Example environment file
├───go.mod                  # Go modules
└───Makefile                # Makefile for task automation
```

## Getting Started

### Prerequisites

- [Go](https://golang.org/)
- [Docker](https://www.docker.com/) (for the database)
- [make](https://www.gnu.org/software/make/)

### Installation

1.  **Clone the repository:**
    ```sh
    git clone <repository-url>
    cd go-base-project
    ```

2.  **Create a `.env` file:**
    Copy `.env.example` to `.env` and change the values if necessary.
    ```sh
    cp .env.example .env
    ```

3.  **Install dependencies:**
    ```sh
    go mod tidy
    ```

### Running

1.  **Start the database:**
    To run PostgreSQL in Docker, execute the command. This will create a `base-project-db` container with the `dadata_v2` database, which will be available on port `5432`.
    ```sh
    docker run --name base-project-db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=dadata_v2 -p 5432:5432 -d postgres
    ```

2.  **Apply migrations:**
    ```sh
    make migrate_up
    ```

3.  **Run the application:**
    ```sh
    go run ./cmd/main.go
    ```

The application will be available at `http://localhost:8080` (or the port specified in `.env`).

## Makefile Commands

- `make all`: Runs code generation, validation (`go vet`, `lint`), and tests (`go test`).
- `make validate`: Runs `go vet`, `lint`, and `go test`.
- `make lint`: Runs static code analysis using `golangci-lint`.
- `make migrate_up`: Applies database migrations.
- `make migrate_new`: Creates a new migration file named "new".