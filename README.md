# LabelZoom Print Agent

[![Docker Hub](https://img.shields.io/badge/Docker%20Hub-labelzoom%2Flz--print--agent--local-blue?logo=docker)](https://hub.docker.com/r/labelzoom/lz-print-agent-local)
[![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A lightweight local print agent that enables browser-based printing to direct thermal printers when using the [LabelZoom Web App](https://www.labelzoom.net/app). This agent acts as a bridge between your web browser and RAW printers on your local network, allowing fast and accurate reproduction of barcode labels in the printer's native language (such as Zebra ZPL).

## 🎯 What It Does

The LabelZoom Print Agent:
- **Receives print jobs** from your web browser via HTTP
- **Forwards print data** to thermal printers via TCP (port 9100)
- **Supports RAW printing** for direct thermal label printers
- **Runs locally** to communicate with printers on your network

## 🚀 Quick Start

### Using Docker (Recommended)

The easiest way to run the print agent is using Docker:

```bash
docker run -d -p 52045:8080 labelzoom/lz-print-agent-local
```

This starts the agent and makes it available at `http://localhost:52045`.

**Why port 52045?** Port 8080 is commonly used by other applications, so we recommend remapping to a higher port (50000+). The default is 52045, but you can use any available port.

### Using Docker Compose

Create a `docker-compose.yml` file:

```yaml
version: '3.8'
services:
  lz-print-agent:
    image: labelzoom/lz-print-agent-local:latest
    ports:
      - "52045:8080"
    restart: unless-stopped
```

Then run:

```bash
docker-compose up -d
```

### Building from Source

If you prefer to build and run the agent locally:

```bash
# Clone the repository
git clone https://github.com/yourusername/lz-print-agent-go.git
cd lz-print-agent-go

# Build the binary
go build -o lz-print-agent-local .

# Run the agent
./lz-print-agent-local
```

The agent will start on port 8080 by default.

## 📋 Requirements

- **Docker** (for containerized deployment) OR
- **Go 1.24+** (for building from source)
- **Network access** to your thermal printers (typically port 9100)
- **Web browser** with access to [LabelZoom Web App](https://www.labelzoom.net/app)

## 🔧 Configuration

The print agent uses the following default settings:

| Setting | Default | Description |
|---------|---------|-------------|
| HTTP Port | `8080` | Port the agent listens on for print jobs |
| Printer Port | `9100` | Standard RAW printing port for thermal printers |
| CORS Origins | Multiple | Allowed origins for web requests |

### Allowed CORS Origins

The agent accepts requests from:
- `https://labelzoom.net`
- `https://www.labelzoom.net`
- `http://local.labelzoom.net`
- `http://localhost`
- `http://localhost:3000`

## 🖨️ Supported Printers

The agent works with any thermal printer that supports RAW printing over TCP/IP, including:
- **Zebra** printers (ZPL language)
- **Datamax** printers
- **SATO** printers
- Other direct thermal/thermal transfer printers with network connectivity

## 🧪 Testing

### Run Tests

```bash
go test -v ./...
```

### Test the Agent

Once running, test the `/ping` endpoint:

```bash
curl http://localhost:52045/ping
```

Expected response:
```json
{"message":"pong"}
```

### Send a Test Print Job

```bash
curl -X POST http://localhost:52045/print \
  -H "Content-Type: application/json" \
  -d '{
    "printerHostname": "192.168.1.100",
    "text": "^XA^FO50,50^ADN,36,20^FDTest Label^FS^XZ"
  }'
```

Replace `192.168.1.100` with your printer's IP address.

## 🏗️ Development

### Project Structure

```
lz-print-agent-go/
├── main.go              # Main application code
├── main_test.go         # Unit tests
├── resources/           # Embedded resources (logo)
├── Dockerfile           # Multi-stage Docker build
├── go.mod               # Go module definition
└── .github/
    └── workflows/       # CI/CD pipeline
```

### Building Docker Image

```bash
docker build -t lz-print-agent .
```

The final image is only **~7.5MB** thanks to:
- Multi-stage build
- Scratch base image
- Statically compiled Go binary

### CI/CD

This project uses GitHub Actions for:
- ✅ Automated testing on push/PR
- ✅ Docker image building
- ✅ Multi-platform builds (amd64, arm64)
- ✅ Automatic publishing to Docker Hub on release tags

## 📦 Docker Hub

Pre-built images are available on Docker Hub:

**🔗 [hub.docker.com/r/labelzoom/lz-print-agent-local](https://hub.docker.com/r/labelzoom/lz-print-agent-local)**

Pull the latest version:

```bash
docker pull labelzoom/lz-print-agent-local:latest
```

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- **LabelZoom Web App**: [www.labelzoom.net/app](https://www.labelzoom.net/app)
- **Docker Hub**: [hub.docker.com/r/labelzoom/lz-print-agent-local](https://hub.docker.com/r/labelzoom/lz-print-agent-local)
- **Issues**: [GitHub Issues](https://github.com/yourusername/lz-print-agent-go/issues)

## 💡 How It Works

```
┌─────────────┐      HTTP POST       ┌──────────────────┐      TCP (9100)      ┌─────────┐
│   Browser   │ ───────────────────> │  Print Agent     │ ───────────────────> │ Printer │
│ (LabelZoom) │   (localhost:52045)  │  (Go Service)    │   (RAW ZPL data)     │ (Zebra) │
└─────────────┘                      └──────────────────┘                      └─────────┘
```

1. User creates a label in the LabelZoom web app
2. Browser sends print job to local agent via HTTP
3. Agent forwards ZPL/RAW data to printer via TCP
4. Printer prints the label

---
**Made with ❤️ for the LabelZoom community**
