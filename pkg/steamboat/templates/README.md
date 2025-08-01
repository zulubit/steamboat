# <<!.ProjectName!>>

A web application built with Steamboat framework.

## Getting Started

### Prerequisites

- Go 1.24 or higher
- SQLite3

### Installation

```bash
# Install dependencies
go mod download

# Run database migrations
go run cmd/cli/main.go migrate

# Start the development server
go run cmd/cli/main.go serve
```

### Available Commands

```bash
# Generate a new model
go run cmd/cli/main.go make model [name]

# Generate a new migration
go run cmd/cli/main.go make migration [name]

# Run migrations
go run cmd/cli/main.go migrate

# Start the server
go run cmd/cli/main.go serve
```

## Project Structure

```
<<!.ProjectName!>>/
├── cmd/
│   ├── api/        # API server entry point
│   └── cli/        # CLI tool
├── internal/
│   ├── database/   # Database models and migrations
│   ├── handlers/   # HTTP handlers
│   ├── middleware/ # HTTP middleware
│   ├── routes/     # Route definitions
│   ├── server/     # Server configuration
│   └── utils/      # Utilities
└── db/             # SQLite database files
```

## Configuration

Configuration is managed through environment variables in the `.env` file:

- `PORT` - Server port (default: 8080)
- `DB_URL` - Database file path
- `APP_ENV` - Application environment
- `SESSION_KEY` - Secret key for session encryption

## License

MIT