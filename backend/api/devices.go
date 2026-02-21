package api

import (
	"net/http"
	"strings"

	"shelly-control/store"
)

// DeviceWithStatus kombiniert Device-Stammdaten mit aktuellem Poll-Status
type DeviceWithStatus struct {
	store.Device
	On      bool    `json:"on"`
	Watt    float64 `json:"watt"`
	Online  bool    `json:"online"`
}

// GET /api/devices → alle Tabs + Devices mit Status
// POST /api/devices → neues Gerät anlegen
func (h *Handler) handleDevices(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listDevices(w, r)
	case http.MethodPost:
		h.createDevice(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "Methode nicht erlaubt")
	}
}

func (h *Handler) listDevices(w http.ResponseWriter, r *http.Request) {
	tabs := h.devices.GetTabs()
	devices := h.devices.GetDevices()
	states := h.poller.GetAllStates()

	enriched := make([]DeviceWithStatus, len(devices))
	for i, d := range devices {
		s, ok := states[d.ID]
		enriched[i] = DeviceWithStatus{
			Device: d,
			On:     s.On,
			Watt:   s.Watt,
			Online: ok && s.Err == nil,
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"tabs":    tabs,
		"devices": enriched,
	})
}

func (h *Handler) createDevice(w http.ResponseWriter, r *http.Request) {
	var d store.Device
	if err := decodeJSON(r, &d); err != nil {
		writeError(w, http.StatusBadRequest, "Ungültiger JSON-Body")
		return
	}
	if d.ID == "" || d.Name == "" || d.IP == "" {
		writeError(w, http.StatusBadRequest, "id, name und ip sind Pflichtfelder")
		return
	}
	if d.Generation != 1 && d.Generation != 3 && d.Generation != 4 {
		writeError(w, http.StatusBadRequest, "generation muss 1, 3 oder 4 sein")
		return
	}
	if err := h.devices.AddDevice(d); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, d)
}

// Routing für /api/devices/{id}[/action]
func (h *Handler) handleDevice(w http.ResponseWriter, r *http.Request) {
	// Pfad: /api/devices/{id} oder /api/devices/{id}/toggle etc.
	path := strings.TrimPrefix(r.URL.Path, "/api/devices/")
	parts := strings.SplitN(path, "/", 2)
	id := parts[0]
	action := ""
	if len(parts) == 2 {
		action = parts[1]
	}

	if id == "" {
		writeError(w, http.StatusBadRequest, "Geräte-ID fehlt")
		return
	}

	switch action {
	case "":
		h.handleDeviceByID(w, r, id)
	case "toggle":
		h.toggleDevice(w, r, id)
	case "status":
		h.deviceStatus(w, r, id)
	case "history":
		h.deviceHistory(w, r, id)
	case "test":
		h.testDevice(w, r, id)
	default:
		writeError(w, http.StatusNotFound, "Unbekannte Aktion")
	}
}

func (h *Handler) handleDeviceByID(w http.ResponseWriter, r *http.Request, id string) {
	switch r.Method {
	case http.MethodGet:
		h.deviceStatus(w, r, id)
	case http.MethodPut:
		h.updateDevice(w, r, id)
	case http.MethodPatch:
		h.patchDevice(w, r, id)
	case http.MethodDelete:
		h.deleteDevice(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "Methode nicht erlaubt")
	}
}

func (h *Handler) deviceStatus(w http.ResponseWriter, r *http.Request, id string) {
	d, ok := h.devices.GetDevice(id)
	if !ok {
		writeError(w, http.StatusNotFound, "Gerät nicht gefunden")
		return
	}
	s, online := h.poller.GetState(id)
	writeJSON(w, http.StatusOK, DeviceWithStatus{
		Device: d,
		On:     s.On,
		Watt:   s.Watt,
		Online: online && s.Err == nil,
	})
}

func (h *Handler) toggleDevice(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Methode nicht erlaubt")
		return
	}
	var body struct {
		On bool `json:"on"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "Ungültiger JSON-Body")
		return
	}
	d, ok := h.devices.GetDevice(id)
	if !ok {
		writeError(w, http.StatusNotFound, "Gerät nicht gefunden")
		return
	}
	if err := h.shelly.SetRelay(d.IP, d.Generation, body.On); err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *Handler) deviceHistory(w http.ResponseWriter, r *http.Request, id string) {
	if _, ok := h.devices.GetDevice(id); !ok {
		writeError(w, http.StatusNotFound, "Gerät nicht gefunden")
		return
	}
	samples := h.buffer.Get(id)
	writeJSON(w, http.StatusOK, samples)
}

func (h *Handler) testDevice(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Methode nicht erlaubt")
		return
	}
	d, ok := h.devices.GetDevice(id)
	if !ok {
		writeError(w, http.StatusNotFound, "Gerät nicht gefunden")
		return
	}
	if err := h.shelly.Test(d.IP, d.Generation); err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{"ok": false, "error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *Handler) updateDevice(w http.ResponseWriter, r *http.Request, id string) {
	var d store.Device
	if err := decodeJSON(r, &d); err != nil {
		writeError(w, http.StatusBadRequest, "Ungültiger JSON-Body")
		return
	}
	if err := h.devices.UpdateDevice(id, d); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *Handler) patchDevice(w http.ResponseWriter, r *http.Request, id string) {
	var fields map[string]interface{}
	if err := decodeJSON(r, &fields); err != nil {
		writeError(w, http.StatusBadRequest, "Ungültiger JSON-Body")
		return
	}
	if err := h.devices.PatchDevice(id, fields); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *Handler) deleteDevice(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.devices.DeleteDevice(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
