package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	input, _ := os.ReadFile("input.txt")

	// Create nested filesystem from commands

	commands := bytes.Split(input, []byte("$"))

	hostDir := directoryOrFile{
		parent:   nil,
		size:     0,
		name:     "",
		children: map[string]*directoryOrFile{},
	}
	activeDir := &hostDir

	for _, command := range commands {
		segments := bytes.Split(command, []byte("\n"))
		input := newCommandInput(segments[0])

		switch input.program {
		case "cd":
			if input.argument == ".." {
				activeDir = activeDir.parent
				continue
			}

			dir := &directoryOrFile{
				parent:   activeDir,
				name:     input.argument,
				size:     0,
				children: map[string]*directoryOrFile{},
			}

			activeDir.children[input.argument] = dir
			activeDir = dir

			continue

		case "ls":
			outputs := segments[1:]

			for _, output := range outputs {

				// Ignore empty output artifacts

				if string(output) == "" {
					continue
				}

				lsOutputToDirectoryOrFile(output, activeDir)
			}
		}
	}

	// Traverse tree to find directory sizes

	var directorySizes []int
	setDirectorySizes(&hostDir, &directorySizes)

	// Collate acceptable dirs and sum their sizes

	var totalSize int

	for _, size := range directorySizes {
		if size <= 100_000 {
			totalSize += size
		}
	}

	// Output result

	fmt.Printf("Total size of directories: %d\n", totalSize)

	// Reset values

	totalSize = 0

	// Construct space to be freed value

	capacity := 70_000_000
	maxSize := capacity - 30_000_000
	totalUsedSpace := directorySizes[len(directorySizes)-1]
	spaceToBeFreed := totalUsedSpace - maxSize

	// Sort ints and retrieve first acceptable size

	sort.Ints(directorySizes)

	for _, size := range directorySizes {
		if size >= spaceToBeFreed {
			totalSize = size
			break
		}
	}

	// Output result

	fmt.Printf("Size of directory to be deleted: %d\n", totalSize)
}

type commandInput struct {
	program  string
	argument string
}

func newCommandInput(inputBytes []byte) commandInput {
	trimmed := bytes.Trim(inputBytes, " ")
	programAndArgument := bytes.Split(trimmed, []byte(" "))
	input := commandInput{
		program: string(programAndArgument[0]),
	}

	if len(programAndArgument) > 1 {
		input.argument = string(programAndArgument[1])
	}

	return input
}

func setDirectorySizes(dir *directoryOrFile, directorySizes *[]int) int {
	var size int

	for _, dirOrFile := range dir.children {
		if !dirOrFile.isDir() {
			size += dirOrFile.size
			continue
		}

		size += setDirectorySizes(dirOrFile, directorySizes)
	}

	*directorySizes = append(*directorySizes, size)

	return size
}

func lsOutputToDirectoryOrFile(lsOutput []byte, parent *directoryOrFile) {
	sizeAndName := bytes.Split(lsOutput, []byte(" "))
	name := string(sizeAndName[1])

	_, ok := parent.children[name]
	if ok {
		return
	}

	var size int

	sizeStr := string(sizeAndName[0])
	if sizeStr != "dir" {
		size, _ = strconv.Atoi(sizeStr)
	}

	dir := &directoryOrFile{
		parent:   parent,
		name:     name,
		size:     size,
		children: map[string]*directoryOrFile{},
	}
	parent.children[name] = dir
}

type directoryOrFile struct {
	parent   *directoryOrFile
	name     string
	size     int
	children map[string]*directoryOrFile
}

func (d directoryOrFile) isDir() bool {
	return d.size == 0
}
