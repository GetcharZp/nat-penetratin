package test

import (
	"log"
	"net"
	"testing"
)

const (
	addr    = "0.0.0.0:22222"
	bufSize = 10
)

// 监听
func TestTcpListen(t *testing.T) {
	log.Println("Tcp Server Running")
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	tcpListen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		t.Fatal(err)
	}
	for {
		tcpConn, err := tcpListen.AcceptTCP()
		if err != nil {
			t.Fatal(err)
		}
		// 读取数据
		data := make([]byte, 0)
		for {
			var buf [bufSize]byte
			n, err := tcpConn.Read(buf[:])
			if err != nil {
				t.Fatal(err)
			}
			data = append(data, buf[:n]...)
			if n < bufSize {
				break
			}
		}
		log.Println("Receive Data : ", string(data))
		// 写数据
		tcpConn.Write([]byte(data))
	}
}

// 创建连接
func TestCreateTcp(t *testing.T) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	// 写数据
	tcpConn.Write([]byte("Client ==> hello world, 你好世界"))
	// 读取数据
	data := make([]byte, 0)
	for {
		var buf [bufSize]byte
		n, err := tcpConn.Read(buf[:])
		if err != nil {
			t.Fatal(err)
		}
		data = append(data, buf[:n]...)
		if n < bufSize {
			break
		}
	}
	log.Println("Send Data : ", string(data))
}
