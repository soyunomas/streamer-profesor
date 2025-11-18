# Makefile para el proyecto emisor-pantalla
#
# IMPORTANTE: Las líneas de comandos (indentadas) DEBEN empezar
# con un carácter de TABULADOR, no con espacios.

# --- Variables de Configuración ---
BINARY_NAME=streamer-profesor
CMD_PATH=./cmd/streamer-profesor

# --- Comandos Principales ---

all: build

## build: Compila el binario para el sistema operativo y arquitectura actual (nativo)
build:
	@echo "==> Compilando para el sistema nativo..."
	@go build -o $(BINARY_NAME) $(CMD_PATH)
	@echo "==> Binario '$(BINARY_NAME)' creado."

## run: Ejecuta la aplicación usando 'go run'
run: tidy
	@echo "==> Ejecutando la aplicación..."
	@go run $(CMD_PATH)

## clean: Elimina los binarios compilados
clean:
	@echo "==> Limpiando binarios..."
	@rm -f $(BINARY_NAME)
	@rm -rf ./build
	@echo "==> Limpieza completada."

## tidy: Sincroniza las dependencias del proyecto
tidy:
	@echo "==> Sincronizando dependencias..."
	@go mod tidy

# --- Comandos de Cross-Compilación ---

## build-all: Compila para las plataformas más comunes
build-all: build-linux build-windows build-macos-intel build-macos-arm
	@echo "==> Todas las compilaciones han finalizado."

## build-linux: Compila para Linux (amd64)
build-linux:
	@echo "==> Compilando para Linux (amd64)..."
	@mkdir -p ./build
	@GOOS=linux GOARCH=amd64 go build -o ./build/$(BINARY_NAME)-linux-amd64 $(CMD_PATH)

## build-windows: Compila para Windows (amd64)
build-windows:
	@echo "==> Compilando para Windows (amd64)..."
	@mkdir -p ./build
	@GOOS=windows GOARCH=amd64 go build -o ./build/$(BINARY_NAME)-windows-amd64.exe $(CMD_PATH)

## build-macos-intel: Compila para macOS (Intel amd64)
build-macos-intel:
	@echo "==> Compilando para macOS (Intel amd64)..."
	@mkdir -p ./build
	@GOOS=darwin GOARCH=amd64 go build -o ./build/$(BINARY_NAME)-macos-amd64 $(CMD_PATH)

## build-macos-arm: Compila para macOS (Apple Silicon arm64)
build-macos-arm:
	@echo "==> Compilando para macOS (Apple Silicon arm64)..."
	@mkdir -p ./build
	@GOOS=darwin GOARCH=arm64 go build -o ./build/$(BINARY_NAME)-macos-arm64 $(CMD_PATH)

# --- Comandos de Cross-Compilación para Raspberry Pi ---

## build-rpi-all: Compila para las arquitecturas comunes de Raspberry Pi
build-rpi-all: build-rpi-64 build-rpi-32
	@echo "==> Compilaciones para Raspberry Pi finalizadas."

## build-rpi-64: Compila para Raspberry Pi (Linux arm64, 64-bit OS)
build-rpi-64:
	@echo "==> Compilando para Raspberry Pi (arm64)..."
	@mkdir -p ./build
	@GOOS=linux GOARCH=arm64 go build -o ./build/$(BINARY_NAME)-linux-arm64 $(CMD_PATH)

## build-rpi-32: Compila para Raspberry Pi (Linux arm, 32-bit OS)
build-rpi-32:
	@echo "==> Compilando para Raspberry Pi (arm)..."
	@mkdir -p ./build
	@GOOS=linux GOARCH=arm go build -o ./build/$(BINARY_NAME)-linux-arm $(CMD_PATH)

# --- Comando de Ayuda ---

## help: Muestra esta ayuda
help:
	@echo "----------------------------------------------------"
	@echo " Comandos disponibles para el proyecto emisor-pantalla:"
	@echo "----------------------------------------------------"
	@echo "  make build             - Compila el binario para el sistema nativo (Linux)."
	@echo "  make run               - Ejecuta la aplicación para pruebas rápidas."
	@echo "  make clean             - Elimina todos los binarios y directorios de compilación."
	@echo "  make tidy              - Sincroniza las dependencias del proyecto."
	@echo ""
	@echo "  make build-all         - Compila para Windows, macOS (Intel/ARM) y Linux (amd64)."
	@echo "  make build-windows     - Compila para Windows (64-bit)."
	@echo "  make build-linux       - Compila para Linux (amd64)."
	@echo "  make build-macos-intel - Compila para macOS (Intel)."
	@echo "  make build-macos-arm   - Compila para macOS (Apple Silicon)."
	@echo ""
	@echo "  make build-rpi-all     - Compila para Raspberry Pi (32 y 64 bits)."
	@echo "  make build-rpi-64      - Compila para Raspberry Pi (64-bit OS)."
	@echo "  make build-rpi-32      - Compila para Raspberry Pi (32-bit OS)."
	@echo ""
	@echo "  make                   - Alias para 'make build'."
	@echo "  make help              - Muestra esta ayuda."

# .PHONY declara targets que no son archivos, evitando conflictos y forzando su ejecución.
.PHONY: all build run clean tidy help build-all build-linux build-windows build-macos-intel build-macos-arm build-rpi-all build-rpi-64 build-rpi-32
