package shelly

import (
	"fmt"
	"net/http"
	"time"
)

type Status struct {
	On   bool
	Watt float64
}

type Client struct {
	http *http.Client
}

func NewClient() *Client {
	return &Client{
		http: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *Client) GetStatus(ip string, generation int) (Status, error) {
	switch generation {
	case 1:
		return getStatusGen1(c.http, ip)
	case 3:
		return getStatusGen3(c.http, ip)
	case 4:
		return getStatusGen4(c.http, ip)
	default:
		return Status{}, fmt.Errorf("unbekannte Generation: %d", generation)
	}
}

func (c *Client) SetRelay(ip string, generation int, on bool) error {
	switch generation {
	case 1:
		return setRelayGen1(c.http, ip, on)
	case 3:
		return setRelayGen3(c.http, ip, on)
	case 4:
		return setRelayGen4(c.http, ip, on)
	default:
		return fmt.Errorf("unbekannte Generation: %d", generation)
	}
}

func (c *Client) Test(ip string, generation int) error {
	_, err := c.GetStatus(ip, generation)
	return err
}
