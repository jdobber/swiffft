package caches

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, val []byte) error
	Has(key string) bool
}
