# IPAM by BananaOps - Project Structure

## Directory Layout

```
.
├── backend/                    # Go backend service
│   ├── cmd/                   # Application entry points
│   │   └── server/           # Main server application
│   │       └── main.go       # Server entry point
│   ├── internal/             # Private application code
│   ├── pkg/                  # Public library code
│   ├── config.example.yaml   # Example configuration file
│   ├── go.mod               # Go module definition
│   └── .gitignore           # Backend-specific ignores
│
├── frontend/                  # React TypeScript frontend
│   ├── src/                  # Source files
│   │   ├── test/            # Test utilities
│   │   │   └── setup.ts    # Test setup
│   │   ├── App.tsx         # Main application component
│   │   ├── main.tsx        # Application entry point
│   │   ├── index.css       # Global styles
│   │   └── vite-env.d.ts   # Vite type definitions
│   ├── index.html           # HTML template
│   ├── package.json         # NPM dependencies
│   ├── tsconfig.json        # TypeScript configuration
│   ├── vite.config.ts       # Vite build configuration
│   ├── .eslintrc.cjs        # ESLint configuration
│   └── .gitignore           # Frontend-specific ignores
│
├── proto/                     # Protocol Buffer definitions
│   └── subnet.proto          # Subnet API definitions
│
├── bin/                       # Compiled binaries (generated)
│   └── ipam-server          # Backend server binary
│
├── Taskfile.yml              # Task runner configuration
├── README.md                 # Project documentation
├── .gitignore               # Root-level ignores
└── PROJECT_STRUCTURE.md     # This file
```

## Key Files

### Backend

- **go.mod**: Defines Go module and dependencies
  - go-ipam: IP address management library
  - protobuf: Protocol Buffers support
  - gorilla/mux: HTTP router
  - mongo-driver: MongoDB support
  - sqlite: SQLite embedded database

- **cmd/server/main.go**: Server entry point (placeholder)

- **config.example.yaml**: Example configuration showing:
  - Server settings (port, host)
  - Database configuration (SQLite/MongoDB)
  - IPAM settings
  - Cloud provider settings

### Frontend

- **package.json**: NPM dependencies including:
  - React 18 + TypeScript
  - Vite (build tool)
  - Axios (HTTP client)
  - FontAwesome (icons)
  - Vitest (testing)
  - fast-check (property-based testing)

- **vite.config.ts**: Vite configuration with:
  - React plugin
  - Path aliases (@/ → src/)
  - Dev server proxy to backend
  - Test configuration

- **tsconfig.json**: TypeScript configuration with strict mode

- **src/App.tsx**: Main application component (placeholder)

- **src/main.tsx**: React application entry point

- **src/index.css**: Global styles with Cyber Minimal color palette

### Protocol Buffers

- **proto/subnet.proto**: Complete API definitions including:
  - Subnet message with all properties
  - LocationType enum (DATACENTER, SITE, CLOUD)
  - CloudInfo for provider metadata
  - SubnetDetails for calculated properties
  - UtilizationInfo for IP usage tracking
  - Request/Response messages for all CRUD operations
  - Error message structure

### Build Tools

- **Taskfile.yml**: Task runner with commands for:
  - Development (dev, dev:backend, dev:frontend)
  - Building (build, build:backend, build:frontend)
  - Testing (test, test:backend, test:frontend)
  - Linting (lint, lint:backend, lint:frontend)
  - Protobuf generation (proto:generate)
  - Docker operations (docker:build, docker:push)
  - Database operations (db:migrate, db:seed)
  - Utilities (clean, install)

## Next Steps

The project structure is now initialized and ready for implementation. The next tasks will:

1. Generate Protobuf code for Go and TypeScript
2. Implement backend service layer with Protobuf
3. Implement REST API Gateway layer
4. Implement database repositories
5. Implement IP calculation service with go-ipam
6. Implement frontend components
7. Implement theme system
8. Wire everything together

## Development Workflow

```bash
# Install dependencies
task install

# Run in development mode
task dev

# Build for production
task build

# Run tests
task test

# Generate Protobuf code
task proto:generate
```

## Dependencies Status

✅ Backend Go module initialized with all required dependencies
✅ Frontend NPM packages installed
✅ Protobuf definitions created
✅ Build tools configured (Taskfile)
✅ Backend compiles successfully
✅ Frontend builds successfully
