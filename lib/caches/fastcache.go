package caches

import "github.com/VictoriaMetrics/fastcache"

type FastCache struct {
	Cache
	Client *fastcache.Cache
}

func NewFastCache() (*FastCache, error) {

	c := FastCache{
		Client: fastcache.New(512 * 1024 * 1024),
	}

	return &c, nil
}

func (c *FastCache) Get(key string) ([]byte, error) {

	var buf []byte
	buf = c.Client.GetBig(nil, []byte(key))
	return buf, nil
}

func (c *FastCache) Set(key string, val []byte) error {

	c.Client.SetBig([]byte(key), val)
	return nil
}

func (c *FastCache) Has(key string) bool {

	return c.Client.Has([]byte(key))

}
