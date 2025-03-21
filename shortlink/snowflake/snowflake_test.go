package snowflake

import (
	"fmt"
	"github.com/jxskiss/base62"
	"sync"
	"testing"
	"time"
)

func TestSnowflake_GenerateId(t *testing.T) {
	start := time.Now().UnixMilli()
	done := make(chan struct{})
	ch := make(chan uint64, 1024)
	hashed := make(map[uint64]struct{})
	go func() {
		for key := range ch {
			if _, ok := hashed[key]; ok {
				fmt.Printf("%d 重复\n", key)
			} else {
				hashed[key] = struct{}{}
				fmt.Println(string(base62.FormatUint(key)))
			}
		}
		fmt.Println(len(hashed))
		done <- struct{}{}
	}()
	snow1 := NewSnowflake(1, 1)
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for y := 0; y < 1000; y++ {
				ch <- snow1.GenerateId()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(ch)
	<-done
	fmt.Printf("handler timestamp %d \n", time.Now().UnixMilli()-start)
}

func TestNewSnowflake(t *testing.T) {
	timestamp := time.Now().UnixMilli()
	fmt.Println(timestamp)
	fmt.Printf("%b \n", timestamp)
}
