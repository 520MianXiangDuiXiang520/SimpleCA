package tools

import (
	"os"
	"path"
	"runtime"
)

func HasThisFile(p string) bool {
	_, currently, _, _ := runtime.Caller(1)
	filename := path.Join(path.Dir(currently), p)
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
