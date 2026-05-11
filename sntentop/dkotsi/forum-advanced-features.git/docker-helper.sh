#!/bin/bash
# ═══════════════════════════════════════════════════════════════════════════
# DOCKER MANAGEMENT HELPER SCRIPT FOR RUNNING FORUM
# ═══════════════════════════════════════════════════════════════════════════
# This script provides easy commands to manage the Dockerized Running Forum
# Usage: ./docker-helper.sh [command]
# ═══════════════════════════════════════════════════════════════════════════

set -e  # Exit on error

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

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

# Function to show usage
show_usage() {
    cat << EOF
═══════════════════════════════════════════════════════════════════════════
Running Forum Docker Helper Script
═══════════════════════════════════════════════════════════════════════════

USAGE: ./docker-helper.sh [COMMAND]

COMMANDS:

  start         Start the forum (builds if needed)
  stop          Stop the forum (keeps data)
  restart       Restart the forum
  rebuild       Rebuild and restart (after code changes)
  logs          Show logs (follow mode)
  status        Show container status
  backup        Backup the database
  restore       Restore database from backup
  clean         Remove everything (including data!)
  shell         Open a shell inside the container
  health        Check container health
  help          Show this help message

EXAMPLES:

  ./docker-helper.sh start        # Start the forum
  ./docker-helper.sh logs         # View logs
  ./docker-helper.sh backup       # Backup database
  ./docker-helper.sh rebuild      # After editing code

ACCESS:
  Frontend: https://localhost:3000
  Backend:  http://localhost:8080

═══════════════════════════════════════════════════════════════════════════
EOF
}

# Check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running! Please start Docker first."
        exit 1
    fi
}

# Start the forum
cmd_start() {
    print_info "Starting Running Forum..."
    docker-compose up -d --build
    print_success "Forum started!"
    print_info "Frontend: https://localhost:3000"
    print_info "Backend:  http://localhost:8080"
    print_info "View logs with: ./docker-helper.sh logs"
}

# Stop the forum
cmd_stop() {
    print_info "Stopping Running Forum..."
    docker-compose down
    print_success "Forum stopped! (Data preserved)"
}

# Restart the forum
cmd_restart() {
    print_info "Restarting Running Forum..."
    docker-compose restart
    print_success "Forum restarted!"
}

# Rebuild and restart
cmd_rebuild() {
    print_warning "Rebuilding forum (this may take a minute)..."
    docker-compose down
    docker-compose build --no-cache
    docker-compose up -d
    print_success "Forum rebuilt and started!"
}

# Show logs
cmd_logs() {
    print_info "Showing logs (Ctrl+C to exit)..."
    docker-compose logs -f --tail=100
}

# Show status
cmd_status() {
    print_info "Container Status:"
    docker-compose ps
    echo ""
    print_info "Volume Status:"
    docker volume ls | grep forum_data || echo "Volume not found!"
    echo ""
    print_info "Health Status:"
    docker inspect --format='{{.State.Health.Status}}' running-forum 2>/dev/null || echo "Container not running"
}

# Backup database
cmd_backup() {
    BACKUP_DIR="./backups"
    BACKUP_FILE="$BACKUP_DIR/forum-backup-$(date +%Y%m%d-%H%M%S).db"
    
    mkdir -p "$BACKUP_DIR"
    
    print_info "Backing up database..."
    docker cp running-forum:/home/appuser/app/data/forum.db "$BACKUP_FILE"
    print_success "Database backed up to: $BACKUP_FILE"
}

# Restore database
cmd_restore() {
    BACKUP_DIR="./backups"
    
    # Find latest backup
    LATEST_BACKUP=$(ls -t "$BACKUP_DIR"/forum-backup-*.db 2>/dev/null | head -1)
    
    if [ -z "$LATEST_BACKUP" ]; then
        print_error "No backup found in $BACKUP_DIR"
        exit 1
    fi
    
    print_warning "This will restore from: $LATEST_BACKUP"
    read -p "Are you sure? (yes/no): " confirm
    
    if [ "$confirm" = "yes" ]; then
        print_info "Restoring database..."
        docker cp "$LATEST_BACKUP" running-forum:/home/appuser/app/data/forum.db
        docker-compose restart
        print_success "Database restored!"
    else
        print_info "Restore cancelled."
    fi
}

# Clean everything
cmd_clean() {
    print_warning "This will DELETE ALL DATA including the database!"
    read -p "Are you absolutely sure? Type 'DELETE' to confirm: " confirm
    
    if [ "$confirm" = "DELETE" ]; then
        print_info "Cleaning up..."
        docker-compose down -v
        docker system prune -af
        print_success "All Docker data removed!"
    else
        print_info "Clean cancelled."
    fi
}

# Open shell in container
cmd_shell() {
    print_info "Opening shell in container (type 'exit' to leave)..."
    docker exec -it running-forum sh
}

# Check health
cmd_health() {
    print_info "Checking health status..."
    
    if docker ps | grep -q running-forum; then
        HEALTH=$(docker inspect --format='{{.State.Health.Status}}' running-forum 2>/dev/null)
        
        if [ "$HEALTH" = "healthy" ]; then
            print_success "Container is HEALTHY ✓"
        elif [ "$HEALTH" = "unhealthy" ]; then
            print_error "Container is UNHEALTHY ✗"
            print_info "Check logs: ./docker-helper.sh logs"
        else
            print_warning "Health status: $HEALTH"
        fi
    else
        print_error "Container is not running!"
    fi
}

# Main script logic
main() {
    # Check if command provided
    if [ $# -eq 0 ]; then
        show_usage
        exit 0
    fi
    
    COMMAND=$1
    
    # Commands that don't need Docker running
    case $COMMAND in
        help)
            show_usage
            exit 0
            ;;
    esac
    
    # Check Docker for all other commands
    check_docker
    
    # Execute command
    case $COMMAND in
        start)
            cmd_start
            ;;
        stop)
            cmd_stop
            ;;
        restart)
            cmd_restart
            ;;
        rebuild)
            cmd_rebuild
            ;;
        logs)
            cmd_logs
            ;;
        status)
            cmd_status
            ;;
        backup)
            cmd_backup
            ;;
        restore)
            cmd_restore
            ;;
        clean)
            cmd_clean
            ;;
        shell)
            cmd_shell
            ;;
        health)
            cmd_health
            ;;
        *)
            print_error "Unknown command: $COMMAND"
            echo ""
            show_usage
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"

