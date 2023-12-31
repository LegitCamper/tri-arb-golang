package platforms

import (
	"github.com/gorilla/websocket"
	"net/url"
)

type Platform struct {
	Sandbox       bool
	WebsocketHost WebsocketHost
	RestHost      RestHost
	// Coins
	// Pairs
	// stables
}
type WebsocketHost struct {
	User       string
	UserPath   string
	UserSubs   []string
	Market     string
	MarketPath string
	MarketSubs []string
	Scheme     string
}
type RestHost struct {
	SandboxApi string
	Api        string
	Scheme     string
	// limits
}

// Ensure my broker structs are passble
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
	// Websocket functions
	GetPlatform() Platform
	Decode([]byte) Response
	Encode(Request) []byte
	Ping(Response) bool
	PongMessage(int) []byte
	PongHandler(int, *websocket.Conn)
	SubscriptionMessage([]string) []byte

	// Rest functions
	MakeUrl(string) url.URL
	DownloadSymbols() []string
}
