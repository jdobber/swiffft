package sources

import (
	"io/ioutil"
)

type FileSource struct {
	Source
	prefix string
}

type FileSourceOptions struct {
	Prefix string `name:"file.prefix" desc:"The basepath for images."`
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
