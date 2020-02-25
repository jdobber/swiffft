package sources

import (
	"errors"
)

type SourceOptions struct {
	Sources []string `name:"sources" default:"file" desc:"An ordered list of sources, eg. 'minio' or 'minio file'."`

	MinioOptions
	FileSourceOptions
}

type Source interface {
	Read(key string) ([]byte, error)
}

func ReadFromSources(sources []Source, key string) ([]byte, error) {

	for _, s := range sources {

		body, err := s.Read(key)
		if err == nil {
			return body, err
		}
	}

	return nil, errors.New("could not read item from configured sources")

}
