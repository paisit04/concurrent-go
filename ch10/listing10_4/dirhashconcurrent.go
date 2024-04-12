package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func FHash(filepath string) []byte {
	file, _ := os.Open(filepath)
	defer file.Close()

	sha := sha256.New()
	io.Copy(sha, file)
	return sha.Sum(nil)
}

func main() {
	dir := os.Args[1]
	files, _ := os.ReadDir(dir)
	sha := sha256.New()
	var prev, next chan int
	for _, file := range files {
		if !file.IsDir() {
			next = make(chan int)
			go func(filename string, prev, next chan int) {
				fPath := filepath.Join(dir, filename)
				hashOnFile := FHash(fPath)
				if prev != nil {
					<-prev
				}
				sha.Write(hashOnFile)
				next <- 0
			}(file.Name(), prev, next)
			prev = next
		}
	}
	<-next
	fmt.Printf("%x\n", sha.Sum(nil))
}
