package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Barrier struct {
	size      int
	waitCount int
	cond      *sync.Cond
}

func NewBarrier(size int) *Barrier {
	condVar := sync.NewCond(&sync.Mutex{})
	return &Barrier{size, 0, condVar}
}

func (b *Barrier) Wait() {
	b.cond.L.Lock()
	b.waitCount++
	if b.waitCount == b.size {
		b.waitCount = 0
		b.cond.Broadcast()
	} else {
		b.cond.Wait()
	}
	b.cond.L.Unlock()
}

const matrixSize = 3

func matrixMultiply(matrixA, matrixB, result *[matrixSize][matrixSize]int) {
	for row := 0; row < matrixSize; row++ {
		for col := 0; col < matrixSize; col++ {
			sum := 0
			for i := 0; i < matrixSize; i++ {
				sum += matrixA[row][i] * matrixB[i][col]
			}
			result[row][col] = sum
		}
	}
}

func generateRandMatrix(matrix *[matrixSize][matrixSize]int) {
	for row := 0; row < matrixSize; row++ {
		for col := 0; col < matrixSize; col++ {
			matrix[row][col] = rand.Intn(10) - 5
		}
	}
}

func rowMultiply(matrixA, matrixB, result *[matrixSize][matrixSize]int, row int, barrier *Barrier) {
	for {
		barrier.Wait()
		for col := 0; col < matrixSize; col++ {
			sum := 0
			for i := 0; i < matrixSize; i++ {
				sum += matrixA[row][i] * matrixB[i][col]
			}
			result[row][col] = sum
		}
		barrier.Wait()
	}
}

func main() {
	var matrixA, matrixB, result [matrixSize][matrixSize]int
	barrier := NewBarrier(matrixSize + 1)

	for row := 0; row < matrixSize; row++ {
		go rowMultiply(&matrixA, &matrixB, &result, row, barrier)
	}

	for i := 0; i < 4; i++ {
		generateRandMatrix(&matrixA)
		generateRandMatrix(&matrixB)
		barrier.Wait()

		barrier.Wait()
		for i := 0; i < matrixSize; i++ {
			fmt.Println(matrixA[i], matrixB[i], result[i])
		}
		fmt.Println()
	}
}
