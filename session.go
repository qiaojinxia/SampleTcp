package main

import (
	"context"
	"fmt"
	"net"
)

type Session struct {
	Uid uint64
	*Reader
	*Writer
	Die      chan struct{}
	userData interface{} //用户的业务数据
}

func NewSession(uid uint64, conn net.Conn, ctx context.Context) (*Session, error) {
	reader, err := NewReader(conn, 1024, 4, 2, ctx)
	if err != nil {
		return nil, err
	}
	writer := NewWriter(conn, ctx)
	return &Session{Uid: uid, Reader: reader, Writer: writer, Die: make(chan struct{}, 2)}, nil
}

func (s *Session) SetUserData(userData interface{}) {
	s.userData = userData
}

func (s *Session) GetUserData() interface{} {
	return s.userData
}

func (s *Session) Run() {
	HandlerAsyncFunc(
		func() error {
			err := s.Writer.Run()
			s.Die <- struct{}{}
			return err
		})
	HandlerAsyncFunc(
		func() error {
			err := s.Reader.Run()
			s.Die <- struct{}{}
			return err
		})

	for {
		select {
		case msg := <-s.Reader.Message:
			fmt.Println(msg)
		case <-s.Die:
			return
		default:
		}
	}
	//序列化数据

}
