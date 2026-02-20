package processor

import (
	"sort"
)

// DirEntry はソート・表示対象のディレクトリエントリ．
type DirEntry struct {
	Path string
	Size int64 // バイト単位
}

// TopN は sizes マップからサイズ上位 n 件を降順で返す．
func TopN(sizes map[string]int64, n int) []DirEntry {
	entries := make([]DirEntry, 0, len(sizes))
	for path, size := range sizes {
		entries = append(entries, DirEntry{Path: path, Size: size})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Size > entries[j].Size
	})

	if n > 0 && n < len(entries) {
		entries = entries[:n]
	}
	return entries
}
