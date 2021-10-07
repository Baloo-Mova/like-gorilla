package server

import (
	"bitmex-gorilla/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)

type Broadcaster struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client

	Data chan models.PriceInfo
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (b *Broadcaster) BroadcastData(message models.PriceInfo) {

	for client := range b.clients {
		for _, item := range client.SymbolsSubscribed {
			if item == message.Symbol {
				client.Write(message)
			}
		}
	}
}

func (b *Broadcaster) Start() {
	for {
		select {
		case client := <-b.register:
			b.clients[client] = true
			break

		case client := <-b.unregister:
			_, ok := b.clients[client]
			if ok {
				delete(b.clients, client)
			}
			break

		case message := <- b.Data:
			b.BroadcastData(message)
		}
	}
}

func (b *Broadcaster) NewClient(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		conn: ws,
	}

	b.register <- client

	client.Read()

	b.unregister <- client
}

func NewBroadcaster(Data chan models.PriceInfo) *Broadcaster {
	return &Broadcaster{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Data: Data,
	}
}
