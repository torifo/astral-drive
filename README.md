# astral-drive

ストレージ内の巨大ディレクトリを高速に検出する CLI ツール．
Go の Worker Pool パターンによる並列スキャンで，数万ディレクトリを数秒で分析できる．

## インストール

[GitHub Releases](https://github.com/torifo/astral-drive/releases/latest) から各 OS 向けのバイナリをダウンロードしてください．
詳細な手順は [INSTALL.md](./INSTALL.md) を参照してください．

```bash
# Linux / macOS（ワンライナー）
curl -fsSL https://raw.githubusercontent.com/torifo/astral-drive/main/scripts/install.sh | bash

# ソースからビルド
git clone https://github.com/torifo/astral-drive
cd astral-drive && go build -o astral-drive .
```

## 使い方

```
astral-drive [flags] [path]

Flags:
  -n int     表示件数 (デフォルト: 20)
  -w int     Worker 数 (デフォルト: CPU コア数 × 2)
  -no-color  カラー出力を無効化
  -h         ヘルプ表示
```

### 例

```bash
# OS に応じた全ドライブを自動スキャン
astral-drive

# /home 配下の上位 10 件を表示
astral-drive -n 10 /home

# Windows で C: ドライブをスキャン
astral-drive C:\

# Worker 数を固定してスキャン
astral-drive -w 4 /var
```

### 出力例

```
   #    SIZE        PATH
  ---  ----------  --------------------------------------------------
    1   186.8 MB  C:\Users\user\logo_list
    2   142.0 MB  C:\project\node_modules\@next\swc-linux-x64-gnu
    3    12.3 MB  C:\Users\user\Downloads

Scanned 21,847 dirs in 0.8s
```

## 開発フロー（Windows 編集 → WSL2 ビルド）

ソースは Windows で編集し，ビルドは WSL2 Ubuntu で行う構成を前提としている．

### ソースを WSL2 へ同期

```powershell
# Windows PowerShell から実行
powershell -ExecutionPolicy Bypass -File scripts\sync-wsl.ps1
```

### WSL2 でビルド・実行

```bash
# WSL2 端末内で実行
cd ~/devops/CLI/astral-drive

make build      # Linux バイナリ
make test       # ユニットテスト
make run ARGS="-n 10 /home"   # 引数付きで実行
make build-win  # Windows バイナリをクロスコンパイル
```

## Makefile ターゲット（WSL2 側）

```bash
make build      # Linux バイナリ
make build-win  # Windows バイナリ (クロスコンパイル)
make build-mac  # macOS バイナリ (クロスコンパイル)
make build-all  # 全プラットフォーム
make test       # ユニットテスト
make run ARGS="-n 10 /home"  # 引数付きで実行
make clean      # 成果物を削除
```

## アーキテクチャ

`spec/design.md` を参照．
