package api

import (
	"net/http"

	"shelly-control/shelly"
	"shelly-control/store"
)

type Handler struct {
	devices *store.DeviceStore
	shelly  *shelly.Client
	buffer  *store.PowerBuffer
	poller  *shelly.Poller
}

func NewHandler(devices *store.DeviceStore, shellyClient *shelly.Client, buffer *store.PowerBuffer, poller *shelly.Poller) *Handler {
	return &Handler{
		devices: devices,
		shelly:  shellyClient,
		buffer:  buffer,
		poller:  poller,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/devices", h.handleDevices)
	mux.HandleFunc("/api/devices/", h.handleDevice)
	mux.HandleFunc("/api/tabs", h.handleTabs)
	mux.HandleFunc("/api/tabs/", h.handleTab)
	mux.HandleFunc("/api/events", h.handleEvents)
	mux.HandleFunc("/api/export", h.handleExport)
	mux.HandleFunc("/api/import", h.handleImport)
}
