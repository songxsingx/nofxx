#!/bin/bash

# pm2.sh - PM2 process manager for NexTrade Trading Bot
# This script provides a user-friendly interface to manage the trading bot with PM2.

# Get the project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_ROOT"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function: Print colored messages
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_header() {
    echo -e "${PURPLE}═══════════════════════════════════════${NC}"
    echo -e "${PURPLE}  🤖 NexTrade Trading Bot - PM2 Manager${NC}"
    echo -e "${PURPLE}═══════════════════════════════════════${NC}"
    echo ""
}

# Function: Check if PM2 is installed
check_pm2() {
    if ! command -v pm2 &> /dev/null; then
        print_error "PM2 未安装，请先安装: npm install -g pm2"
        exit 1
    fi
}

# Function: Ensure log directories exist
ensure_log_dirs() {
    mkdir -p "$PROJECT_ROOT/logs"
    mkdir -p "$PROJECT_ROOT/web/logs"
    print_info "日志目录已创建"
}

# Function: Build backend
build_backend() {
    print_info "正在编译后端..."
    go build -o nextrade
    if [ $? -eq 0 ]; then
        print_success "后端编译完成"
    else
        print_error "后端编译失败"
        exit 1
    fi
}

# Function: Start services
start() {
    print_header
    check_pm2
    ensure_log_dirs
    build_backend
    
    print_info "正在启动服务..."
    
    # Start backend
    pm2 start pm2.config.js --only nextrade-backend
    
    # Wait a moment for backend to start
    sleep 2
    
    # Start frontend
    pm2 start pm2.config.js --only nextrade-frontend
    
    # Save PM2 config
    pm2 save
    
    print_success "服务已启动！"
    echo ""
    print_info "应用查看:"
    print_info "  后端 API: http://localhost:8080"
    print_info "  前端界面: http://localhost:3000"
    echo ""
    print_info "管理命令:"
    print_info "  查看状态: ./pm2.sh status"
    print_info "  查看日志: ./pm2.sh logs"
    print_info "  停止服务: ./pm2.sh stop"
}

# Function: Stop services
stop() {
    print_header
    check_pm2
    
    print_info "正在停止服务..."
    pm2 stop pm2.config.js
    print_success "服务已停止"
}

# Function: Restart services
restart() {
    print_header
    check_pm2
    
    print_info "正在重启服务..."
    pm2 restart pm2.config.js
    print_success "服务已重启"
}

# Function: Show status
status() {
    print_header
    check_pm2
    
    print_info "服务状态:"
    pm2 list
}

# Function: Show logs
logs() {
    check_pm2
    
    if [ -z "$2" ]; then
        print_info "显示所有服务日志 (按 Ctrl+C 退出)..."
        pm2 logs
    else
        case "$2" in
            backend)
                print_info "显示后端日志 (按 Ctrl+C 退出)..."
                pm2 logs nextrade-backend
                ;;
            frontend)
                print_info "显示前端日志 (按 Ctrl+C 退出)..."
                pm2 logs nextrade-frontend
                ;;
            *)
                print_error "未知服务: $2 (支持: backend, frontend)"
                ;;
        esac
    fi
}

# Function: Show help
show_help() {
    print_header
    echo "NexTrade Trading Bot - PM2 管理脚本"
    echo ""
    echo "用法: ./pm2.sh [command]"
    echo ""
    echo "命令:"
    echo "  start     启动服务"
    echo "  stop      停止服务"
    echo "  restart   重启服务"
    echo "  status    查看服务状态"
    echo "  logs      查看日志 (可选: backend/frontend)"
    echo "  help      显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  ./pm2.sh start      # 启动所有服务"
    echo "  ./pm2.sh logs backend  # 查看后端日志"
}

# Main: Command dispatcher
main() {
    case "$1" in
        start)
            start
            ;;
        stop)
            stop
            ;;
        restart)
            restart
            ;;
        status)
            status
            ;;
        logs)
            logs "$1" "$2"
            ;;
        help|"")
            show_help
            ;;
        *)
            print_error "未知命令: $1"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"