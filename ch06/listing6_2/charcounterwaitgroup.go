package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func countLetter(url string, frequency []int, mutex *sync.Mutex) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic("Server returning error status code: " + resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	mutex.Lock()
	for _, b := range body {
		c := strings.ToLower(string(b))
		cIndex := strings.Index(allLetters, c)
		if cIndex >= 0 {
			frequency[cIndex] += 1
		}
	}
	mutex.Unlock()
	fmt.Println("Completed:", url)
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(31)
	mutex := sync.Mutex{}
	var frequency = make([]int, 26)
	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go func() {
			countLetter(url, frequency, &mutex)
			wg.Done()
		}()
	}
	wg.Wait()

	mutex.Lock()
	for i, c := range allLetters {
		fmt.Printf("%c-%d ", c, frequency[i])
	}
	mutex.Unlock()
}
