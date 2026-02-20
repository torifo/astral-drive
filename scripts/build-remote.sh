#!/usr/bin/env bash
# ============================================================
#  build-remote.sh — sync してから SSH でリモートビルドまで一発実行
#
#  実行場所: Windows の Git Bash (プロジェクトルートから)
#    ./scripts/build-remote.sh          # Linux バイナリのみ
#    ./scripts/build-remote.sh all      # 全プラットフォーム
# ============================================================
set -euo pipefail

REMOTE="${REMOTE:-toriforiumu@DELLXPS13}"
REMOTE_DIR="${REMOTE_DIR:-~/devops/CLI/astral-drive}"
TARGET="${1:-build}"   # build / build-win / build-all / test

# 1. ソースを同期
echo "==> [1/3] syncing source..."
bash "$(dirname "$0")/sync.sh"

# 2. リモートでビルド
echo "==> [2/3] building on ${REMOTE} (make ${TARGET})..."
ssh "${REMOTE}" "cd ${REMOTE_DIR} && make ${TARGET}"

# 3. 成果物を取得
echo "==> [3/3] fetching binaries..."
bash "$(dirname "$0")/fetch-bin.sh"

echo ""
echo "==> all done!"
