package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

// Config はスキャナーの実行オプション．
type Config struct {
	Workers int
	Exclude []string // スキップするディレクトリ名（例: ["node_modules", ".git"]）
}

// Result はスキャン完了後の集計結果．
type Result struct {
	Sizes     map[string]int64 // ディレクトリパス → 合計バイト数
	TotalDirs int64
	Elapsed   time.Duration
}

// sizeMap は mutex で保護された map．
type sizeMap struct {
	mu sync.Mutex
	m  map[string]int64
}

func (s *sizeMap) add(path string, size int64) {
	s.mu.Lock()
	s.m[path] += size
	s.mu.Unlock()
}

func (s *sizeMap) toMap() map[string]int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[string]int64, len(s.m))
	for k, v := range s.m {
		out[k] = v
	}
	return out
}

// isExcluded はディレクトリ名が除外リストに含まれるか判定する．
func isExcluded(name string, exclude []string) bool {
	for _, ex := range exclude {
		if name == ex {
			return true
		}
	}
	return false
}

// Run は roots を起点に Worker Pool でディレクトリを並列走査し，Result を返す．
func Run(roots []string, cfg Config) (*Result, error) {
	start := time.Now()

	sizes := &sizeMap{m: make(map[string]int64)}
	dirQueue := make(chan string, 1024)
	var wg sync.WaitGroup
	var totalDirs int64 // atomic で操作する

	// Worker を起動
	for i := 0; i < cfg.Workers; i++ {
		go func() {
			for path := range dirQueue {
				scanDir(path, dirQueue, sizes, &wg, &totalDirs, cfg.Exclude)
				wg.Done()
			}
		}()
	}

	// 起点パスをキューに投入
	for _, root := range roots {
		wg.Add(1)
		dirQueue <- root
	}

	// 全 Worker の完了を待機してチャンネルを閉じる
	wg.Wait()
	close(dirQueue)

	return &Result{
		Sizes:     sizes.toMap(),
		TotalDirs: atomic.LoadInt64(&totalDirs),
		Elapsed:   time.Since(start),
	}, nil
}

// scanDir は 1 ディレクトリを読み取り，サイズを集計してサブディレクトリをキューに追加する．
func scanDir(
	path string,
	dirQueue chan<- string,
	sizes *sizeMap,
	wg *sync.WaitGroup,
	totalDirs *int64,
	exclude []string,
) {
	atomic.AddInt64(totalDirs, 1)
	fmt.Fprintf(os.Stderr, "\r  scanning... %d dirs", atomic.LoadInt64(totalDirs))

	entries, err := os.ReadDir(path)
	if err != nil {
		// アクセス拒否等はスキップ
		return
	}

	var dirSize int64
	for _, entry := range entries {
		// シンボリックリンクは追跡しない
		if entry.Type()&fs.ModeSymlink != 0 {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if entry.IsDir() {
			if isExcluded(entry.Name(), exclude) {
				continue
			}
			subPath := filepath.Join(path, entry.Name())
			wg.Add(1)
			select {
			case dirQueue <- subPath:
			default:
				// チャンネルが満杯の場合はこの goroutine で同期処理
				scanDir(subPath, dirQueue, sizes, wg, totalDirs, exclude)
				wg.Done()
			}
		} else {
			dirSize += info.Size()
		}
	}

	sizes.add(path, dirSize)
}
