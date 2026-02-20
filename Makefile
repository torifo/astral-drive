# ============================================================
#  astral-drive Makefile  (Linux 側で実行する)
# ============================================================
BINARY   := astral-drive
MODULE   := astral-drive
LDFLAGS  := -s -w

.PHONY: all build build-win build-linux run test clean

## デフォルト: Linux バイナリをビルド
all: build

## Linux 向けバイナリ
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	  go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

## Windows 向けバイナリ (クロスコンパイル)
build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 \
	  go build -ldflags "$(LDFLAGS)" -o $(BINARY).exe .

## macOS (Apple Silicon) 向け
build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 \
	  go build -ldflags "$(LDFLAGS)" -o $(BINARY)-mac .

## 全プラットフォーム向けを一括ビルド
build-all: build build-win build-mac

## ローカル実行 (Linux)
run:
	./$(BINARY) $(ARGS)

## テスト
test:
	go test ./internal/... -v

## 成果物の削除
clean:
	rm -f $(BINARY) $(BINARY).exe $(BINARY)-mac
