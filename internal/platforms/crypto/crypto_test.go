package crypto

import (
	"encoding/json"
	"testing"
)

func TestHeartbeatResponse(t *testing.T) {
	heartbeat := `{
		"id": 1587523073344,
		"method": "public/heartbeat",
		"code": 0
	}`
	var output WebsocketResponse
	err := json.Unmarshal([]byte(heartbeat), &output)
	if err != nil {
		t.Fatalf("Failed to decode heartbeat: %s", err)
	}
}

func TestHeartbeatRequest(t *testing.T) {
	heartbeat := WebsocketRequest{Id: 5000, Method: "public/heartbeat"}
	_, err := json.Marshal(heartbeat)
	if err != nil {
		t.Fatalf("Failed to encode heartbeat: %s", err)
	}
}

func TestUnmarshalResponse(t *testing.T) {
	response := `{
	  "id": -1,
	  "code": 0,
	  "method": "subscribe",
	  "result": {
	    "instrument_name": "BTCUSD-PERP",
	    "subscription": "book.BTCUSD-PERP.10",
	    "channel": "book",
	    "depth": 10,
	    "data": [{
	      "asks": [
	        ["50126.000000", "0.400000", "2"],
	        ["50130.000000", "1.279000", "3"],
	        ["50136.000000", "1.279000", "5"],
	        ["50137.000000", "0.800000", "7"],
	        ["50142.000000", "1.279000", "1"],
	        ["50148.000000", "2.892900", "9"],
	        ["50154.000000", "1.279000", "5"],
	        ["50160.000000", "1.133000", "2"],
	        ["50166.000000", "3.090700", "1"],
	        ["50172.000000", "1.279000", "1"]
	      ],
	      "bids": [
	        ["50113.500000", "0.400000", "3"],
	        ["50113.000000", "0.051800", "1"],
	        ["50112.000000", "1.455300", "1"],
	        ["50106.000000", "1.174800", "2"],
	        ["50100.500000", "0.800000", "4"],
	        ["50100.000000", "1.455300", "5"],
	        ["50097.500000", "0.048000", "8"],
	        ["50097.000000", "0.148000", "9"],
	        ["50096.500000", "0.399200", "2"],
	        ["50095.000000", "0.399200", "3"]
	      ],
	      "tt": 1647917462799,
	      "t": 1647917463000,
	      "u": 7845460001
	    }]
	  }
	}`
	var output WebsocketResponse
	err := json.Unmarshal([]byte(response), &output)
	if err != nil {
		t.Fatalf("Failed to decode market response: %s", err)
	}
}
