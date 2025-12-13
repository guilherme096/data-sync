#!/bin/bash

# Startup script for data-sync project
# This script starts all components: data sources, Trino, Go backend, and frontend

set -e  # Exit on error

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

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if a service is ready
wait_for_service() {
    local service_name=$1
    local host=$2
    local port=$3
    local max_attempts=${4:-30}
    local attempt=1

    print_info "Waiting for $service_name to be ready at $host:$port..."

    while [ $attempt -le $max_attempts ]; do
        if nc -z "$host" "$port" 2>/dev/null; then
            print_success "$service_name is ready!"
            return 0
        fi
        echo -n "."
        sleep 2
        attempt=$((attempt + 1))
    done

    print_error "$service_name failed to start within expected time"
    return 1
}

# Function to cleanup on exit
cleanup() {
    print_warning "\nShutting down services..."

    # Kill frontend if running
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
    fi

    # Stop Docker services
    print_info "Stopping Docker containers..."
    docker-compose down 2>/dev/null || true
    cd "$SCRIPT_DIR/data-sources" && docker-compose down 2>/dev/null || true
    cd "$SCRIPT_DIR"

    print_success "Cleanup complete"
}

trap cleanup EXIT INT TERM

# Check if required commands are available
print_info "Checking prerequisites..."
command -v docker >/dev/null 2>&1 || { print_error "docker is not installed. Aborting."; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { print_error "docker-compose is not installed. Aborting."; exit 1; }
command -v npm >/dev/null 2>&1 || { print_error "npm is not installed. Aborting."; exit 1; }
command -v nc >/dev/null 2>&1 || { print_error "nc (netcat) is not installed. Aborting."; exit 1; }
print_success "All prerequisites are available"

# Create logs directory
mkdir -p logs

# Step 1: Create Docker network if it doesn't exist
print_info "========================================="
print_info "Step 1: Creating Docker network..."
print_info "========================================="
if ! docker network inspect trino-network >/dev/null 2>&1; then
    docker network create trino-network
    print_success "Created trino-network"
else
    print_info "trino-network already exists"
fi

# Step 2: Start data sources
print_info "========================================="
print_info "Step 2: Starting data sources..."
print_info "========================================="
cd data-sources
docker-compose up -d
cd ..

# Wait for data sources to be ready
wait_for_service "PostgreSQL" "localhost" "5433" 60
wait_for_service "MySQL" "localhost" "3307" 60
wait_for_service "MongoDB" "localhost" "27017" 60

print_success "All data sources are running and initialized"

# Give databases a moment to fully initialize
sleep 3

# Step 3: Start Trino cluster first
print_info "========================================="
print_info "Step 3: Starting Trino cluster..."
print_info "========================================="
docker-compose up -d trino-coordinator trino-worker-1 trino-worker-2 trino-worker-3

# Wait for Trino to be ready
wait_for_service "Trino" "localhost" "8080" 90

# Additional check for Trino readiness via HTTP
print_info "Verifying Trino is fully operational..."
sleep 5
if curl -s http://localhost:8080/v1/info > /dev/null; then
    print_success "Trino cluster is fully operational"
else
    print_warning "Trino may not be fully ready yet, but continuing..."
fi

# Step 4: Start Go backend after Trino is ready
print_info "========================================="
print_info "Step 4: Starting Go backend..."
print_info "========================================="
docker-compose up -d --build data-sync

# Wait for Go backend to be ready
wait_for_service "Go Backend" "localhost" "8081" 60
print_success "Go backend is running"

# Step 5: Start Frontend
print_info "========================================="
print_info "Step 5: Starting Frontend..."
print_info "========================================="

cd interface/data-sync

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    print_info "Installing frontend dependencies..."
    npm install
fi

# Start frontend in background
print_info "Starting frontend development server..."
npm run dev > ../../logs/frontend.log 2>&1 &
FRONTEND_PID=$!
cd ../..

# Wait a bit for frontend to start
sleep 5

print_success "Frontend is starting (PID: $FRONTEND_PID)"

# Summary
print_info "========================================="
print_success "All services started successfully!"
print_info "========================================="
echo ""
print_info "Services:"
echo "  - PostgreSQL:     localhost:5433"
echo "    - User: testuser, Password: testpass, DB: testdb"
echo "  - MySQL:          localhost:3307"
echo "    - User: testuser, Password: testpass, DB: testdb"
echo "  - MongoDB:        localhost:27017"
echo "    - User: admin, Password: adminpass"
echo "  - Trino:          http://localhost:8080"
echo "  - Go Backend:     http://localhost:8081"
echo "  - Frontend:       http://localhost:5173 (typically)"
echo ""
print_info "Docker containers:"
docker-compose ps
echo ""
cd data-sources
docker-compose ps
cd ..
echo ""
print_info "Frontend log: logs/frontend.log"
echo ""
print_info "Press Ctrl+C to stop all services"
echo ""

# Keep script running and tail logs
tail -f logs/frontend.log 2>/dev/null || wait
