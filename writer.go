package main

import (
	"context"
	"net"
)

type Writer struct {
	Conn    net.Conn
	Message chan []byte     //单次接收的完整数据
	ctx     context.Context //关闭时通知停止读取
}

func NewWriter(conn net.Conn, ctx context.Context) *Writer {
	return &Writer{
		Conn:    conn,
		Message: make(chan []byte, 10),
		ctx:     ctx,
	}
}

func (w *Writer) SendAsync(msg []byte) {
	w.Message <- msg
}

func (w *Writer) Run() error {
	for {
		select {
		case <-w.ctx.Done():
			return nil
		case msg := <-w.Message:
			err := w.SendSync(msg)
			if err != nil {
				return err
			}
		default:
		}
	}
}

func (w *Writer) SendSync(msg []byte) error {
	var err error
	HandlerFunc(func() error {
		_, err = w.Conn.Write(msg)
		return err
	})
	return err
}
