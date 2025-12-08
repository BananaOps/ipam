# IPAM by BananaOps - Frontend

React + TypeScript frontend application for IP Address Management.

## Project Structure

```
src/
├── components/          # Reusable UI components
│   └── Layout.tsx      # Main layout with header and navigation
├── config/             # Configuration and constants
│   └── constants.ts    # App-wide constants (colors, API URLs, etc.)
├── pages/              # Page components (route handlers)
│   ├── SubnetListPage.tsx
│   ├── SubnetDetailPage.tsx
│   ├── CreateSubnetPage.tsx
│   └── EditSubnetPage.tsx
├── proto/              # Generated Protobuf types
│   ├── subnet.ts       # Generated from proto definitions
│   └── example.ts      # Example usage of proto types
├── services/           # API and external service clients
│   └── api.ts          # REST API client (placeholder)
├── test/               # Test configuration
│   └── setup.ts        # Vitest setup
├── types/              # TypeScript type definitions
│   └── index.ts        # Core application types
├── utils/              # Utility functions
│   └── cloudProviderIcons.ts  # Cloud provider icon mapping
├── App.tsx             # Main app component with routing
├── main.tsx            # Application entry point
└── index.css           # Global styles
```

## Available Scripts

- `npm run dev` - Start development server (Vite)
- `npm run build` - Build for production
- `npm run lint` - Run ESLint
- `npm run preview` - Preview production build
- `npm run test` - Run tests with Vitest

## Technology Stack

- **React 18** - UI framework
- **TypeScript 5** - Type safety
- **Vite** - Build tool and dev server
- **React Router 6** - Client-side routing
- **Axios** - HTTP client
- **FontAwesome 6** - Icons
- **Vitest** - Testing framework
- **fast-check** - Property-based testing

## Routing Structure

- `/` - Subnet list page
- `/subnets` - Subnet list page (alias)
- `/subnets/create` - Create new subnet
- `/subnets/:id` - View subnet details
- `/subnets/:id/edit` - Edit subnet

## Development

The frontend is configured to proxy API requests to `http://localhost:8080` during development.

## Next Steps

The following features are placeholders and will be implemented in subsequent tasks:

- Task 10: Theme system (dark/light/auto mode)
- Task 11: Complete API client implementation
- Task 12: Subnet list component with filtering
- Task 13: Subnet detail component
- Task 14: Subnet creation form
- Task 15: Subnet update functionality
- Task 16: Subnet deletion functionality
- Task 17: Cloud provider icon mapping
- Task 18: Error handling UI
- Task 19: Logo design and implementation
- Task 20: Branding and styling (Cyber Minimal design)
