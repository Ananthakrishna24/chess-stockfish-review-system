# Chess Review System - Backend API Requirements

## 🎯 Overview
This document outlines the API requirements for the Chess Review System backend, designed to replace the current frontend-only stockfish.js implementation with a robust server-side analysis engine.

## 📊 Current System Analysis

### Current Data Flow
The existing frontend system processes:
1. **PGN Import** → Parse chess game notation
2. **Position Analysis** → Stockfish evaluation of each position
3. **Move Classification** → Categorize moves (brilliant, great, best, mistake, blunder, etc.)
4. **Game Statistics** → Calculate player accuracy and move distributions
5. **Visualization Data** → Generate evaluation charts and analysis panels

### Core Data Structures Used

#### Game State
```typescript
interface GameState {
  moves: ChessMove[];
  currentMoveIndex: number;
  gameInfo: GameInfo;
  pgn: string;
  startingFen?: string;
}

interface GameInfo {
  white: string;
  black: string;
  whiteRating?: number;
  blackRating?: number;
  result: string; // '1-0', '0-1', '1/2-1/2', or '*'
  date?: string;
  event?: string;
  site?: string;
  opening?: string;
  eco?: string;
}
```

#### Analysis Results
```typescript
interface GameAnalysis {
  moves: MoveAnalysis[];
  whiteStats: PlayerStatistics;
  blackStats: PlayerStatistics;
  openingAnalysis: {
    name: string;
    eco: string;
    accuracy: number;
  };
  gamePhases: {
    opening: number;
    middlegame: number; 
    endgame: number;
  };
  criticalMoments: number[];
  evaluationHistory: EngineEvaluation[];
  phaseAnalysis: {
    openingAccuracy: number;
    middlegameAccuracy: number;
    endgameAccuracy: number;
  };
}

interface PlayerStatistics {
  accuracy: number;
  moveCounts: {
    brilliant: number;
    great: number;
    best: number;
    excellent: number;
    good: number;
    book: number;
    inaccuracy: number;
    mistake: number;
    blunder: number;
    miss: number;
  };
  tacticalMoves?: number;
  forcingMoves?: number;
  criticalMoments?: number;
}
```

---

## 🚀 API Endpoints

### 1. Game Analysis Endpoints

#### POST /api/games/analyze
**Purpose**: Analyze a complete chess game from PGN
**Request**:
```json
{
  "pgn": "string", // Complete PGN notation
  "options": {
    "depth": number, // Engine depth (4-24, default: 15)
    "timePerMove": number, // Max time per position in ms (default: 1000)
    "includeBookMoves": boolean, // Include opening book analysis (default: true)
    "includeTacticalAnalysis": boolean, // Deep tactical pattern analysis (default: true)
    "playerRatings": {
      "white": number, // For accuracy calculation adjustment
      "black": number
    }
  }
}
```

**Response**:
```json
{
  "gameId": "string", // Unique identifier for caching
  "gameInfo": {
    "white": "string",
    "black": "string", 
    "whiteRating": number,
    "blackRating": number,
    "result": "string",
    "date": "string",
    "event": "string",
    "site": "string",
    "opening": "string",
    "eco": "string"
  },
  "analysis": {
    "moves": [
      {
        "moveNumber": number,
        "move": "string", // UCI notation (e.g., "e2e4")
        "san": "string", // Standard algebraic notation (e.g., "e4")
        "fen": "string", // Position after move
        "evaluation": {
          "score": number, // Centipawns from white's perspective
          "depth": number,
          "bestMove": "string",
          "principalVariation": ["string"],
          "nodes": number,
          "time": number,
          "mate": number // Optional: mate in X moves
        },
        "classification": "brilliant|great|best|excellent|good|book|inaccuracy|mistake|blunder|miss",
        "alternativeMoves": [
          {
            "move": "string",
            "evaluation": {
              "score": number,
              "bestMove": "string"
            }
          }
        ],
        "tacticalAnalysis": {
          "patterns": ["fork|pin|skewer|discovery|..."],
          "isForcing": boolean,
          "isTactical": boolean,
          "threatLevel": "low|medium|high|critical",
          "description": "string"
        },
        "comment": "string" // Optional analysis comment
      }
    ],
    "whiteStats": {
      "accuracy": number, // Percentage (0-100)
      "moveCounts": {
        "brilliant": number,
        "great": number,
        "best": number,
        "excellent": number,
        "good": number,
        "book": number,
        "inaccuracy": number,
        "mistake": number,
        "blunder": number,
        "miss": number
      },
      "tacticalMoves": number,
      "forcingMoves": number,
      "criticalMoments": number
    },
    "blackStats": {}, // Same structure as whiteStats
    "openingAnalysis": {
      "name": "string",
      "eco": "string", // ECO code
      "accuracy": number,
      "theory": "string", // Opening theory assessment
      "deviationMove": number // First move out of book
    },
    "gamePhases": {
      "opening": number, // Last move of opening
      "middlegame": number, // Last move of middlegame  
      "endgame": number // First move of endgame
    },
    "phaseAnalysis": {
      "openingAccuracy": number,
      "middlegameAccuracy": number,
      "endgameAccuracy": number
    },
    "criticalMoments": [
      {
        "moveNumber": number,
        "beforeEval": number,
        "afterEval": number,
        "advantage": "white|black",
        "description": "string"
      }
    ],
    "evaluationHistory": [] // Complete evaluation for each position
  },
  "processingTime": number, // Analysis time in seconds
  "timestamp": "string" // ISO 8601
}
```

#### GET /api/games/analyze/{gameId}
**Purpose**: Retrieve cached analysis results
**Response**: Same as POST /api/games/analyze

#### GET /api/games/analyze/{gameId}/progress
**Purpose**: Get real-time analysis progress (for long-running analyses)
**Response**:
```json
{
  "gameId": "string",
  "status": "queued|analyzing|completed|failed",
  "progress": {
    "currentMove": number,
    "totalMoves": number,
    "percentage": number, // 0-100
    "estimatedTimeRemaining": number // seconds
  },
  "error": "string" // If status is "failed"
}
```

### 2. Position Analysis Endpoints

#### POST /api/positions/analyze
**Purpose**: Analyze a single chess position
**Request**:
```json
{
  "fen": "string", // FEN notation of position
  "depth": number, // Engine depth (default: 15)
  "multiPv": number, // Number of best moves to return (default: 3)
  "timeLimit": number // Analysis time limit in ms (default: 5000)
}
```

**Response**:
```json
{
  "fen": "string",
  "evaluation": {
    "score": number,
    "depth": number,
    "bestMove": "string",
    "principalVariation": ["string"],
    "nodes": number,
    "time": number,
    "mate": number
  },
  "alternativeMoves": [
    {
      "move": "string",
      "san": "string",
      "evaluation": {
        "score": number,
        "depth": number
      }
    }
  ],
  "positionInfo": {
    "phase": "opening|middlegame|endgame",
    "material": {
      "white": number,
      "black": number
    },
    "safety": {
      "whiteKing": "safe|exposed|danger",
      "blackKing": "safe|exposed|danger"
    }
  }
}
```

### 3. Opening Database Endpoints

#### GET /api/openings/search
**Purpose**: Get opening information by ECO code or position
**Query Parameters**:
- `eco`: ECO code (e.g., "B00")
- `fen`: Position FEN
- `moves`: Opening moves sequence

**Response**:
```json
{
  "eco": "string",
  "name": "string",
  "variation": "string",
  "moves": ["string"], // Move sequence
  "popularity": number, // Usage frequency percentage
  "statistics": {
    "white": number, // Win percentage
    "draw": number,
    "black": number
  },
  "theory": "string", // Opening assessment
  "keyIdeas": ["string"] // Main strategic ideas
}
```

### 4. Engine Configuration Endpoints

#### GET /api/engine/config
**Purpose**: Get current engine configuration
**Response**:
```json
{
  "version": "string", // Stockfish version
  "features": ["string"], // Supported features
  "limits": {
    "maxDepth": number,
    "maxTime": number,
    "maxNodes": number
  },
  "currentConfig": {
    "threads": number,
    "hash": number, // Hash table size in MB
    "contempt": number,
    "analysisContempt": "off|white|black|both"
  }
}
```

#### POST /api/engine/config
**Purpose**: Update engine configuration
**Request**:
```json
{
  "threads": number,
  "hash": number,
  "contempt": number,
  "analysisContempt": "off|white|black|both"
}
```

### 5. Statistics & Leaderboards

#### GET /api/stats/player/{playername}
**Purpose**: Get historical statistics for a player
**Response**:
```json
{
  "playerName": "string",
  "gamesAnalyzed": number,
  "averageAccuracy": number,
  "ratingRange": {
    "min": number,
    "max": number,
    "current": number
  },
  "recentGames": [
    {
      "gameId": "string",
      "opponent": "string",
      "result": "string",
      "accuracy": number,
      "date": "string"
    }
  ],
  "strengths": ["string"], // e.g., ["tactical", "endgame"]
  "weaknesses": ["string"], // e.g., ["opening", "time management"]
  "improvementSuggestions": ["string"]
}
```

---

## 🔧 Technical Requirements

### Infrastructure
- **Runtime**: Go 1.21+ with Gin or Echo web framework
- **Chess Engine**: Stockfish 16+ binary (system process execution)
- **Database**: Redis for caching, PostgreSQL for persistent data (optional for v1)
- **Queue System**: Go-based worker pools for analysis job management
- **Process Management**: OS process execution for Stockfish engine
- **Real-time**: HTTP polling for progress (WebSocket as future enhancement)

### Performance Requirements
- **Response Time**: 
  - Single position analysis: < 10 seconds
  - Complete game analysis: < 3 minutes (average game)
  - Cached results: < 200ms
- **Throughput**: Support 20+ concurrent analysis requests
- **Caching**: 24-hour cache for complete analyses, 1-hour for positions
- **Process Pool**: Maintain 4-8 Stockfish processes for concurrent analysis

### Security & Rate Limiting
- **Rate Limits**:
  - Game analysis: 10 requests/hour per IP
  - Position analysis: 100 requests/hour per IP
  - Opening lookups: 1000 requests/hour per IP
- **Input Validation**: Strict PGN and FEN validation
- **Resource Protection**: Analysis timeout and resource limits

### Data Storage
- **Analysis Cache**: Store complete analyses for 30 days
- **Game History**: Store game metadata indefinitely
- **Usage Analytics**: Track API usage and performance metrics

---

## 🎨 Integration Points

### Frontend Integration
The backend will replace the current `useStockfish` hook and provide:
1. **HTTP Polling** for real-time analysis progress (every 1-2 seconds)
2. **REST API Client** for game submission and result retrieval
3. **Caching Strategy** to avoid re-analyzing identical games
4. **Error Handling** for network issues and analysis failures

### Expected Frontend Changes
- Replace `stockfish.js` worker with HTTP API calls
- Add polling mechanism for progress updates
- Implement retry logic with exponential backoff
- Update loading states and error handling

---

## 📈 Future Enhancements

### Phase 2 Features (Future)
- **User Accounts**: Personal game history and statistics
- **Game Collections**: Organize and share analyzed games  
- **Opening Database**: ECO code integration
- **WebSocket Support**: Real-time progress updates
- **Advanced Caching**: Persistent storage with PostgreSQL

### Advanced Analysis (Future)
- **Multiple Engines**: Support for Leela Chess Zero, Komodo
- **Cloud Analysis**: Distributed analysis for complex positions
- **Tournament Analysis**: Batch processing for tournament games
- **Advanced Tactics**: Deep tactical pattern recognition

---

## 🚦 Implementation Priority

### Phase 1: Core Go API (Week 1-3)
1. **Go Web Server Setup**
   - Gin/Echo framework setup
   - Basic project structure
   - CORS and middleware configuration

2. **Stockfish Integration**
   - Process pool management
   - UCI protocol communication
   - Position analysis with timeout

3. **Core Endpoints**
   - `POST /api/games/analyze` - Basic game analysis
   - `GET /api/games/analyze/{gameId}/progress` - Progress tracking
   - `GET /api/games/analyze/{gameId}` - Results retrieval

4. **In-Memory Caching**
   - Simple map-based caching for development
   - Analysis result storage

### Phase 2: Production Features (Week 4-5)
1. **Redis Integration** for persistent caching
2. **Rate Limiting** and security middleware
3. **Advanced Move Classification** algorithm
4. **Error Handling** and logging

### Phase 3: Polish & Deploy (Week 6)
1. **Performance Optimization**
2. **Docker Containerization**
3. **API Documentation** (Swagger/OpenAPI)
4. **Testing** and debugging

---

## 💡 Technical Notes

### Move Classification Algorithm
Implement Chess.com's Expected Points Model:
- Convert centipawn evaluations to expected points using sigmoid function
- Adjust for player rating (higher-rated players need smaller advantages)
- Classify moves based on expected points change thresholds
- Special handling for brilliant moves (good sacrifices) and great moves

### Opening Recognition
- Integrate ECO opening database
- Track theoretical vs. practical move accuracy
- Identify first deviation from book theory
- Provide opening-specific advice and statistics

### Tactical Pattern Detection
- Implement pattern recognition for common tactics
- Classify forcing vs. non-forcing moves
- Detect critical moments and missed opportunities
- Generate natural language descriptions

---

## 🐹 Go Implementation Details

### Project Structure
```
chess-backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   ├── analysis.go
│   │   └── health.go
│   ├── services/
│   │   ├── stockfish.go
│   │   ├── chess.go
│   │   └── cache.go
│   ├── models/
│   │   ├── game.go
│   │   └── analysis.go
│   └── middleware/
│       ├── cors.go
│       └── ratelimit.go
├── pkg/
│   ├── uci/
│   │   └── engine.go
│   └── pgn/
│       └── parser.go
├── configs/
│   └── config.go
├── docker/
│   └── Dockerfile
├── scripts/
│   └── install-stockfish.sh
├── go.mod
├── go.sum
└── README.md
```

### Key Go Dependencies
```go
// Web Framework
github.com/gin-gonic/gin v1.9.1

// Chess Libraries
github.com/notnil/chess v1.9.0

// Caching (Phase 2)
github.com/go-redis/redis/v8 v8.11.5

// Configuration
github.com/spf13/viper v1.16.0

// Logging
github.com/sirupsen/logrus v1.9.3

// Rate Limiting
golang.org/x/time/rate

// Testing
github.com/stretchr/testify v1.8.4
```

### Stockfish Process Management
```go
type StockfishPool struct {
    processes   []*StockfishEngine
    available   chan *StockfishEngine
    maxWorkers  int
}

type StockfishEngine struct {
    cmd     *exec.Cmd
    stdin   io.WriteCloser
    stdout  io.ReadCloser
    scanner *bufio.Scanner
    mutex   sync.Mutex
}
```

### Core API Structure
```go
// Handler function
func (h *AnalysisHandler) AnalyzeGame(c *gin.Context) {
    var req AnalyzeGameRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Start async analysis
    gameID := h.analysisService.StartAnalysis(req.PGN, req.Options)
    
    c.JSON(202, gin.H{
        "gameId": gameID,
        "status": "queued"
    })
}
```

### Installation Script
```bash
#!/bin/bash
# scripts/install-stockfish.sh

# Download and install Stockfish binary
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    wget https://stockfishchess.org/files/stockfish_15_linux_x64_avx2.zip
elif [[ "$OSTYPE" == "darwin"* ]]; then
    wget https://stockfishchess.org/files/stockfish_15_mac.zip
fi

unzip stockfish_*.zip
chmod +x stockfish_*
sudo mv stockfish_* /usr/local/bin/stockfish
```

### Docker Configuration
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o chess-backend cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates stockfish
WORKDIR /root/

COPY --from=builder /app/chess-backend .
CMD ["./chess-backend"]
```

This Go-focused implementation provides a solid foundation for your chess review system while being more straightforward to implement and deploy than the Node.js version. 