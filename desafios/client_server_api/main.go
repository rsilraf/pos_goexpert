package main

import (
	"time"

	. "github.com/rsilraf/pos_goexpert/desafios/client_server_api/client"
	. "github.com/rsilraf/pos_goexpert/desafios/client_server_api/server"
)

func main() {
	go Server()
	<-time.After(1 * time.Second)
	Client()
	<-time.After(3 * time.Second)
}
