package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	. "webssh/lib"

	"webssh/Godeps/_workspace/src/github.com/googollee/go-socket.io"
)

func main() {

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		client := NewClient()
		so.On("conn", func(msg string) {
			p := strings.Split(msg, "|")
			if len(p) < 3 {
				so.Emit("conn", false)
				return
			}
			user := p[0]
			pass := p[1]
			host := p[2]
			go client.Connect(host, user, pass)
			startTime := time.Now()

			for !client.IsConnected() {
				if time.Now().Sub(startTime) > time.Second {
					fmt.Println("conn timeout")
					so.Emit("conn", false)
					return
				}
			}

			so.Emit("conn", true)
		})

		so.On("cmd", func(msg string) {
			if client.IsConnected() {
				so.Emit("cmd", client.SendCmd(msg))
			}
		})

		so.On("disconnection", func() {
			if client.IsConnected() {
				client.DisConnect()
			}
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:3003...")
	log.Fatal(http.ListenAndServe(":3003", nil))
}
