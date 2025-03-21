package singleflight

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestGroup_Do(t *testing.T) {
	group := &Group{}
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			result, err := group.Do("test_do", func() (any, error) {
				fmt.Println("execute")
				time.Sleep(1 * time.Second)
				return index, nil
			})
			if err != nil {
				panic(err)
			}
			fmt.Println(result)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
