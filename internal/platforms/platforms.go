package platforms

import "github.com/gorilla/websocket"

type Host struct {
	User       string
	UserPath   string
	UserSubs   []string
	Market     string
	MarketPath string
	MarketSubs []string
	Scheme     string
}

type Platform struct {
	Host Host
	// Coins
	// Pairs
	// stables
}

// ensure my broker structs are passble
type Response interface {
	GetMethod() string
	GetId() int
}
type Request interface {
	ToJson() []byte
	Timestamp() int64
	AddTimestamp()
}

// Interface to allow several Platform types (eg. Crypto.com)
type Broker interface {
	GetPlatform() Platform
	Decode([]byte) Response
	Encode(Request) []byte
	Ping(Response) bool
	PongMessage(int) []byte
	PongHandler(int, *websocket.Conn)
	SubscriptionMessage([]string) []byte
}
