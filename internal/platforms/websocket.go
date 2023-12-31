package platforms

import (
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// Create one of a few platforms
func platform_websocket(b Broker, u url.URL, subscriptions []string, channel chan<- Response) { //*websocket.Conn
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	time.Sleep(time.Second)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer c.Close()

	// subscribe
	{
		err := c.WriteMessage(websocket.TextMessage, b.SubscriptionMessage(subscriptions))
		if err != nil {
			log.Println("write:", err)
			return
		}
	}

	for {
		message_type, message, err := c.ReadMessage()
		if err != nil {
			log.Println("ReadMessage() error:", err)
			panic("")
		}
		log.Printf("got: %s", message) // TODO: Remove

		if message_type == websocket.TextMessage {
			d_message := b.Decode(message)
			// If ping is seen handle pong
			if b.Ping(d_message) {
				b.PongHandler(d_message.GetId(), c)
			} else {
				channel <- d_message
			}
		}
	}
}

type Platform_api struct {
	User_conn   chan Response
	Market_conn chan Response
	Platform    Platform
}

// platform_websocket Handler
func Handler(b Broker) Platform_api {
	host := b.GetPlatform().Host

	user_conn := make(chan Response)
	market_conn := make(chan Response)

	go platform_websocket(b, url.URL{Scheme: host.Scheme, Host: host.User, Path: host.UserPath}, host.UserSubs, user_conn)
	go platform_websocket(b, url.URL{Scheme: host.Scheme, Host: host.Market, Path: host.MarketPath}, host.MarketSubs, market_conn)

	return Platform_api{
		User_conn:   user_conn,
		Market_conn: market_conn,
		Platform:    Platform{Host: host},
	}
}
