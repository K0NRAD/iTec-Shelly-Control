package shelly

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type gen1Status struct {
	Relays []struct {
		IsOn bool `json:"ison"`
	} `json:"relays"`
	Meters []struct {
		Power float64 `json:"power"`
	} `json:"meters"`
}

func getStatusGen1(client *http.Client, ip string) (Status, error) {
	resp, err := client.Get(fmt.Sprintf("http://%s/status", ip))
	if err != nil {
		return Status{}, fmt.Errorf("Gen1 status: %w", err)
	}
	defer resp.Body.Close()

	var s gen1Status
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return Status{}, fmt.Errorf("Gen1 status decode: %w", err)
	}
	if len(s.Relays) == 0 || len(s.Meters) == 0 {
		return Status{}, fmt.Errorf("Gen1 status: unerwartete Antwortstruktur")
	}
	return Status{On: s.Relays[0].IsOn, Watt: s.Meters[0].Power}, nil
}

func setRelayGen1(client *http.Client, ip string, on bool) error {
	state := "off"
	if on {
		state = "on"
	}
	resp, err := client.Get(fmt.Sprintf("http://%s/relay/0?turn=%s", ip, state))
	if err != nil {
		return fmt.Errorf("Gen1 relay: %w", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Gen1 relay: HTTP %d", resp.StatusCode)
	}
	return nil
}
