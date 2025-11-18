# 🖥️ Emisor de Pantalla

Herramienta ligera para transmitir el escritorio a otros dispositivos en la misma red local (LAN) utilizando un navegador web.

![Captura de la aplicación](screenshot.png)

## Características

*   **Sin instalaciones extra:** El receptor solo necesita un navegador web moderno.
*   **Multiplataforma y Arquitectura:** Funciona nativamente en Windows, Linux, macOS (Intel y Apple Silicon) y dispositivos ARM como Raspberry Pi.
*   **Feedback visual:** Muestra un indicador naranja al hacer clic para facilitar el seguimiento.
*   **Interactivo:** El receptor puede dibujar sobre la transmisión (lápiz digital).
*   **Configurable:** Consumo de recursos ajustable (calidad y FPS).

## Caso de Uso Principal

Este proyecto fue concebido para su uso en aulas con **PDIs (Pizarras Digitales Interactivas)**.

Permite al docente enviar la señal de su ordenador a la PDI simplemente abriendo una URL en el navegador de la pizarra. Esto elimina la necesidad de instalar aplicaciones de terceros en la pantalla, gestionar cables HDMI largos o depender de dongles propietarios. Al funcionar en Raspberry Pi, también puede convertir cualquier proyector antiguo en un receptor inalámbrico económico.

## 📥 Descarga

Puedes descargar los binarios compilados en la sección de **[Releases](https://github.com/soyunomas/emisor-pantalla/releases)**.

Solo descarga el archivo correspondiente a tu arquitectura:
*   **Windows** (`.exe`)
*   **Linux** (`amd64`, `arm64`)
*   **macOS** (`Intel`, `M1/M2`)
*   **Raspberry Pi** (ARM)

## ⚙️ Uso

### Linux / macOS / Raspberry Pi
Abre una terminal en la carpeta de descarga y ejecuta:

```bash
chmod +x streamer-profesor  # Solo la primera vez para dar permisos
./streamer-profesor
```

### Windows
1.  Abre la carpeta donde descargaste el archivo `.exe`.
2.  Escribe `cmd` en la barra de direcciones del explorador de archivos y pulsa Enter.
3.  Escribe el nombre del archivo (ej: `streamer-profesor-windows.exe`) y pulsa Enter.
4.  **Importante:** Si aparece el Firewall de Windows, permite el acceso en redes privadas.

---

Una vez iniciado, la terminal mostrará la dirección (ej: `http://192.168.1.35:8080`) que debes introducir en el navegador de la PDI o del alumno.

### Configuración (Flags)

Puedes ajustar el rendimiento añadiendo parámetros al comando:

*   `-port`: Puerto del servidor (defecto: 8080).
*   `-quality`: Calidad de compresión JPEG (1-100).
*   `-fps`: Límite de cuadros por segundo.

**Ejemplo para redes WiFi inestables:**
```bash
./streamer-profesor -quality=30 -fps=10
```

**Ejemplo para máxima calidad (red cableada):**
```bash
./streamer-profesor -quality=80 -fps=25
```

## 🛠️ Compilación desde código fuente

Si prefieres compilarlo tú mismo desde el código fuente:

### Requisitos (Linux)
Necesario para capturar pantalla y eventos de entrada:
```bash
sudo apt update
sudo apt install libx11-dev libx11-xcb-dev libxtst-dev libxcb1-dev libxkbcommon-dev libxkbcommon-x11-dev
```

### Comandos
Usa el `Makefile` incluido:

```bash
make build      # Compilar para el sistema actual
make build-all  # Compilar para Windows, Mac y Linux
make clean      # Limpiar
```

## 📄 Licencia

Este proyecto se distribuye bajo la licencia MIT.
