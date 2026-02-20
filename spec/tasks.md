# タスクリスト

## フェーズ 1: 基盤構築

- [x] リポジトリ初期化と `go.mod` の作成 (`module astral-drive`, Go 1.21+)
- [x] ディレクトリ構造の作成 (`main.go`, `internal/env/`, `internal/scanner/`, `internal/processor/`, `internal/ui/`)
- [x] `main.go` の雛形作成（`flag` パッケージによる `-n`, `-w`, `-no-color` フラグ解析）

## フェーズ 2: コア実装

- [x] `internal/env/env.go`: `runtime.GOOS` による OS 判定と，初期ターゲット（Windows ならドライブ一覧，その他は `/`）の取得実装
- [x] `internal/scanner/scanner.go`: シングルスレッドでの再帰スキャン（ベースライン）の実装
- [x] `internal/scanner/scanner.go`: Worker Pool と Channel を用いた並列スキャンへの拡張
- [x] `internal/scanner/scanner.go`: エラーハンドリング（アクセス拒否フォルダのスキップ，シンボリックリンク非追跡）の追加
- [x] `internal/processor/processor.go`: サイズ降順ソートと上位 N 件抽出の実装

## フェーズ 3: 表示・進捗

- [x] `internal/ui/ui.go`: カラム整形，単位変換（B/KB/MB/GB/TB），ANSI カラー出力の実装
- [x] 進捗表示（スキャン中のディレクトリ数を標準エラー出力に `\r` で逐次表示）の実装
- [x] 完了後サマリー（経過時間，スキャンディレクトリ総数）の表示実装

## フェーズ 4: 品質確保

- [x] `internal/processor` のユニットテスト作成（ソート正確性，件数上限の境界値）
- [x] `internal/ui` のユニットテスト作成（単位変換の正確性）
- [x] Windows / Linux のクロスコンパイル確認 (`GOOS=linux go build ./...`, `GOOS=windows go build ./...`)
- [x] 出力フォーマットの最終確認と `README.md` の整備
