/**
 * Copyright 2015 @ to2.net.
 * name : tcp_coder
 * author : jarryliu
 * date : 2015-11-23 18:59
 * description :
 * history :
 */
package tcpserve

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

// 分包与解包
func encodeContent(content string) ([]byte, error) {
	// 读取消息的长度
	var length int32 = int32(len(content))
	var pkg *bytes.Buffer = new(bytes.Buffer)
	// 写入消息头
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	// 写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, []byte(content))
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

func decodeContent(reader *bufio.Reader) (string, error) {
	// 获取读取的内容长度
	var length int
	lenBytes, err := reader.Peek(4)
	err = binary.Read(bytes.NewBuffer(lenBytes), binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}
	// 判断是否是完整的包
	if reader.Buffered() < 4+length {
		return "", err
	}
	// 获取完整包的内容
	var bytes []byte = make([]byte, 4+length)
	_, err = reader.Read(bytes)
	if err != nil {
		return "", err
	}
	return string(bytes[4:]), nil
}
