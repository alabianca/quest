package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const INDEX_COMMAND = "index"

type readResult struct {
	data []byte
	ext  string
	path string
}

func findMatchingFiles(dir, pattern string) <-chan string {
	outChannel := make(chan string)

	go func(_dir, _p string) {
		defer close(outChannel)

		filepath.WalkDir(_dir, func(path string, d fs.DirEntry, err error) error {
			_, fileName := filepath.Split(path)
			if matched, _ := filepath.Match(pattern, fileName); matched {
				outChannel <- path
			}
			return nil
		})
	}(dir, pattern)

	return outChannel
}

func readFiles(files <-chan string) <-chan readResult {
	outChannel := make(chan readResult)

	go func(_files <-chan string) {
		defer close(outChannel)

		for file := range _files {
			b, err := os.ReadFile(file)
			if err != nil {
				fmt.Printf("Error reading %s\n", file)
				continue
			}

			outChannel <- readResult{data: b, ext: filepath.Ext(file), path: file}

			//switch filepath.Ext(file) {
			//case ".json":
			//	var payload map[string]interface{}
			//	json.Unmarshal(b, &payload)
			//	fmt.Println(payload["ID"])
			//}
		}
	}(files)

	return outChannel
}

func merge(cs ...<-chan readResult) <-chan readResult {
	var wg sync.WaitGroup
	out := make(chan readResult)

	output := func(c <-chan readResult) {
		for res := range c {
			out <- res
		}
		wg.Done()
	}

	wg.Add(len(cs))

	for i, c := range cs {
		fmt.Printf("Starting worker %d\n", i+1)
		go output(c)
	}

	// go routine to close the out channel
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func index(p string) int64 {
	pattern := p
	// account of macOS "~" shortcut
	if strings.HasPrefix(p, "~/") {
		dirname, _ := os.UserHomeDir()
		pattern = path.Join(dirname, pattern[2:])
	}

	files := findMatchingFiles(filepath.Split(pattern))

	// start multiple read workers
	c1 := readFiles(files)
	c2 := readFiles(files)
	c3 := readFiles(files)
	c4 := readFiles(files)
	c5 := readFiles(files)
	c6 := readFiles(files)
	c7 := readFiles(files)
	c8 := readFiles(files)
	c9 := readFiles(files)
	c10 := readFiles(files)
	c11 := readFiles(files)
	c12 := readFiles(files)
	c13 := readFiles(files)
	c14 := readFiles(files)
	c15 := readFiles(files)
	c16 := readFiles(files)

	var numFiles int64
	for range merge(c1, c2, c3, c4, c5, c6, c7, c8, c9, c10, c11, c12, c13, c14, c15, c16) {
		numFiles++
	}

	return numFiles
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Incorrect Arguments. Require at least 1 command")
		os.Exit(1)
	}

	command := args[0]

	switch command {
	case INDEX_COMMAND:
		fmt.Println("Indexing...")
		before := time.Now()
		x := index(args[1])
		fmt.Printf("Took %f seconds to index %d files\n", time.Now().Sub(before).Seconds(), x)
	}
}
