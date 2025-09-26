
## como executar

o arquivo principal é **thread_indexer.go**. Ele foi construido como modulo go e indexa os arquivos dentro da pasta **exemplo**.

```shell
go run thread_indexer.go
```



## Resultado

A construção do indexador tem como base duas funções addToFIlelist e addFileToIndexer. A primeira  função addToFilelist adiciona o path de todos os arquivos dentro   de um diretório à uma variável global chamada filelist 

```go
var fileList []string
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
```

A segunda recebe o path de um arquivo a adiciona suas palavras ao indexador global

```go
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
```

A carga de trabalho de cada thread é definida separando a lista filelist em listas menores que serão a carga de trabalho de cada thread. Assim, se existem 100 arquivos na filelist e 5 threads cada thread ficará responsavel por indexar 20 arquivos.

```go
workload := len(fileList) / num
loadlist := [][]string{}
```

Cada thread então aplica a função addFileToIndexer nos arquivos que estão na sua carga de trabalho.

```go
	
	for _, load :=range loadlist {
				wg.Add(1)
				go func(files []string) {
					defer wg.Done()
					for _, file := range files {
						addFileToIndexer(file)
					}
				}(load)
			}
```

Nos experimentos comparamos o tempo de execução das threads para finalizar a tarefa. Para cada numero thread executamos a tarefa 50 vezes. Abaixo estão os resultados com a média e desvio padrão do tempo que as threads usaram para construir o indexador. 

![](https://github.com/gabriel-ferreira-da-silva/listas-concorrentes/blob/main/atv1/thread_means_stddev.png?raw=true)
