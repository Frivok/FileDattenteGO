package main

import "github.com/gorilla/websocket"

func Dial() {
	websocket.DefaultDialer.Dial("ws://localhost:4000/ws", nil)
}

type Consommateur interface {
	Commencer() error
}
