#!/bin/bash

# Docker build and run script for Go Echo API
# This script helps build and run the application with Docker

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    print_warning "Docker Compose is not installed. Using Docker directly."
    USE_COMPOSE=false
else
    USE_COMPOSE=true
fi

# Function to build Docker image
build_image() {
    print_status "Building Docker image..."
    docker build -t go-echo-api .
    print_status "Docker image built successfully!"
}

# Function to run with Docker Compose
run_with_compose() {
    print_status "Starting services with Docker Compose..."
    docker-compose up -d
    print_status "Services started successfully!"
    print_status "API will be available at http://localhost:8080"
    print_status "Database will be available at localhost:5432"
    print_status "Use 'docker-compose logs -f app' to view logs"
}

# Function to run with Docker directly
run_with_docker() {
    print_status "Starting container with Docker..."
    
    # Check if config.yaml exists
    if [ ! -f "config.yaml" ]; then
        print_warning "config.yaml not found. Creating from example..."
        cp config.yaml.example config.yaml
    fi
    
    docker run -d \
        --name go-echo-api \
        -p 8080:8080 \
        -v "$(pwd)/config.yaml:/root/config.yaml" \
        -v "$(pwd)/database/migrations:/root/database/migrations" \
        go-echo-api
    
    print_status "Container started successfully!"
    print_status "API will be available at http://localhost:8080"
    print_status "Use 'docker logs go-echo-api' to view logs"
}

# Main script logic
case "${1:-build}" in
    "build")
        build_image
        ;;
    "run")
        if [ "$USE_COMPOSE" = true ]; then
            run_with_compose
        else
            run_with_docker
        fi
        ;;
    "compose")
        if [ "$USE_COMPOSE" = true ]; then
            build_image
            run_with_compose
        else
            print_error "Docker Compose is not available"
            exit 1
        fi
        ;;
    "docker")
        build_image
        run_with_docker
        ;;
    "clean")
        print_status "Cleaning up Docker resources..."
        docker-compose down 2>/dev/null || true
        docker rm -f go-echo-api 2>/dev/null || true
        docker rmi go-echo-api 2>/dev/null || true
        print_status "Cleanup completed!"
        ;;
    *)
        echo "Usage: $0 {build|run|compose|docker|clean}"
        echo ""
        echo "Commands:"
        echo "  build   - Build Docker image only"
        echo "  run     - Build and run (auto-detects Docker Compose)"
        echo "  compose - Build and run with Docker Compose"
        echo "  docker  - Build and run with Docker directly"
        echo "  clean   - Clean up Docker resources"
        exit 1
        ;;
esac 