package sources

type Source interface {
	Read(key string) ([]byte, error)
}
