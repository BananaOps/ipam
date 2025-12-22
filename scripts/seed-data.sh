#!/bin/bash

# IPAM K6 Data Seeding Script
# Usage: ./seed-data.sh [clean]
# 
# Options:
#   clean    - Delete existing data before seeding
#   (none)   - Add data to existing database (handles duplicates gracefully)

set -e

# Configuration
BASE_URL="http://localhost:8080"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
K6_SCRIPT="$SCRIPT_DIR/k6-seed-data.js"

# Check if k6 is installed
if ! command -v k6 &> /dev/null; then
    echo "❌ k6 is not installed. Please install k6 first:"
    echo "   https://k6.io/docs/getting-started/installation/"
    exit 1
fi

# Check if backend is running
echo "🔍 Checking if backend is running..."
if ! curl -s "$BASE_URL/api/v1/subnets" > /dev/null; then
    echo "❌ Backend is not running or not accessible at $BASE_URL"
    echo "   Please start the backend first with: task dev:backend"
    exit 1
fi

echo "✅ Backend is accessible"

# Set environment variables
export BASE_URL="$BASE_URL"

# Check if clean option is provided
if [[ "$1" == "clean" ]]; then
    echo "🧹 Running with CLEAN_FIRST=true (will delete existing data)"
    export CLEAN_FIRST="true"
else
    echo "📝 Running in append mode (will handle duplicates gracefully)"
    export CLEAN_FIRST="false"
fi

# Run the K6 script
echo "🚀 Starting data seeding..."
echo "   Script: $K6_SCRIPT"
echo "   Base URL: $BASE_URL"
echo "   Clean first: $CLEAN_FIRST"
echo ""

k6 run "$K6_SCRIPT"

echo ""
echo "🎉 Seeding completed!"
echo "📋 You can now:"
echo "   • View the frontend at http://localhost:3000 or http://localhost:3001"
echo "   • Check VPC/Subnet badges and relationships"
echo "   • Click on VPCs to see their child subnets"
