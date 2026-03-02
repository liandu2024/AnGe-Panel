#!/usr/bin/env bash
set -euo pipefail

REPO="liandu2024/AnGe-Panel"
APP_NAME="ange-panel"
INSTALL_DIR="/opt/${APP_NAME}"
SERVICE_FILE="/etc/systemd/system/${APP_NAME}.service"

need_cmd() { command -v "$1" >/dev/null 2>&1 || { echo "Missing required command: $1"; exit 1; }; }
need_cmd curl
need_cmd tar
need_cmd systemctl

ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64) ASSET_ARCH="amd64" ;;
  aarch64|arm64) ASSET_ARCH="arm64" ;;
  *) echo "Unsupported arch: $ARCH"; exit 1 ;;
 esac

API="https://api.github.com/repos/${REPO}/releases/latest"
echo "[1/6] Fetching latest release info from ${REPO}..."
JSON="$(curl -fsSL "$API")"
TAG="$(echo "$JSON" | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]\+\)".*/\1/p' | head -n1)"
if [[ -z "${TAG}" ]]; then
  echo "Failed to detect tag_name from GitHub API."
  exit 1
fi

ASSET_NAME="${APP_NAME}_${TAG#v}_linux_${ASSET_ARCH}.tar.gz"
DOWNLOAD_URL="$(echo "$JSON" | grep -E '"browser_download_url"' | grep -F "${ASSET_NAME}" | head -n1 | sed -n 's/.*"browser_download_url"[[:space:]]*:[[:space:]]*"\([^"]\+\)".*/\1/p')"

if [[ -z "${DOWNLOAD_URL}" ]]; then
  echo "Could not find asset ${ASSET_NAME} in latest release (${TAG})."
  echo "Make sure the release workflow uploaded it."
  exit 1
fi

echo "[2/6] Downloading ${ASSET_NAME}..."
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT
curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/${ASSET_NAME}"

echo "[3/6] Installing to ${INSTALL_DIR}..."
sudo mkdir -p "$INSTALL_DIR"
sudo tar -xzf "$TMP_DIR/${ASSET_NAME}" -C "$INSTALL_DIR" --strip-components=1
sudo chmod +x "$INSTALL_DIR/${APP_NAME}" || true
# Place web assets to /opt/<app>/web so backend can serve it
if [[ -d "$INSTALL_DIR/web" ]]; then
  sudo chown -R "$APP_NAME":"$APP_NAME" "$INSTALL_DIR/web" || true
fi

# system user
if ! id -u "$APP_NAME" >/dev/null 2>&1; then
  echo "[4/6] Creating user ${APP_NAME}..."
  sudo useradd --system --no-create-home --shell /usr/sbin/nologin "$APP_NAME" || true
fi
sudo chown -R "$APP_NAME":"$APP_NAME" "$INSTALL_DIR"

# install systemd unit
echo "[5/6] Installing systemd service..."
sudo tee "$SERVICE_FILE" >/dev/null <<EOF
[Unit]
Description=AnGe-Panel
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=${APP_NAME}
Group=${APP_NAME}
WorkingDirectory=${INSTALL_DIR}
ExecStart=${INSTALL_DIR}/${APP_NAME}
Restart=always
RestartSec=2
LimitNOFILE=1048576

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable --now "${APP_NAME}.service"

echo "[6/6] Done. Status:"
systemctl --no-pager --full status "${APP_NAME}.service" || true

echo "\nOpen: http://localhost:3002"
