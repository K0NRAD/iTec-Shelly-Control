GOOS   = windows
GOARCH = amd64
BINARY = shelly-control.exe

.PHONY: build dev dist clean

# Produktions-Build für Windows
build:
	cd frontend && npm run build
	rm -rf backend/static/*
	cp -r frontend/dist/* backend/static/
	cd backend && CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY) .

# Lokale Entwicklung auf macOS
dev:
	cd frontend && npm run dev & \
	cd backend && go run . --dev; \
	kill %1 2>/dev/null; wait

# Alias für Windows-Distribution
dist: build
	@echo "✓ $(BINARY) bereit für Deployment auf Windows"

# macOS-Binary zum lokalen Testen (ohne Frontend-Build)
run:
	cd backend && go run .

clean:
	rm -f backend/$(BINARY) backend/shelly-control
	rm -rf backend/static/*
	rm -rf frontend/dist
