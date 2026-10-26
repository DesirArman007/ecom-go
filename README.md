# Ecom API

A Go e-commerce API backed by PostgreSQL. The application is packaged with Docker and can be started with Docker Compose.

## Project structure

```text
cmd/                              Application entrypoint and HTTP server
internal/                         Application and database code
internal/adapters/postgresql/    SQL queries, generated code, and migrations
Dockerfile                        Multi-stage Go container build
docker-compose.yml                Backend and PostgreSQL services
Jenkinsfile                       Jenkins build, deploy, and health-check pipeline
```

## Requirements

- Go 1.25.4 or later
- Docker Desktop or Docker Engine with Docker Compose
- Git

## Configuration

Create a local environment file from the example:

```powershell
Copy-Item .env.example .env
```

Edit `.env` and set a secure value for `POSTGRES_PASSWORD`:

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=change-me
POSTGRES_DB=ecom
```

## Run with Docker Compose

Build and start the backend and PostgreSQL containers:

```powershell
docker compose up --build -d
```

Check the service status:

```powershell
docker compose ps
```

The API is available at:

```text
http://localhost:8000
```

Check the health endpoint:

```powershell
Invoke-WebRequest http://localhost:8000/health
```

The backend connects to PostgreSQL internally using `postgres:5432`. PostgreSQL is intentionally not published to the host, so there is no `5432:5432` port mapping.

View logs with:

```powershell
docker compose logs -f backend
docker compose logs -f postgres
```

Stop the services without deleting the database volume:

```powershell
docker compose down
```

To stop the services and delete the stored PostgreSQL data:

```powershell
docker compose down -v
```

## Run the Go application locally

For local Go development, use a PostgreSQL instance reachable from your computer and set the connection string in the environment:

```powershell
$env:GOOSE_DBSTRING = "user=postgres password=postgres host=localhost port=5432 dbname=ecom sslmode=disable"
go run ./cmd
```

The local server listens on port `8080`:

```text
http://localhost:8080/health
```

The default Docker Compose setup keeps PostgreSQL private to the Docker network. If you need to run the Go server on the host while using the Compose database, use a separate development Compose override that temporarily publishes PostgreSQL to localhost.

## Database migrations

Migration files are located in:

```text
internal/adapters/postgresql/migrations/
```

Apply them with the Goose CLI using the same database connection string configured for your environment. Do not commit database passwords or other secrets.



## API endpoints

- `GET /health`
- `GET /products`
- `GET /products/{id}`
- `POST /products`
- `POST /orders`
- `GET /orders/{id}`
