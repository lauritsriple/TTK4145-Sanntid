package main

import (
	. "fmt"
	"runtime"
)

func func1(queue chan int, finish chan int) {
	temp:=0
	for j := 0; j < 1000000; j++ {
		temp= <- queue
		temp ++
		queue <- temp
	}
	finish <- 1
}

func func2(queue chan int, finish chan int) {
	temp:=0
	for j := 0; j < 1000000; j++ {
		temp= <- queue
		temp --
		queue <- temp
	}
	finish <- 1
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	queue := make(chan int,1)
	finish := make(chan int,2)
	queue <- 0
	go func1 (queue,finish)
	go func2 (queue,finish)

	<- finish
	<- finish
	Println("Hello from main")
	Println(<-queue)
}
