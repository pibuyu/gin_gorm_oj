package test

import (
	"fmt"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	ch := make(chan int)
	//创建channel的时候就可以defer关闭，避免通道不关泄露内存
	defer close(ch)
	go func() {
		ch <- 3 + 4
	}()
	i := <-ch
	fmt.Println(i)
}

func sum(ch chan int, s []int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	ch <- sum
}
func TestBlocking(t *testing.T) {
	ch := make(chan int)
	s := []int{1, 2, 3, 4, 5, 6}

	go sum(ch, s[:len(s)/2])
	go sum(ch, s[len(s)/2:])

	x, y := <-ch, <-ch
	fmt.Println(x, y, x+y)
}

// select结构，类似于switch，选择一个case处理
func fibonacci(c, quit chan int) {
	x, y := 0, 1
	// 在select外面加上for循环一直执行，直到遇到return
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}
func TestSelect(t *testing.T) {
	ch := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-ch)
		}
		quit <- 0
	}()
	fibonacci(ch, quit)
}

// 直接迭代遍历channel
func TestIterChannel(t *testing.T) {
	go func() {
		time.Sleep(1 * time.Hour)
	}()
	c := make(chan int)
	go func() {
		for i := 0; i < 10; i = i + 1 {
			c <- i
		}
		close(c)
	}()
	//直接对channel进行遍历
	for i := range c {
		fmt.Println(i)
	}
	fmt.Println("Finished")
}

// 给select里加上一个处理超时的case，避免一直阻塞等待
func TestTimeOut(t *testing.T) {
	ch := make(chan string, 1)
	go func() {
		time.Sleep(time.Second * 2)
		ch <- "result 1"
	}()
	select {
	case res := <-ch:
		fmt.Println(res)
	case <-time.After(time.Second * 1):
		fmt.Println("timeout 1")
	}
}
