package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
)

type Tab struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Order int    `json:"order"`
}

type Device struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IP          string `json:"ip"`
	Generation  int    `json:"generation"`
	TabID       string `json:"tab_id"`
	Description string `json:"description"`
}

type devicesFile struct {
	Tabs    []Tab    `json:"tabs"`
	Devices []Device `json:"devices"`
}

type DeviceStore struct {
	mu       sync.RWMutex
	path     string
	tabs     []Tab
	devices  []Device
}

func NewDeviceStore(path string) *DeviceStore {
	return &DeviceStore{path: path}
}

func (s *DeviceStore) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		s.tabs = []Tab{}
		s.devices = []Device{}
		return nil
	}
	if err != nil {
		return fmt.Errorf("devices.json lesen: %w", err)
	}

	var f devicesFile
	if err := json.Unmarshal(data, &f); err != nil {
		return fmt.Errorf("devices.json parsen: %w", err)
	}
	s.tabs = f.Tabs
	s.devices = f.Devices
	return nil
}

func (s *DeviceStore) save() error {
	f := devicesFile{Tabs: s.tabs, Devices: s.devices}
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

// Tabs

func (s *DeviceStore) GetTabs() []Tab {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Tab, len(s.tabs))
	copy(result, s.tabs)
	sort.Slice(result, func(i, j int) bool { return result[i].Order < result[j].Order })
	return result
}

func (s *DeviceStore) AddTab(tab Tab) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range s.tabs {
		if t.ID == tab.ID {
			return fmt.Errorf("Tab-ID '%s' existiert bereits", tab.ID)
		}
	}
	s.tabs = append(s.tabs, tab)
	return s.save()
}

func (s *DeviceStore) UpdateTab(id string, update Tab) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, t := range s.tabs {
		if t.ID == id {
			s.tabs[i].Name = update.Name
			return s.save()
		}
	}
	return fmt.Errorf("Tab '%s' nicht gefunden", id)
}

func (s *DeviceStore) DeleteTab(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, d := range s.devices {
		if d.TabID == id {
			return fmt.Errorf("Tab '%s' enthält noch Geräte", id)
		}
	}
	for i, t := range s.tabs {
		if t.ID == id {
			s.tabs = append(s.tabs[:i], s.tabs[i+1:]...)
			return s.save()
		}
	}
	return fmt.Errorf("Tab '%s' nicht gefunden", id)
}

// Devices

func (s *DeviceStore) GetDevices() []Device {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Device, len(s.devices))
	copy(result, s.devices)
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}

func (s *DeviceStore) GetDevice(id string) (Device, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, d := range s.devices {
		if d.ID == id {
			return d, true
		}
	}
	return Device{}, false
}

func (s *DeviceStore) AddDevice(d Device) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, existing := range s.devices {
		if existing.ID == d.ID {
			return fmt.Errorf("Geräte-ID '%s' existiert bereits", d.ID)
		}
	}
	s.devices = append(s.devices, d)
	return s.save()
}

func (s *DeviceStore) UpdateDevice(id string, update Device) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, d := range s.devices {
		if d.ID == id {
			update.ID = id
			s.devices[i] = update
			return s.save()
		}
	}
	return fmt.Errorf("Gerät '%s' nicht gefunden", id)
}

func (s *DeviceStore) PatchDevice(id string, fields map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, d := range s.devices {
		if d.ID == id {
			if tabID, ok := fields["tab_id"].(string); ok {
				s.devices[i].TabID = tabID
			}
			if name, ok := fields["name"].(string); ok {
				s.devices[i].Name = name
			}
			if ip, ok := fields["ip"].(string); ok {
				s.devices[i].IP = ip
			}
			if desc, ok := fields["description"].(string); ok {
				s.devices[i].Description = desc
			}
			if gen, ok := fields["generation"].(float64); ok {
				s.devices[i].Generation = int(gen)
			}
			return s.save()
		}
	}
	return fmt.Errorf("Gerät '%s' nicht gefunden", id)
}

func (s *DeviceStore) DeleteDevice(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, d := range s.devices {
		if d.ID == id {
			s.devices = append(s.devices[:i], s.devices[i+1:]...)
			return s.save()
		}
	}
	return fmt.Errorf("Gerät '%s' nicht gefunden", id)
}

func (s *DeviceStore) ReplaceAll(tabs []Tab, devices []Device) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tabs = tabs
	s.devices = devices
	return s.save()
}
