package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func pacLog(lines ...string) {
	lines = append(lines, "----------------")
	result := strings.Join(lines, "\n")
	fmt.Println(result)
}

func fixPath(pathToFix string) string {
	pathToFix = path.Clean(pathToFix)
	if path.IsAbs(pathToFix) {
		return pathToFix
	}
	cwd, err := os.Getwd()
	if err != nil {
		return pathToFix
	}
	return path.Join(cwd, pathToFix)
}
