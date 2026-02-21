package api

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"shelly-control/store"
)

type importDevice struct {
	Name        string `json:"name"`
	IP          string `json:"ip"`
	Generation  int    `json:"generation"`
	Tab         string `json:"tab"`
	Description string `json:"description"`
}

func (h *Handler) handleExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Methode nicht erlaubt")
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}
	scope := r.URL.Query().Get("scope")

	tabs := h.devices.GetTabs()
	devices := h.devices.GetDevices()

	// Tab-Name-Lookup
	tabNames := make(map[string]string)
	for _, t := range tabs {
		tabNames[t.ID] = t.Name
	}

	// Nach Scope filtern
	var filtered []store.Device
	for _, d := range devices {
		if scope == "" || scope == "all" || d.TabID == scope {
			filtered = append(filtered, d)
		}
	}

	date := time.Now().Format("2006-01-02")
	scopeLabel := "all"
	if scope != "" && scope != "all" {
		scopeLabel = scope
	}

	switch format {
	case "csv":
		filename := fmt.Sprintf("shelly-export-%s-%s.csv", scopeLabel, date)
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
		cw := csv.NewWriter(w)
		cw.Write([]string{"name", "ip", "generation", "tab", "description"})
		for _, d := range filtered {
			cw.Write([]string{d.Name, d.IP, strconv.Itoa(d.Generation), tabNames[d.TabID], d.Description})
		}
		cw.Flush()

	default: // json
		filename := fmt.Sprintf("shelly-export-%s-%s.json", scopeLabel, date)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
		out := make([]importDevice, len(filtered))
		for i, d := range filtered {
			out[i] = importDevice{
				Name:        d.Name,
				IP:          d.IP,
				Generation:  d.Generation,
				Tab:         tabNames[d.TabID],
				Description: d.Description,
			}
		}
		json.NewEncoder(w).Encode(out)
	}
}

type importPreviewRow struct {
	Name        string `json:"name"`
	IP          string `json:"ip"`
	Generation  int    `json:"generation"`
	Tab         string `json:"tab"`
	Description string `json:"description"`
	Status      string `json:"status"` // "new", "update", "unchanged", "error"
	Error       string `json:"error,omitempty"`
}

func (h *Handler) handleImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Methode nicht erlaubt")
		return
	}

	format := r.URL.Query().Get("format")
	mode := r.URL.Query().Get("mode") // add, overwrite, replace
	if mode == "" {
		mode = "add"
	}

	var items []importDevice
	var parseErr error

	switch format {
	case "csv":
		items, parseErr = parseCSV(r)
	default:
		items, parseErr = parseJSON(r)
	}

	if parseErr != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Parse-Fehler: %v", parseErr))
		return
	}

	existingDevices := h.devices.GetDevices()
	existingTabs := h.devices.GetTabs()

	// Index für schnelle Suche
	deviceByName := make(map[string]store.Device)
	for _, d := range existingDevices {
		deviceByName[d.Name] = d
	}
	tabByName := make(map[string]store.Tab)
	for _, t := range existingTabs {
		tabByName[t.Name] = t
	}

	// Preview aufbauen + Validierung
	preview := make([]importPreviewRow, len(items))
	for i, item := range items {
		row := importPreviewRow{
			Name:        item.Name,
			IP:          item.IP,
			Generation:  item.Generation,
			Tab:         item.Tab,
			Description: item.Description,
		}
		// Validierung
		if item.Name == "" {
			row.Status = "error"
			row.Error = "name fehlt"
		} else if item.IP == "" {
			row.Status = "error"
			row.Error = "ip fehlt"
		} else if item.Generation != 1 && item.Generation != 3 && item.Generation != 4 {
			row.Status = "error"
			row.Error = "generation muss 1, 3 oder 4 sein"
		} else if existing, exists := deviceByName[item.Name]; exists {
			if existing.IP == item.IP && existing.Generation == item.Generation && existing.Description == item.Description {
				row.Status = "unchanged"
			} else {
				row.Status = "update"
			}
		} else {
			row.Status = "new"
		}
		preview[i] = row
	}

	// Dry-run? Query-Parameter preview=true
	if r.URL.Query().Get("preview") == "true" {
		writeJSON(w, http.StatusOK, preview)
		return
	}

	// Ausführen
	switch mode {
	case "replace":
		// Alles löschen und neu aufbauen
		newTabs, newDevices := buildStoreData(items, tabByName)
		if err := h.devices.ReplaceAll(newTabs, newDevices); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	default: // add / overwrite
		for _, item := range items {
			if item.Name == "" || item.IP == "" || (item.Generation != 1 && item.Generation != 3) {
				continue
			}
			// Tab sicherstellen
			tab, tabExists := tabByName[item.Tab]
			if !tabExists && item.Tab != "" {
				tab = store.Tab{
					ID:    slugify(item.Tab),
					Name:  item.Tab,
					Order: len(tabByName) + 1,
				}
				h.devices.AddTab(tab)
				tabByName[item.Tab] = tab
			}

			d := store.Device{
				ID:          slugify(item.Name),
				Name:        item.Name,
				IP:          item.IP,
				Generation:  item.Generation,
				TabID:       tab.ID,
				Description: item.Description,
			}

			if _, exists := deviceByName[item.Name]; exists && mode == "overwrite" {
				h.devices.UpdateDevice(d.ID, d)
			} else if _, exists := deviceByName[item.Name]; !exists {
				h.devices.AddDevice(d)
			}
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"ok":      true,
		"preview": preview,
	})
}

func parseJSON(r *http.Request) ([]importDevice, error) {
	var items []importDevice
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}

func parseCSV(r *http.Request) ([]importDevice, error) {
	reader := csv.NewReader(r.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("CSV leer oder nur Header")
	}

	// Header-Mapping
	header := records[0]
	idx := make(map[string]int)
	for i, h := range header {
		idx[strings.ToLower(strings.TrimSpace(h))] = i
	}

	var items []importDevice
	for _, row := range records[1:] {
		get := func(key string) string {
			i, ok := idx[key]
			if !ok || i >= len(row) {
				return ""
			}
			return strings.TrimSpace(row[i])
		}
		gen, _ := strconv.Atoi(get("generation"))
		items = append(items, importDevice{
			Name:        get("name"),
			IP:          get("ip"),
			Generation:  gen,
			Tab:         get("tab"),
			Description: get("description"),
		})
	}
	return items, nil
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func buildStoreData(items []importDevice, existingTabsByName map[string]store.Tab) ([]store.Tab, []store.Device) {
	tabOrder := len(existingTabsByName)
	tabsByName := make(map[string]store.Tab)
	for k, v := range existingTabsByName {
		tabsByName[k] = v
	}

	var devices []store.Device
	for _, item := range items {
		if item.Name == "" || item.IP == "" {
			continue
		}
		if _, exists := tabsByName[item.Tab]; !exists && item.Tab != "" {
			tabOrder++
			tabsByName[item.Tab] = store.Tab{
				ID:    slugify(item.Tab),
				Name:  item.Tab,
				Order: tabOrder,
			}
		}
		tab := tabsByName[item.Tab]
		devices = append(devices, store.Device{
			ID:          slugify(item.Name),
			Name:        item.Name,
			IP:          item.IP,
			Generation:  item.Generation,
			TabID:       tab.ID,
			Description: item.Description,
		})
	}

	tabs := make([]store.Tab, 0, len(tabsByName))
	for _, t := range tabsByName {
		tabs = append(tabs, t)
	}
	return tabs, devices
}
