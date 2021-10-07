package bitmex

import (
	"bitmex-gorilla/models"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Client struct {
	WssEndpoint string
	PingTimer   *time.Timer
	DataStream  chan models.PriceInfo
	connection  *websocket.Conn
	pingTest    bool
}

func NewClient(WssEndpoint string) *Client {
	return &Client{
		WssEndpoint: WssEndpoint,
		DataStream:  make(chan models.PriceInfo, 40),
	}
}

func (client *Client) Start() {
	var err error
	client.connection, _, err = websocket.DefaultDialer.Dial(client.WssEndpoint, nil)
	if err != nil {
		log.Println("dial:", err)
	}
	client.pingTest = false
	client.startTimer()

  	log.Println("Connected")
	go client.listener()

	client.connection.WriteJSON(models.BitmexRequestModel{
		Op:   "subscribe",
		Args: []string{"instrument"},
	})

	log.Println("Subscribe")

}

func (client *Client) startTimer() {
	if client.PingTimer == nil {
		client.PingTimer = time.NewTimer(5 * time.Second)
		return
	}

	client.PingTimer.Reset(5 * time.Second)
}

func (client *Client) listener() {
	for {
		var data models.BitmexResponseModel
		err := client.connection.ReadJSON(&data)
		if err != nil {
			log.Println(err)
			return
		}

		client.PingTimer.Reset(5 * time.Second)
		client.pingTest = false

		for _, item := range data.Data {
			if item.LastPrice != 0 {
				client.DataStream <- models.PriceInfo{
					TimeStamp: item.Timestamp.String(),
					Symbol:    item.Symbol,
					Price:     item.LastPrice,
				}
			}
		}
	}
}

func (client *Client) restart() {
	for {
		select {
		case <-client.PingTimer.C:
			if client.pingTest {
				log.Println(time.Now(), "No connection Restart WSS to Bitmex")
				client.finish()
				client.Start()
			}

			client.pingTest = true
			client.PingTimer.Reset(5 * time.Second)
		}
	}
}

func (client *Client) finish() {

	//TODO: Need to close connection more gracefully

	client.PingTimer.Stop()
	client.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	client.connection.Close()
}

func (client *Client) Stop() {
	client.finish()
}
