package caches

type CacheOptions struct {
	UseCache bool `default:"false" name:"cache.activate" desc:"Activate the cache."`

}

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, val []byte) error
	Has(key string) bool
}
