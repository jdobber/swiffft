package sources

import (
	"io/ioutil"
)

type FileSource struct {
	Source
	prefix string
}

func NewFileSource(prefix string) (*FileSource, error) {

	c := FileSource{
		prefix: prefix,
	}

	return &c, nil
}

func (c *FileSource) Read(key string) ([]byte, error) {

	return ioutil.ReadFile(c.prefix + "/" + key)

}
