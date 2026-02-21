# Shelly Control

Web-basierte Oberfläche zur zentralen Verwaltung und Steuerung von Shelly Smart Plugs im lokalen Netzwerk. Unterstützt Shelly Gen 1, Gen 3 und Gen 4 Geräte.

## Features

- **Geräteübersicht** – Alle Shelly Plugs auf einen Blick mit aktuellem Status (ein/aus, Erreichbarkeit)
- **Schalten** – Geräte direkt aus der Web-UI ein- und ausschalten
- **Leistungsanzeige** – Aktuelle Stromaufnahme (Watt) mit Live-Sparkline je Gerät
- **Tabs** – Geräte in Gruppen (z. B. Racks) organisieren
- **Polling** – Konfigurierbares Intervall für automatische Statusabfragen
- **Import / Export** – Gerätekonfiguration als JSON oder CSV sichern und wiederherstellen
- **Geräteverwaltung** – Geräte und Tabs direkt in der UI anlegen, bearbeiten und löschen
- **Dark Mode** – Automatisch per `prefers-color-scheme`, manuell umschaltbar
- **Shelly-Generationen** – Gen 1 (Plug S), Gen 3 (Plug S MTR), Gen 4 (1PM Gen4)

## Tech Stack

| Bereich       | Technologie                         |
|---------------|-------------------------------------|
| Frontend      | Svelte 5 (Runes), Vite, Bulma CSS   |
| Backend       | Go 1.23, stdlib `net/http`          |
| Build         | Makefile, eingebettete Static Files |
| Zielplattform | Windows amd64, macOS amd64 / arm64  |

## Projektstruktur

```
shelly-control/
├── .github/
│   └── workflows/
│       └── release.yml   # CI/CD: Windows + macOS Builds
├── backend/
│   ├── api/              # HTTP-Handler (REST-Endpunkte)
│   ├── shelly/           # Shelly-Client (Gen1, Gen3, Gen4), Poller
│   ├── store/            # DeviceStore, PowerBuffer
│   ├── config.json       # Konfiguration (Port, Polling, History)
│   ├── devices.json      # Gerät- und Tab-Definitionen
│   └── main.go
└── frontend/
    └── src/
        ├── components/   # DeviceCard, TabBar, Sparkline, Modals, …
        ├── stores/       # Svelte Stores (devices, editMode, theme)
        └── styles/       # Globale CSS-Variablen und Utilities
```

## Konfiguration

**`backend/config.json`**
```json
{
  "port": 8080,
  "polling_interval_seconds": 5,
  "power_history_minutes": 5
}
```

**`backend/devices.json`** – Geräte und Tabs (kann auch über die UI verwaltet werden):
```json
{
  "tabs": [
    { "id": "rack1", "name": "Rack 1", "order": 1 }
  ],
  "devices": [
    {
      "id": "s1r1",
      "name": "S1R1",
      "ip": "192.168.1.50",
      "generation": 1,
      "tab_id": "rack1",
      "description": "Dell PowerEdge R640"
    }
  ]
}
```

## Shelly-Generationen

| Generation | Gerät           | API                          |
|------------|-----------------|------------------------------|
| 1          | Plug S          | `GET /status`, `GET /relay/0`|
| 3          | Plug S MTR      | RPC: `Switch.GetStatus`      |
| 4          | 1PM Gen4        | RPC: `Switch.GetStatus`      |

## API-Endpunkte

| Methode | Pfad                      | Beschreibung                       |
|---------|---------------------------|------------------------------------|
| GET     | `/api/devices`            | Alle Geräte mit Status             |
| POST    | `/api/devices`            | Neues Gerät anlegen                |
| PUT     | `/api/devices/{id}`       | Gerät aktualisieren                |
| PATCH   | `/api/devices/{id}`       | Teilupdate (z. B. Tab-Zuweisung)   |
| DELETE  | `/api/devices/{id}`       | Gerät entfernen                    |
| POST    | `/api/devices/{id}/toggle`| Relay schalten                     |
| POST    | `/api/devices/{id}/test`  | Verbindung testen                  |
| GET     | `/api/devices/{id}/history` | Leistungshistorie (Rolling Buffer)|
| POST    | `/api/tabs`               | Neuen Tab anlegen                  |
| PUT     | `/api/tabs/{id}`          | Tab aktualisieren                  |
| DELETE  | `/api/tabs/{id}`          | Tab entfernen                      |
| GET     | `/api/events`             | SSE-Stream für Live-Updates        |
| GET     | `/api/export`             | Konfiguration exportieren          |
| POST    | `/api/import`             | Konfiguration importieren          |

## Entwicklung

### Voraussetzungen

- [Go 1.23+](https://go.dev/dl/)
- [Node.js 22+](https://nodejs.org/)

### Setup

```bash
cd frontend && npm install
```

### Dev-Modus starten (Frontend + Backend parallel)

```bash
make dev
```

Öffne [http://localhost:5173](http://localhost:5173) im Browser (Vite Dev Server mit HMR).
Das Backend läuft auf Port 8080.

## Build & Deployment

### Lokal bauen (Windows-Binary)

```bash
make build
```

Erzeugt `backend/shelly-control.exe` mit eingebettetem Frontend — ein einzelnes, portables Binary.

### Ausführen auf dem Zielsystem (Windows)

```
shelly-control.exe
```

Öffne im Browser: `http://localhost:8080`

### Release via GitHub Actions

Ein neues Release wird automatisch gebaut und veröffentlicht, sobald ein Tag gepusht wird:

```bash
git tag v1.0.0
git push origin v1.0.0
```

Der Workflow erstellt drei Binaries als Release Assets:

| Datei                              | Plattform        |
|------------------------------------|------------------|
| `shelly-control-windows-amd64.exe` | Windows 10/11    |
| `shelly-control-macos-amd64`       | macOS Intel      |
| `shelly-control-macos-arm64`       | macOS Apple Silicon |

### Aufräumen

```bash
make clean
```
