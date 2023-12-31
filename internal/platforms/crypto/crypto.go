package crypto

// https://exchange-docs.crypto.com/exchange/v1/rest-ws/index.html

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"tri-arb/internal/platforms"
)

type Request struct {
	Id      int            `json:"id"`
	Method  string         `json:"method"`
	Params  *RequestParams `json:"params,omitempty"`
	Api_key string         `json:"spi_key,omitempty"`
	Sig     string         `json:"sig,omitempty"`
	Nonce   int64          `json:"nonce,omitempty"`
}
type RequestParams struct {
	Channels             []string `json:"channels,omitempty"`
	BookSubscriptionType string   `json:"book_subscription_type,omitempty"`
	BookUpdateFrequency  int      `json:"book_update_frequency,omitempty"`
}

func (r Request) ToJson() []byte {
	s, err := json.Marshal(&r)
	if err != nil {
		log.Println("error:", err)
	}
	return s
}
func (r Request) Timestamp() int64 {
	return time.Now().UnixMilli()
}
func (r Request) AddTimestamp() {
	r.Nonce = r.Timestamp()
}

type Response struct {
	Id       int              `json:"id"`
	Method   string           `json:"method"`
	Result   *ResponseResults `json:"reslult,omitempty"`
	Code     int              `json:"code,omitempty"`
	Message  string           `json:"message,omitempty"`
	Original string           `json:"original,omitempty"`
}
type ResponseResults struct {
	Instrument_name string                 `json:"instrument_name,omitempty"`
	Subscription    string                 `json:"subscription,omitempty"`
	Channel         string                 `json:"channel,omitempty"`
	Depth           uint                   `json:"depth,omitempty"`
	Data            []*ResponseResultsData `json:"data"`
}
type ResponseResultsData struct {
	Asks [][3]float32 `json:"asks,omitempty"`
	Bids [][3]float32 `json:"bids,omitempty"`
	T    int          `json:"t,omitempty"`
	Tt   int          `json:"tt,omitempty"`
	U    int          `json:"U,omitempty"`
}

func (r Response) GetMethod() string {
	return r.Method
}

func (r Response) GetId() int {
	return r.Id
}

type Crypto platforms.Platform

func New(sandbox bool) Crypto {
	var host platforms.Host
	if !sandbox {
		host = platforms.Host{
			User:       "stream.crypto.com",
			UserPath:   "/exchange/v1/user",
			UserSubs:   []string{""},
			Market:     "stream.crypto.com",
			MarketPath: "/exchange/v1/market",
			MarketSubs: []string{"book.BTCUSD-PERP.50"},
			Scheme:     "wss",
		}
	} else {
		host = platforms.Host{
			User:       "uat-stream.3ona.co",
			UserPath:   "/exchange/v1/user",
			UserSubs:   []string{""},
			Market:     "uat-stream.3ona.co",
			MarketPath: "/exchange/v1/market",
			MarketSubs: []string{"book.BTCUSD-PERP.50"},
			Scheme:     "wss",
		}
	}

	return Crypto{Host: host}
}

func (c Crypto) GetPlatform() platforms.Platform {
	return platforms.Platform{Host: c.Host}
}

func (c Crypto) Encode(r platforms.Request) []byte {
	return []byte(r.ToJson())
}

func (c Crypto) Decode(b []byte) platforms.Response {
	var response Response
	err := json.Unmarshal([]byte(b), &response)
	if err != nil {
		log.Println("error:", err)
	}
	return response
}

func (c Crypto) Ping(r platforms.Response) bool {
	return r.GetMethod() == "public/heartbeat"
}

func (c Crypto) PongMessage(id int) []byte {
	return Request{Id: id, Method: "public/respond-heartbeat"}.ToJson()
}

func (c Crypto) PongHandler(id int, ws *websocket.Conn) {
	err := ws.WriteMessage(websocket.TextMessage, c.PongMessage(id))
	if err != nil {
		log.Println("write:", err)
		panic("")
	}
}

func (c Crypto) SubscriptionMessage(list []string) []byte {
	return Request{
		Id:     -1,
		Method: "subscribe",
		Params: &RequestParams{
			Channels:             list,
			BookSubscriptionType: "SNAPSHOT",
			BookUpdateFrequency:  100,
		},
	}.ToJson()
}
