// adapted from here: https://github.com/pkg4go/rewrite/blob/master/rewrite.go

package middleware

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

type RewriteOptions struct {
	Rewrites []string `name:"rewrites" desc:"An ordered list of rewrite rules, eg. ':key::/new/:key/path'."`
}

type Rule struct {
	Pattern string
	To      string
	*regexp.Regexp
}

var regfmt = regexp.MustCompile(`:[^/#?()\.\\]+`)

func NewRule(pattern, to string) (*Rule, error) {
	pattern = regfmt.ReplaceAllStringFunc(pattern, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
	})

	reg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &Rule{
		pattern,
		to,
		reg,
	}, nil
}

func (r *Rule) Rewrite(from string) (bool, string) {

	if !r.MatchString(from) {
		return false, ""
	}

	to := path.Clean(r.Replace(from))

	return true, to
}

func (r *Rule) Replace(from string) string {
	if !hit("\\$|\\:", r.To) {
		return r.To
	}

	regFrom := regexp.MustCompile(r.Pattern)
	match := regFrom.FindStringSubmatchIndex(from)

	result := regFrom.ExpandString([]byte(""), r.To, from, match)

	str := string(result[:])

	if hit("\\:", str) {
		return r.replaceNamedParams(from, str)
	}

	return str
}

var urlreg = regexp.MustCompile(`:[^/#?()\.\\]+|\(\?P<[a-zA-Z0-9]+>.*\)`)

func (r *Rule) replaceNamedParams(from, to string) string {
	fromMatches := r.FindStringSubmatch(from)

	if len(fromMatches) > 0 {
		for i, name := range r.SubexpNames() {
			if len(name) > 0 {
				to = strings.Replace(to, ":"+name, fromMatches[i], -1)
			}
		}
	}

	return to
}

func NewRewriteHandler(rules []string) *RewriteHandler {
	var h RewriteHandler

	for _, rule := range rules {
		s := strings.Split(rule, "::")
		r, e := NewRule(s[0], s[1])
		if e != nil {
			panic(e)
		}

		h.rules = append(h.rules, r)
	}
	return &h
}

type RewriteHandler struct {
	rules []*Rule
}

func (h *RewriteHandler) ApplyRules(from string) (bool, string) {
	for _, r := range h.rules {
		ok, to := r.Rewrite(from)
		if ok {
			return ok, to
		}
	}
	return false, ""
}

func hit(pattern, str string) bool {
	r, e := regexp.MatchString(pattern, str)
	if e != nil {
		return false
	}

	return r
}
