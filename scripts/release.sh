#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "用法: $0 vX.Y.Z"; exit 1;
fi
VER=$1
git fetch --all --tags
echo "Tagging $VER on main HEAD..."
git checkout main
git pull --ff-only
git tag -a "$VER" -m "Release $VER"
git push origin "$VER"
echo "✅ 已创建 tag $VER 并推送到远端"

