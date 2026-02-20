package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) handleEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "SSE nicht unterstützt")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	ch := h.poller.Subscribe()
	defer h.poller.Unsubscribe(ch)

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case state, ok := <-ch:
			if !ok {
				return
			}
			data, _ := json.Marshal(map[string]interface{}{
				"id":     state.ID,
				"on":     state.On,
				"watt":   state.Watt,
				"online": state.Err == nil,
			})
			fmt.Fprintf(w, "event: device_update\ndata: %s\n\n", data)
			flusher.Flush()
		}
	}
}
