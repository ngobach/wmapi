package server

import (
	"fmt"
	"github.com/ngobach/wmapi/config"
	"net/http"
)

type MyHandler struct{}

func (m *MyHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	panic("implement me")
}

func configServer() {
	mainHandler := &MyHandler{}
	http.Handle("/*", mainHandler)
}

func StartServer() error {
	configServer()
	addr := fmt.Sprintf("0.0.0.0:%d", config.Port)
	fmt.Println("Server will listen at", addr)
	return http.ListenAndServe(addr, nil)
}
