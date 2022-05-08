package test

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"testing"
)

const (
	serverAddr = "0.0.0.0:22300"
	tunnelAddr = "0.0.0.0:22301"
	BufSize    = 10
)

// server
func TestServer(t *testing.T) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
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
		b := make([]byte, 0)
		// 读取数据
		for {
			var buf [BufSize]byte
			n, err := tcpConn.Read(buf[:])
			if err != nil {
				t.Fatal(err)
			}
			b = append(b, buf[:n]...)
			if n < BufSize {
				break
			}
		}
		// 写数据
		i, err := strconv.Atoi(string(b))
		if err != nil {
			t.Fatal(err)
		}
		i = i + 2
		tcpConn.Write([]byte(strconv.Itoa(i)))
	}
}

// client
func TestClient(t *testing.T) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", tunnelAddr)
	if err != nil {
		t.Fatal(err)
	}
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	// 写数据
	tcpConn.Write([]byte("1500"))
	// 读取数据
	b := make([]byte, 0)
	for {
		var buf [BufSize]byte
		n, err := tcpConn.Read(buf[:])
		if err != nil {
			t.Fatal(err)
		}
		b = append(b, buf[:n]...)
		if n < BufSize {
			break
		}
	}
	fmt.Println(string(b))
}

// tunnel
func TestTunnel(t *testing.T) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", tunnelAddr)
	if err != nil {
		t.Fatal(err)
	}
	tcpListen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		t.Fatal(err)
	}
	for {
		// client tcp Conn
		clientTcpConn, err := tcpListen.AcceptTCP()
		if err != nil {
			t.Fatal(err)
		}

		// 获取用户传递过来的数据
		//b := make([]byte, 0)
		//for {
		//	var buf [BufSize]byte
		//	n, err := clientTcpConn.Read(buf[:])
		//	if err != nil {
		//		t.Fatal(err)
		//	}
		//	b = append(b, buf[:n]...)
		//	if n < BufSize {
		//		break
		//	}
		//}

		// 与服务端创建连接
		addr, err := net.ResolveTCPAddr("tcp", serverAddr)
		if err != nil {
			t.Fatal(err)
		}
		serverTcpConn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			t.Fatal(err)
		}

		//serverTcpConn.Write(b)

		// 获取服务端响应过来的数据
		//b2 := make([]byte, 0)
		//for {
		//	var buf [BufSize]byte
		//	n, err := serverTcpConn.Read(buf[:])
		//	if err != nil {
		//		t.Fatal(err)
		//	}
		//	b2 = append(b2, buf[:n]...)
		//	if n < BufSize {
		//		break
		//	}
		//}
		//
		//clientTcpConn.Write(b2)

		go io.Copy(serverTcpConn, clientTcpConn)
		go io.Copy(clientTcpConn, serverTcpConn)
	}
}
