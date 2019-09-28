package cmd

import (
	"encoding/json"
	"fmt"

	sources "github.com/jdobber/swiffft/lib/sources"
	caches "github.com/jdobber/swiffft/lib/caches"
	"github.com/wrfly/gua"
)

type CommandOptions struct {
	Config   string `desc:"Path to a valid config file"`
	Host     string `default:"127.0.0.1:8080" desc:"Bind the server to this host and port"`
	Endpoint string `default:"http://127.0.0.1:8080/iiif" desc:"Use this endpoint in profiles"`

	sources.SourceOptions
	caches.CacheOptions
}

func Init() CommandOptions {

	args := new(CommandOptions)
	if err := gua.Parse(args); err != nil {
		panic(err)
	}

	bs, _ := json.MarshalIndent(args, "", "  ")
	fmt.Printf("%s\n", bs)

	return *args

}
