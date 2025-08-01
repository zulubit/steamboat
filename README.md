```
      _    _
   __|_|__|_|__
 _|____________|__
|o o o o o o o o /
`~'`~'`~'`~'`~'`~'
```

# Steamboat Framework

A Go web framework for building modern web applications with built-in session management, middleware, and CLI tooling.

v0.0.1

## Features

- 🔥 **Fast Development** - Get started with a single command
- 🔐 **Built-in Sessions** - Encrypted cookie-based sessions
- 🛡️ **Security First** - CORS, Rate limiting, Request ID tracking
- 📦 **CLI Tooling** - Project generation, migrations, and more
- 🧪 **Test Ready** - Comprehensive test suite included
- 🐳 **Docker Ready** - Dockerfile and docker-compose included

## Quick Start

### Install the CLI

```bash
go install github.com/zulubit/steamboat/cmd/steamboat@latest
```

### Create a New Project

```bash
steamboat create myproject
cd myproject
go mod tidy
```

### Run Your Project

```bash
# Run migrations
go run cmd/cli/main.go migrate

# Start the server
go run cmd/cli/main.go serve
```

## CLI Commands

- `steamboat create [name]` - Create a new project
- `steamboat make model [name]` - Generate a model
- `steamboat make migration [name]` - Generate a migration
- `steamboat migrate` - Run migrations
- `steamboat serve` - Start the development server

## Framework Structure

Generated projects include:

- **Session Management** - Encrypted, secure sessions
- **Middleware Stack** - CORS, rate limiting, compression, logging
- **Database Layer** - SQLite with migrations
- **Testing** - Full test coverage
- **Docker Support** - Ready for containerization

## Development

```bash
# Build the CLI
go build -o steamboat cmd/cli/main.go

# Test template generation
./steamboat create testproject
```
