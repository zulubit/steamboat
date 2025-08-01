# Steamboat CLI

The Steamboat CLI is a standalone tool for managing Steamboat applications.

## Building for Distribution

The CLI is completely decoupled from the main application and can be distributed as a standalone binary.

### Using GoReleaser (Recommended)

```bash
# Install GoReleaser
brew install goreleaser/tap/goreleaser

# Create a release
goreleaser release --snapshot --clean

# For a real release (requires git tag)
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
goreleaser release
```

### Manual Build

```bash
# Build for current platform
go build -o steamboat cmd/cli/main.go

# Build with version info
go build -ldflags "-X 'steamboat/pkg/steamboat/cmd.Version=1.0.0'" -o steamboat cmd/cli/main.go
```

## Features

- **Project Creation**: `steamboat create [project-name]`
- **Database Migrations**: `steamboat migrate`
- **Model Generation**: `steamboat make model [name]`
- **Migration Generation**: `steamboat make migration [name]`
- **Development Server**: `steamboat serve`

## Environment Variables

The CLI uses a `.env` file in the current directory for configuration:

```env
PORT=8080
DB_URL=./db/test.db
APP_ENV=development
SESSION_KEY=your-secret-key-here
```

## Distribution

The CLI is self-contained and only requires:
- No external dependencies beyond what's in go.mod
- Reads `.env` from the current working directory
- All templates are embedded in the binary