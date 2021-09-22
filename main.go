package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func run() {
	num_writer := 4
	num_reader := 4

	var reader_stop, writer_stop int32
	var count int32

	ch := make(chan int, 4)

	writer_wg := sync.WaitGroup{}
	for i := 0; i < num_writer; i++ {
		writer_wg.Add(1)
		go func() {
			for atomic.LoadInt32(&writer_stop) == 0 {
				for j := 0; j < 10; j++ {
					ch <- j
				}
				atomic.AddInt32(&count, 10)
			}
			writer_wg.Done()
		}()
	}

	reader_wg := sync.WaitGroup{}
	for i := 0; i < num_reader; i++ {
		reader_wg.Add(1)
		go func() {
			for atomic.LoadInt32(&reader_stop) == 0 {
				for j := 0; j < 10; j++ {
					<-ch
				}
			}
			reader_wg.Done()
		}()
	}

	timer := time.NewTimer(time.Second * 10)
	ticker := time.NewTicker(time.Second)
	var last_c int32
	done := false

	for !done {
		select {
		case <-ticker.C:
			c := atomic.LoadInt32(&count)
			fmt.Printf("Count: %d Count/s: %d\n", c, c-last_c)
			last_c = c
		case <-timer.C:
			done = true
			fmt.Println("done")
			break
		}
	}
	c := atomic.LoadInt32(&count)
	fmt.Printf("Count: %d Count/s: %f\n", c, float32(c/10))
	atomic.StoreInt32(&writer_stop, 1)
	writer_wg.Wait()
	atomic.StoreInt32(&reader_stop, 1)
	reader_wg.Wait()
}

func main() {
	run()
}
