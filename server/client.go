package server

import (
	"bitmex-gorilla/models"
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Client struct {
	SymbolsSubscribed []string
	conn              *websocket.Conn
}

func (c *Client) Read() {
	defer c.conn.Close()

	for {
		var request models.SubscriptionStatusRequest
		_, message, err := c.conn.ReadMessage()

		if string(message) == "PING" {
			c.conn.WriteMessage(websocket.TextMessage, []byte("PONG"))
			continue
		}

		json.Unmarshal(message, &request)
		if err != nil {
			break
		}

		switch request.Action {
		case "subscribe":
			c.SymbolsSubscribed = request.Symbols
			break
		case "unsubscribe":
			c.SymbolsSubscribed = nil
		default:
			c.conn.WriteMessage(websocket.TextMessage, []byte("Undefined command"))
		}
	}
}

func (c *Client) Write(message models.PriceInfo) error {
	return c.conn.WriteJSON(message)
}
