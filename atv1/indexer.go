package main

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "strings"
	"time"
	"log"
    "myproject/plotcalc"
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

	for i := 0; i < 50; i++ {
		fileList = []string{}
		indexador = make(map[string][]string)

		now := time.Now()
		addToFileList("exemplo")
		fmt.Printf("Time to build file list: %v (Files found %d)\n", time.Since(now), len(fileList))

		now = time.Now()
		for _, file := range fileList {
			addFileToIndexer(file)
		}
		elapsed := time.Since(now)
		data = append(data, elapsed)
		fmt.Printf("Time to index files: %v\n", elapsed)
	}

	// Write stats to file
	f, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	mean := plotcalc.Mean(data)
	stddev := plotcalc.StdDev(data)
    
    plotcalc.PlotDurations(data, "durations.png")

	writer.WriteString(fmt.Sprintf("Mean: %v\n", mean))
	writer.WriteString(fmt.Sprintf("Std Dev: %v\n", stddev))
	writer.WriteString("Data:\n")
	for _, d := range data {
		writer.WriteString(fmt.Sprintf("%v\n", d))
	}
	writer.Flush()

	// Plot graph
	if err := plotcalc.PlotDurations(data, "durations.png"); err != nil {
		log.Fatal(err)
	}
}