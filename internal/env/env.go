package env

import (
	"os"
	"runtime"
)

// Detect は実行 OS に応じたスキャン起点パスの一覧を返す．
// Windows の場合はマウント済みドライブ（C:\, D:\ 等），
// それ以外は "/" を返す．
func Detect() []string {
	if runtime.GOOS == "windows" {
		return detectWindows()
	}
	return []string{"/"}
}

func detectWindows() []string {
	var drives []string
	for c := 'A'; c <= 'Z'; c++ {
		path := string(c) + `:\`
		if _, err := os.Stat(path); err == nil {
			drives = append(drives, path)
		}
	}
	return drives
}
