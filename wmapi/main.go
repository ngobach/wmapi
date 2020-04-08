package main

import (
	"github.com/ngobach/wmapi/server"
)

func main() {
	err := server.StartServer()
	panic(err)
}
