package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
)

//使用了length field based frame decoder
const (
	PackageLengthBytes   = 4
	HeaderLengthBytes    = 2
	ProtocolVersionBytes = 2
	OperationBytes       = 4
	SequenceIDBytes      = 4

	HeaderLength = PackageLengthBytes + HeaderLengthBytes + ProtocolVersionBytes + OperationBytes + SequenceIDBytes

	//Body
)

//解码器
func Depack(buffer []byte) []byte {
	length := len(buffer)

	var i int
	data := make([]byte, 32)

	for i = 0; i < length; i++ {
		if length < i+HeaderLength {
			break
		}

		messageLength := ByteToInt(buffer[i : i+PackageLengthBytes])
		if length < i+HeaderLength+messageLength {
			break
		}

		site := i + PackageLengthBytes
		headerLength := ByteToInt(buffer[site : site+HeaderLengthBytes])
		site += HeaderLengthBytes

		protocolVersion := ByteToInt16(buffer[site : site+ProtocolVersionBytes])
		site += ProtocolVersionBytes

		operation := ByteToInt(buffer[site : site+OperationBytes])
		site += OperationBytes

		sequenceID := ByteToInt(buffer[site : site+SequenceIDBytes])
		site += SequenceIDBytes

		fmt.Printf("packageLength: %d, headerLength: %d , protocolVersion: %d, operation: %d, sequenceID: %d \n", messageLength, headerLength, protocolVersion, operation, sequenceID)

		data = buffer[i+HeaderLength : i+HeaderLength+messageLength]
		break
	}

	if i == length {
		return make([]byte, 0)
	}

	return data
}

//字节转换成整形
func ByteToInt(n []byte) int {
	bytesbuffer := bytes.NewBuffer(n)
	var x int32
	binary.Read(bytesbuffer, binary.BigEndian, &x)

	return int(x)
}

func ByteToInt16(n []byte) int {
	bytesbuffer := bytes.NewBuffer(n)
	var x int16
	binary.Read(bytesbuffer, binary.BigEndian, &x)

	return int(x)
}

func main() {
	netListen, err := net.Listen("tcp", "localhost:7373")
	CheckErr(err)
	defer netListen.Close()

	Log("Waiting for client ...") //启动后，等待客户端访问。
	for {
		conn, err := netListen.Accept() //监听客户端
		if err != nil {
			Log(conn.RemoteAddr().String(), "发了了错误：", err)
			continue
		}
		Log(conn.RemoteAddr().String(), "tcp connection success")
		go handleConnection(conn)
	}
}

//连接处理
func handleConnection(conn net.Conn) {
	//缓冲区，存储被截断的数据
	tmpBuffer := make([]byte, 0)
	//接收解包
	readerChannel := make(chan []byte, 10000)
	go reader(readerChannel)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), "connection error: ", err)
			return
		}

		tmpBuffer = Depack(append(tmpBuffer, buffer[:n]...))
		readerChannel <- tmpBuffer //接收的信息写入通道
	}
	defer conn.Close()
}

//获取通道数据
func reader(readerchannel chan []byte) {
	for {
		select {
		case data := <-readerchannel:
			Log(string(data)) //打印通道内的信息
		}
	}
}

//日志处理
func Log(v ...interface{}) {
	log.Println(v...)
}

//错误处理
func CheckErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
