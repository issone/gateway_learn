package unpack

import (
	"encoding/binary"
	"errors"
	"io"
)

const MsgHeader = "12345678"

func Encode(bytesBuffer io.Writer, content string) error {
	//定义数据格式: msg_header + content_len + content

	//数据长度 8 + 4 + content_len

	// 写入固定消息头
	if err := binary.Write(bytesBuffer, binary.BigEndian, []byte(MsgHeader)); err != nil {
		return err
	}
	// 写入消息内容的长度，将内容的长度转int32,占4个字节
	clen := int32(len([]byte(content)))
	if err := binary.Write(bytesBuffer, binary.BigEndian, clen); err != nil {
		return err
	}
	// 写入实际消息内容
	if err := binary.Write(bytesBuffer, binary.BigEndian, []byte(content)); err != nil {
		return err
	}
	return nil

}

func Decode(bytesBuffer io.Reader) (bodyBuf []byte, err error) {
	MagicBuf := make([]byte, len(MsgHeader))
	if _, err := io.ReadFull(bytesBuffer, MagicBuf); err != nil {
		return nil, err
	}
	if string(MagicBuf) != MsgHeader {
		return nil, errors.New("msg_header error")
	}
	lengthBuf := make([]byte, 4)
	if _, err = io.ReadFull(bytesBuffer, lengthBuf); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lengthBuf)
	bodyBuf = make([]byte, length)
	if _, err = io.ReadFull(bytesBuffer, bodyBuf); err != nil {
		return nil, err
	}
	return bodyBuf, err

}
