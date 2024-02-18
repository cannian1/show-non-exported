package _runtime

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func TestGetProcID(t *testing.T) {
	numProcs := runtime.GOMAXPROCS(0)
	if numProcs > runtime.NumCPU() {
		t.Skip("unreliable with high GOMAXPROCS")
	}
	seen := make([]int64, numProcs)
	// 确保所有 goroutine 同时开始执行，避免由于 goroutine 启动时间不同步而导致的竞态条件
	start := make(chan struct{})
	var wg sync.WaitGroup
	for range seen {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for i := 0; i < 1e6; i++ {
				atomic.AddInt64(&seen[GetProcID()], 1)
			}
		}()
	}
	close(start)
	wg.Wait()
	for i, n := range seen {
		if n == 0 {
			t.Fatalf("did not see proc id %d (got: %v)", i, seen)
		}
		fmt.Println(i, n)
	}
}

func TestGetGoID(t *testing.T) {
	for i := 0; i < 100; i++ {
		go func(i int) {
			goid := GetGoID()
			t.Logf("goroutine %d, goid: %d", i, goid)
		}(i)
	}
}

func BenchmarkGetGoID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetGoID()
	}
}
