#!/bin/bash

echo "Modifying go.mod"
go mod edit -replace github.com/jdobber/go-iiif-mod=$(pwd)/../go-iiif-mod

echo "DONE"
echo -e "Now build the test program:\n\ngo build -o bin/test main.go\n"
echo -e "Then run:\n\n./bin/swiffft -h\n"