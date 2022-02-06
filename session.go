package main

import (
	"context"
	"log"
	"net"
)

type Session struct {
	Uid uint64
	*Reader
	*Writer
}

func NewUser(uid uint64, conn net.Conn, ctx context.Context) *Session {
	reader, err := NewReader(conn, 1024, 4, 2, ctx)
	writer := NewWriter(conn, ctx)
	if err != nil {
		log.Fatal(err)
	}
	return &Session{Uid: uid, Reader: reader, Writer: writer}
}

func (u *Session) Run() {
	HandlerAsyncFunc(
		func() error {
			err := u.Writer.Run()
			return err
		},
	)
	HandlerAsyncFunc(
		func() error {
			err := u.Reader.Run()
			return err
		})

	//序列化数据

}
