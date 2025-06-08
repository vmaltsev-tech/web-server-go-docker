#!/bin/bash

# Setup script for Go Web Server Monitoring
# This script initializes monitoring infrastructure

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"
}

# Check if Docker and Docker Compose are installed
check_dependencies() {
    log "Checking dependencies..."
    
    if ! command -v docker &> /dev/null; then
        error "Docker is not installed"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        error "Docker Compose is not installed"
        exit 1
    fi
    
    log "Dependencies check passed ✓"
}

# Create necessary directories
create_directories() {
    log "Creating monitoring directories..."
    
    local dirs=(
        "monitoring/prometheus/data"
        "monitoring/grafana/data"
        "monitoring/alertmanager/data"
        "monitoring/grafana/plugins"
    )
    
    for dir in "${dirs[@]}"; do
        mkdir -p "$dir"
        log "Created directory: $dir"
    done
}

# Set correct permissions
set_permissions() {
    log "Setting permissions..."
    
    # Grafana needs specific user permissions
    sudo chown -R 472:472 monitoring/grafana/data 2>/dev/null || {
        warn "Could not set Grafana permissions. You may need to run: sudo chown -R 472:472 monitoring/grafana/data"
    }
    
    # Prometheus data directory
    sudo chown -R 65534:65534 monitoring/prometheus/data 2>/dev/null || {
        warn "Could not set Prometheus permissions. You may need to run: sudo chown -R 65534:65534 monitoring/prometheus/data"
    }
    
    # Alertmanager data directory
    sudo chown -R 65534:65534 monitoring/alertmanager/data 2>/dev/null || {
        warn "Could not set Alertmanager permissions. You may need to run: sudo chown -R 65534:65534 monitoring/alertmanager/data"
    }
}

# Validate configuration files
validate_configs() {
    log "Validating configuration files..."
    
    # Check if prometheus config exists
    if [[ ! -f "monitoring/prometheus/config/prometheus.yml" ]]; then
        error "Prometheus configuration file not found"
        exit 1
    fi
    
    # Check if grafana datasource config exists
    if [[ ! -f "monitoring/grafana/config/provisioning/datasources/prometheus.yml" ]]; then
        error "Grafana datasource configuration not found"
        exit 1
    fi
    
    log "Configuration validation passed ✓"
}

# Display helpful information
display_info() {
    log "Monitoring setup completed successfully!"
    echo ""
    echo "📊 Services will be available at:"
    echo "   • Grafana:     http://localhost:3000 (admin/admin123)"
    echo "   • Prometheus:  http://localhost:9090"
    echo "   • Alertmanager: http://localhost:9093 (if enabled)"
    echo "   • Web Server:  http://localhost:8080"
    echo ""
    echo "🚀 To start monitoring stack:"
    echo "   docker-compose up -d"
    echo ""
    echo "📈 To view logs:"
    echo "   docker-compose logs -f"
    echo ""
    echo "🛑 To stop monitoring stack:"
    echo "   docker-compose down"
    echo ""
}

# Main execution
main() {
    log "Starting monitoring setup..."
    
    check_dependencies
    create_directories
    set_permissions
    validate_configs
    display_info
    
    log "Setup completed! 🎉"
}

# Execute main function
main "$@" 