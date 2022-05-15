package test

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"nat-pernetration/helper"
	"net"
	"sync"
	"testing"
)

const (
	ControlServerAddr = "0.0.0.0:8080"
	RequestServerAddr = "0.0.0.0:8081"
	KeepAliveStr      = "KeepAlive\n"
)

var wg sync.WaitGroup
var clientConn *net.TCPConn

// 服务端
func TestUserServer(t *testing.T) {
	wg.Add(1)
	// 监听控制中心
	go ControlServer()
	// 监听用户的请求
	go RequestServer()
	wg.Wait()
}

func ControlServer() {
	tcpListener, err := helper.CreateListen(ControlServerAddr)
	if err != nil {
		panic(err)
	}
	log.Printf("ControlServer is listening on %s\n", ControlServerAddr)
	for {
		clientConn, err = tcpListener.AcceptTCP()
		if err != nil {
			return
		}
		go helper.KeepAlive(clientConn)
	}
}

func RequestServer() {
	tcpListener, err := helper.CreateListen(RequestServerAddr)
	if err != nil {
		panic(err)
	}
	log.Printf("RequestServer is listening on %s\n", RequestServerAddr)
	for {
		conn, err := tcpListener.AcceptTCP()
		if err != nil {
			return
		}
		go io.Copy(clientConn, conn)
		go io.Copy(conn, clientConn)
	}
}

// 客户端
func TestUserClient(t *testing.T) {
	conn, err := helper.CreateConn(ControlServerAddr)
	if err != nil {
		log.Printf("[连接失败] %s", err)
	}
	for {
		s, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("Get Data Error:%v", err)
			continue
		}
		log.Printf("Get Data: %v", s)
		_, err = conn.Write([]byte("I Get\n"))
		if err != nil {
			log.Printf("Send Data Error:%v", err)
		}
	}
}

// 用户端
func TestUserRequestClient(t *testing.T) {
	conn, err := helper.CreateConn(RequestServerAddr)
	if err != nil {
		log.Printf("[连接失败] %s", err)
	}
	_, err = conn.Write([]byte("中文 \n"))
	if err != nil {
		log.Printf("[发送失败] %s", err)
	}
	s, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Printf("[接收失败] %s", err)
	}
	fmt.Println(s)
}
