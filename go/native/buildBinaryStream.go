package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const MaxPacketSize = 4096

/*
构造二进制数据流，便于传输
*/

type message struct {
	length uint32 // 消息长度
	id     uint32 // 消息ID
	data   []byte // 消息内容

}

// Pack 封包方法(压缩数据 将消息打包成二进制数据以便于传输)
func Pack(msg message) ([]byte, error) {
	//创建一个存放bytes字节的缓冲区
	/*
		缓冲区是一个用于存储数据的临时存储区域。
		这里因为消息的数据长度、ID和数据长度是不确定的，所以需要使用缓冲区，动态的存储这些数据
	*/
	dataBuff := bytes.NewBuffer([]byte{})

	/*
		LittleEndian是一种字节序，它指定了在多字节数据类型(如int、float等)的存储中，最低有效字节(即最右边的字节)先存储在内存中。
		这与BigEndian相反，后者将最高有效字节(即最左边的字节)存储在内存中。
	*/
	fmt.Printf("pack msg process begin. msg.DataLen^%d|msg.MsgID^%d|msg.Data^%v|dataBuff^%v\n",
		msg.length, msg.id, msg.data, dataBuff)
	//写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.length); err != nil {
		return nil, err
	}
	//if err := binary.Write(dataBuff, binary.BigEndian, msg.length); err != nil {
	//	return nil, err
	//}
	fmt.Printf("insert dataLen to pack res. dataBuff|%v\n", dataBuff.Bytes())

	//写msgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.id); err != nil {
		return nil, err
	}
	fmt.Printf("insert msgId to pack res. dataBuff|%v\n", dataBuff.Bytes())

	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.data); err != nil {
		return nil, err
	}
	fmt.Printf("insert data to pack res. dataBuff|%v\n", dataBuff.Bytes())
	return dataBuff.Bytes(), nil
}

// unpack 拆包(解压数据) 将二进制数据还原成 msg
func unpack(binaryData []byte) (*message, error) {

	//创建一个从输入获取二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)
	fmt.Printf("dataBuff|%v\n", dataBuff)

	/*
		拆包的时候是分两次过程的，第二次依赖第一次的dataLen结果。
		所以Unpack只能解压出包头head的内容，得到msgId和dataLen，
		之后调用者再根据dataLen继续从io流中读取body中的数据
	*/

	//只解压head的信息，得到dataLen|msgId
	msg := &message{}

	/*
		读dataLen
		这里能正确把DataLen的值读到 msg.length 字段中，跟 msg.length 的类型有关
		它是 uint32 类型的，所以会把前4个字节的值读出来赋值给 msg.length
	*/
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.length); err != nil {
		return msg, err
	}
	fmt.Printf("dataBuff|%v^msg.length|%d\n", dataBuff, msg.length)

	//读msgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.id); err != nil {
		return nil, err
	}
	fmt.Printf("dataBuff|%v^msg.DataId|%d\n", dataBuff, msg.id)

	//判断dataLen的长度是否超出我们允许的最大包长度
	if msg.length > MaxPacketSize {
		return nil, errors.New("too large msg data recieved")
	}

	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从buffer读取一次数据
	return msg, nil
}

func main() {
	msg := message{
		id:     100,
		length: 8,
		data:   []byte("hello goku"),
	}
	packInfo, err := Pack(msg)
	fmt.Printf("[binary]|%v\n[str]%s\n", packInfo, string(packInfo))
	fmt.Println(err)

	// 拆包
	msgUnpack, _ := unpack(packInfo)
	fmt.Println(msgUnpack)
}
