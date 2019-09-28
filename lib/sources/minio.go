package sources

import (
	"io/ioutil"
	"os"

	minio "github.com/minio/minio-go/v6"
)

type MinioSource struct {
	Source
	Client *minio.Client
	Opts   MinioOptions
}

type MinioOptions struct {
	Endpoint string `name:"minio.endpoint" desc:"The Minio endpoint to use."`
	Bucket   string `name:"minio.bucket" desc:"The Minio bucket to use."`
	UseSSL   bool   `name:"minio.usessl" desc:"Use ssl for Minio."`
}

func NewMinioSource(opts MinioOptions) (*MinioSource, error) {

	// Initialize minio client object.
	minioClient, err := minio.New(
		opts.Endpoint,
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		opts.UseSSL)
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

	object, err := c.Client.GetObject(c.Opts.Bucket, key, minio.GetObjectOptions{})
	defer object.Close()

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(object)
}
