package main

import (
	"fmt"
	"os"
)

func inExplored(needle Point, haystack []Point) bool {
	for _, n := range haystack {
		if n.Col == needle.Col && n.Row == needle.Row {
			return true
		}
	}
	return false
}

func emptyTmp() {
	const directory = "./tmp/"
	dir, err := os.Open(directory)
	if err != nil {
		panic(fmt.Errorf("error opening temp directory: %w", err))
	}
	defer dir.Close()
	filesToDelete, err := dir.ReadDir(0)
	for _, file := range filesToDelete {
		fullPath := directory + file.Name()
		err = os.Remove(fullPath)
		if err != nil {
			panic(fmt.Errorf("error deleting file %s: %w", fullPath, err))
		}
	}
}

func abs(i int) int {
	if i < 0 {
		return -i
	} else {
		return i
	}
}
