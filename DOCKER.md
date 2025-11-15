# Docker Setup

This application is fully Dockerized with both the backend API and frontend web app.

## Quick Start

### Using Docker Compose (Recommended)

Start both the API and web frontend:

```bash
docker compose up -d
```

Stop the services:

```bash
docker compose down
```

View logs:

```bash
docker compose logs -f
```

### Access the Application

- **Frontend**: http://localhost
- **API**: http://localhost:9090

### Rebuild After Code Changes

```bash
docker compose up -d --build
```

## Individual Services

### Backend API Only

Build:
```bash
docker build -t words-api .
```

Run:
```bash
docker run -d \
  -p 9090:9090 \
  -v $(pwd)/data:/root/data \
  --name words-api \
  words-api
```

### Frontend Only

Build:
```bash
cd web
docker build -t words-web .
```

Run:
```bash
docker run -d \
  -p 80:80 \
  --name words-web \
  words-web
```

## Architecture

```
┌─────────────────┐
│   User Browser  │
└────────┬────────┘
         │ http://localhost
         ▼
┌─────────────────┐
│  Nginx (Port 80)│
│  React Frontend │
└────────┬────────┘
         │ /api/* proxied to backend
         ▼
┌─────────────────┐
│ Go API (9090)   │
│ + SQLite DB     │
└─────────────────┘
```

## Data Persistence

The SQLite database is stored in a Docker volume mapped to `./data/` on the host machine. This ensures your data persists even if containers are removed.

## Health Checks

Both services include health checks:
- **API**: Checks if the server responds to a simple word lookup
- **Web**: Checks if nginx is serving content

View health status:
```bash
docker compose ps
```

## Environment Variables

You can customize the API by setting environment variables in `docker-compose.yml`:

```yaml
environment:
  - DATABASE_PATH=/root/data/words.db
  - PORT=9090
```

## Development vs Production

### Development
For local development, it's recommended to run services natively:
- Backend: `go run cmd/api/main.go`
- Frontend: `cd web && pnpm run dev`

### Production
Use Docker Compose for production deployments:
```bash
docker compose up -d
```

## Troubleshooting

### Port Already in Use
If ports 80 or 9090 are already in use, modify the port mappings in `docker-compose.yml`:

```yaml
ports:
  - "8080:80"    # Maps host port 8080 to container port 80
  - "9091:9090"  # Maps host port 9091 to container port 9090
```

### Database Issues
To reset the database:
```bash
docker compose down
rm -rf data/
docker compose up -d
```

### View Container Logs
```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f api
docker compose logs -f web
```

### Rebuild from Scratch
```bash
docker compose down
docker compose build --no-cache
docker compose up -d
```

## Resource Requirements

- **Memory**: ~200MB (API) + ~50MB (Web)
- **Disk**: ~100MB (images) + database size
- **CPU**: Minimal (suitable for single-core systems)
