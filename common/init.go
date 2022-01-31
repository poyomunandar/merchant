package common

import (
	"os"
	"strings"
)

func init() {
	cwd, _ := os.Getwd()
	roots := strings.Split(cwd, string(os.PathSeparator))
	path := strings.Join(roots, "/")
	os.Args[0] = path
}
