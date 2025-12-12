#!/bin/bash

# Stop script for data-sync project
# This script stops all running services

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Function to print colored messages
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_info "Stopping all services..."

# Stop main docker-compose services (Trino + Go backend)
print_info "Stopping Trino and Go backend..."
docker-compose down

# Stop data sources
print_info "Stopping data sources..."
cd data-sources
docker-compose down
cd ..

# Kill any running npm processes for the frontend
print_info "Stopping frontend (if running)..."
pkill -f "vite" 2>/dev/null || true

# Remove the Docker network if it exists
print_info "Removing Docker network..."
docker network rm trino-network 2>/dev/null || true

print_success "All services stopped!"
