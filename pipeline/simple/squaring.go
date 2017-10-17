package main

import (
	"fmt"
	"sync"
)

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			fmt.Printf("n: %v\n", n)
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(v int, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			fmt.Printf("%v: n * n: %v\n", v, n)
			out <- n * n
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	fmt.Println("Simple")
	// Set up the pipeline.
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	for n := range sq(2, sq(1, gen(nums...))) {
		fmt.Println(n)
	}

	// Fan out
	fmt.Println("Fan out")
	in := gen(nums...)

	c1 := sq(1, in)
	c2 := sq(2, in)
	c3 := sq(3, in)
	// Consume the merged output from c1 and c2.
	for n := range merge(c1, c2, c3) {
		fmt.Println(n) // 4 then 9, or 9 then 4
	}

}
