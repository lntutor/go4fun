package main

import (
	"fmt"
	"runtime"
)

type Vector []float64

// Apply the operation to v[i], v[i+1] ... up to v[n-1].
func (v Vector) DoSome(i, n int, u Vector, c chan int) {
	fmt.Println(i, n)
	for ; i < n; i++ {
		v[i] += u[i]
	}
	c <- 1
}

var numCPU = runtime.NumCPU()

func (a Vector) DoAll(u Vector) {
	fmt.Println("numOfCPU = ", numCPU)
	c := make(chan int, numCPU)
	for i := 0; i < numCPU; i++ {
		go a.DoSome(i*len(a)/numCPU, (i+1)*len(a)/numCPU, u, c)
	}
	// Drain the channel
	for i := 0; i < numCPU; i++ {
		<-c
	}
}

func main() {
	var v Vector
	v = []float64{1, 2, 3, 3, 4}
	var u Vector
	u = []float64{1, 2, 3, 3, 4}
	v.DoAll(u)
	fmt.Printf("u.DoAll(v) = %+v\n", v)
}
