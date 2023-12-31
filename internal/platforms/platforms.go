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

type WebsocketResponse interface {
	GetMethod() string
	GetId() int
	// this is a list and could be remvoing data
	GetMarketData() interface{} // this works, but you have to coerce the struct
}
type WebsocketRequest interface {
	ToJson() []byte
	Timestamp() int64
	AddTimestamp()
}

type RestResponse interface{}
type RestRequest interface {
	ToJson() []byte
}

type PlatformApi struct {
	User_conn   chan WebsocketResponse
	Market_conn chan WebsocketResponse
	Platform    Platform
	MarketData  PlatformMarketData
}
type PlatformMarketData struct {
	Symbol map[string]*PlatformMarketDataSymbol
}
type PlatformMarketDataSymbol struct {
	Asks [][2]float64
	Bids [][2]float64
	Id   int
}

// Interface to allow several Platform types (eg. Crypto.com)
type Broker interface {
	// Websocket functions
	GetPlatform() Platform
	Decode([]byte) WebsocketResponse
	Encode(WebsocketRequest) []byte
	Ping(WebsocketResponse) bool
	PongMessage(int) []byte
	PongHandler(int, *websocket.Conn)
	SubscriptionMessage([]string) []byte
	ProccessMarketData(*PlatformApi)

	// Rest functions
	MakeUrl(string) url.URL
	DownloadSymbols() []string
	DecodeSymbols(string)
}
