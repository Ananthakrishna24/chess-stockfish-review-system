# Chess Review System - API Documentation

## 🚀 Overview

This document provides comprehensive documentation for the Chess Review System backend API. The API provides chess game analysis, position evaluation, opening database lookups, and player statistics tracking.

**Base URL**: `http://localhost:8080/api`

## 📊 Endpoints Summary

### Game Analysis
- `POST /api/games/analyze` - Analyze a complete chess game
- `GET /api/games/analyze/{gameId}` - Get analysis results
- `GET /api/games/analyze/{gameId}/progress` - Get analysis progress

### Position Analysis
- `POST /api/positions/analyze` - Analyze a single position

### Opening Database
- `GET /api/openings/search` - Search openings by criteria
- `GET /api/openings/{eco}` - Get opening by ECO code
- `GET /api/openings` - Get all openings
- `GET /api/openings/categories` - Get openings by ECO categories

### Player Statistics
- `GET /api/stats/player/{playername}` - Get player statistics
- `GET /api/stats/player/{playername}/games` - Get player game history
- `GET /api/stats/players` - Get all tracked players
- `GET /api/stats/leaderboard` - Get player rankings

### Engine Configuration
- `GET /api/engine/config` - Get engine configuration
- `POST /api/engine/config` - Update engine configuration

### System Health
- `GET /api/health` - Health check
- `GET /health` - Simple health check

---

## 🎯 Game Analysis Endpoints

### POST /api/games/analyze

Analyzes a complete chess game from PGN notation.

**Request Body:**
```json
{
  "pgn": "1. e4 e5 2. Nf3 Nc6 3. Bb5 a6...",
  "options": {
    "depth": 15,
    "timePerMove": 1000,
    "includeBookMoves": true,
    "includeTacticalAnalysis": true,
    "playerRatings": {
      "white": 1500,
      "black": 1600
    }
  }
}
```

**Response (202 Accepted):**
```json
{
  "gameId": "abc123def456",
  "status": "queued",
  "message": "Analysis started successfully"
}
```

### GET /api/games/analyze/{gameId}

Retrieves completed analysis results.

**Response (200 OK):**
```json
{
  "gameId": "abc123def456",
  "gameInfo": {
    "white": "Player1",
    "black": "Player2",
    "whiteRating": 1500,
    "blackRating": 1600,
    "result": "1-0",
    "date": "2024.01.15",
    "event": "Casual Game",
    "opening": "Italian Game",
    "eco": "C50"
  },
  "analysis": {
    "moves": [
      {
        "moveNumber": 1,
        "move": "e2e4",
        "san": "e4",
        "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
        "evaluation": {
          "score": 25,
          "depth": 15,
          "bestMove": "e7e5",
          "principalVariation": ["e7e5", "g1f3"],
          "nodes": 150000,
          "time": 1000
        },
        "classification": "book",
        "alternativeMoves": [
          {
            "move": "e7e5",
            "san": "e5",
            "evaluation": {
              "score": 20,
              "depth": 15
            }
          }
        ],
        "tacticalAnalysis": {
          "patterns": [],
          "isForcing": false,
          "isTactical": false,
          "threatLevel": "low",
          "description": "Opening move controlling the center"
        }
      }
    ],
    "whiteStats": {
      "accuracy": 85.2,
      "moveCounts": {
        "brilliant": 0,
        "great": 2,
        "best": 15,
        "excellent": 8,
        "good": 5,
        "book": 3,
        "inaccuracy": 2,
        "mistake": 1,
        "blunder": 0,
        "miss": 0
      },
      "tacticalMoves": 4,
      "forcingMoves": 6,
      "criticalMoments": 2
    },
    "blackStats": {
      "accuracy": 78.9,
      "moveCounts": {
        "brilliant": 1,
        "great": 1,
        "best": 12,
        "excellent": 6,
        "good": 7,
        "book": 3,
        "inaccuracy": 3,
        "mistake": 2,
        "blunder": 1,
        "miss": 0
      }
    },
    "openingAnalysis": {
      "name": "Italian Game",
      "eco": "C50",
      "accuracy": 92.1,
      "theory": "Classical opening focusing on rapid development",
      "deviationMove": 8
    },
    "gamePhases": {
      "opening": 10,
      "middlegame": 25,
      "endgame": 35
    },
    "phaseAnalysis": {
      "openingAccuracy": 90.5,
      "middlegameAccuracy": 82.3,
      "endgameAccuracy": 88.7
    },
    "criticalMoments": [
      {
        "moveNumber": 15,
        "beforeEval": 50,
        "afterEval": -120,
        "advantage": "black",
        "description": "Missed tactical opportunity"
      }
    ],
    "evaluationHistory": []
  },
  "processingTime": 45.2,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### GET /api/games/analyze/{gameId}/progress

Gets real-time analysis progress.

**Response (200 OK):**
```json
{
  "gameId": "abc123def456",
  "status": "analyzing",
  "progress": {
    "currentMove": 15,
    "totalMoves": 40,
    "percentage": 37.5,
    "estimatedTimeRemaining": 30
  }
}
```

---

## 🎯 Position Analysis Endpoints

### POST /api/positions/analyze

Analyzes a single chess position.

**Request Body:**
```json
{
  "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
  "depth": 15,
  "multiPv": 3,
  "timeLimit": 5000
}
```

**Response (200 OK):**
```json
{
  "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
  "evaluation": {
    "score": 25,
    "depth": 15,
    "bestMove": "e7e5",
    "principalVariation": ["e7e5", "g1f3", "b8c6"],
    "nodes": 250000,
    "time": 5000
  },
  "alternativeMoves": [
    {
      "move": "e7e5",
      "san": "e5",
      "evaluation": {
        "score": 25,
        "depth": 15
      }
    },
    {
      "move": "c7c5",
      "san": "c5",
      "evaluation": {
        "score": 15,
        "depth": 15
      }
    }
  ],
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

---

## 🎯 Opening Database Endpoints

### GET /api/openings/search

Search openings by various criteria.

**Query Parameters:**
- `eco` - ECO code (e.g., "B20")
- `fen` - Position FEN
- `moves` - Move sequence (space-separated)
- `name` - Opening name (fuzzy search)

**Examples:**
- `/api/openings/search?eco=B20`
- `/api/openings/search?name=sicilian`
- `/api/openings/search?moves=e4 c5`

**Response (200 OK):**
```json
{
  "results": [
    {
      "eco": "B20",
      "name": "Sicilian Defense",
      "variation": "",
      "moves": ["e4", "c5"],
      "popularity": 16.8,
      "statistics": {
        "white": 52.3,
        "draw": 23.1,
        "black": 24.6
      },
      "theory": "Most popular response to e4, creates imbalanced positions",
      "keyIdeas": [
        "Asymmetrical structure",
        "Counterplay",
        "Sharp positions"
      ]
    }
  ],
  "count": 1
}
```

### GET /api/openings/{eco}

Get specific opening by ECO code.

**Response (200 OK):**
```json
{
  "eco": "B20",
  "name": "Sicilian Defense",
  "variation": "",
  "moves": ["e4", "c5"],
  "popularity": 16.8,
  "statistics": {
    "white": 52.3,
    "draw": 23.1,
    "black": 24.6
  },
  "theory": "Most popular response to e4, creates imbalanced positions",
  "keyIdeas": [
    "Asymmetrical structure",
    "Counterplay",
    "Sharp positions"
  ]
}
```

### GET /api/openings

Get all openings in the database.

**Response (200 OK):**
```json
{
  "results": [
    // Array of all opening objects
  ],
  "count": 12
}
```

### GET /api/openings/categories

Get openings grouped by ECO categories.

**Response (200 OK):**
```json
{
  "categories": {
    "A": [
      // Flank openings (A00-A99)
    ],
    "B": [
      // Semi-open games (B00-B99)
    ],
    "C": [
      // Open games (C00-C99)
    ],
    "D": [
      // Closed games (D00-D99)
    ],
    "E": [
      // Indian defenses (E00-E99)
    ]
  },
  "total": 5
}
```

---

## 🎯 Player Statistics Endpoints

### GET /api/stats/player/{playername}

Get comprehensive player statistics.

**Response (200 OK):**
```json
{
  "playerName": "magnus_carlsen",
  "gamesAnalyzed": 25,
  "averageAccuracy": 92.3,
  "ratingRange": {
    "min": 2800,
    "max": 2850,
    "current": 2835
  },
  "recentGames": [
    {
      "gameId": "game123",
      "opponent": "hikaru_nakamura",
      "result": "win",
      "accuracy": 94.2,
      "date": "2024.01.15",
      "opening": "Queen's Gambit",
      "eco": "D06"
    }
  ],
  "strengths": [
    "endgame technique",
    "positional play",
    "opening preparation"
  ],
  "weaknesses": [
    "time management"
  ],
  "improvementSuggestions": [
    "Practice tactical puzzles daily",
    "Study complex endgame positions"
  ],
  "phasePerformance": {
    "openingAccuracy": 95.1,
    "middlegameAccuracy": 91.2,
    "endgameAccuracy": 96.8
  },
  "openingRepertoire": {
    "D06": {
      "eco": "D06",
      "name": "Queen's Gambit",
      "games": 8,
      "accuracy": 93.5,
      "results": {
        "wins": 6,
        "draws": 2,
        "losses": 0
      }
    }
  },
  "tacticalStats": {
    "totalTacticalMoves": 45,
    "totalForcingMoves": 78,
    "totalCriticalMoments": 12,
    "brilliantMoves": 3,
    "blunderRate": 0.2
  },
  "lastUpdated": "2024-01-15T10:30:00Z"
}
```

### GET /api/stats/player/{playername}/games

Get player's game history.

**Response (200 OK):**
```json
{
  "playerName": "magnus_carlsen",
  "games": [
    {
      "gameId": "game123",
      "opponent": "hikaru_nakamura",
      "result": "win",
      "accuracy": 94.2,
      "date": "2024.01.15",
      "opening": "Queen's Gambit",
      "eco": "D06"
    }
  ],
  "totalGames": 25
}
```

### GET /api/stats/players

Get list of all tracked players.

**Response (200 OK):**
```json
{
  "players": [
    "magnus_carlsen",
    "hikaru_nakamura",
    "fabiano_caruana"
  ],
  "count": 3
}
```

### GET /api/stats/leaderboard

Get player rankings by accuracy.

**Query Parameters:**
- `limit` - Number of players to return (default: 10, max: 100)

**Response (200 OK):**
```json
{
  "rankings": [
    {
      "playerName": "magnus_carlsen",
      "gamesAnalyzed": 25,
      "averageAccuracy": 92.3,
      "currentRating": 2835
    },
    {
      "playerName": "hikaru_nakamura",
      "gamesAnalyzed": 18,
      "averageAccuracy": 89.7,
      "currentRating": 2780
    }
  ],
  "count": 2,
  "limit": 10
}
```

---

## 🎯 Engine Configuration Endpoints

### GET /api/engine/config

Get current engine configuration.

**Response (200 OK):**
```json
{
  "version": "Stockfish 16",
  "features": ["UCI", "MultiPV", "Hash", "Threads"],
  "limits": {
    "maxDepth": 24,
    "maxTime": 30000,
    "maxNodes": 10000000
  },
  "currentConfig": {
    "threads": 4,
    "hash": 128,
    "contempt": 0,
    "analysisContempt": "off"
  }
}
```

### POST /api/engine/config

Update engine configuration.

**Request Body:**
```json
{
  "threads": 8,
  "hash": 256,
  "contempt": 10,
  "analysisContempt": "white"
}
```

**Response (200 OK):**
```json
{
  "message": "Configuration updated successfully",
  "config": {
    "threads": 8,
    "hash": 256,
    "contempt": 10,
    "analysisContempt": "white"
  }
}
```

---

## 🎯 System Health Endpoints

### GET /api/health

Comprehensive health check.

**Response (200 OK):**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "services": {
    "stockfish": "healthy",
    "cache": "healthy",
    "database": "healthy"
  },
  "uptime": 3600,
  "version": "1.0.0"
}
```

### GET /health

Simple health check.

**Response (200 OK):**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## 🚨 Error Responses

All endpoints return consistent error responses:

### 400 Bad Request
```json
{
  "error": "Invalid request format",
  "details": "PGN cannot be empty"
}
```

### 404 Not Found
```json
{
  "error": "Analysis not found",
  "details": "Analysis job not found"
}
```

### 429 Too Many Requests
```json
{
  "error": "Rate limit exceeded",
  "details": "Too many requests, please try again later"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error",
  "details": "Stockfish engine unavailable"
}
```

---

## 🔧 Rate Limits

**⚠️ RATE LIMITING CURRENTLY DISABLED**

Rate limiting has been disabled for development purposes. All endpoints accept unlimited requests.

*To re-enable rate limiting, uncomment the middleware line in `cmd/server/main.go`:*
```go
// router.Use(middleware.RateLimit(cfg.RateLimit))
```

When enabled, the following limits will apply:
- **Game Analysis**: 10,000 requests/hour per IP
- **Position Analysis**: 100,000 requests/hour per IP  
- **Opening Lookups**: 1,000,000 requests/hour per IP
- **Player Statistics**: 500,000 requests/hour per IP
- **General Endpoints**: 1,000,000 requests/hour per IP

---

## 📝 Move Classifications

The analysis engine classifies moves using the following categories:

- **Brilliant** (!!): Exceptional moves, often sacrificial
- **Great** (!): Very strong moves
- **Best**: Engine's top choice
- **Excellent**: Near-optimal moves
- **Good**: Solid, reasonable moves
- **Book**: Opening theory moves
- **Inaccuracy** (?!): Slightly suboptimal
- **Mistake** (?): Clear error
- **Blunder** (??): Major error
- **Miss**: Missed opportunity

---

## 🎯 Game Phases

Games are automatically divided into phases:

- **Opening**: First 10-15 moves (theory-based)
- **Middlegame**: Complex tactical/strategic phase
- **Endgame**: Simplified positions with few pieces

---

## 🔍 Tactical Patterns

The system recognizes common tactical patterns:

- Fork, Pin, Skewer
- Discovery, Double Attack
- Deflection, Decoy
- Clearance, Interference
- Zugzwang, Stalemate tricks

---

## 🚀 Getting Started

1. **Start the server**: `./server`
2. **Analyze a game**: POST to `/api/games/analyze` with PGN
3. **Check progress**: GET `/api/games/analyze/{gameId}/progress`
4. **Get results**: GET `/api/games/analyze/{gameId}`
5. **View player stats**: GET `/api/stats/player/{playername}`

---

## 📊 Example Workflow

```bash
# 1. Analyze a game
curl -X POST http://localhost:8080/api/games/analyze \
  -H "Content-Type: application/json" \
  -d '{"pgn": "1. e4 e5 2. Nf3 Nc6..."}'

# 2. Check progress
curl http://localhost:8080/api/games/analyze/abc123/progress

# 3. Get results
curl http://localhost:8080/api/games/analyze/abc123

# 4. Search openings
curl "http://localhost:8080/api/openings/search?eco=B20"

# 5. View player stats
curl http://localhost:8080/api/stats/player/magnus_carlsen
```

This API provides a complete chess analysis solution with comprehensive game evaluation, opening database integration, and player performance tracking. 