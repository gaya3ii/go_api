# go_api

A REST API built with Go and PostgreSQL, containerized with Docker.

## Prerequisites

- [Go 1.24+](https://golang.org/dl/)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

---

## Running locally (without Docker)

```bash
# Download dependencies
go mod download

# Run the server
go run ./cmd/app
```

---

## Docker

### Build the image

```bash
docker build -t go-api .
```

### Run the container (standalone)

Requires a running PostgreSQL instance. Pass DB credentials via environment variables:

```bash
docker run -p 8080:8080 \
  -e DB_HOST=<host> \
  -e DB_USER=postgres \
  -e DB_PASSWORD=password \
  -e DB_NAME=goapi \
  -e DB_PORT=5432 \
  go-api
```

---

## Docker Compose

Starts the API and a PostgreSQL database together.

Requires a `.env` file in this directory with `DB_PASSWORD` set — Compose reads
it automatically to fill in `${DB_PASSWORD}` in [docker-compose.yml](docker-compose.yml). No
sourcing/exporting needed, just the file being present.

```bash
# Build and start all services
docker compose up --build

# Start in detached (background) mode
docker compose up --build -d

# Stop all services
docker compose down

# Stop and remove volumes (wipes the database)
docker compose down -v

# View logs
docker compose logs -f

# View logs for a specific service
docker compose logs -f app
docker compose logs -f postgres
```

---

## API Endpoints

| Method | Path           | Description       |
|--------|----------------|-------------------|
| GET    | /users         | List all users    |
| POST   | /users         | Create a user     |
| GET    | /users/get     | Get user by ID    |
| PUT    | /users/update  | Update a user     |
| DELETE | /users/delete  | Delete a user     |

Query parameter `id` (UUID) is required for `/users/get`, `/users/update`, and `/users/delete`.

### Example requests

```bash
# Create a user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","age":30}'

# Get all users
curl http://localhost:8080/users

# Get user by ID
curl "http://localhost:8080/users/get?id=<uuid>"

# Update a user
curl -X PUT "http://localhost:8080/users/update?id=<uuid>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Smith","email":"alice@example.com","age":31}'

# Delete a user
curl -X DELETE "http://localhost:8080/users/delete?id=<uuid>"
```

---

## Environment Variables

| Variable      | Default     | Description              |
|---------------|-------------|---------------------------|
| `DB_HOST`     | `localhost` | PostgreSQL host          |
| `DB_USER`     | — required  | PostgreSQL user          |
| `DB_PASSWORD` | — required  | PostgreSQL password      |
| `DB_NAME`     | — required  | PostgreSQL database name |
| `DB_PORT`     | `5432`      | PostgreSQL port          |

For `DB_USER`, `DB_PASSWORD`, and `DB_NAME` when running locally (outside Docker/Kubernetes), export the values from
`.env` first:

```bash
set -a; source .env; set +a
```
