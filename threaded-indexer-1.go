package main

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "strings"
	"time"
	"sync"
)

var (
	counter int
	mu      sync.Mutex
)


type Indexador struct {
    indice map[string][]string
}

var fileList []string
var indexador = make(map[string][]string)

func addToFileList(filePath string) {
	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}else{
			fileList = append(fileList, path)
		}
		return nil
	})

	if err != nil {
        fmt.Printf("Error walking directory: %v\n", err)
        return
    }
}

func addFileToIndexer(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)
		for _, word := range words {
			mu.Lock()
			indexador[word] = append(indexador[word], filePath)
			mu.Unlock()
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
}

func main() {
	now := time.Now()
    addToFileList("exemplo")
	fmt.Printf("Time taken to file to build file list: %v\n", time.Since(now))

    fmt.Println("*****************************")
    fmt.Println("Files found ", len(fileList))
    fmt.Println("File List:", fileList)
	now = time.Now()
	
	for _, file := range fileList {
        go addFileToIndexer(file)
    }

    fmt.Println("\n\n*****************************")
	fmt.Printf("Time taken to file to build indexer: %v\n", time.Since(now))
    fmt.Println("*****************************")
    fmt.Println("\n\n*****************************")
	fmt.Println("Index:")
	for word, files := range indexador {
		fmt.Printf("%s: %v\n", word, files)
	}
}