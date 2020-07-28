package caches

type CacheOptions struct {
	UseCache bool `default:"false" name:"cache.activate" desc:"Activate the cache."`
	InfoJSON bool `default:"true" name:"cache.infojson" desc:"Cache info.json responses."`
	Tiles    bool `default:"true" name:"cache.tiles" desc:"Cache tiles."`
	Size     int  `default:"128" name:"cache.size" desc:"The size of the cache in MB."`
}

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, val []byte) error
	Has(key string) bool
}
