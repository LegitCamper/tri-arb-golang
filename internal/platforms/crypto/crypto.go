package crypto

// https://exchange-docs.crypto.com/exchange/v1/rest-ws/index.html

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"tri-arb/internal/platforms"
)

// Websocket implementations

type WebsocketRequest struct {
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

func (r WebsocketRequest) ToJson() []byte {
	s, err := json.Marshal(&r)
	if err != nil {
		log.Println("error:", err)
	}
	return s
}
func (r WebsocketRequest) Timestamp() int64 {
	return time.Now().UnixMilli()
}
func (r WebsocketRequest) AddTimestamp() {
	r.Nonce = r.Timestamp()
}

type WebsocketResponse struct {
	Id       int                      `json:"id"`
	Method   string                   `json:"method"`
	Result   WebsocketResponseResults `json:"result,omitempty"`
	Code     int                      `json:"code,omitempty"`
	Message  string                   `json:"message,omitempty"`
	Original string                   `json:"original,omitempty"`
}
type WebsocketResponseResults struct {
	Instrument_name string                         `json:"instrument_name,omitempty"`
	Subscription    string                         `json:"subscription,omitempty"`
	Channel         string                         `json:"channel,omitempty"`
	Depth           uint                           `json:"depth,omitempty"`
	Data            []WebsocketResponseResultsData `json:"data"`
}
type WebsocketResponseResultsData struct {
	Update *WebsocketResponseResultsData `json:"update"`
	Asks   [][3]string                   `json:"asks"`
	Bids   [][3]string                   `json:"bids"`
	T      int                           `json:"t"`
	Tt     int                           `json:"tt"`
	U      int                           `json:"u"`
	Pu     int                           `json:"pu,omitempty"`

	// Not in actual response - just for simplicity
	Symbol string
}

func (r WebsocketResponse) GetMethod() string {
	return r.Method
}

func (r WebsocketResponse) GetId() int {
	return r.Id
}

func (r WebsocketResponse) GetMarketData() interface{} {
	if len(r.Result.Data) >= 1 {
		// this is a list!!!! but seems to only ever have 1 element
		r.Result.Data[0].Symbol = r.Result.Instrument_name
		return r.Result.Data[0]
	}
	return nil
}

type Crypto platforms.Platform

func New(sandbox bool) Crypto {
	var websocket_host platforms.WebsocketHost
	var rest_host platforms.RestHost
	if !sandbox {
		websocket_host = platforms.WebsocketHost{
			User:       "stream.crypto.com",
			UserPath:   "/exchange/v1/user",
			UserSubs:   []string{""},
			Market:     "stream.crypto.com",
			MarketPath: "/exchange/v1/market",
			MarketSubs: []string{"book.BTCUSD-PERP.50"},
			Scheme:     "wss",
		}
		rest_host = platforms.RestHost{Scheme: "https://", Api: "api.crypto.com"}
	} else {
		websocket_host = platforms.WebsocketHost{
			User:       "uat-stream.3ona.co",
			UserPath:   "/exchange/v1/user",
			UserSubs:   []string{""},
			Market:     "uat-stream.3ona.co",
			MarketPath: "/exchange/v1/market",
			MarketSubs: []string{"book.BTCUSD-PERP.50"},
			Scheme:     "wss",
		}
		rest_host = platforms.RestHost{Scheme: "https://", Api: "uat-api.3ona.co"}
	}

	return Crypto{
		Sandbox:       sandbox,
		WebsocketHost: websocket_host,
		RestHost:      rest_host,
	}
}

func (c Crypto) GetPlatform() platforms.Platform {
	return platforms.Platform{
		WebsocketHost: c.WebsocketHost,
		RestHost:      c.RestHost,
	}
}

func (c Crypto) Encode(r platforms.WebsocketRequest) []byte {
	return []byte(r.ToJson())
}

func (c Crypto) Decode(b []byte) platforms.WebsocketResponse {
	var response WebsocketResponse
	err := json.Unmarshal([]byte(b), &response)
	if err != nil {
		log.Println("error:", err)
	}
	return response
}

func (c Crypto) Ping(r platforms.WebsocketResponse) bool {
	return r.GetMethod() == "public/heartbeat"
}

func (c Crypto) PongMessage(id int) []byte {
	return WebsocketRequest{Id: id, Method: "public/respond-heartbeat"}.ToJson()
}

func (c Crypto) PongHandler(id int, ws *websocket.Conn) {
	err := ws.WriteMessage(websocket.TextMessage, c.PongMessage(id))
	if err != nil {
		log.Println("write:", err)
		panic("")
	}
}

func (c Crypto) SubscriptionMessage(list []string) []byte {
	return WebsocketRequest{
		Id:     -1,
		Method: "subscribe",
		Params: &RequestParams{
			Channels:             list,
			BookSubscriptionType: "SNAPSHOT_AND_UPDATE",
			BookUpdateFrequency:  10,
		},
	}.ToJson()
}

// Rest implementations

type restRequest struct {
	Id     int         `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params,omitempty"`
	ApiKey string      `json:"api_key,omitempty"`
	Sig    string      `json:"sig,omitempty"`
	Nonce  int         `json:"nonce"`
}

func (r restRequest) ToJson() []byte {
	s, err := json.Marshal(&r)
	if err != nil {
		log.Println("error:", err)
	}
	return s
}

type restResponse struct {
	Result restResponseResult `json:"result,omitempty"`
}
type restResponseResult struct {
	Currencymap map[string]*restResponseCurrencyMap `json:"currency_map,omitempty"`
}
type restResponseCurrencyMap struct {
	Fullname    string             `json:"full_name,omitempty"`
	NetworkList []*restNetworkList `json:"network_list,omitempty"`
}
type restNetworkList struct {
	Networkid string `json:"network_id,omitempty"`
	// WithdrawalFee        null   `json:"withdrawal_fee,omitempty"`
	WithdrawEnabled      bool    `json:"withdraw_endabled,omitempty"`
	MinWithdrawalAmount  float64 `json:"min_withdrawal_amount,omitempty"`
	DepositEnabled       bool    `json:"deposit_enabled,omitempty"`
	ConfirmationRequired int     `json:"confirmation_required ,omitempty"`
}

func (c Crypto) MakeUrl(s string) url.URL {
	return url.URL{Scheme: c.RestHost.Scheme, Host: c.RestHost.Api, Path: s}
}

func (c Crypto) DownloadSymbols() []string {
	log.Println(platforms.RestGet(c, "/exchange/v1/private/get-currency-networks", restRequest{}))
	return []string{""}
}
func (c Crypto) DecodeSymbols(s string) platforms.RestResponse {
	var response restResponse
	{
		err := json.Unmarshal([]byte(s), &response)
		if err != nil {
			log.Println("error:", err)
		}
	}
	return response
}

func (c Crypto) ProccessMarketData(platform *platforms.PlatformApi) {
	// Runs for the lifetime of the program
	for data := range platform.Market_conn {
		new_market_data := data.GetMarketData()
		if new_market_data == nil {
			return
		}
		r, ok := new_market_data.(WebsocketResponseResultsData) // Coerce interface back into ResponseResultsData
		if !ok {
			log.Println("Error: invalid type assertion")
			panic("")
		}
		symbol := r.Symbol
		// There are two types of data here snapshots and updates
		if r.Update == nil {
			// snapshots
			market_data := *platform.MarketData.Symbol[symbol]
			for _, e := range r.Asks {
				float0, _ := strconv.ParseFloat(strings.TrimSpace(e[0]), 64)
				float1, _ := strconv.ParseFloat(strings.TrimSpace(e[1]), 64)
				new_e := [2]float64{float0, float1}
				market_data.Asks = append(market_data.Asks, new_e)
			}
			for _, e := range r.Bids {
				float0, _ := strconv.ParseFloat(strings.TrimSpace(e[0]), 64)
				float1, _ := strconv.ParseFloat(strings.TrimSpace(e[1]), 64)
				new_e := [2]float64{float0, float1}
				market_data.Bids = append(market_data.Bids, new_e)
			}
			market_data.Id = r.U
		} else {
			// updates
			if r.Pu == platform.MarketData.Symbol[symbol].Id {
				market_data := *platform.MarketData.Symbol[symbol]
				for _, e := range r.Update.Asks {
					float0, _ := strconv.ParseFloat(strings.TrimSpace(e[0]), 64)
					float1, _ := strconv.ParseFloat(strings.TrimSpace(e[1]), 64)
					new_e := [2]float64{float0, float1}
					market_data.Asks = append(market_data.Asks, new_e)
				}
				for _, e := range r.Update.Bids {
					float0, _ := strconv.ParseFloat(strings.TrimSpace(e[0]), 64)
					float1, _ := strconv.ParseFloat(strings.TrimSpace(e[1]), 64)
					new_e := [2]float64{float0, float1}
					market_data.Bids = append(market_data.Bids, new_e)
				}
				market_data.Id = r.U
			} else {
				// TODO: download book via rest here
				log.Println("Existing id and new id dont match. Removing broken book")
				platform.MarketData = platforms.PlatformMarketData{Symbol: make(map[string]*platforms.PlatformMarketDataSymbol)}
			}
		}
	}
	panic("Websocket Market Channel closed")
}
