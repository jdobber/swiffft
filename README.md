# swiffft
A lightweight and fast IIIF server.

## setup

```sh
go build -o bin/swiffft main.go
./bin/swiffft -config ./config.json -sources "file" -file.prefix ./fixtures -cache.activate true -cache.tiles true -cache.size 128
```