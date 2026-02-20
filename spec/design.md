# 設計書

## アーキテクチャ概要

Worker Pool Pattern を採用し，複数の Worker が Channel からディレクトリパスを受け取り，`os.ReadDir` を実行する並行構造とする．
メインゴルーチンがルートパスをキューに投入し，Worker が再帰的に処理する BFS（幅優先）方式でファイルシステムを走査する．

```
[main] --> [Environment] --> [Scanner (Worker Pool)] --> [Processor] --> [UI]
                                      |
                          +-----------+-----------+
                          |           |           |
                       Worker1    Worker2    Worker3
                          |           |           |
                          +-----------+-----------+
                                      |
                               dirQueue (chan string)
```

## 主要モジュール

### Environment (`internal/env/`)

- `runtime.GOOS` を用いた実行環境（Windows / Linux / macOS）の特定．
- Windows の場合: 固定ドライブレター（A〜Z）に対して `os.Stat` でマウント済みドライブを列挙する．
- Linux / macOS の場合: ルート（`/`）を起点として利用する．
- 戻り値: `[]string` 形式のスキャン起点パス一覧．

### Scanner (`internal/scanner/`)

- `sync/atomic` による並行安全なサイズ集計．
- `sync.WaitGroup` で全 Worker の完了を待機．
- アクセス拒否（`os.ErrPermission`）などのエラーを無視してスキップ．
- シンボリックリンクは `DirEntry.Type().IsDir()` で判定し，追跡しない．
- 各ディレクトリの合計サイズ（バイト）を `sync.Map` で並行安全に保持．

**Worker Pool 構成:**

| パラメータ         | 値                         |
|--------------------|----------------------------|
| デフォルト Worker 数 | `runtime.NumCPU() * 2`   |
| Channel バッファサイズ | `1024`                  |
| Channel の型       | `chan string`（ディレクトリパス） |

### Processor (`internal/processor/`)

- `sync.Map` から `[]DirEntry` スライスへ変換．
- `sort.Slice` でサイズ降順にソート．
- 上位 N 件を抽出（デフォルト N=20）．
- 戻り値: `[]DirEntry{Path string, Size int64}`．

### UI (`internal/ui/`)

- ターミナル向けカラム整形表示（左寄せパス，右寄せサイズ）．
- サイズ単位の自動変換: バイト → KB / MB / GB / TB（1024 基数）．
- 進捗表示: スキャン中のディレクトリ数を標準エラー出力に逐次表示（`\r` で上書き）．
- カラー出力: `NO_COLOR` 環境変数が未設定かつ TTY 接続時に ANSI エスケープコードを使用．

## データフロー

```
起動
 └─ Environment.Detect()
      └─ []string (ドライブ/ルートパス)
           └─ Scanner.Run(paths, workerN)
                 ├─ dirQueue (chan string) に起点パスを投入
                 │    ├─ Worker[0]: ReadDir → subDirs → dirQueue, sizes[path] += fileSize
                 │    ├─ Worker[1]: 同上
                 │    └─ Worker[N]: 同上
                 └─ map[string]int64 (path → totalSize)
                      └─ Processor.TopN(sizes, N)
                           └─ []DirEntry (ソート済み上位 N 件)
                                └─ UI.Render(entries)
                                     └─ 標準出力へ表示
```

## ディレクトリ構造

```
astral-drive/
├── main.go                  # エントリーポイント，CLI フラグ解析
├── internal/
│   ├── env/
│   │   └── env.go           # OS 判定，ターゲットパス列挙
│   ├── scanner/
│   │   ├── scanner.go       # Worker Pool，ディレクトリ走査
│   │   └── scanner_test.go
│   ├── processor/
│   │   ├── processor.go     # ソート，上位 N 件抽出
│   │   └── processor_test.go
│   └── ui/
│       ├── ui.go            # 表示フォーマット，単位変換
│       └── ui_test.go
├── spec/
│   ├── design.md
│   ├── requirements.md
│   └── tasks.md
├── go.mod
└── README.md
```

## CLI インターフェース設計

```
astral-drive [flags] [path]

Flags:
  -n int     表示件数 (デフォルト: 20)
  -w int     Worker 数 (デフォルト: CPU コア数 × 2)
  -no-color  カラー出力を無効化
  -h         ヘルプ表示

使用例:
  astral-drive              # OS に応じた全ドライブを自動スキャン
  astral-drive -n 10 /home  # /home 配下の上位 10 件を表示
  astral-drive C:\          # Windows で C: ドライブをスキャン
  astral-drive -w 4 /var    # Worker 数を 4 に固定してスキャン
```

**出力例:**

```
 #   SIZE      PATH
---  --------  --------------------------------------------------
  1   45.2 GB  /home/user/.local/share/containers
  2   12.8 GB  /home/user/Downloads
  3    8.3 GB  /var/cache/apt
  ...

Scanned 1,204,832 dirs in 18.4s
```

## 主要型定義

```go
// DirEntry はソート・表示対象のディレクトリエントリ．
type DirEntry struct {
    Path string
    Size int64 // バイト単位
}

// ScanResult はスキャン結果の集計．
type ScanResult struct {
    Entries    []DirEntry
    TotalDirs  int64
    Elapsed    time.Duration
}

// Config は実行時オプションを保持する．
type Config struct {
    TopN     int
    Workers  int
    NoColor  bool
    RootPath string
}
```
