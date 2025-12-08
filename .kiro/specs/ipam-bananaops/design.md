# Design Document - IPAM by BananaOps

## Overview

IPAM by BananaOps est une application full-stack de gestion d'adresses IP avec une architecture moderne séparant clairement le backend (Go), la couche de données (MongoDB/SQLite embarqué) et le frontend (React/TypeScript). La communication entre frontend et backend utilise Protobuf pour garantir performance et typage fort. Le système est conçu pour être extensible, permettant l'intégration future de récupération dynamique d'IPs depuis les cloud providers.

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Frontend (React/TS)                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Subnet List  │  │ Subnet Detail│  │ Theme System │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
                      REST API (JSON)
                            │
┌─────────────────────────────────────────────────────────────┐
│                    Backend Service (Go)                      │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              REST API Gateway Layer                   │   │
│  │  (JSON ↔ Protobuf conversion, HTTP routing)          │   │
│  └──────────────────────────────────────────────────────┘   │
│                            │                                 │
│                   Internal Protobuf                          │
│                            │                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Service Layer│  │ IP Calculator│  │ Cloud Module │      │
│  │  (Protobuf)  │  │  (go-ipam)   │  │  Interface   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
                            │
┌─────────────────────────────────────────────────────────────┐
│              Embedded Database Layer                         │
│         MongoDB Embedded  OR  SQLite Embedded                │
└─────────────────────────────────────────────────────────────┘
```

### Technology Stack

- **Backend**: Go 1.24+, go-ipam module, Protobuf (protoc-gen-go), grpc-gateway or custom REST proxy
- **Database**: MongoDB embedded (via go.mongodb.org/mongo-driver) OR SQLite embedded (via modernc.org/sqlite)
- **Frontend**: React 18+, TypeScript 5+, FontAwesome 6+, Axios for HTTP requests
- **Communication**: REST API (JSON) over HTTP for frontend-backend, Protocol Buffers (proto3) internally
- **Build Tools**: Go modules, npm/yarn, protoc compiler

## Components and Interfaces

### Backend Components

#### 1. REST API Gateway Layer
Responsable de la gestion des requêtes HTTP REST, conversion JSON ↔ Protobuf, et routage.

```go
type RESTGateway struct {
    serviceLayer *ServiceLayer
}

// REST Endpoints (JSON)
// POST /api/v1/subnets - Create subnet
// GET /api/v1/subnets - List subnets with filters
// GET /api/v1/subnets/{id} - Get subnet details
// PUT /api/v1/subnets/{id} - Update subnet
// DELETE /api/v1/subnets/{id} - Delete subnet

// Conversion functions
func (g *RESTGateway) jsonToProtobuf(jsonData []byte, msg proto.Message) error
func (g *RESTGateway) protobufToJSON(msg proto.Message) ([]byte, error)
```

#### 2. Service Layer (Internal Protobuf)
Couche de logique métier utilisant Protobuf en interne pour le typage fort.

```go
type ServiceLayer struct {
    ipService    *IPService
    subnetRepo   SubnetRepository
    cloudManager *CloudProviderManager
}

// Internal methods using Protobuf messages
func (s *ServiceLayer) CreateSubnet(ctx context.Context, req *pb.CreateSubnetRequest) (*pb.CreateSubnetResponse, error)
func (s *ServiceLayer) ListSubnets(ctx context.Context, req *pb.ListSubnetsRequest) (*pb.ListSubnetsResponse, error)
func (s *ServiceLayer) GetSubnet(ctx context.Context, req *pb.GetSubnetRequest) (*pb.GetSubnetResponse, error)
func (s *ServiceLayer) UpdateSubnet(ctx context.Context, req *pb.UpdateSubnetRequest) (*pb.UpdateSubnetResponse, error)
func (s *ServiceLayer) DeleteSubnet(ctx context.Context, req *pb.DeleteSubnetRequest) (*pb.DeleteSubnetResponse, error)
```

#### 3. IP Service Layer
Utilise go-ipam pour tous les calculs et validations IP.

```go
type IPService struct {
    ipam ipam.Ipamer
}

func (s *IPService) CalculateSubnetDetails(cidr string) (*SubnetDetails, error)
func (s *IPService) ValidateIPAddress(ip string) error
func (s *IPService) CalculateUtilization(subnet *Subnet) (float64, error)
```

#### 4. Repository Layer
Abstraction pour l'accès aux données, supportant MongoDB et SQLite.

```go
type SubnetRepository interface {
    Create(ctx context.Context, subnet *Subnet) error
    FindByID(ctx context.Context, id string) (*Subnet, error)
    FindAll(ctx context.Context, filters SubnetFilters) ([]*Subnet, error)
    Update(ctx context.Context, subnet *Subnet) error
    Delete(ctx context.Context, id string) error
}

type MongoSubnetRepository struct { /* ... */ }
type SQLiteSubnetRepository struct { /* ... */ }
```

#### 5. Cloud Provider Module
Interface extensible pour l'intégration future des cloud providers.

```go
type CloudProvider interface {
    GetName() string
    FetchSubnets(ctx context.Context, credentials CloudCredentials) ([]*CloudSubnet, error)
    GetRegions() []string
}

type CloudProviderManager struct {
    providers map[string]CloudProvider
}

// Future implementations: AWSProvider, AzureProvider, GCPProvider, ScalewayProvider, OVHProvider
```

### Frontend Components

#### 1. API Client Service
Service HTTP pour communiquer avec le backend REST API.

```typescript
class APIClient {
  private baseURL: string;
  private axiosInstance: AxiosInstance;

  // Subnet operations
  async createSubnet(data: CreateSubnetRequest): Promise<Subnet>;
  async listSubnets(filters: SubnetFilters): Promise<SubnetListResponse>;
  async getSubnet(id: string): Promise<Subnet>;
  async updateSubnet(id: string, data: UpdateSubnetRequest): Promise<Subnet>;
  async deleteSubnet(id: string): Promise<void>;
  
  // Error handling
  private handleError(error: AxiosError): APIError;
}

interface CreateSubnetRequest {
  cidr: string;
  name: string;
  description?: string;
  location: string;
  locationType: LocationType;
  cloudInfo?: CloudInfo;
}

interface SubnetListResponse {
  subnets: Subnet[];
  totalCount: number;
}
```

#### 2. Theme System
Gestion des thèmes avec React Context et synchronisation système.

```typescript
type Theme = 'dark' | 'light' | 'auto';

interface ThemeContextValue {
  theme: Theme;
  effectiveTheme: 'dark' | 'light';
  setTheme: (theme: Theme) => void;
}

// Hook personnalisé pour détecter le thème système
function useSystemTheme(): 'dark' | 'light';
```

#### 3. Subnet List Component
Affichage de la liste des sous-réseaux avec filtres.

```typescript
interface SubnetListProps {
  filters: SubnetFilters;
  onFilterChange: (filters: SubnetFilters) => void;
}

interface SubnetFilters {
  location?: string;
  cloudProvider?: CloudProviderType;
  searchQuery?: string;
}
```

#### 4. Subnet Detail Component
Affichage détaillé d'un sous-réseau avec calculs et utilisation.

```typescript
interface SubnetDetailProps {
  subnetId: string;
}

interface SubnetInfo {
  address: string;
  netmask: string;
  wildcard: string;
  network: string;
  type: string;
  broadcast: string;
  hostMin: string;
  hostMax: string;
  hostsPerNet: number;
  isPublic: boolean;
  utilizationPercent: number;
}
```

#### 5. Cloud Provider Icon Component
Rendu des icônes FontAwesome pour chaque provider.

```typescript
const CLOUD_PROVIDER_ICONS: Record<CloudProviderType, IconDefinition> = {
  aws: faAws,
  azure: faMicrosoft,
  gcp: faGoogle,
  scaleway: faCloud,
  ovh: faServer,
};
```

## Data Models

### Subnet Model

```protobuf
message Subnet {
  string id = 1;
  string cidr = 2;
  string name = 3;
  string description = 4;
  string location = 5;
  LocationType location_type = 6;
  CloudInfo cloud_info = 7;
  SubnetDetails details = 8;
  UtilizationInfo utilization = 9;
  int64 created_at = 10;
  int64 updated_at = 11;
}

enum LocationType {
  DATACENTER = 0;
  SITE = 1;
  CLOUD = 2;
}

message CloudInfo {
  string provider = 1;
  string region = 2;
  string account_id = 3;
}

message SubnetDetails {
  string address = 1;
  string netmask = 2;
  string wildcard = 3;
  string network = 4;
  string type = 5;
  string broadcast = 6;
  string host_min = 7;
  string host_max = 8;
  int32 hosts_per_net = 9;
  bool is_public = 10;
}

message UtilizationInfo {
  int32 total_ips = 1;
  int32 allocated_ips = 2;
  float utilization_percent = 3;
}
```

### API Messages

```protobuf
message CreateSubnetRequest {
  string cidr = 1;
  string name = 2;
  string description = 3;
  string location = 4;
  LocationType location_type = 5;
  CloudInfo cloud_info = 6;
}

message CreateSubnetResponse {
  Subnet subnet = 1;
  Error error = 2;
}

message ListSubnetsRequest {
  string location_filter = 1;
  string cloud_provider_filter = 2;
  string search_query = 3;
  int32 page = 4;
  int32 page_size = 5;
}

message ListSubnetsResponse {
  repeated Subnet subnets = 1;
  int32 total_count = 2;
  Error error = 3;
}

message GetSubnetRequest {
  string id = 1;
}

message GetSubnetResponse {
  Subnet subnet = 1;
  Error error = 2;
}

message Error {
  string code = 1;
  string message = 2;
}
```

### Database Schema

#### MongoDB Collections

```javascript
// subnets collection
{
  _id: ObjectId,
  id: String (UUID),
  cidr: String,
  name: String,
  description: String,
  location: String,
  locationType: String,
  cloudInfo: {
    provider: String,
    region: String,
    accountId: String
  },
  details: {
    address: String,
    netmask: String,
    wildcard: String,
    network: String,
    type: String,
    broadcast: String,
    hostMin: String,
    hostMax: String,
    hostsPerNet: Number,
    isPublic: Boolean
  },
  utilization: {
    totalIps: Number,
    allocatedIps: Number,
    utilizationPercent: Number
  },
  createdAt: Date,
  updatedAt: Date
}

// Indexes
db.subnets.createIndex({ location: 1 })
db.subnets.createIndex({ "cloudInfo.provider": 1 })
db.subnets.createIndex({ cidr: 1 }, { unique: true })
```

#### SQLite Schema

```sql
CREATE TABLE subnets (
  id TEXT PRIMARY KEY,
  cidr TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  location TEXT,
  location_type TEXT,
  cloud_provider TEXT,
  cloud_region TEXT,
  cloud_account_id TEXT,
  address TEXT,
  netmask TEXT,
  wildcard TEXT,
  network TEXT,
  type TEXT,
  broadcast TEXT,
  host_min TEXT,
  host_max TEXT,
  hosts_per_net INTEGER,
  is_public BOOLEAN,
  total_ips INTEGER,
  allocated_ips INTEGER,
  utilization_percent REAL,
  created_at INTEGER,
  updated_at INTEGER
);

CREATE INDEX idx_subnets_location ON subnets(location);
CREATE INDEX idx_subnets_cloud_provider ON subnets(cloud_provider);
CREATE INDEX idx_subnets_cidr ON subnets(cidr);
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*


### Property Reflection

After analyzing all acceptance criteria, several properties can be consolidated to avoid redundancy:

- Properties 4.2, 10.2, and 10.3 all relate to Protobuf communication and can be combined into a single comprehensive property
- Properties 1.3 and 6.3 both relate to cloud provider metadata display and can be combined
- Properties 2.1 and 2.4 both relate to subnet detail display and can be combined

### Core Properties

**Property 1: Location filter correctness**
*For any* collection of subnets and any location filter value, applying the filter should return only subnets where the location matches the filter criteria (datacenter, site, or cloud provider).
**Validates: Requirements 1.2**

**Property 2: Cloud provider metadata completeness**
*For any* cloud provider subnet, the rendered display should include the provider-specific icon, region, and account ID.
**Validates: Requirements 1.3, 1.4, 6.3**

**Property 3: Subnet detail completeness**
*For any* subnet, the detail view should display all required fields: address, netmask, wildcard, network, type, broadcast, hostMin, hostMax, hostsPerNet, public/private classification, utilization percentage, and progress bar.
**Validates: Requirements 2.1, 2.2, 2.4**

**Property 4: Invalid subnet error handling**
*For any* invalid subnet data (malformed CIDR, invalid IP), the system should return a structured error message indicating the validation failure.
**Validates: Requirements 2.5, 8.5**

**Property 5: Theme persistence round-trip**
*For any* theme selection (dark, light, auto), setting the theme and reloading the application should restore the same theme preference.
**Validates: Requirements 3.5**

**Property 6: System theme synchronization**
*For any* system theme change when in auto mode, the application theme should update to match the system preference within a reasonable time.
**Validates: Requirements 3.4**

**Property 7: REST to Protobuf conversion consistency**
*For any* API request or response, converting from JSON to Protobuf and back to JSON should preserve all data fields and maintain type correctness.
**Validates: Requirements 4.2, 10.2, 10.3**

**Property 8: Database persistence round-trip**
*For any* subnet with all its properties and metadata, storing it to the database and then retrieving it should return an equivalent subnet object.
**Validates: Requirements 5.3**

**Property 9: Cloud provider error resilience**
*For any* cloud provider failure or unavailability, the system should handle the error gracefully without crashing and provide meaningful error information.
**Validates: Requirements 6.5**

**Property 10: CIDR validation correctness**
*For any* IP address and CIDR notation input, the validation should correctly accept valid formats and reject invalid formats according to IP addressing standards.
**Validates: Requirements 8.1**

**Property 11: Subnet property calculation completeness**
*For any* valid CIDR input, creating a subnet should automatically calculate and store all properties: address, netmask, wildcard, network, type, broadcast, hostMin, hostMax, hostsPerNet, and public/private classification.
**Validates: Requirements 8.2**

**Property 12: Subnet update recalculation**
*For any* subnet update that changes the CIDR, all dependent properties should be recalculated and the updated values should be persisted correctly.
**Validates: Requirements 8.3**

**Property 13: Subnet deletion completeness**
*For any* existing subnet, after deletion, attempting to retrieve that subnet should return a not-found error.
**Validates: Requirements 8.4**

**Property 14: Utilization calculation accuracy**
*For any* subnet with a known number of total IPs and allocated IPs, the calculated utilization percentage should equal (allocated / total) * 100.
**Validates: Requirements 9.1**

**Property 15: High utilization visual indication**
*For any* subnet with utilization above a threshold (e.g., 80%), the display should include visual indicators (color, icon, or styling) to highlight the high utilization status.
**Validates: Requirements 9.3**

**Property 16: Structured error responses**
*For any* communication error or API failure, the error response should be a structured Protobuf message containing an error code and descriptive message.
**Validates: Requirements 10.5**

## Error Handling

### Backend Error Handling

1. **Input Validation Errors**
   - Invalid CIDR notation: Return error code `INVALID_CIDR` with details
   - Invalid IP address: Return error code `INVALID_IP` with details
   - Missing required fields: Return error code `MISSING_FIELD` with field name

2. **Database Errors**
   - Connection failures: Retry with exponential backoff, return `DB_CONNECTION_ERROR`
   - Duplicate subnet: Return error code `DUPLICATE_SUBNET`
   - Not found: Return error code `SUBNET_NOT_FOUND`

3. **Cloud Provider Errors**
   - Provider unavailable: Log error, return cached data if available, return `PROVIDER_UNAVAILABLE`
   - Authentication failure: Return error code `PROVIDER_AUTH_FAILED`
   - Rate limiting: Implement backoff, return `PROVIDER_RATE_LIMITED`

4. **Protobuf Serialization Errors**
   - Malformed message: Return error code `INVALID_MESSAGE_FORMAT`
   - Unknown message type: Return error code `UNKNOWN_MESSAGE_TYPE`

### Frontend Error Handling

1. **API Communication Errors**
   - Network timeout: Display user-friendly message, offer retry
   - Server error (5xx): Display error message, log details
   - Client error (4xx): Display validation feedback

2. **Theme System Errors**
   - Invalid theme preference: Fallback to system default
   - Storage unavailable: Use in-memory theme, warn user

3. **Rendering Errors**
   - Missing subnet data: Display placeholder or empty state
   - Invalid utilization value: Display 0% with warning icon

### Error Response Format

All errors follow a consistent Protobuf structure:

```protobuf
message Error {
  string code = 1;           // Machine-readable error code
  string message = 2;        // Human-readable error message
  map<string, string> details = 3;  // Additional context
  int64 timestamp = 4;       // Error occurrence time
}
```

## Testing Strategy

### Unit Testing

**Backend Unit Tests (Go)**
- Test IP calculation functions with go-ipam for various CIDR notations
- Test repository CRUD operations with mock database
- Test Protobuf serialization/deserialization
- Test cloud provider interface implementations
- Test error handling for invalid inputs

**Frontend Unit Tests (TypeScript/Jest)**
- Test theme context and hook behavior
- Test subnet list filtering logic
- Test utilization percentage calculations
- Test icon mapping for cloud providers
- Test error boundary components

### Property-Based Testing

**Property-Based Testing Library**: For Go, we will use `gopter` (https://github.com/leanovate/gopter). For TypeScript/React, we will use `fast-check` (https://github.com/dubzzz/fast-check).

**Configuration**: Each property-based test should run a minimum of 100 iterations to ensure adequate coverage of the input space.

**Test Tagging**: Each property-based test must include a comment tag in this exact format:
```
// Feature: ipam-bananaops, Property {number}: {property_text}
```

**Backend Property Tests (Go/gopter)**
- Property 1: Location filter correctness - Generate random subnets and filter values, verify filtering
- Property 4: Invalid subnet error handling - Generate invalid CIDR/IP inputs, verify errors
- Property 7: Protobuf communication consistency - Generate random subnet data, verify round-trip
- Property 8: Database persistence round-trip - Generate random subnets, verify store/retrieve
- Property 10: CIDR validation correctness - Generate valid/invalid CIDR strings, verify validation
- Property 11: Subnet property calculation completeness - Generate valid CIDRs, verify all properties calculated
- Property 12: Subnet update recalculation - Generate subnet updates, verify recalculation
- Property 13: Subnet deletion completeness - Generate and delete subnets, verify removal
- Property 14: Utilization calculation accuracy - Generate subnets with various allocations, verify percentage
- Property 16: Structured error responses - Generate error conditions, verify response structure

**Frontend Property Tests (TypeScript/fast-check)**
- Property 2: Cloud provider metadata completeness - Generate random cloud subnets, verify display
- Property 3: Subnet detail completeness - Generate random subnets, verify all fields displayed
- Property 5: Theme persistence round-trip - Generate theme selections, verify persistence
- Property 6: System theme synchronization - Simulate system theme changes, verify updates
- Property 9: Cloud provider error resilience - Generate provider failures, verify graceful handling
- Property 15: High utilization visual indication - Generate high utilization values, verify indicators

**Integration Property Tests**
- End-to-end subnet creation flow with random valid inputs
- End-to-end filtering with random filter combinations
- Theme switching across multiple sessions

### Integration Testing

- Test complete subnet CRUD flow through API
- Test database initialization for both MongoDB and SQLite
- Test theme persistence across browser sessions
- Test Protobuf communication between frontend and backend
- Test error propagation from backend to frontend

### Test Data Generators

**Go Generators (gopter)**
```go
// Generate valid CIDR notations
func GenValidCIDR() gopter.Gen

// Generate invalid CIDR notations
func GenInvalidCIDR() gopter.Gen

// Generate subnet with all properties
func GenSubnet() gopter.Gen

// Generate cloud provider types
func GenCloudProvider() gopter.Gen
```

**TypeScript Generators (fast-check)**
```typescript
// Generate valid subnet objects
const subnetArbitrary: fc.Arbitrary<Subnet>

// Generate theme values
const themeArbitrary: fc.Arbitrary<Theme>

// Generate utilization percentages
const utilizationArbitrary: fc.Arbitrary<number>

// Generate cloud provider types
const cloudProviderArbitrary: fc.Arbitrary<CloudProviderType>
```

## Deployment and Configuration

### Backend Configuration

Configuration via environment variables or config file:

```yaml
server:
  port: 8080
  host: "0.0.0.0"

database:
  type: "sqlite"  # or "mongodb"
  path: "./data/ipam.db"  # for SQLite
  # connection_string: "mongodb://localhost:27017"  # for MongoDB

ipam:
  default_allocation_size: 256

cloud_providers:
  enabled: false  # Enable when cloud modules are implemented
  sync_interval: "5m"
```

### Frontend Configuration

Build-time configuration:

```typescript
const config = {
  apiBaseUrl: process.env.REACT_APP_API_URL || 'http://localhost:8080',
  defaultTheme: 'auto',
  colors: {
    darkPrimary: '#0A1A2F',
    cyanAccent: '#0EA5E9',
    lightGray: '#F3F4F6',
    white: '#FFFFFF',
  },
};
```

### Development Workflow with Taskfile

Le projet utilise [Task](https://taskfile.dev) comme task runner pour simplifier les commandes de développement.

**Taskfile.yml Structure:**

```yaml
version: '3'

vars:
  BACKEND_DIR: backend
  FRONTEND_DIR: frontend
  PROTO_DIR: proto
  DOCKER_REGISTRY: ghcr.io/bananaops
  IMAGE_NAME: ipam-bananaops

tasks:
  # Development tasks
  dev:
    desc: Run both backend and frontend in development mode
    deps: [dev:backend, dev:frontend]

  dev:backend:
    desc: Run backend in development mode with hot reload
    dir: "{{.BACKEND_DIR}}"
    cmds:
      - air # Using cosmtrek/air for Go hot reload

  dev:frontend:
    desc: Run frontend in development mode
    dir: "{{.FRONTEND_DIR}}"
    cmds:
      - npm run dev

  # Build tasks
  build:
    desc: Build both backend and frontend
    deps: [build:backend, build:frontend]

  build:backend:
    desc: Build backend binary
    dir: "{{.BACKEND_DIR}}"
    cmds:
      - task: proto:generate
      - go build -o ../bin/ipam-server cmd/server/main.go

  build:frontend:
    desc: Build frontend for production
    dir: "{{.FRONTEND_DIR}}"
    cmds:
      - npm run build

  # Protobuf generation
  proto:generate:
    desc: Generate Protobuf code for Go and TypeScript
    cmds:
      - protoc --go_out=. --go_opt=paths=source_relative {{.PROTO_DIR}}/*.proto
      - protoc --ts_out={{.FRONTEND_DIR}}/src/proto {{.PROTO_DIR}}/*.proto

  # Test tasks
  test:
    desc: Run all tests
    deps: [test:backend, test:frontend]

  test:backend:
    desc: Run backend tests
    dir: "{{.BACKEND_DIR}}"
    cmds:
      - go test -v ./...

  test:backend:coverage:
    desc: Run backend tests with coverage
    dir: "{{.BACKEND_DIR}}"
    cmds:
      - go test -v -coverprofile=coverage.out ./...
      - go tool cover -html=coverage.out -o coverage.html

  test:frontend:
    desc: Run frontend tests
    dir: "{{.FRONTEND_DIR}}"
    cmds:
      - npm run test

  # Lint tasks
  lint:
    desc: Run linters for backend and frontend
    deps: [lint:backend, lint:frontend]

  lint:backend:
    desc: Run Go linter
    dir: "{{.BACKEND_DIR}}"
    cmds:
      - golangci-lint run

  lint:frontend:
    desc: Run ESLint
    dir: "{{.FRONTEND_DIR}}"
    cmds:
      - npm run lint

  # Docker tasks
  docker:build:
    desc: Build Docker images
    cmds:
      - docker build -t {{.DOCKER_REGISTRY}}/{{.IMAGE_NAME}}-backend:latest -f docker/Dockerfile.backend .
      - docker build -t {{.DOCKER_REGISTRY}}/{{.IMAGE_NAME}}-frontend:latest -f docker/Dockerfile.frontend .

  docker:push:
    desc: Push Docker images to registry
    cmds:
      - docker push {{.DOCKER_REGISTRY}}/{{.IMAGE_NAME}}-backend:latest
      - docker push {{.DOCKER_REGISTRY}}/{{.IMAGE_NAME}}-frontend:latest

  # Database tasks
  db:migrate:
    desc: Run database migrations
    dir: "{{.BACKEND_DIR}}"
    cmds:
      - go run cmd/migrate/main.go

  db:seed:
    desc: Seed database with sample data
    dir: "{{.BACKEND_DIR}}"
    cmds:
      - go run cmd/seed/main.go

  # Clean tasks
  clean:
    desc: Clean build artifacts
    cmds:
      - rm -rf bin/
      - rm -rf {{.FRONTEND_DIR}}/dist/
      - rm -rf {{.BACKEND_DIR}}/coverage.*

  # Install dependencies
  install:
    desc: Install all dependencies
    deps: [install:backend, install:frontend]

  install:backend:
    desc: Install backend dependencies
    dir: "{{.BACKEND_DIR}}"
    cmds:
      - go mod download

  install:frontend:
    desc: Install frontend dependencies
    dir: "{{.FRONTEND_DIR}}"
    cmds:
      - npm install
```

**Common Task Commands:**

```bash
# Development
task dev                    # Run full stack in dev mode
task dev:backend           # Run only backend
task dev:frontend          # Run only frontend

# Building
task build                 # Build everything
task build:backend         # Build backend binary
task build:frontend        # Build frontend assets

# Testing
task test                  # Run all tests
task test:backend          # Run backend tests
task test:backend:coverage # Run backend tests with coverage
task test:frontend         # Run frontend tests

# Linting
task lint                  # Lint everything
task lint:backend          # Lint Go code
task lint:frontend         # Lint TypeScript/React code

# Docker
task docker:build          # Build Docker images
task docker:push           # Push images to registry

# Protobuf
task proto:generate        # Generate Protobuf code

# Database
task db:migrate            # Run migrations
task db:seed               # Seed sample data

# Utilities
task clean                 # Clean build artifacts
task install               # Install all dependencies
```

### Kubernetes Deployment with Skaffold

Le projet utilise [Skaffold](https://skaffold.dev) pour le développement et le déploiement sur Kubernetes.

**skaffold.yaml Configuration:**

```yaml
apiVersion: skaffold/v4beta6
kind: Config
metadata:
  name: ipam-bananaops

build:
  artifacts:
    - image: ipam-backend
      context: .
      docker:
        dockerfile: docker/Dockerfile.backend
      sync:
        manual:
          - src: "backend/**/*.go"
            dest: /app
    
    - image: ipam-frontend
      context: .
      docker:
        dockerfile: docker/Dockerfile.frontend
      sync:
        manual:
          - src: "frontend/src/**/*"
            dest: /app/src

  tagPolicy:
    gitCommit: {}

deploy:
  kubectl:
    manifests:
      - k8s/namespace.yaml
      - k8s/configmap.yaml
      - k8s/secret.yaml
      - k8s/backend-deployment.yaml
      - k8s/backend-service.yaml
      - k8s/frontend-deployment.yaml
      - k8s/frontend-service.yaml
      - k8s/ingress.yaml

profiles:
  # Development profile
  - name: dev
    activation:
      - command: dev
    build:
      artifacts:
        - image: ipam-backend
          docker:
            dockerfile: docker/Dockerfile.backend.dev
        - image: ipam-frontend
          docker:
            dockerfile: docker/Dockerfile.frontend.dev
    deploy:
      kubectl:
        manifests:
          - k8s/dev/*.yaml
    portForward:
      - resourceType: service
        resourceName: ipam-backend
        port: 8080
        localPort: 8080
      - resourceType: service
        resourceName: ipam-frontend
        port: 3000
        localPort: 3000

  # Staging profile
  - name: staging
    build:
      tagPolicy:
        sha256: {}
    deploy:
      kubectl:
        manifests:
          - k8s/staging/*.yaml

  # Production profile
  - name: prod
    build:
      tagPolicy:
        sha256: {}
      artifacts:
        - image: ipam-backend
          docker:
            dockerfile: docker/Dockerfile.backend
            buildArgs:
              GO_VERSION: "1.21"
        - image: ipam-frontend
          docker:
            dockerfile: docker/Dockerfile.frontend
            buildArgs:
              NODE_VERSION: "20"
    deploy:
      kubectl:
        manifests:
          - k8s/prod/*.yaml
```

**Kubernetes Manifests Structure:**

```
k8s/
├── namespace.yaml              # Namespace definition
├── configmap.yaml              # Application configuration
├── secret.yaml                 # Secrets (database credentials, etc.)
├── backend-deployment.yaml     # Backend deployment
├── backend-service.yaml        # Backend service
├── frontend-deployment.yaml    # Frontend deployment
├── frontend-service.yaml       # Frontend service
├── ingress.yaml                # Ingress for external access
├── dev/
│   ├── backend-deployment.yaml # Dev-specific backend config
│   └── frontend-deployment.yaml # Dev-specific frontend config
├── staging/
│   ├── backend-deployment.yaml
│   ├── frontend-deployment.yaml
│   └── hpa.yaml                # Horizontal Pod Autoscaler
└── prod/
    ├── backend-deployment.yaml
    ├── frontend-deployment.yaml
    ├── hpa.yaml
    └── pdb.yaml                # Pod Disruption Budget
```

**Example Backend Deployment (k8s/backend-deployment.yaml):**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ipam-backend
  namespace: ipam-bananaops
  labels:
    app: ipam-backend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: ipam-backend
  template:
    metadata:
      labels:
        app: ipam-backend
    spec:
      containers:
      - name: backend
        image: ipam-backend
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: SERVER_PORT
          value: "8080"
        - name: DATABASE_TYPE
          valueFrom:
            configMapKeyRef:
              name: ipam-config
              key: database.type
        - name: DATABASE_PATH
          value: "/data/ipam.db"
        volumeMounts:
        - name: data
          mountPath: /data
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: ipam-data
```

**Skaffold Commands:**

```bash
# Development
skaffold dev                    # Continuous development with hot reload
skaffold dev --port-forward     # Dev mode with port forwarding

# Build
skaffold build                  # Build and tag images
skaffold build --push           # Build and push to registry

# Deploy
skaffold deploy                 # Deploy to current kubectl context
skaffold deploy -p staging      # Deploy to staging
skaffold deploy -p prod         # Deploy to production

# Run (build + deploy)
skaffold run                    # Build and deploy once
skaffold run -p prod            # Build and deploy to production

# Debug
skaffold debug                  # Run with debugging enabled

# Cleanup
skaffold delete                 # Delete deployed resources
```

**Docker Configuration:**

**docker/Dockerfile.backend:**
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git protobuf-dev

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source code
COPY backend/ ./
COPY proto/ ../proto/

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /ipam-server cmd/server/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /ipam-server .

EXPOSE 8080

CMD ["./ipam-server"]
```

**docker/Dockerfile.frontend:**
```dockerfile
# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY frontend/package*.json ./
RUN npm ci

# Copy source code
COPY frontend/ ./

# Build
RUN npm run build

# Runtime stage
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

### Local Development Setup

**Prerequisites:**
- Go 1.21+
- Node.js 20+
- Docker
- kubectl
- Task (task runner)
- Skaffold
- Minikube or kind (for local Kubernetes)

**Quick Start:**

```bash
# 1. Install dependencies
task install

# 2. Generate Protobuf code
task proto:generate

# 3. Start local development
task dev

# 4. Run tests
task test

# 5. Build for production
task build
```

**Kubernetes Local Development:**

```bash
# 1. Start local Kubernetes cluster
minikube start
# or
kind create cluster

# 2. Start Skaffold in dev mode
skaffold dev

# 3. Access the application
# Backend: http://localhost:8080
# Frontend: http://localhost:3000
```

### CI/CD Integration

The project can be integrated with CI/CD pipelines using Task and Skaffold:

**GitHub Actions Example:**
```yaml
name: CI/CD

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: arduino/setup-task@v1
      - name: Run tests
        run: task test

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: arduino/setup-task@v1
      - name: Build
        run: task build

  deploy:
    needs: build
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: google-github-actions/setup-gcloud@v1
      - name: Deploy with Skaffold
        run: skaffold run -p prod
```

## Future Enhancements

1. **Cloud Provider Integration**
   - Implement AWS provider using AWS SDK
   - Implement Azure provider using Azure SDK
   - Implement GCP provider using Google Cloud SDK
   - Implement Scaleway provider using Scaleway SDK
   - Implement OVH provider using OVH API

2. **Advanced Features**
   - IP allocation tracking within subnets
   - Subnet hierarchy and parent-child relationships
   - IPAM conflict detection
   - Historical utilization tracking
   - Export/import functionality
   - Multi-tenancy support

3. **Performance Optimizations**
   - Caching layer for frequently accessed subnets
   - Pagination for large subnet lists
   - WebSocket support for real-time updates
   - Database query optimization

4. **Security Enhancements**
   - Authentication and authorization
   - Role-based access control (RBAC)
   - Audit logging
   - Encryption at rest and in transit
