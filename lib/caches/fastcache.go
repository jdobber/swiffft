package caches

import (
	"errors"

	"github.com/VictoriaMetrics/fastcache"
)

type FastCache struct {
	Cache
	Client *fastcache.Cache
	Size   int
}

func NewFastCache(size int) (*FastCache, error) {

	c := FastCache{
		Client: fastcache.New(size * 1024 * 1024),
		Size:   size,
	}

	return &c, nil
}

func (c *FastCache) Get(key string) ([]byte, error) {

	var buf []byte
	buf = c.Client.GetBig(nil, []byte(key))
	if len(buf) > 0 {
		//log.Printf("Cache [HIT] %s", key)
		return buf, nil
	}
	//log.Printf("Cache [MISS] %s", key)
	return buf, errors.New("empty value from cache")

}

func (c *FastCache) Set(key string, val []byte) error {

	c.Client.SetBig([]byte(key), val)
	return nil
}

func (c *FastCache) Has(key string) bool {

	return c.Client.Has([]byte(key))

}
