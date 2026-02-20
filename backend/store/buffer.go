package store

import (
	"sync"
	"time"
)

type PowerSample struct {
	Timestamp time.Time `json:"timestamp"`
	Watt      float64   `json:"watt"`
}

type PowerBuffer struct {
	mu             sync.RWMutex
	historyMinutes int
	data           map[string][]PowerSample
}

func NewPowerBuffer(historyMinutes int) *PowerBuffer {
	return &PowerBuffer{
		historyMinutes: historyMinutes,
		data:           make(map[string][]PowerSample),
	}
}

func (b *PowerBuffer) Add(deviceID string, watt float64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	cutoff := time.Now().Add(-time.Duration(b.historyMinutes) * time.Minute)
	sample := PowerSample{Timestamp: time.Now(), Watt: watt}

	samples := b.data[deviceID]
	// Alte Samples entfernen
	start := 0
	for start < len(samples) && samples[start].Timestamp.Before(cutoff) {
		start++
	}
	samples = append(samples[start:], sample)
	b.data[deviceID] = samples
}

func (b *PowerBuffer) Get(deviceID string) []PowerSample {
	b.mu.RLock()
	defer b.mu.RUnlock()

	samples := b.data[deviceID]
	result := make([]PowerSample, len(samples))
	copy(result, samples)
	return result
}

func (b *PowerBuffer) Latest(deviceID string) (float64, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	samples := b.data[deviceID]
	if len(samples) == 0 {
		return 0, false
	}
	return samples[len(samples)-1].Watt, true
}
