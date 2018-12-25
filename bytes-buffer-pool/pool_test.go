package bytes_buffer_pool

import (
	"bytes"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"testing"
	"time"
)

func init() {
	go http.ListenAndServe(":9000", nil)
}

func Test_Contrast_Time(t *testing.T) {
	b := time.Now().UnixNano()
	for i := 0; i < 10000; i++ {
		_ = bytes.NewBuffer(make([]byte, 0, 100000))
	}
	e := time.Now().UnixNano()
	fmt.Println("no pool cost time ", (e-b)/int64(time.Microsecond))
	pool := sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 100000))
		},
	}
	pb := time.Now().UnixNano()
	for i := 0; i < 10000; i++ {
		pool.Put(pool.Get())
	}
	pe := time.Now().UnixNano()
	fmt.Println("pool cost time (put back before get) ", (pe-pb)/int64(time.Microsecond))
	pb1 := time.Now().UnixNano()
	for i := 0; i < 10000; i++ {
		_ = pool.Get()
	}
	pe1 := time.Now().UnixNano()
	fmt.Println("pool cost time (get and no put back) ", (pe1-pb1)/int64(time.Microsecond))
}
