package sources

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/jdobber/swiffft/lib/middleware"
	minio "github.com/minio/minio-go/v6"
)

type MinioSource struct {
	Source
	Client         *minio.Client
	Opts           MinioOptions
	RewriteHandler *middleware.RewriteHandler
}

type MinioOptions struct {
	Endpoint string   `name:"minio.endpoint" desc:"The Minio endpoint to use."`
	Bucket   string   `name:"minio.bucket" desc:"The Minio bucket to use."`
	UseSSL   bool     `name:"minio.usessl" desc:"Use ssl for Minio."`
	Rewrites []string `name:"minio.rewrites" desc:"An ordered list of rewrite rules, eg. ':key::/new/:key/path'."`
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
		Client:         minioClient,
		Opts:           opts,
		RewriteHandler: middleware.NewRewriteHandler(opts.Rewrites),
	}

	return &c, nil
}

func (c *MinioSource) Read(key string) (*SourceInfo, error) {

	ok, to := c.RewriteHandler.ApplyRules(key)
	if !ok {
		to = key
	}

	sourceInfo := SourceInfo{}

	info, err := c.Client.StatObject(c.Opts.Bucket, to, minio.StatObjectOptions{})
	if err != nil {
		return nil, err

	}

	sourceInfo.LastModified = info.LastModified.Format(time.RFC1123Z)

	object, err := c.Client.GetObject(c.Opts.Bucket, to, minio.GetObjectOptions{})
	defer object.Close()
	if err != nil {
		return nil, err
	}

	sourceInfo.Payload, err = ioutil.ReadAll(object)
	if err != nil {
		return nil, err
	}

	return &sourceInfo, nil
}
