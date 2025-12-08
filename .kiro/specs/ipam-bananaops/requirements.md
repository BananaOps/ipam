# Requirements Document

## Introduction

IPAM by BananaOps est une application de gestion d'adresses IP (IP Address Management) moderne avec une architecture full-stack. L'application permet de gérer, visualiser et analyser des sous-réseaux IP avec support pour les environnements on-premise et cloud. Le système utilise Go pour le backend avec calcul IP natif, MongoDB ou SQLite embarqué pour le stockage, et React/TypeScript pour une interface utilisateur moderne avec support multi-thèmes.

## Glossary

- **IPAM System**: Le système complet de gestion d'adresses IP incluant backend, base de données et frontend
- **Backend Service**: Le service Go qui gère la logique métier, les calculs IP et la communication Protobuf
- **Frontend Application**: L'application React/TypeScript qui fournit l'interface utilisateur
- **Subnet**: Un sous-réseau IP avec ses propriétés (adresse, masque, plage d'hôtes)
- **Cloud Provider**: Un fournisseur de services cloud (AWS, Azure, Google Cloud, Scaleway, OVH)
- **Theme System**: Le système de gestion des thèmes visuels (dark, light, auto)
- **go-ipam Module**: Le module Go utilisé pour la validation et le calcul des adresses IP
- **Protobuf API**: L'API de communication utilisant Protocol Buffers entre backend et frontend
- **Embedded Database**: La base de données embarquée (MongoDB ou SQLite) configurée au démarrage

## Requirements

### Requirement 1

**User Story:** En tant qu'administrateur réseau, je veux visualiser tous mes sous-réseaux avec des filtres, afin de pouvoir rapidement localiser et gérer mes ressources IP.

#### Acceptance Criteria

1. WHEN a user accesses the subnet listing page THEN the IPAM System SHALL display all subnets with their basic information
2. WHEN a user applies a location filter THEN the IPAM System SHALL display only subnets matching the selected datacenter, site, or cloud provider
3. WHEN displaying cloud provider subnets THEN the IPAM System SHALL show the region and associated account for each subnet
4. WHEN rendering cloud provider entries THEN the Frontend Application SHALL display the appropriate FontAwesome icon for each provider (AWS, Azure, Google Cloud, Scaleway, OVH)
5. WHEN the subnet list is empty THEN the IPAM System SHALL display an appropriate empty state message

### Requirement 2

**User Story:** En tant qu'administrateur réseau, je veux voir les détails complets d'un sous-réseau, afin de comprendre son utilisation et ses caractéristiques techniques.

#### Acceptance Criteria

1. WHEN a user selects a subnet THEN the IPAM System SHALL display the subnet address, netmask, wildcard, network address, type, broadcast address, HostMin, HostMax, and Hosts/Net count
2. WHEN displaying subnet details THEN the IPAM System SHALL indicate whether the subnet is public or private
3. WHEN calculating subnet information THEN the Backend Service SHALL use the go-ipam module for all IP calculations and validations
4. WHEN displaying subnet utilization THEN the Frontend Application SHALL show a percentage indicator and progress bar representing IP address usage
5. WHEN subnet data is invalid THEN the IPAM System SHALL display appropriate error messages

### Requirement 3

**User Story:** En tant qu'utilisateur, je veux choisir un thème visuel adapté à mes préférences, afin d'améliorer mon confort d'utilisation.

#### Acceptance Criteria

1. WHEN a user selects dark mode THEN the Frontend Application SHALL apply the dark theme using the color palette (Bleu nuit #0A1A2F, Bleu cyan #0EA5E9)
2. WHEN a user selects light mode THEN the Frontend Application SHALL apply the light theme using the color palette (Gris clair #F3F4F6, Blanc pur #FFFFFF)
3. WHEN a user selects auto mode THEN the Frontend Application SHALL synchronize with the system theme preference
4. WHEN the system theme changes in auto mode THEN the Frontend Application SHALL update the theme automatically
5. WHEN a theme is selected THEN the Frontend Application SHALL persist the user preference across sessions

### Requirement 4

**User Story:** En tant que développeur, je veux une architecture backend modulaire avec communication Protobuf, afin de garantir performance et extensibilité.

#### Acceptance Criteria

1. WHEN the Backend Service starts THEN it SHALL initialize the go-ipam module for IP management
2. WHEN the Backend Service receives a request THEN it SHALL communicate using Protobuf message format
3. WHEN processing IP calculations THEN the Backend Service SHALL use the go-ipam module for validation and computation
4. WHEN the Backend Service starts THEN it SHALL load the configured embedded database (MongoDB or SQLite)
5. WHEN database configuration is invalid THEN the Backend Service SHALL fail startup with a clear error message

### Requirement 5

**User Story:** En tant qu'administrateur système, je veux choisir entre MongoDB embarqué et SQLite embarqué, afin d'adapter le stockage à mes besoins de déploiement.

#### Acceptance Criteria

1. WHEN the Backend Service starts with MongoDB configuration THEN it SHALL initialize an embedded MongoDB instance
2. WHEN the Backend Service starts with SQLite configuration THEN it SHALL initialize an embedded SQLite database
3. WHEN storing subnet data THEN the Embedded Database SHALL persist all subnet properties and metadata
4. WHEN querying subnet data THEN the Embedded Database SHALL return results efficiently with proper indexing
5. WHEN the database configuration changes THEN the Backend Service SHALL require a restart to apply the new configuration

### Requirement 6

**User Story:** En tant qu'administrateur réseau, je veux une architecture extensible pour les cloud providers, afin de pouvoir intégrer la récupération automatique d'IPs à l'avenir.

#### Acceptance Criteria

1. WHEN the Backend Service is structured THEN it SHALL include a modular cloud provider interface for future extensions
2. WHEN adding a new cloud provider THEN the IPAM System SHALL support the provider without modifying core logic
3. WHEN displaying cloud subnets THEN the Frontend Application SHALL render provider-specific icons and metadata
4. WHEN cloud provider modules are added THEN the Backend Service SHALL support dynamic IP retrieval from provider APIs (AWS, Azure, Google Cloud, Scaleway, OVH)
5. WHEN a cloud provider is unavailable THEN the IPAM System SHALL handle the error gracefully and display cached data

### Requirement 7

**User Story:** En tant qu'utilisateur, je veux une interface moderne avec un style "Cyber Minimal", afin d'avoir une expérience professionnelle et agréable.

#### Acceptance Criteria

1. WHEN the Frontend Application renders THEN it SHALL use FontAwesome icons throughout the interface
2. WHEN displaying the application THEN the Frontend Application SHALL present a modern, minimalist, tech-focused design
3. WHEN showing the branding THEN the Frontend Application SHALL display a logo consistent with the Cyber Minimal style and color palette
4. WHEN rendering UI components THEN the Frontend Application SHALL maintain visual consistency across all pages
5. WHEN the application loads THEN the Frontend Application SHALL display the IPAM by BananaOps branding prominently

### Requirement 8

**User Story:** En tant qu'administrateur réseau, je veux créer et gérer des sous-réseaux, afin de maintenir mon inventaire IP à jour.

#### Acceptance Criteria

1. WHEN a user creates a subnet THEN the Backend Service SHALL validate the IP address and CIDR notation using go-ipam
2. WHEN a user creates a subnet THEN the IPAM System SHALL calculate and store all subnet properties automatically
3. WHEN a user updates a subnet THEN the Backend Service SHALL recalculate affected properties and persist changes
4. WHEN a user deletes a subnet THEN the IPAM System SHALL remove the subnet and update related data
5. WHEN subnet operations fail THEN the IPAM System SHALL provide clear error messages indicating the cause

### Requirement 9

**User Story:** En tant qu'administrateur réseau, je veux suivre l'utilisation des adresses IP dans mes sous-réseaux, afin d'anticiper les besoins en capacité.

#### Acceptance Criteria

1. WHEN calculating subnet utilization THEN the Backend Service SHALL determine the percentage of allocated IP addresses
2. WHEN displaying utilization THEN the Frontend Application SHALL show a visual progress bar with percentage
3. WHEN a subnet reaches high utilization THEN the Frontend Application SHALL use visual indicators to highlight the status
4. WHEN utilization data is updated THEN the IPAM System SHALL reflect changes in real-time
5. WHEN no IPs are allocated THEN the Frontend Application SHALL display zero percent utilization

### Requirement 10

**User Story:** En tant que développeur, je veux une API Protobuf bien définie, afin d'assurer une communication efficace et typée entre frontend et backend.

#### Acceptance Criteria

1. WHEN defining API contracts THEN the Backend Service SHALL use Protobuf message definitions for all endpoints
2. WHEN the Frontend Application makes requests THEN it SHALL serialize data using Protobuf format
3. WHEN the Backend Service responds THEN it SHALL return Protobuf-encoded messages
4. WHEN API schemas change THEN the IPAM System SHALL regenerate client and server code from Protobuf definitions
5. WHEN communication errors occur THEN the Protobuf API SHALL provide structured error responses
