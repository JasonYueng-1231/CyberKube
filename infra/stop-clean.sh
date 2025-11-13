#!/usr/bin/env bash
set -euo pipefail
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
# 停止并清理容器、网络与数据卷（包含 MySQL/Redis 数据）
"$SCRIPT_DIR/compose.sh" down -v --remove-orphans

