package platforms

import (
	"encoding/json"
	"log"

	"github.com/go-resty/resty/v2"
)

func RestGet(b Broker, s string) string {
	client := resty.New().
		SetJSONUnmarshaler(json.Unmarshal)

	url := b.MakeUrl(s)
	resp, err := client.R().
		EnableTrace().
		Get(url.Scheme + url.Host + url.Path)
	if err != nil {
		log.Println("error:", err)
	}
	log.Printf("resp: %s", resp)
	return resp.String()

	// var response rest
	// err := json.Unmarshal([]byte(resp), &response)
	// if err != nil {
	// 	log.Println("error:", err)
	// }

}
