package main

import (
	"os"
	"strings"
)

func cleanDirectory(dir string) (cleanDir string) {
	if strings.LastIndex(dir, string(os.PathSeparator)) != len(dir) {
		return dir + string(os.PathSeparator)
	}

	return dir
}

func cleanFilename(filename string) (cleanFilename string) {
	cleanFilename = strings.Replace(filename, "(", "_", -1)
	cleanFilename = strings.Replace(cleanFilename, ")", "_", -1)

	return
}
