package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch_bool1 := make(chan bool)
	ch_bool2 := make(chan bool)
	ch_bool3 := make(chan bool)

	rand.Seed(time.Now().UnixNano())

	go producer(ch1)

	go consumer(1, ch1, ch_bool1)
	go consumer(2, ch1, ch_bool2)
	go consumer(3, ch1, ch_bool3)

	<-ch_bool1
	<-ch_bool2
	<-ch_bool3

	defer fmt.Println("main ... over ...")
}

func producer(ch1 chan int) {
	for i := 1; i <= 10; i++ {
		ch1 <- i
		fmt.Println("生产蛋糕，标号为：", i)
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
	defer close(ch1)
}

func consumer(num int, ch1 chan int, ch chan bool) {
	for data := range ch1 {
		pre := strings.Repeat("_____", num)
		fmt.Printf("%s %d号购买%d号蛋糕 \n", pre, num, data)
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
	ch <- true
	defer close(ch)
}
