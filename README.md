# IPAM by BananaOps

A modern IP Address Management (IPAM) application with a full-stack architecture featuring Go backend, React/TypeScript frontend, and Protocol Buffers for communication.

## Project Structure

```
.
├── backend/              # Go backend service
│   ├── cmd/             # Application entry points
│   │   └── server/      # Main server
│   └── go.mod           # Go module definition
├── frontend/            # React TypeScript frontend
│   ├── src/            # Source files
│   ├── package.json    # NPM dependencies
│   └── vite.config.ts  # Vite configuration
├── proto/              # Protocol Buffer definitions
│   └── subnet.proto    # Subnet API definitions
└── Taskfile.yml        # Task runner configuration
```

## Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- npm or yarn
- Protocol Buffers compiler (protoc)
- Task (https://taskfile.dev)

## Quick Start

### Install Dependencies

```bash
# Install all dependencies (backend + frontend)
task install

# Or install separately
task install:backend
task install:frontend
```

### Development

```bash
# Run both backend and frontend in development mode
task dev

# Or run separately
task dev:backend
task dev:frontend
```

The backend will be available at `http://localhost:8080` and the frontend at `http://localhost:3000`.

### Build

```bash
# Build both backend and frontend
task build

# Or build separately
task build:backend
task build:frontend
```

### Testing

```bash
# Run all tests
task test

# Run backend tests with coverage
task test:backend:coverage

# Run frontend tests
task test:frontend
```

### Generate Protobuf Code

```bash
task proto:generate
```

## Technology Stack

### Backend
- **Language**: Go 1.21+
- **IP Management**: go-ipam
- **Database**: MongoDB embedded OR SQLite embedded
- **API**: Protocol Buffers (proto3)
- **HTTP Router**: Gorilla Mux

### Frontend
- **Framework**: React 18+
- **Language**: TypeScript 5+
- **Build Tool**: Vite
- **Icons**: FontAwesome 6+
- **HTTP Client**: Axios
- **Testing**: Vitest + fast-check (property-based testing)

### Communication
- REST API (JSON) over HTTP for frontend-backend
- Protocol Buffers internally for type safety

## Available Tasks

Run `task --list` to see all available tasks:

- `task dev` - Run full stack in development mode
- `task build` - Build everything
- `task test` - Run all tests
- `task lint` - Lint all code
- `task clean` - Clean build artifacts
- `task proto:generate` - Generate Protobuf code

## Configuration

Backend configuration is loaded from environment variables or `config.yaml`:

```yaml
server:
  port: 8080
  host: "0.0.0.0"

database:
  type: "sqlite"  # or "mongodb"
  path: "./data/ipam.db"
```

## License

Copyright © 2024 BananaOps
