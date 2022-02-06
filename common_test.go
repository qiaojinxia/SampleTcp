package main

import (
	"fmt"
	"net"
	"testing"
	"time"
)

//测试 sync.Pool 和 直接申请内存 性能
//go test -bench BenchmarkBytes .\common_test.go .\common.go

func BenchmarkBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		obj := make([]byte, 1024)
		_ = obj
	}
}

func BenchmarkBytesWithCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		obj := cacheBytes.Get().(*[]byte)
		_ = obj
		cacheBytes.Put(obj)
	}
}

func Test_client(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	for i := 0; i < 100; i++ {
		data := PackageMsg([]byte(fmt.Sprintf("caomaoboy的第%d条消息", i)))
		_, err = conn.Write(data)
		if err != nil {
			panic(err)
		}
	}

	time.Sleep(time.Second * 10)
}
