package helper

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"nat-pernetration/conf"
	"nat-pernetration/define"
	"net"
	"time"
)

type UserClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// CreateListen 监听
func CreateListen(serverAddr string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		return nil, err
	}
	tcpListen, err := net.ListenTCP("tcp", tcpAddr)
	return tcpListen, err
}

// CreateConn 创建连接
func CreateConn(serverAddr string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	return conn, err
}

// KeepAlive 设置连接保活
func KeepAlive(conn *net.TCPConn) {
	for {
		_, err := conn.Write([]byte(define.KeepAliveStr))
		if err != nil {
			log.Printf("[KeepAlive] Error %s", err)
			return
		}
		time.Sleep(time.Second * 3)
	}
}

// GetDataFromConnection 获取Connection中的数据
func GetDataFromConnection(bufSize int, conn *net.TCPConn) ([]byte, error) {
	b := make([]byte, 0)
	for {
		buf := make([]byte, bufSize)
		n, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}
		b = append(b, buf[:n]...)
		if n < bufSize {
			break
		}
	}
	return b, nil
}

// GetServerConf 解析 server.yaml
func GetServerConf() (*conf.Server, error) {
	s := new(conf.Server)
	b, err := ioutil.ReadFile("./conf/server.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, s)
	if err != nil {
		return nil, err
	}
	return s, err
}

var myKey = []byte("nat-penetration-key")

// GenerateToken
// 生成 token
func GenerateToken(name string) (string, error) {
	UserClaim := &UserClaims{
		Username:       name,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyseToken
// 解析 token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token Error:%v", err)
	}
	return userClaim, nil
}
