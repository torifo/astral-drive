package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"astral-drive/internal/env"
	"astral-drive/internal/processor"
	"astral-drive/internal/scanner"
	"astral-drive/internal/ui"
)

func main() {
	topN := flag.Int("n", 20, "表示件数")
	workers := flag.Int("w", runtime.NumCPU()*2, "Worker 数")
	noColor := flag.Bool("no-color", false, "カラー出力を無効化")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: astral-drive [flags] [path]\n\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n  astral-drive\n  astral-drive -n 10 /home\n  astral-drive C:\\\n")
	}
	flag.Parse()

	// スキャン対象パスの決定
	var roots []string
	if flag.NArg() > 0 {
		root := flag.Arg(0)
		if _, err := os.Stat(root); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		roots = []string{root}
	} else {
		roots = env.Detect()
	}

	// スキャン実行
	cfg := scanner.Config{
		Workers: *workers,
	}
	result, err := scanner.Run(roots, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// 上位 N 件を抽出
	entries := processor.TopN(result.Sizes, *topN)

	// 表示
	uiCfg := ui.Config{
		NoColor: *noColor,
	}
	ui.Render(entries, result, uiCfg)
}
