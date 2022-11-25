package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const delimiter = ","

var client = &http.Client{Timeout: 10 * time.Second}

func main() {
	min, _ := strconv.Atoi(os.Getenv("MIN"))
	max, _ := strconv.Atoi(os.Getenv("MAX"))
	step, _ := strconv.Atoi(os.Getenv("STEP"))
	_ = os.WriteFile(os.Getenv("OUTPUT"), generateCsv(testLoad(min, max, step, os.Getenv("METHOD"), os.Getenv("URL"), os.Getenv("LOAD"))), 0666)
}

func generateCsv(data map[int][]int) []byte {
	headers := make([]int, 0, len(data))
	for header, latencies := range data {
		sort.Ints(latencies)
		data[header] = latencies
		headers = append(headers, header)
	}
	sort.Ints(headers)

	csv := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(headers)), delimiter), "[]")
	columnLength := headers[len(headers)-1]
	for i := 0; i < columnLength; i++ {
		csv += "\n"
		for _, header := range headers {
			if len(data[header]) > i {
				csv += fmt.Sprintf("%d%s", data[header][i], delimiter)
			} else {
				csv += delimiter
			}
		}
	}
	return []byte(csv)
}

func testLoad(min, max, step int, method, url, load string) map[int][]int {
	completeResults := make(map[int][]int)
	for requests := min; requests <= max; requests += step {
		log.Print("Nr. of requests: ", requests)
		resultsChan := make(chan int)
		errorsChan := make(chan error)
		start := time.Now()
		for i := 0; i < requests; i++ {
			go executeAndMeasureRequest(resultsChan, errorsChan, method, url, load)
		}
		var partialResults []int
		for i := 0; i < requests; i++ {
			select {
			case result := <-resultsChan:
				partialResults = append(partialResults, result)
			case err := <-errorsChan:
				log.Print(err)
			}
		}
		log.Print("Overall time needed: ", time.Now().UnixMilli()-start.UnixMilli(), "millis")
		completeResults[requests] = partialResults
		time.Sleep(time.Second)
	}
	return completeResults
}

func executeAndMeasureRequest(resultsChan chan<- int, errorsChan chan<- error, method, url, load string) {
	req, err := http.NewRequest(method, url, strings.NewReader(load))
	if err != nil {
		errorsChan <- err
		return
	}
	start := time.Now().UnixMilli()
	response, err := client.Do(req)
	if err != nil {
		errorsChan <- err
		return
	}
	latency := int(time.Now().UnixMilli() - start)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
	} else {
		log.Print(string(body))
	}
	resultsChan <- latency
}
