package api

import (
	"net/http"
	"strings"

	"shelly-control/store"
)

func (h *Handler) handleTabs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Methode nicht erlaubt")
		return
	}
	var t store.Tab
	if err := decodeJSON(r, &t); err != nil {
		writeError(w, http.StatusBadRequest, "Ungültiger JSON-Body")
		return
	}
	if t.ID == "" || t.Name == "" {
		writeError(w, http.StatusBadRequest, "id und name sind Pflichtfelder")
		return
	}
	if err := h.devices.AddTab(t); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, t)
}

func (h *Handler) handleTab(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/tabs/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Tab-ID fehlt")
		return
	}
	switch r.Method {
	case http.MethodPut:
		var t store.Tab
		if err := decodeJSON(r, &t); err != nil {
			writeError(w, http.StatusBadRequest, "Ungültiger JSON-Body")
			return
		}
		if err := h.devices.UpdateTab(id, t); err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})

	case http.MethodDelete:
		if err := h.devices.DeleteTab(id); err != nil {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})

	default:
		writeError(w, http.StatusMethodNotAllowed, "Methode nicht erlaubt")
	}
}
