package syncmap_test

import (
	"testing"

	"github.com/sapslaj/morbius/syncmap"
)

func TestSyncMapWorks(t *testing.T) {
	m, err := syncmap.NewMap[string, string]()
	if err != nil {
		t.Fatalf("syncmap.NewMap returned error: %v", err)
	}

	_, ok := m.Load("test1")
	if ok {
		t.Error("Load on undefined key returned ok")
	}

	m.Store("test2", "test2")
	value, ok := m.Load("test2")
	if !ok {
		t.Error("Load on defined key did not return ok")
	}
	if value != "test2" {
		t.Error("Load did not return expected value")
	}

	m.Range(func(key, value string) bool {
		if key != "test2" || value != "test2" {
			t.Error("Range did not return expected results")
		}
		return true
	})
}
