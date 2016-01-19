package main

import (
	. "fmt"
	"runtime"
)

var i int = 0

func func1() {
	for j := 0; j < 1000000; j++ {
		i += 1
	}
}
func func2() {
	for j := 0; j < 1000000; j++ {
		i -= 1
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	go func1()
	go func2()
	Println("Hello from main")
	Println(i)
}
