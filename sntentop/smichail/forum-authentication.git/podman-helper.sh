#!/bin/bash
# ═══════════════════════════════════════════════════════════════════════════
# PODMAN MANAGEMENT HELPER SCRIPT FOR RUNNING forum-authentication
# Compatible ONLY with Podman + podman-compose
# ═══════════════════════════════════════════════════════════════════════════

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_info()    { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
print_error()   { echo -e "${RED}[ERROR]${NC} $1"; }

show_usage() {
    cat << EOF
════════════════════════════════════════════════════════════
Running forum-authentication (Podman Helper)
════════════════════════════════════════════════════════════

USAGE: ./podman-helper.sh [COMMAND]

COMMANDS:
  start        Start services (build if needed)
  stop         Stop services
  restart      Restart services
  rebuild      Rebuild images & restart
  logs         Show logs
  status       Show container status
  backup       Backup SQLite database
  restore      Restore latest backup
  clean        Remove everything (including data)
  shell        Open shell inside backend container
  health       Show container health
  help         Show this usage message

EOF
}

# Ensure Podman is installed
check_podman() {
    if ! podman info > /dev/null 2>&1; then
        print_error "Podman is not available! Install Podman first."
        exit 1
    fi
}

# Helper functions
cmd_start() {
    print_info "Starting Running forum-authentication..."
    podman-compose up -d --build
    print_success "Started!"
}

cmd_stop() {
    print_info "Stopping..."
    podman-compose down
    print_success "Stopped."
}

cmd_restart() {
    print_info "Restarting..."
    podman-compose restart
    print_success "Restarted!"
}

cmd_rebuild() {
    print_warning "Rebuilding images..."
    podman-compose down
    podman-compose build --no-cache
    podman-compose up -d
    print_success "Rebuilt and started!"
}

cmd_logs() {
    print_info "Logs (Ctrl+C to exit)..."
    podman-compose logs -f --tail=100
}

cmd_status() {
    print_info "Container Status:"
    podman-compose ps
}

cmd_backup() {
    BACKUP_DIR="./backups"
    mkdir -p "$BACKUP_DIR"
    BACKUP_FILE="$BACKUP_DIR/forum-authentication-backup-$(date +%Y%m%d-%H%M%S).db"

    print_info "Backing up SQLite DB..."
    podman cp running-forum-authentication:/home/appuser/app/data/forum-authentication.db "$BACKUP_FILE"
    print_success "Backup saved → $BACKUP_FILE"
}

cmd_restore() {
    BACKUP_DIR="./backups"
    LATEST=$(ls -t $BACKUP_DIR/*.db 2>/dev/null | head -1)

    if [ -z "$LATEST" ]; then
        print_error "No backup found."
        exit 1
    fi

    print_warning "Restoring: $LATEST"
    read -p "Type YES to restore: " confirm

    if [ "$confirm" = "YES" ]; then
        podman cp "$LATEST" running-forum-authentication:/home/appuser/app/data/forum-authentication.db
        podman-compose restart
        print_success "Restored."
    else
        print_info "Cancelled."
    fi
}

cmd_clean() {
    print_warning "This deletes EVERYTHING including DB!"
    read -p "Type DELETE to confirm: " confirm
    if [ "$confirm" = "DELETE" ]; then
        podman-compose down -v
        podman system prune -af
        print_success "Environment cleaned!"
    else
        print_info "Cancelled."
    fi
}

cmd_shell() {
    print_info "Opening shell..."
    podman exec -it running-forum-authentication sh
}

cmd_health() {
    print_info "Health status:"
    podman inspect --format '{{.State.Health.Status}}' running-forum-authentication || echo "Not running"
}

# Main
main() {
    if [ $# -eq 0 ]; then
        show_usage
        exit 0
    fi

    check_podman

    case $1 in
        start)   cmd_start ;;
        stop)    cmd_stop ;;
        restart) cmd_restart ;;
        rebuild) cmd_rebuild ;;
        logs)    cmd_logs ;;
        status)  cmd_status ;;
        backup)  cmd_backup ;;
        restore) cmd_restore ;;
        clean)   cmd_clean ;;
        shell)   cmd_shell ;;
        health)  cmd_health ;;
        help)    show_usage ;;
        *)
            print_error "Unknown command: $1"
            exit 1 ;;
    esac
}

main "$@"

