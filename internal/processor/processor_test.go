package processor

import (
	"testing"
)

func TestTopN_Sort(t *testing.T) {
	sizes := map[string]int64{
		"/a": 100,
		"/b": 300,
		"/c": 200,
	}
	entries := TopN(sizes, 10)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Size != 300 || entries[0].Path != "/b" {
		t.Errorf("expected first entry /b(300), got %s(%d)", entries[0].Path, entries[0].Size)
	}
	if entries[1].Size != 200 {
		t.Errorf("expected second entry size 200, got %d", entries[1].Size)
	}
	if entries[2].Size != 100 {
		t.Errorf("expected third entry size 100, got %d", entries[2].Size)
	}
}

func TestTopN_Limit(t *testing.T) {
	sizes := map[string]int64{
		"/a": 100,
		"/b": 300,
		"/c": 200,
		"/d": 400,
		"/e": 50,
	}
	entries := TopN(sizes, 3)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	// 上位 3 件は /d(400), /b(300), /c(200)
	if entries[0].Size != 400 {
		t.Errorf("expected top size 400, got %d", entries[0].Size)
	}
	if entries[2].Size != 200 {
		t.Errorf("expected 3rd size 200, got %d", entries[2].Size)
	}
}

func TestTopN_Empty(t *testing.T) {
	entries := TopN(map[string]int64{}, 10)
	if len(entries) != 0 {
		t.Errorf("expected empty, got %d", len(entries))
	}
}

func TestTopN_NGreaterThanLen(t *testing.T) {
	sizes := map[string]int64{"/a": 1, "/b": 2}
	entries := TopN(sizes, 100)
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}
