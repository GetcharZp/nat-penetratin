package main

import (
	"encoding/json"
	"log"
	"nat-pernetration/define"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		q := request.URL.Query()
		b, err := json.Marshal(q)
		if err != nil {
			log.Printf("Marshal Error:%v", err)
		}
		writer.Write(b)
	})
	log.Println("本地服务已启动 ", define.LocalServerAddr)
	http.ListenAndServe(define.LocalServerAddr, nil)
}
