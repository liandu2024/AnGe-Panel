#!/usr/bin/env sh
set -eu

APP_DIR="/app"
DATA_DIR="${ANGE_DATA_DIR:-/data}"

# Create persistent dirs
mkdir -p "${DATA_DIR}/conf" "${DATA_DIR}/database" "${DATA_DIR}/uploads" "${DATA_DIR}/runtime/temp"

# If this is a fresh volume, initialize config/db from image defaults/seed
# 1) conf
if [ -z "$(ls -A "${DATA_DIR}/conf" 2>/dev/null || true)" ]; then
  if [ -d "${APP_DIR}/conf" ]; then
    cp -a "${APP_DIR}/conf/." "${DATA_DIR}/conf/" || true
  fi
fi

# 2) database
if [ ! -f "${DATA_DIR}/database/database.db" ]; then
  if [ -f "${APP_DIR}/seed/database/database.db" ]; then
    cp -a "${APP_DIR}/seed/database/database.db" "${DATA_DIR}/database/database.db" || true
  fi
fi

# 3) uploads (only if empty)
if [ -z "$(ls -A "${DATA_DIR}/uploads" 2>/dev/null || true)" ] || { [ "$(ls -A "${DATA_DIR}/uploads" | wc -l | tr -d ' ')" = "1" ] && [ -f "${DATA_DIR}/uploads/.gitkeep" ]; }; then
  if [ -d "${APP_DIR}/seed/uploads" ]; then
    cp -a "${APP_DIR}/seed/uploads/." "${DATA_DIR}/uploads/" || true
  fi
fi

# Link runtime paths expected by the app (relative ./conf ./database ./uploads ./runtime)
link_dir() {
  src="$1"; dst="$2"
  if [ -e "${src}" ] && [ ! -L "${src}" ]; then
    # leave existing real dir in place
    return
  fi
  rm -f "${src}" 2>/dev/null || true
  ln -s "${dst}" "${src}"
}

link_dir "${APP_DIR}/conf" "${DATA_DIR}/conf"
link_dir "${APP_DIR}/database" "${DATA_DIR}/database"
link_dir "${APP_DIR}/uploads" "${DATA_DIR}/uploads"
link_dir "${APP_DIR}/runtime" "${DATA_DIR}/runtime"

exec "$@"
