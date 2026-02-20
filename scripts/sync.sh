#!/usr/bin/env bash
# ============================================================
#  sync.sh — Windows 側のソースを Linux ビルドマシンへ転送する
#
#  実行場所: Windows の Git Bash (プロジェクトルートから)
#    ./scripts/sync.sh
#
#  転送先を変えたい場合:
#    REMOTE=user@host ./scripts/sync.sh
# ============================================================
set -euo pipefail

REMOTE="${REMOTE:-toriforiumu@DELLXPS13}"
REMOTE_DIR="${REMOTE_DIR:-~/devops/CLI/astral-drive}"

echo "==> syncing to ${REMOTE}:${REMOTE_DIR}"

rsync -avz --progress \
  --exclude='.git/' \
  --exclude='*.exe' \
  --exclude='astral-drive' \
  --exclude='astral-drive-mac' \
  ./ "${REMOTE}:${REMOTE_DIR}/"

echo "==> sync done"
echo "    Next: ssh ${REMOTE} 'cd ${REMOTE_DIR} && make build'"
