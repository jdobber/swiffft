package sources

import (
	"bytes"
	"encoding/gob"
	"errors"
)

type SourceOptions struct {
	Sources []string `name:"sources" default:"file" desc:"An ordered list of sources, eg. 'minio' or 'minio file'."`

	MinioOptions
	FileSourceOptions
}

// SourceInfo ...
type SourceInfo struct {
	Payload      []byte
	LastModified string
	ETag         string
}

type Source interface {
	Read(key string) (*SourceInfo, error)
}

func ReadFromSources(sources []Source, key string) (*SourceInfo, error) {

	for _, s := range sources {

		payload, err := s.Read(key)
		if err == nil {
			return payload, err
		}
	}

	return nil, errors.New("could not read item from configured sources")

}

// Pack ...
func Pack(sourceInfo *SourceInfo) ([]byte, error) {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(sourceInfo)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}

// Unpack ...
func Unpack(b []byte) (*SourceInfo, error) {
	var sourceInfo SourceInfo
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&sourceInfo)
	if err != nil {
		return nil, err
	}
	return &sourceInfo, nil
}
