package main

import (
	"encoding/json"
	"fmt"
	"github.com/alabianca/quest/trie"
	"io/fs"
	"math"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const IndexCommand = "index"

type document struct {
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

func readFiles(files <-chan string) <-chan document {
	outChannel := make(chan document)

	go func(_files <-chan string) {
		defer close(outChannel)

		for file := range _files {
			b, err := os.ReadFile(file)
			if err != nil {
				fmt.Printf("Error reading %s\n", file)
				continue
			}

			outChannel <- document{data: b, ext: filepath.Ext(file), path: file}

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

func merge(cs ...<-chan document) <-chan document {
	var wg sync.WaitGroup
	out := make(chan document)

	output := func(c <-chan document) {
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

func processDocument(t *trie.Trie, doc document, params map[string][]string) {
	switch doc.ext {
	case ".json":
		var payload map[string]interface{}
		err := json.Unmarshal(doc.data, &payload)
		if err != nil {
			return
		}
		if fieldsToIndex, ok := params[FieldsParam]; ok && len(fieldsToIndex) > 0 {
			// only index specific fields
			for _, field := range fieldsToIndex {
				v, ok := payload[field]
				if !ok {
					continue
				}
				// only index strings
				switch v.(type) {
				case string:
					fmt.Printf("Inserting %v\n", v)
					t.Insert(fmt.Sprintf("%v", v))
				}
			}
		} else {
			// index all fields
			// TODO
		}
	}
}

func index(params map[string][]string, concurrency int) (int64, error) {
	patternParam, ok := params[PatternParam]
	if !ok {
		return 0, fmt.Errorf("-p is required when indexing")
	}
	pattern := patternParam[0]
	if pattern == "" {
		pattern = "."
	}
	c := int(math.Min(float64(runtime.NumCPU()), float64(concurrency)))
	if c < 1 {
		c = 1
	}
	// account of macOS "~" shortcut
	if strings.HasPrefix(pattern, "~/") {
		dirname, _ := os.UserHomeDir()
		pattern = path.Join(dirname, pattern[2:])
	}

	t := trie.New()

	files := findMatchingFiles(filepath.Split(pattern))

	// start multiple read workers
	channels := make([]<-chan document, c)
	for i := 0; i < c; i++ {
		channels[i] = readFiles(files)
	}

	var numFiles int64
	for doc := range merge(channels...) {
		processDocument(t, doc, params)
		numFiles++
	}

	return numFiles, nil
}

const PatternParam = "-p"
const FieldsParam = "-f"

func isParam(s string) bool {
	return s == PatternParam || s == FieldsParam
}

func parseParams(rawParams []string) map[string][]string {
	pList := make(map[string][]string)
	var lastParamSeen string
	for _, s := range rawParams {
		if isParam(s) {
			lastParamSeen = s
			pList[s] = make([]string, 0)
		} else if p, ok := pList[lastParamSeen]; ok {
			pList[lastParamSeen] = append(p, s)

		}
	}

	return pList
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Incorrect Arguments. Require at least 1 command")
		os.Exit(1)
	}

	command := args[0]

	switch command {
	case IndexCommand:
		params := parseParams(args[1:])
		before := time.Now()
		x, _ := index(params, 16)
		fmt.Printf("Took %f seconds to index %d files\n", time.Now().Sub(before).Seconds(), x)
	}
}
