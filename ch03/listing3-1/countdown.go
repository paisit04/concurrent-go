package main

import (
	"fmt"
	"time"
)

func main() {
	count := 5
	go countdown(&count)
	for count > 0 {
		fmt.Println(count)
		time.Sleep(time.Millisecond * 500)
	}
}

func countdown(seconds *int) {
	for *seconds > 0 {
		time.Sleep(time.Second * 1)
		*seconds -= 1
	}
}
