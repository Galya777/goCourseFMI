package week6

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func multiplex(ctx context.Context, inputs []<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	wg.Add(len(inputs))

	for _, ch := range inputs {
		go func(c <-chan int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-c:
					if !ok {
						return
					}
					select {
					case out <- v:
					case <-ctx.Done():
						return
					}
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func generator(ctx context.Context, power int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		for i := 1; i <= 1000; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(time.Duration(r.Intn(1000)) * time.Millisecond)
				out <- int(math.Pow(float64(i), float64(power)))
			}
		}
	}()

	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel след 1 минута
	time.AfterFunc(1*time.Minute, cancel)

	var channels []<-chan int
	for i := 1; i <= 5; i++ {
		channels = append(channels, generator(ctx, i))
	}

	out := multiplex(ctx, channels)

	// Main goroutine чете резултатите
	for v := range out {
		fmt.Println(v)
	}

	// Даваме време на всички goroutines да приключат
	time.Sleep(500 * time.Millisecond)

	fmt.Println("Active goroutines:", runtime.NumGoroutine())
}
