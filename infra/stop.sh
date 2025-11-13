#!/usr/bin/env bash
set -euo pipefail
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
# 仅停止并移除容器与网络，保留数据卷
"$SCRIPT_DIR/compose.sh" down --remove-orphans
