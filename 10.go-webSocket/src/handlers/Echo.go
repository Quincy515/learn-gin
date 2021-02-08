package handlers

import (
	"learn-websocket/src/core"
	"log"
	"net/http"
)

func Echo(w http.ResponseWriter, r *http.Request) {
	client, err := core.Upgrader.Upgrade(w, r, nil) // 升级
	if err != nil {
		log.Println()
	} else {
		core.ClientMap.Store(client)
	}
}
