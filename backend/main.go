package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"shelly-control/api"
	"shelly-control/shelly"
	"shelly-control/store"
)

//go:embed all:static
var staticFiles embed.FS

type Config struct {
	Port                  int `json:"port"`
	PollingIntervalSeconds int `json:"polling_interval_seconds"`
	PowerHistoryMinutes   int `json:"power_history_minutes"`
}

func loadConfig(path string) Config {
	cfg := Config{
		Port:                  8080,
		PollingIntervalSeconds: 5,
		PowerHistoryMinutes:   10,
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("config.json nicht gefunden, verwende Standardwerte: %v", err)
		return cfg
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("config.json ungültig: %v", err)
	}
	return cfg
}

func main() {
	devMode := flag.Bool("dev", false, "Entwicklungsmodus (CORS für Vite Dev Server)")
	flag.Parse()

	cfg := loadConfig("config.json")

	deviceStore := store.NewDeviceStore("devices.json")
	if err := deviceStore.Load(); err != nil {
		log.Fatalf("devices.json konnte nicht geladen werden: %v", err)
	}

	shellyClient := shelly.NewClient()
	powerBuffer := store.NewPowerBuffer(cfg.PowerHistoryMinutes)

	poller := shelly.NewPoller(shellyClient, deviceStore, powerBuffer, cfg.PollingIntervalSeconds)
	poller.Start()

	mux := http.NewServeMux()

	apiHandler := api.NewHandler(deviceStore, shellyClient, powerBuffer, poller)
	apiHandler.RegisterRoutes(mux)

	// Static Files ausliefern
	var staticFS http.FileSystem
	if *devMode {
		// Im Dev-Modus direkt aus dem Dateisystem (für HMR via Vite)
		staticFS = http.Dir("static")
	} else {
		sub, err := fs.Sub(staticFiles, "static")
		if err != nil {
			log.Fatalf("Static files einbetten fehlgeschlagen: %v", err)
		}
		staticFS = http.FS(sub)
	}

	fileServer := http.FileServer(staticFS)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// SPA: Alle nicht-API Routen → index.html
		if r.URL.Path != "/" {
			_, err := staticFS.Open(r.URL.Path)
			if err != nil {
				r.URL.Path = "/"
			}
		}
		if *devMode {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		fileServer.ServeHTTP(w, r)
	})

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Shelly Control läuft auf http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
