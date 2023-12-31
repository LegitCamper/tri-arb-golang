package platforms

import (
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// Create one of a few platforms
func platform_websocket(b Broker, u url.URL, subscriptions []string, channel chan<- WebsocketResponse) { //*websocket.Conn
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

// platform_websocket Handler
func Handler(b Broker) PlatformApi {
	websocket_host := b.GetPlatform().WebsocketHost

	user_conn := make(chan WebsocketResponse)
	market_conn := make(chan WebsocketResponse)

	go platform_websocket(b, url.URL{Scheme: websocket_host.Scheme, Host: websocket_host.User, Path: websocket_host.UserPath}, websocket_host.UserSubs, user_conn)
	go platform_websocket(b, url.URL{Scheme: websocket_host.Scheme, Host: websocket_host.Market, Path: websocket_host.MarketPath}, websocket_host.MarketSubs, market_conn)

	platform := PlatformApi{
		User_conn:   user_conn,
		Market_conn: market_conn,
		Platform:    Platform{WebsocketHost: websocket_host},
		MarketData:  PlatformMarketData{Symbol: make(map[string]*PlatformMarketDataSymbol)},
	}

	return platform
}
