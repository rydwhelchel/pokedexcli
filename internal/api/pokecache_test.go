package api

import (
	"fmt"
	"testing"
	"time"
)

// need to try unit testing Go code a bit more
func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

// Tests that accessing an entry updates the time created
// and thus staves off the reaper
func TestUpdateCacheEntry(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = 3 * time.Millisecond
	const key = "dont fear"
	const value = "the reaper"
	cache := NewCache(baseTime)
	cache.Add(key, []byte(value))

	time.Sleep(waitTime)
	cache.Update(key)
	time.Sleep(waitTime)
	_, ok := cache.Get(key)
	if ok == false {
		t.Error("The entry should still be alive in the map")
	}

	time.Sleep(waitTime + (2 * time.Millisecond))
	_, ok = cache.Get(key)
	if ok {
		t.Error("The entry should have been reaped")
	}
}
