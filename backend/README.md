# Chess Analysis Backend - Phase 1

A Go-based REST API for chess game analysis using Stockfish engine.

## 🚀 Features

- **Game Analysis**: Complete PGN analysis with move classification and statistics
- **Position Analysis**: Single position evaluation with multiple variations
- **Real-time Progress**: HTTP polling for analysis progress tracking
- **Engine Configuration**: Configurable Stockfish settings
- **Caching**: In-memory caching for improved performance
- **Rate Limiting**: Configurable API usage limits

## 📋 Requirements

- Go 1.21 or higher
- Stockfish 16 chess engine

## 🛠️ Installation

### 1. Install Stockfish

Run the provided installation script:

```bash
chmod +x scripts/install-stockfish.sh
./scripts/install-stockfish.sh
```

Or install manually:
- **Ubuntu/Debian**: `sudo apt-get install stockfish`
- **macOS**: `brew install stockfish`
- **Manual**: Download from [Stockfish releases](https://stockfishchess.org/download/)

### 2. Install Dependencies

```bash
go mod download
```

### 3. Build and Run

```bash
# Build
go build -o chess-backend cmd/server/main.go

# Run
./chess-backend
```

Or run directly:

```bash
go run cmd/server/main.go
```

## 🔧 Configuration

Configure via environment variables:

```bash
# Server settings
export SERVER_PORT=8080
export APP_MODE=debug  # or "release"

# Engine settings
export ENGINE_BINARY_PATH=stockfish
export ENGINE_MAX_WORKERS=4
export ENGINE_DEFAULT_DEPTH=15
export ENGINE_THREADS=1
export ENGINE_HASH_SIZE_MB=128

# Rate limiting
export RATE_LIMIT_GAME_ANALYSIS_PER_HOUR=10
export RATE_LIMIT_POSITION_ANALYSIS_PER_HOUR=100
```

## 🌐 API Endpoints

### Game Analysis

#### Start Game Analysis
```http
POST /api/games/analyze
Content-Type: application/json

{
  "pgn": "1. e4 e5 2. Nf3 Nc6...",
  "options": {
    "depth": 15,
    "timePerMove": 1000,
    "includeBookMoves": true,
    "includeTacticalAnalysis": true,
    "playerRatings": {
      "white": 1500,
      "black": 1400
    }
  }
}
```

Response:
```json
{
  "gameId": "abc123...",
  "status": "queued",
  "message": "Analysis started successfully"
}
```

#### Get Analysis Progress
```http
GET /api/games/analyze/{gameId}/progress
```

Response:
```json
{
  "gameId": "abc123...",
  "status": "analyzing",
  "progress": {
    "currentMove": 15,
    "totalMoves": 40,
    "percentage": 37.5,
    "estimatedTimeRemaining": 25
  }
}
```

#### Get Analysis Results
```http
GET /api/games/analyze/{gameId}
```

Response:
```json
{
  "gameId": "abc123...",
  "gameInfo": {
    "white": "Player 1",
    "black": "Player 2",
    "result": "1-0"
  },
  "analysis": {
    "moves": [...],
    "whiteStats": {
      "accuracy": 85.2,
      "moveCounts": {
        "brilliant": 1,
        "great": 3,
        "best": 8,
        "excellent": 5,
        "good": 2,
        "inaccuracy": 2,
        "mistake": 1,
        "blunder": 0
      }
    },
    "blackStats": {...},
    "criticalMoments": [...],
    "evaluationHistory": [...]
  },
  "processingTime": 23.5,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Position Analysis

#### Analyze Position
```http
POST /api/positions/analyze
Content-Type: application/json

{
  "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
  "depth": 18,
  "multiPv": 3,
  "timeLimit": 5000
}
```

Response:
```json
{
  "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
  "evaluation": {
    "score": 25,
    "depth": 18,
    "bestMove": "e7e5",
    "principalVariation": ["e7e5", "g1f3", "b8c6"],
    "nodes": 2847264,
    "time": 1250
  },
  "alternativeMoves": [...],
  "positionInfo": {
    "phase": "opening",
    "material": {
      "white": 39,
      "black": 39
    },
    "safety": {
      "whiteKing": "safe",
      "blackKing": "safe"
    }
  }
}
```

### Engine Configuration

#### Get Engine Config
```http
GET /api/engine/config
```

#### Update Engine Config
```http
POST /api/engine/config
Content-Type: application/json

{
  "threads": 2,
  "hash": 256,
  "contempt": 0,
  "analysisContempt": "off"
}
```

### Health & Stats

```http
GET /health
GET /api/health
GET /api/stats
```

## 📊 Rate Limits

Default limits per IP address:
- Game analysis: 10 requests/hour
- Position analysis: 100 requests/hour  
- Other endpoints: 1000 requests/hour

Rate limit headers:
- `X-RateLimit-Limit`: Maximum requests allowed
- `X-RateLimit-Type`: Type of rate limit applied
- `X-RateLimit-Remaining`: Requests remaining
- `X-RateLimit-Reset`: Reset time (Unix timestamp)

## 🧪 Testing

### Basic Test

```bash
# Health check
curl http://localhost:8080/health

# Position analysis
curl -X POST http://localhost:8080/api/positions/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
    "depth": 15
  }'

# Game analysis
curl -X POST http://localhost:8080/api/games/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "pgn": "1. e4 e5 2. Nf3 Nc6 3. Bb5 a6"
  }'
```

### Sample PGN for Testing

```pgn
[Event "Test Game"]
[Site "Local"]
[Date "2024.01.15"]
[Round "1"]
[White "Player1"]
[Black "Player2"]
[Result "1-0"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 
6. Re1 b5 7. Bb3 d6 8. c3 O-O 9. h3 Nb8 10. d4 Nbd7 
11. c4 c6 12. cxb5 axb5 13. Nc3 Bb7 14. Bg5 b4 15. Nb1 h6 
16. Bh4 c5 17. dxe5 Nxe4 18. Bxe7 Qxe7 19. exd6 Qf6 20. Nbd2 1-0
```

## 🔍 Troubleshooting

### Common Issues

1. **Stockfish not found**
   ```
   Failed to initialize Stockfish service: failed to create engine 0: failed to start engine: executable file not found in $PATH
   ```
   Solution: Install Stockfish or set `ENGINE_BINARY_PATH` environment variable

2. **Port already in use**
   ```
   Failed to start server: listen tcp :8080: bind: address already in use
   ```
   Solution: Set different port with `SERVER_PORT=8081` or kill process using port 8080

3. **Rate limit exceeded**
   ```json
   {"error": "Rate limit exceeded", "message": "Too many game_analysis requests. Limit: 10 per hour"}
   ```
   Solution: Wait for rate limit reset or increase limits via environment variables

### Debug Mode

Enable debug logging:
```bash
export APP_MODE=debug
go run cmd/server/main.go
```

## 🚧 Phase 1 Limitations

Current implementation includes:
- ✅ Basic game and position analysis
- ✅ In-memory caching
- ✅ Rate limiting
- ✅ Progress tracking
- ✅ Engine configuration

Not yet implemented (planned for Phase 2):
- ❌ Redis caching
- ❌ PostgreSQL persistence  
- ❌ Advanced tactical analysis
- ❌ Opening database integration
- ❌ WebSocket real-time updates
- ❌ User accounts and authentication

## 📈 Performance

Typical performance on modern hardware:
- Position analysis (depth 15): ~1-3 seconds
- Game analysis (40 moves): ~30-60 seconds
- Concurrent analysis: 4 games simultaneously
- Memory usage: ~50-100MB baseline + ~10MB per active analysis

## 🤝 Development

### Project Structure

```
backend/
├── cmd/server/           # Application entry point
├── configs/              # Configuration management
├── internal/
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # HTTP middleware
│   ├── models/          # Data structures
│   └── services/        # Business logic
├── pkg/uci/             # UCI engine communication
└── scripts/             # Setup and utility scripts
```

### Adding New Features

1. Define models in `internal/models/`
2. Implement business logic in `internal/services/`
3. Add HTTP handlers in `internal/handlers/`
4. Update routes in `cmd/server/main.go`

## 📝 License

This project is part of the Chess Review System implementation.

## 🔗 Related

- [Stockfish Chess Engine](https://stockfishchess.org/)
- [Chess.com Analysis](https://www.chess.com/analysis)
- [Frontend Integration Guide](../README.md) 