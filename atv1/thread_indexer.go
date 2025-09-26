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
	"myproject/plotcalc"

	"image/color"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
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

type BarError struct {
	Xs, Ys, Errors []float64
}

func (b BarError) Len() int                 { return len(b.Xs) }
func (b BarError) XY(i int) (x, y float64) { return b.Xs[i], b.Ys[i] }

func (b BarError) YError(i int) (low, high float64) {
	return b.Ys[i] - b.Errors[i], b.Ys[i] + b.Errors[i]
}


func main() {
	numthreads := []int{1,2,3,4,5,6,7,8,9,10,11,12}
	eval := map[int][]time.Duration{}
    
	for _, num:= range numthreads {
		var data []time.Duration
	
		for i:=0; i<50; i++{

			fileList = []string{}
			indexador = make(map[string][]string)
			
			var wg sync.WaitGroup
			now := time.Now()
			addToFileList("exemplo")
			elapsed := time.Since(now)
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
		}

		eval[num] = data
	}


    f, err := os.Create("output2.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    writer := bufio.NewWriter(f)

    for k, v := range eval {

		mean := plotcalc.Mean(v)
		stddev := plotcalc.StdDev(v)

		_, err := writer.WriteString(fmt.Sprintf("num threads: %v\n", k))
		_, err = writer.WriteString(fmt.Sprintf("Mean: %v\n", mean))
		_, err = writer.WriteString(fmt.Sprintf("StdDev: %v\n", stddev))

        if err != nil {
            log.Fatal(err)
        }
    }
    writer.Flush()

	keys := make([]int, 0, 1+len(eval))
	for k := range eval {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	values := make(plotter.Values, len(keys))
	labels := make([]string, len(keys))
	for i, k := range keys {
		mean := plotcalc.Mean(eval[k]).Seconds() * 1000 // convert to ms
		values[i] = mean
		labels[i] = fmt.Sprintf("%d", k)
	}

	p := plot.New()
	p.Title.Text = "Mean Indexing Time by Number of Threads"
	p.Y.Label.Text = "Time (ms)"

	bars, err := plotter.NewBarChart(values, vg.Points(20))
	if err != nil {
		log.Fatal(err)
	}
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = color.RGBA{R: 100, G: 150, B: 255, A: 255} // choose a color
	p.Add(bars)

	p.NominalX(labels...)

	if err := p.Save(8*vg.Inch, 4*vg.Inch, "thread_means.png"); err != nil {
		log.Fatal(err)
	}

	//values := make(plotter.Values, len(keys))
	errs := make([]float64, len(keys))
	xs := make([]float64, len(keys))
	//labels := make([]string, len(keys))

	for i, k := range keys {
		mean := plotcalc.Mean(eval[k]).Seconds() * 1000 // ms
		std := plotcalc.StdDev(eval[k]).Seconds() * 1000 // ms

		values[i] = mean
		errs[i] = std
		xs[i] = float64(i)
		labels[i] = fmt.Sprintf("%d", k)
	}

	p = plot.New()
	p.Title.Text = "Mean Indexing Time by Number of Threads"
	p.Y.Label.Text = "Time (ms)"

	// Bars
	bars, err = plotter.NewBarChart(values, vg.Points(20))
	if err != nil {
		log.Fatal(err)
	}
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = color.RGBA{R: 100, G: 150, B: 255, A: 255}
	p.Add(bars)

	// X-axis labels
	p.NominalX(labels...)

	// Error bars
	eb, err := plotter.NewYErrorBars(BarError{Xs: xs, Ys: values, Errors: errs})
	if err != nil {
		log.Fatal(err)
	}
	p.Add(eb)

	// Save
	if err := p.Save(8*vg.Inch, 4*vg.Inch, "thread_means_stddev.png"); err != nil {
		log.Fatal(err)
	}

}