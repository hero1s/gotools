package lru

import (
	"github.com/hero1s/gotools/log"
	"testing"
)

func TestNewCache(t *testing.T) {
	cache := NewCache(2)

	cache.Put("1", "one")
	log.Info("%v", cache.Get("1"))

	log.Warning("===============")
	cache.Put("2", "two")
	log.Info("%v", cache.Get("1"))

	log.Warning("===============")
	cache.Put("3", "three")
	log.Info("%v", cache.Get("2"))
	log.Info("%v", cache.Get("3"))
	log.Info("%v", cache.Get("3"))
	log.Info("%v", cache.Get("1"))

	log.Warning("===============")
	cache.Put("2", "two")
	log.Info("%v", cache.Get("3"))
	log.Info("%v", cache.Get("1"))
}
