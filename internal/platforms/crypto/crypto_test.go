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
	var output Response
	err := json.Unmarshal([]byte(heartbeat), &output)
	if err != nil {
		t.Fatalf("Failed to decode heartbeat: %s", err)
	}
}

func TestHeartbeatRequest(t *testing.T) {
	heartbeat := Request{Id: 5000, Method: "public/heartbeat"}
	_, err := json.Marshal(heartbeat)
	if err != nil {
		t.Fatalf("Failed to encode heartbeat: %s", err)
	}
}
