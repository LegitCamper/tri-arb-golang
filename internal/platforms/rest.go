package platforms

import (
	"encoding/json"
	"log"

	"github.com/go-resty/resty/v2"
)

func RestGet(b Broker, path string, req RestRequest) RestResponse {
	client := resty.New().
		SetJSONUnmarshaler(json.Unmarshal)

	url := b.MakeUrl(path)
	resp, err := client.R().
		EnableTrace().
		Get(url.Scheme + url.Host + url.Path)
	if err != nil {
		log.Println("error:", err)
	}
	log.Printf("resp: %s", resp)

	return b.DecodeSymbols(resp.String())
}
