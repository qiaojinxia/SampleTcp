package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
)

type Reader struct {
	Conn         net.Conn
	Buff         []byte          //数据接收缓冲区
	Start        int             //数据读取开始位置
	End          int             //数据读取结束位置
	BuffLen      int             //数据接收缓冲区大小
	HeaderLen    int             //包头长度
	LengthOffset int             //指示包体长度的字段在包头中的位置（2字节）
	Message      chan string     //单次接收的完整数据
	ctx          context.Context //关闭时通知停止读取
}

func NewReader(conn net.Conn, maxBufferSize, headerLen, lengthOffset int, ctx context.Context) (*Reader, error) {
	if lengthOffset+2 > headerLen {
		return nil, fmt.Errorf("incorrect 'headerLen' or 'lengthOffset'")
	}
	return &Reader{
		Conn:         conn,
		Buff:         make([]byte, maxBufferSize),
		Start:        0,
		End:          0,
		BuffLen:      maxBufferSize,
		HeaderLen:    headerLen,
		LengthOffset: lengthOffset,
		Message:      make(chan string, 10),
		ctx:          ctx,
	}, nil
}

func (r *Reader) Run() (err error) {
	defer close(r.Message)
	err = r.read()
	if err != nil {
		return fmt.Errorf("read data error:%v", err)
	}
	return
}

//读取tcp数据流
func (r *Reader) read() error {
	for {
		select {
		case <-r.ctx.Done():
			return nil
		default:
		}
		r.move()
		if r.End == r.BuffLen {
			//缓冲区的宽度容纳不了一条消息的长度
			return fmt.Errorf("message is too large:%v", r)
		}
		length, err := r.Conn.Read(r.Buff[r.End:])
		if err != nil {
			return err
		}
		r.End += length
		r.readFromBuff()
	}
}

//前移上一次未处理完的数据
func (r *Reader) move() {
	if r.Start == 0 {
		return
	}
	copy(r.Buff, r.Buff[r.Start:r.End])
	r.End -= r.Start
	r.Start = 0
}

//读取buff中的单条数据
func (r *Reader) readFromBuff() {
	if r.End-r.Start < r.HeaderLen {
		//包头的长度不够，继续接收
		return
	}
	//读取包头数据
	headerData := r.Buff[r.Start:(r.Start + r.HeaderLen)]

	//读取包体的长度(2字节)
	bodyLen := binary.BigEndian.Uint16(headerData[r.LengthOffset : r.LengthOffset+2])
	if r.End-r.Start-r.HeaderLen < int(bodyLen) {
		//包体的长度不够，继续接收
		return
	}
	//读取包体数据
	bodyData := r.Buff[(r.Start + r.HeaderLen) : r.Start+r.HeaderLen+int(bodyLen)]
	//把body的数据包用通道传递出去
	r.Message <- string(bodyData)
	//每读完一次数据 start 后移
	r.Start += r.HeaderLen + int(headerData[r.HeaderLen-1])
	r.readFromBuff()
}
