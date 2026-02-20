package shelly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type gen3SwitchStatus struct {
	Output bool    `json:"output"`
	APower float64 `json:"apower"`
}

func getStatusGen3(client *http.Client, ip string) (Status, error) {
	resp, err := client.Get(fmt.Sprintf("http://%s/rpc/Switch.GetStatus?id=0", ip))
	if err != nil {
		return Status{}, fmt.Errorf("Gen3 status: %w", err)
	}
	defer resp.Body.Close()

	var s gen3SwitchStatus
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return Status{}, fmt.Errorf("Gen3 status decode: %w", err)
	}
	return Status{On: s.Output, Watt: s.APower}, nil
}

func setRelayGen3(client *http.Client, ip string, on bool) error {
	body, _ := json.Marshal(map[string]interface{}{"id": 0, "on": on})
	resp, err := client.Post(
		fmt.Sprintf("http://%s/rpc/Switch.Set", ip),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("Gen3 relay: %w", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Gen3 relay: HTTP %d", resp.StatusCode)
	}
	return nil
}
