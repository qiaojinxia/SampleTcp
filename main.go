package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	Init()
	Start(fmt.Sprintf("127.0.0.1:%s", AppConfig.Port))
}

var (
	quit = make(chan os.Signal, 1)
)

func Start(addr string) {
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGKILL)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Println("server start...")
	ctx, _ := context.WithCancel(context.Background())
	go func() {
		<-quit
		close(quit)
		log.Println("server stop...")
		ctx.Done()     //控制Handle 不再接受新的请求
		process.Wait() //等待所有已有的请求处理完毕
		listen.Close() //关闭tcp监听
	}()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}
		if isChanClose(quit) {
			return
		}
		go handle(ctx, conn)
	}
}

func handle(ctx context.Context, conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
		log.Printf("close Conn %s", conn.RemoteAddr().String())
		conn.Close()
	}()
	reader, err := NewReader(conn, 1024, 4, 2, ctx)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			select {
			case msg := <-reader.Message:
				fmt.Println("收到消息:", msg)
			default:
			}
		}
	}()
	err = reader.Do()
	if err != nil {
		log.Fatal(err)
	}

}
