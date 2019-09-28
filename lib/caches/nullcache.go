package caches

import (
	"errors"
)

type NullCache struct {
	Cache
}

func NewNullCache() (*NullCache, error) {

	c := NullCache{}
	return &c, nil
}

func (c *NullCache) Get(key string) ([]byte, error) {

	return nil, errors.New("Method not implemented.")
}

func (c *NullCache) Set(key string, val []byte) error {

	return nil
}

func (c *NullCache) Has(key string) bool {

	return false

}
