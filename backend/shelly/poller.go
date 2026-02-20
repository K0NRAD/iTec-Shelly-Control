package shelly

import (
	"log"
	"sync"
	"time"

	"shelly-control/store"
)

type DeviceState struct {
	ID   string
	On   bool
	Watt float64
	Err  error
}

type Poller struct {
	client          *Client
	deviceStore     *store.DeviceStore
	powerBuffer     *store.PowerBuffer
	intervalSeconds int

	mu          sync.RWMutex
	states      map[string]DeviceState
	subscribers []chan DeviceState
	stopCh      chan struct{}
}

func NewPoller(client *Client, ds *store.DeviceStore, pb *store.PowerBuffer, intervalSeconds int) *Poller {
	return &Poller{
		client:          client,
		deviceStore:     ds,
		powerBuffer:     pb,
		intervalSeconds: intervalSeconds,
		states:          make(map[string]DeviceState),
		stopCh:          make(chan struct{}),
	}
}

func (p *Poller) Start() {
	go func() {
		p.poll()
		ticker := time.NewTicker(time.Duration(p.intervalSeconds) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				p.poll()
			case <-p.stopCh:
				return
			}
		}
	}()
}

func (p *Poller) Stop() {
	close(p.stopCh)
}

func (p *Poller) poll() {
	devices := p.deviceStore.GetDevices()
	var wg sync.WaitGroup
	for _, d := range devices {
		wg.Add(1)
		go func(dev store.Device) {
			defer wg.Done()
			status, err := p.client.GetStatus(dev.IP, dev.Generation)
			state := DeviceState{ID: dev.ID, Err: err}
			if err == nil {
				state.On = status.On
				state.Watt = status.Watt
				p.powerBuffer.Add(dev.ID, status.Watt)
			} else {
				log.Printf("Poll %s (%s): %v", dev.Name, dev.IP, err)
			}
			p.mu.Lock()
			p.states[dev.ID] = state
			subs := make([]chan DeviceState, len(p.subscribers))
			copy(subs, p.subscribers)
			p.mu.Unlock()

			for _, ch := range subs {
				select {
				case ch <- state:
				default:
				}
			}
		}(d)
	}
	wg.Wait()
}

func (p *Poller) GetState(id string) (DeviceState, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	s, ok := p.states[id]
	return s, ok
}

func (p *Poller) GetAllStates() map[string]DeviceState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	result := make(map[string]DeviceState, len(p.states))
	for k, v := range p.states {
		result[k] = v
	}
	return result
}

func (p *Poller) Subscribe() chan DeviceState {
	ch := make(chan DeviceState, 32)
	p.mu.Lock()
	p.subscribers = append(p.subscribers, ch)
	p.mu.Unlock()
	return ch
}

func (p *Poller) Unsubscribe(ch chan DeviceState) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, s := range p.subscribers {
		if s == ch {
			p.subscribers = append(p.subscribers[:i], p.subscribers[i+1:]...)
			close(ch)
			return
		}
	}
}
