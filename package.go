package main

import "encoding/binary"

//TCP包解析格式 指令号:数据包长度:版本号:操作指令

func PackageMsg(body []byte) (res []byte) {
	buff := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(buff, 1234)
	res = append(res, buff...)
	bLen := len(body)
	buff = buff[0:]
	binary.BigEndian.PutUint16(buff, uint16(bLen))
	res = append(res, buff...)
	res = append(res, body...)
	return
}
