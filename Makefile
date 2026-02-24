.PHONY: build build-macos dev dist run clean

# Produktions-Build für Windows (amd64)
build:
	cd frontend && npm run build
	rm -rf backend/static/*
	cp -r frontend/dist/* backend/static/
	cd backend && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o shelly-control-windows-amd64.exe .
	@echo "✓ shelly-control-windows-amd64.exe bereit"

# macOS-Builds (Intel + Apple Silicon)
build-macos:
	cd frontend && npm run build
	rm -rf backend/static/*
	cp -r frontend/dist/* backend/static/
	cd backend && CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o shelly-control-macos-amd64 .
	cd backend && CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o shelly-control-macos-arm64 .
	@echo "✓ shelly-control-macos-amd64 (Intel) bereit"
	@echo "✓ shelly-control-macos-arm64 (Apple Silicon) bereit"

# Alle Plattformen bauen
dist: build build-macos
	@echo "✓ Alle Binaries bereit"

# Lokale Entwicklung auf macOS
dev:
	cd frontend && npm run dev & \
	cd backend && go run . --dev; \
	kill %1 2>/dev/null; wait

# macOS-Binary für lokales Testen (ohne Frontend-Build)
run:
	cd backend && go run .

clean:
	rm -f backend/shelly-control-windows-amd64.exe
	rm -f backend/shelly-control-macos-amd64
	rm -f backend/shelly-control-macos-arm64
	rm -f backend/shelly-control
	rm -rf backend/static/*
	rm -rf frontend/dist
