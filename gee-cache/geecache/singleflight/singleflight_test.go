package singleflight

import (
	"sync"
	"testing"
)

func TestDo(t *testing.T) {
	var g Group
	var wg sync.WaitGroup
	requestCount, callCount := 10, 0
	for i := 0; i < requestCount; i++ {
		wg.Add(1)
		go func() {
			v, err := g.Do("key", func() (interface{}, error) {
				for j := 0; j < 10000000; j++ {
					callCount++
					callCount--
				}
				callCount++
				return "bar", nil
			})
			if v != "bar" || err != nil {
				t.Errorf("Do v = %v, error = %v", v, err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	t.Logf("requestCount = %v, callCount = %v", requestCount, callCount)

}
