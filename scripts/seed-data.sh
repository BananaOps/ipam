#!/bin/bash

# IPAM Data Seeding Script
# This script uses k6 to populate the IPAM with sample data

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default configuration
BASE_URL=${BASE_URL:-"http://localhost:8081"}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
K6_SCRIPT="$SCRIPT_DIR/k6-seed-data.js"

echo -e "${BLUE}🚀 IPAM Data Seeding Script${NC}"
echo -e "${BLUE}=============================${NC}"
echo ""

# Check if k6 is installed
if ! command -v k6 &> /dev/null; then
    echo -e "${RED}❌ k6 is not installed.${NC}"
    echo -e "${YELLOW}Please install k6 from: https://k6.io/docs/getting-started/installation/${NC}"
    echo ""
    echo -e "${YELLOW}Quick install options:${NC}"
    echo -e "${YELLOW}  macOS: brew install k6${NC}"
    echo -e "${YELLOW}  Ubuntu/Debian: sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69${NC}"
    echo -e "${YELLOW}                 echo 'deb https://dl.k6.io/deb stable main' | sudo tee /etc/apt/sources.list.d/k6.list${NC}"
    echo -e "${YELLOW}                 sudo apt-get update && sudo apt-get install k6${NC}"
    exit 1
fi

# Check if the backend is running
echo -e "${BLUE}🔍 Checking if IPAM backend is running at ${BASE_URL}...${NC}"
if ! curl -s -f "${BASE_URL}/api/v1/subnets" > /dev/null; then
    echo -e "${RED}❌ IPAM backend is not running or not accessible at ${BASE_URL}${NC}"
    echo -e "${YELLOW}Please make sure:${NC}"
    echo -e "${YELLOW}  1. The backend is running (task dev:backend)${NC}"
    echo -e "${YELLOW}  2. The backend is accessible at ${BASE_URL}${NC}"
    echo -e "${YELLOW}  3. Or set BASE_URL environment variable: BASE_URL=http://your-server:port $0${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Backend is running and accessible${NC}"
echo ""

# Run the k6 script
echo -e "${BLUE}🌱 Starting data seeding with k6...${NC}"
echo -e "${BLUE}Script: ${K6_SCRIPT}${NC}"
echo -e "${BLUE}Target: ${BASE_URL}${NC}"
echo ""

# Export BASE_URL for k6 script
export BASE_URL

# Run k6
if k6 run "$K6_SCRIPT"; then
    echo ""
    echo -e "${GREEN}🎉 Data seeding completed successfully!${NC}"
    echo -e "${GREEN}You can now access your IPAM at: ${BASE_URL%:*}:3000${NC}"
else
    echo ""
    echo -e "${RED}❌ Data seeding failed. Check the output above for details.${NC}"
    exit 1
fi
