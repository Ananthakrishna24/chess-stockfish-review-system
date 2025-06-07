# Setup Guide - Chess Analysis Backend

## 📋 Prerequisites Installation

### 1. Install Go

#### Option A: Using Package Manager (Ubuntu/Pop OS)
```bash
sudo apt update
sudo apt install golang-go

# Verify installation
go version
```

#### Option B: Install Latest Go (Recommended)
```bash
# Download and install Go 1.21+
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH=$PATH:/usr/local/go/bin

# Reload shell or run:
source ~/.bashrc

# Verify installation
go version
```

### 2. Install Stockfish

```bash
# Option A: Using package manager
sudo apt install stockfish

# Option B: Using our script (recommended)
chmod +x scripts/install-stockfish.sh
./scripts/install-stockfish.sh

# Verify installation
stockfish quit
```

## 🚀 Quick Start

### 1. Install Dependencies
```bash
cd backend
go mod download
```

### 2. Run the Server
```bash
# Development mode
go run cmd/server/main.go

# Or build and run
go build -o chess-backend cmd/server/main.go
./chess-backend
```

### 3. Test the API
```bash
# Health check
curl http://localhost:8080/health

# Position analysis test
curl -X POST http://localhost:8080/api/positions/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
    "depth": 12
  }'
```

## 🔧 Environment Configuration

Create a `.env` file or set environment variables:

```bash
# Server Configuration
export SERVER_PORT=8080
export APP_MODE=debug

# Engine Configuration  
export ENGINE_BINARY_PATH=/usr/local/bin/stockfish  # or just "stockfish"
export ENGINE_MAX_WORKERS=4
export ENGINE_DEFAULT_DEPTH=15
export ENGINE_THREADS=1
export ENGINE_HASH_SIZE_MB=128

# Rate Limiting
export RATE_LIMIT_GAME_ANALYSIS_PER_HOUR=10
export RATE_LIMIT_POSITION_ANALYSIS_PER_HOUR=100
```

## 🧪 Development Workflow

### Build and Test
```bash
# Clean dependencies
go mod tidy

# Format code
go fmt ./...

# Build
go build -o chess-backend cmd/server/main.go

# Run tests (when implemented)
go test ./...

# Run with custom config
ENGINE_BINARY_PATH=stockfish SERVER_PORT=8081 go run cmd/server/main.go
```

### Project Structure
```
backend/
├── cmd/server/main.go           # Server entry point
├── configs/config.go            # Configuration management
├── internal/
│   ├── handlers/               # HTTP handlers
│   │   ├── analysis.go         # Game/position analysis endpoints
│   │   └── health.go           # Health check endpoints
│   ├── middleware/             # HTTP middleware
│   │   └── ratelimit.go        # Rate limiting
│   ├── models/                 # Data structures
│   │   ├── analysis.go         # Analysis request/response models
│   │   └── game.go             # Game state and job models
│   └── services/               # Business logic
│       ├── analysis.go         # Analysis orchestration
│       ├── cache.go            # In-memory caching
│       ├── chess.go            # PGN parsing and chess logic
│       └── stockfish.go        # Stockfish engine management
├── pkg/uci/                    # UCI engine communication
│   └── engine.go               # Stockfish UCI interface
├── scripts/
│   └── install-stockfish.sh    # Stockfish installation script
├── go.mod                      # Go module definition
└── README.md                   # Documentation
```

## 🔍 Troubleshooting

### Common Issues

1. **Go not found**
   ```
   Command 'go' not found
   ```
   Solution: Install Go following instructions above

2. **Stockfish not found**
   ```
   Failed to initialize Stockfish service
   ```
   Solution: Install Stockfish and ensure it's in PATH

3. **Port in use**
   ```
   bind: address already in use
   ```
   Solution: Use different port with `SERVER_PORT=8081`

4. **Permission denied on script**
   ```
   Permission denied: ./scripts/install-stockfish.sh
   ```
   Solution: `chmod +x scripts/install-stockfish.sh`

### Debug Mode

Run with debug logging:
```bash
APP_MODE=debug go run cmd/server/main.go
```

### Testing Endpoints

#### 1. Health Check
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### 2. Position Analysis
```bash
curl -X POST http://localhost:8080/api/positions/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
    "depth": 12
  }'
```

#### 3. Game Analysis
```bash
curl -X POST http://localhost:8080/api/games/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "pgn": "[Event \"Test Game\"]\n[White \"Player1\"]\n[Black \"Player2\"]\n[Result \"*\"]\n\n1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 *"
  }'
```

#### 4. Check Analysis Progress
```bash
# Replace GAME_ID with actual ID from previous response
curl http://localhost:8080/api/games/analyze/GAME_ID/progress
```

## 🚧 Known Limitations (Phase 1)

- In-memory only caching (no persistence)
- Basic move classification algorithm
- No opening database integration
- HTTP polling for progress (no WebSockets)
- No user authentication
- Simple rate limiting

## 📈 Next Steps

After successful setup:

1. Test all API endpoints
2. Integrate with frontend
3. Plan Phase 2 enhancements:
   - Redis caching
   - PostgreSQL persistence
   - Advanced analysis features
   - WebSocket support

## 🆘 Getting Help

If you encounter issues:
1. Check the logs for error messages
2. Verify Stockfish installation: `stockfish quit`
3. Ensure Go is properly installed: `go version`
4. Check if ports are available: `lsof -i :8080`
5. Try running with debug mode enabled 