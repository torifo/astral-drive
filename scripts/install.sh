#!/usr/bin/env bash
# install.sh - GitHub Releases から astral-drive をインストールする
# curl -fsSL https://raw.githubusercontent.com/torifo/astral-drive/main/scripts/install.sh | bash

set -euo pipefail

REPO="torifo/astral-drive"   # ← GitHub リポジトリ名に書き換える
BIN_DIR="${BIN_DIR:-/usr/local/bin}"
BIN_NAME="astral-drive"

# OS / アーキテクチャを判定
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  *)       echo "unsupported arch: $ARCH"; exit 1 ;;
esac

# 最新タグを取得
TAG="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' | sed 's/.*"tag_name": *"\(.*\)".*/\1/')"

ASSET="${BIN_NAME}-${OS}-${ARCH}"
URL="https://github.com/${REPO}/releases/download/${TAG}/${ASSET}"

echo "==> installing ${BIN_NAME} ${TAG} (${OS}/${ARCH})"
curl -fsSL "$URL" -o "/tmp/${BIN_NAME}"
chmod +x "/tmp/${BIN_NAME}"
sudo mv "/tmp/${BIN_NAME}" "${BIN_DIR}/${BIN_NAME}"

echo "==> installed to ${BIN_DIR}/${BIN_NAME}"
${BIN_NAME} -h
