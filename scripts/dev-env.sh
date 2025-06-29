#!/bin/bash

# Development Environment Setup Script for Vapi CLI
# This script helps developers quickly switch between environments

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

show_usage() {
    echo "Vapi CLI Development Environment Setup"
    echo ""
    echo "Usage: $0 <command> [arguments]"
    echo ""
    echo "Commands:"
    echo "  setup <env>     Set up environment (production|staging|development)"
    echo "  status          Show current environment status"
    echo "  reset           Reset to production environment"
    echo "  local           Set up for local development"
    echo ""
    echo "Examples:"
    echo "  $0 setup staging"
    echo "  $0 local"
    echo "  $0 status"
}

setup_environment() {
    local env="$1"
    
    case "$env" in
        production|prod)
            log_info "Setting up production environment"
            unset VAPI_ENV
            unset VAPI_API_BASE_URL
            unset VAPI_DASHBOARD_URL
            go run . config set environment production
            ;;
        staging|stage)
            log_info "Setting up staging environment"
            export VAPI_ENV=staging
            go run . config set environment staging
            ;;
        development|dev|local)
            log_info "Setting up development environment"
            export VAPI_ENV=development
            export VAPI_API_BASE_URL=http://localhost:3000
            export VAPI_DASHBOARD_URL=http://localhost:3001
            go run . config set environment development
            ;;
        *)
            log_error "Unknown environment: $env"
            echo "Valid environments: production, staging, development"
            exit 1
            ;;
    esac
    
    log_success "Environment set to: $env"
    show_current_status
}

show_current_status() {
    echo ""
    log_info "Current Environment Status:"
    
    # Show config
    go run . config get environment
    go run . config get base_url
    go run . config get dashboard_url
    
    # Show relevant environment variables
    echo ""
    if [[ -n "$VAPI_ENV" ]]; then
        echo "VAPI_ENV: $VAPI_ENV"
    fi
    if [[ -n "$VAPI_API_BASE_URL" ]]; then
        echo "VAPI_API_BASE_URL: $VAPI_API_BASE_URL"
    fi
    if [[ -n "$VAPI_DASHBOARD_URL" ]]; then
        echo "VAPI_DASHBOARD_URL: $VAPI_DASHBOARD_URL"
    fi
}

setup_local_development() {
    log_info "Setting up local development environment"
    
    # Set environment variables for current session
    export VAPI_ENV=development
    export VAPI_API_BASE_URL=http://localhost:3000
    export VAPI_DASHBOARD_URL=http://localhost:3001
    
    # Update CLI config
    go run . config set environment development
    
    log_success "Local development environment configured"
    log_warning "Make sure your local Vapi server is running on localhost:3000"
    log_warning "Make sure your local dashboard is running on localhost:3001"
    
    show_current_status
    
    echo ""
    log_info "To persist these settings in your shell, add to your ~/.bashrc or ~/.zshrc:"
    echo "export VAPI_ENV=development"
    echo "export VAPI_API_BASE_URL=http://localhost:3000"
    echo "export VAPI_DASHBOARD_URL=http://localhost:3001"
}

reset_environment() {
    log_info "Resetting to production environment"
    
    # Clear environment variables
    unset VAPI_ENV
    unset VAPI_API_BASE_URL
    unset VAPI_DASHBOARD_URL
    
    # Reset config
    go run . config set environment production
    
    log_success "Reset to production environment"
    show_current_status
}

main() {
    cd "$PROJECT_ROOT"
    
    case "${1:-}" in
        setup)
            if [[ -z "${2:-}" ]]; then
                log_error "Environment required"
                show_usage
                exit 1
            fi
            setup_environment "$2"
            ;;
        status)
            show_current_status
            ;;
        local)
            setup_local_development
            ;;
        reset)
            reset_environment
            ;;
        help|--help|-h)
            show_usage
            ;;
        "")
            log_warning "No command specified"
            show_usage
            exit 1
            ;;
        *)
            log_error "Unknown command: $1"
            show_usage
            exit 1
            ;;
    esac
}

main "$@" 