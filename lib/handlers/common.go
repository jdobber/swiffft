package handlers

import (
	"strings"

	iiifparser "github.com/jdobber/go-iiif-mod/lib/parser"

	echo "github.com/labstack/echo/v4"
	"github.com/whosonfirst/go-sanitize"
)

func NewIIIFQueryParser(c echo.Context) (*iiifparser.IIIFQueryParser, error) {

	opts := sanitize.DefaultOptions()
	vars := make(map[string]string)

	for _, name := range c.ParamNames() {
		if name == "quality" {
			a := strings.Split(c.Param(name), ".")
			vars["quality"] = a[0]
			vars["format"] = a[1]
		} else {
			vars[name] = c.Param(name)
		}

	}

	p := iiifparser.IIIFQueryParser{
		Opts: opts,
		Vars: vars,
	}

	return &p, nil
}
