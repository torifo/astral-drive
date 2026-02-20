package ui

import (
	"fmt"
	"os"

	"astral-drive/internal/processor"
	"astral-drive/internal/scanner"
)

// Config は UI の表示オプション．
type Config struct {
	NoColor bool
}

const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
)

// Render はスキャン結果を整形してターミナルに出力する．
func Render(entries []processor.DirEntry, result *scanner.Result, cfg Config) {
	// 進捗表示をクリア
	fmt.Fprintf(os.Stderr, "\r\033[K")

	if len(entries) == 0 {
		fmt.Println("No directories found.")
		return
	}

	color := func(c, s string) string {
		if cfg.NoColor {
			return s
		}
		return c + s + colorReset
	}

	// ヘッダー
	fmt.Printf("  %s  %-10s  %s\n",
		color(colorCyan, " # "),
		color(colorCyan, "SIZE      "),
		color(colorCyan, "PATH"),
	)
	fmt.Printf("  %s  %s  %s\n", "---", "----------", "--------------------------------------------------")

	// エントリー
	for i, e := range entries {
		sizeStr := formatSize(e.Size)
		fmt.Printf("  %3d  %-10s  %s\n",
			i+1,
			color(colorYellow, fmt.Sprintf("%10s", sizeStr)),
			e.Path,
		)
	}

	// サマリー
	fmt.Printf("\n%s\n",
		color(colorGreen, fmt.Sprintf("Scanned %s dirs in %.1fs",
			formatCount(result.TotalDirs),
			result.Elapsed.Seconds(),
		)),
	)
}

// formatSize はバイト数を人間が読みやすい単位文字列に変換する．
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"KB", "MB", "GB", "TB", "PB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// formatCount はディレクトリ数をカンマ区切りの文字列に変換する．
func formatCount(n int64) string {
	s := fmt.Sprintf("%d", n)
	result := []byte{}
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}
