package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Cache struct {
	value   sync.Map
	expires int64
}

func (c *Cache) Expired(time int64) bool {
	if c.expires == 0 {
		return false
	}
	return time > c.expires
}

func (c *Cache) Get(key string) string {
	if c.Expired(time.Now().UnixNano()) {
		log.Printf("%s has expired", key)
		return ""
	}

	v, ok := c.value.Load(key)
	var s string
	if ok {
		s, ok = v.(string)
		if !ok {
			log.Printf("%s does not exists", key)
			return ""
		}
	}
	return s
}

func (c *Cache) Put(key string, value string, expired int64) {
	c.value.Store(key, value)
	c.expires = expired
}

var cache = &Cache{}

func main() {
	fk := "first-key"
	sk := "second-key"

	cache.Put(fk, "first-value", time.Now().Add(2*time.Second).UnixNano())
	s := cache.Get(fk)
	fmt.Println(cache.Get(fk))

	time.Sleep(5 * time.Second)

	s = cache.Get(fk)
	if len(s) == 0 {
		cache.Put(sk, "second-value", time.Now().Add(100*time.Second).UnixNano())
	}
	fmt.Println(cache.Get(sk))
}
