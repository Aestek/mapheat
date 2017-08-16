package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

func server(addr string, evts <-chan *Evt) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	var connectionsLock sync.RWMutex
	connections := map[*websocket.Conn]struct{}{}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		connectionsLock.Lock()
		connections[conn] = struct{}{}
		connectionsLock.Unlock()

		conn.SetCloseHandler(func(code int, msg string) error {
			log.Printf("client closed with (%d) %s\n", code, msg)

			connectionsLock.Lock()
			delete(connections, conn)
			connectionsLock.Unlock()

			return nil
		})
	})

	go func() {
		for evt := range evts {
			bytes, _ := json.Marshal(evt)

			for conn := range connections {
				err := conn.WriteMessage(websocket.TextMessage, bytes)
				if err != nil {
					log.Printf("err writting msg: %s\n", err)
					connectionsLock.Lock()
					delete(connections, conn)
					connectionsLock.Unlock()
				}
			}
		}
	}()

	http.Handle("/", http.FileServer(http.Dir("public/")))

	log.Println("serving on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
