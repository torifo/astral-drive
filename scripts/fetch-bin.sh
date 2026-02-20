#!/usr/bin/env bash
# ============================================================
#  fetch-bin.sh — Linux でビルドした成果物を Windows 側へ取得する
#
#  実行場所: Windows の Git Bash (プロジェクトルートから)
#    ./scripts/fetch-bin.sh
# ============================================================
set -euo pipefail

REMOTE="${REMOTE:-toriforiumu@DELLXPS13}"
REMOTE_DIR="${REMOTE_DIR:-~/devops/CLI/astral-drive}"
OUT_DIR="${OUT_DIR:-./dist}"

mkdir -p "$OUT_DIR"

echo "==> fetching binaries from ${REMOTE}:${REMOTE_DIR}"

# Linux バイナリ
scp "${REMOTE}:${REMOTE_DIR}/astral-drive"     "${OUT_DIR}/astral-drive-linux" 2>/dev/null && \
  echo "  [ok] astral-drive-linux" || echo "  [skip] astral-drive-linux (not built)"

# Windows バイナリ (Linux 側でクロスコンパイルしたもの)
scp "${REMOTE}:${REMOTE_DIR}/astral-drive.exe" "${OUT_DIR}/astral-drive.exe"   2>/dev/null && \
  echo "  [ok] astral-drive.exe"   || echo "  [skip] astral-drive.exe (not built)"

echo "==> done. binaries in ${OUT_DIR}/"
