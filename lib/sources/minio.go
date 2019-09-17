package sources

import (
	"io/ioutil"

	minio "github.com/minio/minio-go/v6"
)

type MinioSource struct {
	Source
	Client *minio.Client
}

func NewMinioSource() (*MinioSource, error) {

	// Use a secure connection.
	ssl := false

	// Initialize minio client object.
	minioClient, err := minio.New("localhost:4000", "n365-test", "n365-test", ssl)
	if err != nil {
		return nil, err
	}

	c := MinioSource{
		Client: minioClient,
	}

	return &c, nil
}

func (c *MinioSource) Read(key string) ([]byte, error) {

	object, err := c.Client.GetObject("iiif", key, minio.GetObjectOptions{})
	defer object.Close()

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(object)
}
