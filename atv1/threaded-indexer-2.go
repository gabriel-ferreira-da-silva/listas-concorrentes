package main

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "strings"
	"time"
	"sync"
	"log"
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
	numthreads := []int{2,3,4,5,6,7,8,9,10,11,12}
	eval := map[int]time.Duration{}
    var data []time.Duration
	
	for _, num:= range numthreads {
		fileList = []string{}
		indexador = make(map[string][]string)
		fmt.Println("\n\n*****************************")
		fmt.Printf("Running with %d threads\n", num)
		
		var wg sync.WaitGroup

		now := time.Now()
		addToFileList("exemplo")
		elapsed := time.Since(now)
		fmt.Printf("Time taken to file to build file list: %v\n", elapsed)

		fmt.Println("Files found ", len(fileList))
		fmt.Println("File List:", fileList)

		workload := len(fileList) / num
		loadlist := [][]string{}


		for i := 0; i < num; i++ {
			start := i * workload
			end := start + workload
			if(i==num-1){
				end = len(fileList)
			}
			loadlist = append(loadlist, fileList[start:end])
		}

		now = time.Now()

		for _, load :=range loadlist {
			wg.Add(1)
			go func(files []string) {
				defer wg.Done()
				for _, file := range files {
					addFileToIndexer(file)
				}
			}(load)
		}

		wg.Wait()
		elapsed = time.Since(now)
        data = append(data, elapsed)


		fmt.Printf("Number of files indexed: %d\n", len(fileList))
		fmt.Println("number of threads:", num)
		fmt.Printf("Time taken to file to build indexer: %v\n", elapsed)
		fmt.Println("*****************************\n\n")
		eval[num] = elapsed
	}

	fmt.Println("\n\nEvaluation:")
	for k, v := range eval {
		fmt.Printf("Threads: %d -> Time: %v\n", k, v)
	}



    f, err := os.Create("output2.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    writer := bufio.NewWriter(f)
    for _, d := range data {
        _, err := writer.WriteString(fmt.Sprintf("%v\n", d))
        if err != nil {
            log.Fatal(err)
        }
    }
    writer.Flush()
}