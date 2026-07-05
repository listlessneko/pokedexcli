package pokecache

import (
	"testing"
	"time"
	"bytes"
)

func TestAddAndGet(t *testing.T) {
	cache := NewCache(5 * time.Second)
	url := "https://example.com"
	og_b := []byte("test data")
	cache.Add(url, og_b)

	new_b, ok := cache.Get(url)
	if !ok {
		t.Errorf("expected cache hit, received miss")
	}

	if !bytes.Equal(new_b, og_b) {
		t.Errorf("expected %v, received %v", og_b, new_b)
	}
}
