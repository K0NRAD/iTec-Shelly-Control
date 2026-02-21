package shelly

import "net/http"

// Gen4 (z.B. Shelly 1PM Gen4) nutzt dieselbe RPC-API wie Gen3.
func getStatusGen4(client *http.Client, ip string) (Status, error) {
	return getStatusGen3(client, ip)
}

func setRelayGen4(client *http.Client, ip string, on bool) error {
	return setRelayGen3(client, ip, on)
}
