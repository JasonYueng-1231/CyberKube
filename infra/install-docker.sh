#!/usr/bin/env bash
set -euo pipefail

ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64) ARCH=x86_64 ;;
  aarch64|arm64) ARCH=aarch64 ;;
  *) echo "[install] 不支持的架构: $ARCH" >&2; exit 1 ;;
esac

fetch() { # url dest
  if command -v curl >/dev/null 2>&1; then curl -fSL "$1" -o "$2"; else wget -qO "$2" "$1"; fi
}

echo "[install] 停止旧 docker（如存在）"
if command -v systemctl >/dev/null 2>&1; then systemctl stop docker || true; fi
pkill dockerd 2>/dev/null || true

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

install_docker_static() {
  local VER=$1
  local URL="https://download.docker.com/linux/static/stable/${ARCH}/docker-${VER}.tgz"
  echo "[install] 下载 Docker 静态二进制: $VER"
  fetch "$URL" "$TMPDIR/docker.tgz"
  echo "[install] 解压安装"
  tar -xzf "$TMPDIR/docker.tgz" -C "$TMPDIR"
  install -m 0755 "$TMPDIR"/docker/* /usr/local/bin/
}

# 1) 安装 Docker 静态二进制，带回退版本
if ! install_docker_static 24.0.9 2>/dev/null; then
  echo "[install] 回退到 23.0.6"
  install_docker_static 23.0.6
fi

# 2) 安装 Compose v2 插件
echo "[install] 安装 Docker Compose v2 插件"
mkdir -p /usr/local/lib/docker/cli-plugins
COMPOSE_VER=v2.24.5
fetch "https://github.com/docker/compose/releases/download/${COMPOSE_VER}/docker-compose-linux-${ARCH}" \
  /usr/local/lib/docker/cli-plugins/docker-compose
chmod +x /usr/local/lib/docker/cli-plugins/docker-compose

# 3) 配置 systemd 服务（如可用），否则后台拉起 dockerd
if command -v systemctl >/dev/null 2>&1; then
  echo "[install] 配置并启动 systemd 服务"
  cat >/etc/systemd/system/docker.service <<'UNIT'
[Unit]
Description=Docker Engine
After=network-online.target
Wants=network-online.target
[Service]
Type=notify
ExecStart=/usr/local/bin/dockerd
ExecReload=/bin/kill -s HUP $MAINPID
Restart=always
LimitNOFILE=1048576
LimitNPROC=1048576
[Install]
WantedBy=multi-user.target
UNIT
  systemctl daemon-reload
  systemctl enable --now docker
else
  echo "[install] 无 systemd，后台启动 dockerd"
  nohup /usr/local/bin/dockerd >/var/log/dockerd.log 2>&1 & disown || true
  sleep 3
fi

echo "[install] 版本校验"
/usr/local/bin/docker --version
/usr/local/bin/docker compose version || true

echo "[install] 完成。可在项目根目录执行: bash infra/start.sh"

