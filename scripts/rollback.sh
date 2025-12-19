#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "用法: $0 vX.Y.Z"; exit 1;
fi
VER=$1
git fetch --all --tags
echo "回滚到 $VER（创建 hotfix 分支）..."
git checkout -b "hotfix/revert-to-$VER" "$VER"
echo "现在位于 hotfix/revert-to-$VER 分支。建议以 PR 方式回滚，避免强推 main。"

