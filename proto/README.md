# Protocol Buffers Definitions

This directory contains the Protocol Buffer definitions for the IPAM by BananaOps API.

## Structure

- `subnet.proto` - Main API definitions including:
  - Subnet message with all properties
  - API request/response messages (Create, List, Get, Update, Delete)
  - Error message structure
  - Enums for LocationType

## Code Generation

This project uses [Buf](https://buf.build) for Protocol Buffer code generation.

### Prerequisites

- [Buf CLI](https://buf.build/docs/installation) installed

### Generate Code

To generate Go and TypeScript code from the proto files:

```bash
# From project root
buf generate

# Or using Task
task proto:generate
```

### Generated Files

- **Go**: `backend/subnet.pb.go`
- **TypeScript**: `frontend/src/proto/subnet.ts`

### Configuration Files

- `buf.yaml` (in proto directory) - Buf module configuration
- `buf.gen.yaml` (in project root) - Code generation configuration
- `buf.work.yaml` (in project root) - Workspace configuration

## Making Changes

1. Edit the `.proto` files in this directory
2. Run `buf generate` to regenerate code
3. Commit both the `.proto` files and generated code

## Validation

Buf automatically validates proto files against best practices:

```bash
# Lint proto files
buf lint

# Check for breaking changes
buf breaking --against '.git#branch=main'
```
