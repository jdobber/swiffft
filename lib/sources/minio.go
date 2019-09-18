package sources

import (
	"io/ioutil"

	minio "github.com/minio/minio-go/v6"
)

type MinioSource struct {
	Source
	Client *minio.Client
	Opts   MinioOptions
}

type MinioOptions struct {
	MinioEndpoint  string `help:The Minio endpoint to use."`
	MinioBucket    string `help:"The bucket to use."`
	MinioUseSSL    bool   `help:"Use ssl."`
	MinioAccessKey string `arg:"env:MINIO_ACCESS_KEY" help:"read the access key from env var MINIO_ACCESS_KEY"`
	MinioSecretKey string `arg:"env:MINIO_SECRET_KEY"`
}

func MinioDefaultOptions() MinioOptions {
	return MinioOptions{
		MinioEndpoint: "localhost:4000",
		MinioBucket:   "iiif",
		MinioUseSSL:   false,
	}
}

func NewMinioSource(opts MinioOptions) (*MinioSource, error) {

	// Initialize minio client object.
	minioClient, err := minio.New(
		opts.MinioEndpoint,
		opts.MinioAccessKey,
		opts.MinioSecretKey,
		opts.MinioUseSSL)
	if err != nil {
		return nil, err
	}

	c := MinioSource{
		Client: minioClient,
		Opts:   opts,
	}

	return &c, nil
}

func (c *MinioSource) Read(key string) ([]byte, error) {

	object, err := c.Client.GetObject(c.Opts.MinioBucket, key, minio.GetObjectOptions{})
	defer object.Close()

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(object)
}
