package main

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "strings"
	"time"
	"log"
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
			indexador[word] = append(indexador[word], filePath)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
}

func main() {
    var data []time.Duration

    for i := 0; i < 5; i++ {
	    fileList = []string{}
		indexador = make(map[string][]string)
		
        now := time.Now()
        addToFileList("exemplo") // your function
        fmt.Printf("Time taken to file to build file list: %v\n", time.Since(now))
        fmt.Println("Files found ", len(fileList))

        now = time.Now()
        for _, file := range fileList {
            addFileToIndexer(file) // your function
        }
        elapsed := time.Since(now)
        data = append(data, elapsed)
        fmt.Printf("Time taken to index files: %v\n", elapsed)
    }

    f, err := os.Create("output.txt") // creates or truncates
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    writer := bufio.NewWriter(f)
    for _, d := range data {
        _, err := writer.WriteString(fmt.Sprintf("%v\n", d)) // convert duration to string
        if err != nil {
            log.Fatal(err)
        }
    }
    writer.Flush()
}