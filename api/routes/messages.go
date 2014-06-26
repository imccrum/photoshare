package routes

import (
	"encoding/json"
	"github.com/igm/pubsub"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
)

var pub pubsub.Publisher

type Message struct {
	UserName string `json:"username"`
	PhotoID  int64  `json:"photoID"`
	Type     string `json:"type"`
}

func sendMessage(msg *Message) {
	pub.Publish(msg)
}

func messageHandler(session sockjs.Session) {
	var closedSession = make(chan struct{})
	go func() {
		reader, _ := pub.SubChannel(nil)
		for {
			select {
			case <-closedSession:
				return
			case msg := <-reader:
				msg = msg.(*Message)
				if body, err := json.Marshal(msg); err == nil {
					if err = session.Send(string(body)); err != nil {
						log.Println(err)
						return
					}
				}
			}
		}
	}()
	for {
		if msg, err := session.Recv(); err == nil {
			pub.Publish(msg)
			continue
		}
		break
	}
	close(closedSession)
}
