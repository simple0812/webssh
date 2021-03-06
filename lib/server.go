package lib

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/googollee/go-socket.io"
)

func InitServer() *socketio.Server {
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
			go client.Connect(host, user, pass, func(client *Client, err error) {
				if err != nil {
					fmt.Println(err.Error())
					so.Emit("conn", false)
					return
				}
				so.Emit("conn", true)
				go func() {
					for {
						if !client.IsConnected() {
							break
						}
						time.Sleep(200 * time.Millisecond)
						so.Emit("cmd", client.GetOutFile())
					}
				}()
			})
		})

		so.On("cmd", func(msg string) {
			if client.IsConnected() {
				if msg == "quit" {
					client.DisConnect()
					so.Emit("cmd", "")
					return
				}

				client.SendCmd(msg)
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

	return server
}
