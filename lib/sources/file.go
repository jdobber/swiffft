package sources

import (
	"io/ioutil"
	"os"
	"time"
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

func (c *FileSource) Read(key string) (*SourceInfo, error) {

	filename := c.prefix + "/" + key

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sourceInfo := SourceInfo{}

	if info, err := f.Stat(); err == nil {

		sourceInfo.LastModified = info.ModTime().Format(time.RFC1123Z)
	}

	sourceInfo.Payload, err = ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return &sourceInfo, nil

}
