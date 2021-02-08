package main

import (
	"github.com/gorilla/websocket"
	"learn-websocket/src/core"
	"learn-websocket/src/handlers"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/echo", handlers.Echo)

	http.HandleFunc("/sendall", func(w http.ResponseWriter, req *http.Request) {
		msg:=req.URL.Query().Get("msg")
		core.ClientMap.SendAll(msg)
		w.Write([]byte("OK"))
	})


	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
