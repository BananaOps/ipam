# Protobuf Setup Documentation

## Overview

This document describes the Protobuf setup for IPAM by BananaOps using Buf.

## What Was Implemented

### 1. Protobuf Definitions (proto/subnet.proto)

Complete API contract definitions including:

- **Subnet Message**: Main data structure with all subnet properties
  - Basic info: id, cidr, name, description, location
  - Location type enum (DATACENTER, SITE, CLOUD)
  - Cloud provider information
  - Calculated subnet details
  - Utilization tracking

- **API Request/Response Messages**:
  - CreateSubnetRequest/Response
  - ListSubnetsRequest/Response (with filtering support)
  - GetSubnetRequest/Response
  - UpdateSubnetRequest/Response
  - DeleteSubnetRequest/Response

- **Error Message Structure**: Structured error responses with code, message, details, and timestamp

### 2. Buf Configuration

Three configuration files were created:

- **buf.yaml** (in proto/): Module configuration with linting and breaking change detection
- **buf.gen.yaml**: Code generation configuration for Go and TypeScript
- **buf.work.yaml**: Workspace configuration

### 3. Code Generation

Successfully generated:

- **Go code**: `backend/subnet.pb.go` - Complete Go structs and methods
- **TypeScript code**: `frontend/src/proto/subnet.ts` - TypeScript interfaces and runtime support

### 4. Integration

- Updated `Taskfile.yml` to use `buf generate` command
- Added `@protobuf-ts/runtime` dependency to frontend package.json
- Created `.eslintignore` to exclude generated proto files from linting
- Created test file `backend/proto_test.go` to verify Go code generation
- Created example file `frontend/src/proto/example.ts` to demonstrate TypeScript usage

## Usage

### Generate Protobuf Code

```bash
# Using buf directly
buf generate

# Using Task
task proto:generate
```

### Validate Proto Files

```bash
# Lint proto files
buf lint

# Check for breaking changes
buf breaking --against '.git#branch=main'
```

### Using Generated Code

#### Go Example

```go
import proto "github.com/bananaops/ipam-bananaops"

subnet := &proto.Subnet{
    Id:          "subnet-123",
    Cidr:        "10.0.0.0/24",
    Name:        "Production Subnet",
    LocationType: proto.LocationType_DATACENTER,
}
```

#### TypeScript Example

```typescript
import { Subnet, LocationType } from './proto/subnet';

const subnet: Subnet = {
    id: 'subnet-123',
    cidr: '10.0.0.0/24',
    name: 'Production Subnet',
    locationType: LocationType.DATACENTER,
    // ... other fields
};
```

## Testing

### Backend Tests

```bash
cd backend
go test -v
```

Tests verify:
- Protobuf message creation
- Field assignments
- Request/response message structures

### Frontend Type Checking

```bash
cd frontend
npx tsc --noEmit
```

## Requirements Validated

This implementation satisfies:

- **Requirement 10.1**: Protobuf message definitions for all endpoints
- **Requirement 10.4**: API schema regeneration from Protobuf definitions

## Next Steps

The Protobuf API contracts are now ready for:

1. Backend service layer implementation (Task 3)
2. REST API Gateway layer implementation (Task 4)
3. Frontend API client implementation (Task 11)
