# Shelly Control

Web-basierte Oberfläche zur zentralen Verwaltung und Steuerung von Shelly Smart Plugs im lokalen Netzwerk. Unterstützt Shelly Gen 1, Gen 3 und Gen 4 Geräte.

## Features

- **Geräteübersicht** – Alle Shelly Plugs auf einen Blick mit aktuellem Status (ein/aus, Erreichbarkeit)
- **Schalten** – Geräte direkt aus der Web-UI ein- und ausschalten
- **Leistungsanzeige** – Aktuelle Stromaufnahme (Watt) mit Live-Sparkline je Gerät 
- **Tabs** – Geräte in Gruppen (z. B. Racks) organisieren
- **Polling** – Konfigurierbares Intervall für automatische Statusabfragen
- **Import / Export** – Gerätekonfiguration als JSON sichern und wiederherstellen
- **Geräteverwaltung** – Geräte und Tabs direkt in der UI anlegen, bearbeiten und löschen
- **Gen 1 & Gen 3 & Gen 4** – Unterstützung für die Shelly-Generationen

## Tech Stack

| Bereich   | Technologie              |
|-----------|--------------------------|
| Frontend  | Svelte 5, Vite, Bulma CSS |
| Backend   | Go 1.23, stdlib `net/http` |
| Build     | Makefile, eingebettete Static Files (`embed.FS`) |
| Zielplattform | Windows (amd64) |

## Projektstruktur

```
shelly-control/
├── backend/
│   ├── api/          # HTTP-Handler (REST-Endpunkte)
│   ├── shelly/       # Shelly-Client (Gen1 & Gen3), Poller
│   ├── store/        # DeviceStore, PowerBuffer
│   ├── config.json   # Konfiguration (Port, Polling, History)
│   ├── devices.json  # Gerät- und Tab-Definitionen
│   └── main.go
└── frontend/
    └── src/
        ├── components/   # DeviceCard, TabBar, Modals, ...
        ├── stores/       # Svelte Stores
        └── styles/
```

## Konfiguration

**`backend/config.json`**
```json
{
  "port": 8080,
  "polling_interval_seconds": 5,
  "power_history_minutes": 10
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

## API-Endpunkte

| Methode | Pfad              | Beschreibung                        |
|---------|-------------------|-------------------------------------|
| GET     | `/api/devices`    | Alle Geräte mit Status              |
| POST    | `/api/devices`    | Neues Gerät anlegen                 |
| PUT     | `/api/devices/{id}` | Gerät aktualisieren               |
| DELETE  | `/api/devices/{id}` | Gerät entfernen                   |
| GET     | `/api/tabs`       | Alle Tabs                           |
| POST    | `/api/tabs`       | Neuen Tab anlegen                   |
| PUT     | `/api/tabs/{id}`  | Tab aktualisieren                   |
| DELETE  | `/api/tabs/{id}`  | Tab entfernen                       |
| GET     | `/api/events`     | SSE-Stream für Live-Updates         |
| GET     | `/api/export`     | Konfiguration als JSON exportieren  |
| POST    | `/api/import`     | Konfiguration importieren           |

## Entwicklung

### Voraussetzungen

- [Go 1.23+](https://go.dev/dl/)
- [Node.js 20+](https://nodejs.org/)

### Setup

```bash
# Frontend-Abhängigkeiten installieren
cd frontend && npm install
```

### Dev-Modus starten (Frontend + Backend parallel)

```bash
make dev
```

Öffne [http://localhost:5173](http://localhost:5173) im Browser (Vite Dev Server mit HMR).
Das Backend läuft auf Port 8080.

## Build & Deployment

### Windows-Binary erstellen

```bash
make build
```

Erzeugt `backend/shelly-control.exe` mit eingebettetem Frontend.
Die Datei ist ein einzelnes, portables Binary – kein Webserver oder separate Dateien nötig.

### Ausführen auf dem Zielsystem (Windows)

```
shelly-control.exe
```

Öffne im Browser: `http://localhost:8080`

### Aufräumen

```bash
make clean
```
