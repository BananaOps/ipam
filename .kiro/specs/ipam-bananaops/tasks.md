# Implementation Plan - IPAM by BananaOps

- [x] 1. Initialize project structure and dependencies
  - Create Go module for backend with go-ipam dependency
  - Initialize React TypeScript project for frontend
  - Set up Protobuf definitions directory
  - Configure build tools and scripts
  - _Requirements: 4.1, 4.4_

- [x] 2. Define Protobuf API contracts
  - Create Subnet message definition with all fields
  - Create API request/response messages (Create, List, Get, Update, Delete)
  - Create Error message structure
  - Generate Go and TypeScript code from proto files
  - _Requirements: 10.1, 10.4_

- [ ]* 2.1 Write property test for JSON to Protobuf conversion
  - **Property 7: REST to Protobuf conversion consistency**
  - **Validates: Requirements 4.2, 10.2, 10.3**

- [x] 3. Implement backend service layer with Protobuf
  - Create ServiceLayer structure with internal Protobuf operations
  - Implement CreateSubnet service method
  - Implement ListSubnets service method with filtering
  - Implement GetSubnet service method
  - Implement UpdateSubnet service method
  - Implement DeleteSubnet service method
  - Use Protobuf messages for all internal communication
  - _Requirements: 4.2, 10.1_

- [x] 4. Implement REST API Gateway layer
  - Create RESTGateway structure with HTTP routing
  - Implement JSON to Protobuf conversion functions
  - Implement Protobuf to JSON conversion functions
  - Create REST endpoint handlers (POST, GET, PUT, DELETE)
  - Wire REST handlers to ServiceLayer methods
  - Add request validation and error handling
  - Implement CORS support for frontend
  - _Requirements: 4.2, 10.2, 10.3_

- [ ]* 4.1 Write unit tests for JSON/Protobuf conversion
  - Test conversion accuracy for all message types
  - Test error handling for malformed JSON
  - **Validates: Requirements 10.2, 10.3**

- [x] 5. Implement backend database layer
  - Create SubnetRepository interface
  - Implement SQLite repository with CRUD operations
  - Implement MongoDB repository with CRUD operations
  - Add database configuration loading
  - Create database initialization logic with config selection
  - _Requirements: 5.1, 5.2, 5.5_

- [ ]* 5.1 Write property test for database persistence
  - **Property 8: Database persistence round-trip**
  - **Validates: Requirements 5.3**

- [x] 6. Implement IP calculation service with go-ipam
  - Create IPService with go-ipam integration
  - Implement CIDR validation function
  - Implement subnet details calculation (address, netmask, wildcard, network, broadcast, hostMin, hostMax, hostsPerNet)
  - Implement public/private classification logic
  - Implement utilization percentage calculation
  - _Requirements: 2.3, 4.1, 4.3_

- [ ]* 6.1 Write property test for CIDR validation
  - **Property 10: CIDR validation correctness**
  - **Validates: Requirements 8.1**

- [ ]* 6.2 Write property test for subnet property calculation
  - **Property 11: Subnet property calculation completeness**
  - **Validates: Requirements 8.2**

- [ ]* 6.3 Write property test for utilization calculation
  - **Property 14: Utilization calculation accuracy**
  - **Validates: Requirements 9.1**

- [ ] 7. Wire service layer with repositories and IP service
  - Connect ServiceLayer to SubnetRepository
  - Connect ServiceLayer to IPService
  - Implement business logic for subnet operations
  - Add validation in service methods
  - Implement error handling and error message creation
  - _Requirements: 8.1, 8.2, 8.3, 8.4, 1.2_

- [ ]* 7.1 Write property test for invalid subnet error handling
  - **Property 4: Invalid subnet error handling**
  - **Validates: Requirements 2.5, 8.5**

- [ ]* 7.2 Write property test for subnet update recalculation
  - **Property 12: Subnet update recalculation**
  - **Validates: Requirements 8.3**

- [ ]* 7.3 Write property test for subnet deletion
  - **Property 13: Subnet deletion completeness**
  - **Validates: Requirements 8.4**

- [ ]* 7.4 Write property test for structured error responses
  - **Property 16: Structured error responses**
  - **Validates: Requirements 10.5**

- [ ] 8. Implement cloud provider module interface
  - Create CloudProvider interface definition
  - Create CloudProviderManager with provider registry
  - Add cloud provider types (AWS, Azure, GCP, Scaleway, OVH)
  - Implement error handling for provider unavailability
  - _Requirements: 6.1, 6.2, 6.5_

- [ ]* 8.1 Write property test for cloud provider error resilience
  - **Property 9: Cloud provider error resilience**
  - **Validates: Requirements 6.5**

- [ ] 9. Set up frontend project structure
  - Create React app with TypeScript template
  - Install dependencies (FontAwesome, Axios for HTTP requests)
  - Set up routing structure
  - Configure build and development scripts
  - _Requirements: 7.1_

- [ ] 10. Implement theme system
  - Create ThemeContext with dark/light/auto support
  - Implement useSystemTheme hook for system preference detection
  - Implement theme persistence to localStorage
  - Create theme provider component
  - Define color palette constants
  - Apply theme styling to root components
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5_

- [ ]* 10.1 Write property test for theme persistence
  - **Property 5: Theme persistence round-trip**
  - **Validates: Requirements 3.5**

- [ ]* 10.2 Write property test for system theme synchronization
  - **Property 6: System theme synchronization**
  - **Validates: Requirements 3.4**

- [ ] 11. Create REST API client service
  - Create APIClient class with Axios
  - Implement HTTP methods for all endpoints (POST, GET, PUT, DELETE)
  - Add request/response interceptors for error handling
  - Implement TypeScript interfaces for request/response types
  - Add authentication headers support (for future use)
  - Configure base URL and timeout settings
  - _Requirements: 10.2, 10.3_

- [ ]* 11.1 Write unit tests for API client
  - Test request formatting
  - Test error handling
  - **Validates: Requirements 10.2, 10.3**

- [ ] 12. Implement subnet list component
  - Create SubnetList component with table/grid layout
  - Implement location filter UI (datacenter, site, cloud)
  - Implement cloud provider filter with FontAwesome icons
  - Add search functionality
  - Display basic subnet information (CIDR, name, location)
  - Show cloud provider icon, region, and account for cloud subnets
  - Implement empty state display
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_

- [ ]* 12.1 Write property test for location filter
  - **Property 1: Location filter correctness**
  - **Validates: Requirements 1.2**

- [ ]* 12.2 Write property test for cloud provider metadata display
  - **Property 2: Cloud provider metadata completeness**
  - **Validates: Requirements 1.3, 1.4, 6.3**

- [ ] 13. Implement subnet detail component
  - Create SubnetDetail component with detailed layout
  - Display subnet name and description
  - Display all subnet properties (address, netmask, wildcard, network, type, broadcast, hostMin, hostMax, hostsPerNet)
  - Display public/private classification
  - Implement utilization percentage display
  - Create progress bar component for utilization
  - Add visual indicators for high utilization
  - _Requirements: 2.1, 2.2, 2.4, 9.2, 9.3, 9.5_

- [ ]* 13.1 Write property test for subnet detail completeness
  - **Property 3: Subnet detail completeness**
  - **Validates: Requirements 2.1, 2.2, 2.4**

- [ ]* 13.2 Write property test for high utilization visual indication
  - **Property 15: High utilization visual indication**
  - **Validates: Requirements 9.3**

- [ ] 14. Implement subnet creation form
  - Create form component with CIDR input
  - Add name and description fields
  - Add location field
  - Add cloud provider selection with region and account fields
  - Implement form validation
  - Connect to API create endpoint via APIClient
  - Display success/error feedback
  - _Requirements: 8.1, 8.2, 8.5_

- [ ] 15. Implement subnet update functionality
  - Create edit form component
  - Pre-populate form with existing subnet data
  - Implement update submission
  - Handle recalculation on CIDR change
  - Display success/error feedback
  - _Requirements: 8.3, 8.5_

- [ ] 16. Implement subnet deletion functionality
  - Add delete button to subnet detail view
  - Implement confirmation dialog
  - Connect to API delete endpoint via APIClient
  - Handle navigation after deletion
  - Display success/error feedback
  - _Requirements: 8.4, 8.5_

- [ ] 17. Create cloud provider icon mapping
  - Define FontAwesome icon constants for each provider
  - Create CloudProviderIcon component
  - Implement icon selection logic based on provider type
  - Apply consistent styling
  - _Requirements: 1.4, 6.3_

- [ ] 18. Implement error handling UI
  - Create error boundary component
  - Create error message display component
  - Implement toast/notification system for errors
  - Add retry functionality for failed requests
  - Display user-friendly error messages
  - _Requirements: 2.5, 8.5_

- [ ] 19. Design and implement logo
  - Create logo following Cyber Minimal style
  - Use color palette (Bleu nuit, Bleu cyan, Gris clair, Blanc)
  - Export logo in multiple formats (SVG, PNG)
  - Integrate logo into application header
  - _Requirements: 7.3, 7.5_

- [ ] 20. Apply branding and styling
  - Implement Cyber Minimal design system
  - Apply color palette consistently across components
  - Style all UI components with modern, tech-focused aesthetic
  - Ensure visual consistency across all pages
  - Add responsive design for mobile/tablet
  - _Requirements: 7.2, 7.4_

- [ ] 21. Implement backend server startup
  - Create main server entry point
  - Load configuration from file/environment
  - Initialize database based on configuration
  - Initialize go-ipam module
  - Wire REST Gateway to Service Layer
  - Start HTTP server with REST API routes
  - Add graceful shutdown handling
  - _Requirements: 4.1, 4.4, 5.1, 5.2_

- [ ]* 21.1 Write unit test for database configuration validation
  - Test invalid configuration handling
  - **Validates: Requirements 4.5**

- [ ] 22. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 23. Create Taskfile for development workflow
  - Create Taskfile.yml at project root
  - Add development tasks (dev, dev:backend, dev:frontend)
  - Add build tasks (build, build:backend, build:frontend)
  - Add test tasks (test, test:backend, test:frontend, test:backend:coverage)
  - Add lint tasks (lint, lint:backend, lint:frontend)
  - Add proto:generate task for Protobuf code generation
  - Add docker tasks (docker:build, docker:push)
  - Add database tasks (db:migrate, db:seed)
  - Add utility tasks (clean, install)
  - _Requirements: 10.4_

- [ ] 24. Create Docker configuration
  - Create docker/Dockerfile.backend for production backend image
  - Create docker/Dockerfile.frontend for production frontend image
  - Create docker/Dockerfile.backend.dev for development backend image
  - Create docker/Dockerfile.frontend.dev for development frontend image
  - Create docker/nginx.conf for frontend nginx configuration
  - Create .dockerignore files
  - _Requirements: 10.4_

- [ ] 25. Create Kubernetes manifests
  - Create k8s/namespace.yaml
  - Create k8s/configmap.yaml for application configuration
  - Create k8s/secret.yaml for sensitive data
  - Create k8s/backend-deployment.yaml with health checks
  - Create k8s/backend-service.yaml
  - Create k8s/frontend-deployment.yaml
  - Create k8s/frontend-service.yaml
  - Create k8s/ingress.yaml for external access
  - Create k8s/pvc.yaml for persistent storage
  - _Requirements: 10.4_

- [ ] 26. Create Skaffold configuration
  - Create skaffold.yaml at project root
  - Configure build artifacts for backend and frontend
  - Configure kubectl deployment
  - Create dev profile with hot reload and port forwarding
  - Create staging profile with appropriate settings
  - Create prod profile with optimized builds
  - Configure file sync for development
  - _Requirements: 10.4_

- [ ] 27. Create environment-specific Kubernetes configs
  - Create k8s/dev/ directory with development configs
  - Create k8s/staging/ directory with staging configs (HPA)
  - Create k8s/prod/ directory with production configs (HPA, PDB)
  - Configure resource limits for each environment
  - Configure replicas for each environment
  - _Requirements: 10.4_

- [ ] 28. Write documentation
  - Create README with project overview
  - Document REST API endpoints with examples
  - Document configuration options
  - Create user guide for frontend features
  - Document development setup with Task commands
  - Document Kubernetes deployment with Skaffold
  - Create CONTRIBUTING.md with development workflow
  - Document local development setup (prerequisites, quick start)
  - _Requirements: All_

- [ ] 29. Final checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

