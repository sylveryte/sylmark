package data

import (
	"strings"
	"sylmark/lsp"
)

type Target string  // is like # Some heading
type GTarget string // is full GLink target
type Tag string

func GetGTarget(heading string, uri lsp.DocumentURI) (gtarget GTarget, ok bool) {
	filename := uri.GetFileName()
	splits := strings.Split(filename, ".md")
	if len(splits) < 1 {
		return "", false
	}

	return GTarget(splits[0] + "#" + heading), true
}
