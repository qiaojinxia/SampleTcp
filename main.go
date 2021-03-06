package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
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
	log.Info("server start...")
	ctx, _ := context.WithCancel(context.Background())
	go func() {
		<-quit
		close(quit)
		log.Info("server stop...")
		ctx.Done()     //控制Handle 不再接受新的请求
		process.Wait() //等待所有已有的请求处理完毕
		listen.Close() //关闭tcp监听
	}()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Error(err)
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
			log.Error(err)
		}
		log.Infof("close Remote Conn %s", conn.RemoteAddr())
		conn.Close()
		atomic.AddInt64(&OnlineSession, -1)
	}()
	num := atomic.AddUint64(&GlobalUid, 1)
	sess, err := NewSession(num, conn, ctx)
	if err != nil {
		log.Errorf(" %s Init Session Failed", conn.RemoteAddr())
		return
	}
	atomic.AddInt64(&OnlineSession, 1)
	log.Infof("now user %d online  ", OnlineSession)
	sess.Run()
}
