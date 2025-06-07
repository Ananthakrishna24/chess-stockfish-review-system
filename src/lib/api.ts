import { GameAnalysis, EngineEvaluation, AnalysisProgress } from '@/types/analysis';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

export interface GameAnalysisRequest {
  pgn: string;
  options?: {
    depth?: number;
    timePerMove?: number;
    includeBookMoves?: boolean;
    includeTacticalAnalysis?: boolean;
    playerRatings?: {
      white?: number;
      black?: number;
    };
  };
}

export interface GameAnalysisResponse {
  gameId: string;
  status: 'queued' | 'analyzing' | 'completed' | 'failed';
  message: string;
}

export interface GameAnalysisResult {
  gameId: string;
  gameInfo: {
    white: string;
    black: string;
    whiteRating?: number;
    blackRating?: number;
    result: string;
    date: string;
    event?: string;
    opening?: string;
    eco?: string;
  };
  analysis: GameAnalysis;
  processingTime: number;
  timestamp: string;
}

export interface PositionAnalysisRequest {
  fen: string;
  depth?: number;
  multiPv?: number;
  timeLimit?: number;
}

export interface PositionAnalysisResult {
  fen: string;
  evaluation: EngineEvaluation;
  alternativeMoves?: Array<{
    move: string;
    san: string;
    evaluation: EngineEvaluation;
  }>;
  positionInfo?: {
    phase: string;
    material: { white: number; black: number };
    safety: { whiteKing: string; blackKing: string };
  };
}

export interface PlayerStatsResult {
  playerName: string;
  gamesAnalyzed: number;
  averageAccuracy: number;
  ratingRange?: {
    min: number;
    max: number;
    current: number;
  };
  recentGames: Array<{
    gameId: string;
    opponent: string;
    result: string;
    accuracy: number;
    date: string;
    opening: string;
    eco: string;
  }>;
  strengths: string[];
  weaknesses: string[];
  improvementSuggestions: string[];
  phasePerformance: {
    openingAccuracy: number;
    middlegameAccuracy: number;
    endgameAccuracy: number;
  };
  openingRepertoire: Record<string, {
    eco: string;
    name: string;
    games: number;
    accuracy: number;
    results: { wins: number; draws: number; losses: number };
  }>;
  tacticalStats: {
    totalTacticalMoves: number;
    totalForcingMoves: number;
    totalCriticalMoments: number;
    brilliantMoves: number;
    blunderRate: number;
  };
  lastUpdated: string;
}

export interface Opening {
  eco: string;
  name: string;
  variation?: string;
  moves: string[];
  popularity: number;
  statistics: {
    white: number;
    draw: number;
    black: number;
  };
  theory?: string;
  keyIdeas?: string[];
}

class ApiError extends Error {
  constructor(public status: number, message: string, public details?: string) {
    super(message);
    this.name = 'ApiError';
  }
}

class ApiClient {
  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    
    const response = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    });

    if (!response.ok) {
      let errorMessage = `HTTP ${response.status}: ${response.statusText}`;
      let errorDetails: string | undefined;
      
      try {
        const errorData = await response.json();
        errorMessage = errorData.error || errorMessage;
        errorDetails = errorData.details;
      } catch {
        // If we can't parse error JSON, use the default message
      }
      
      throw new ApiError(response.status, errorMessage, errorDetails);
    }

    return response.json();
  }

  // Game Analysis APIs
  async analyzeGame(request: GameAnalysisRequest): Promise<GameAnalysisResponse> {
    return this.request<GameAnalysisResponse>('/games/analyze', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  async getGameAnalysis(gameId: string): Promise<GameAnalysisResult> {
    return this.request<GameAnalysisResult>(`/games/analyze/${gameId}`);
  }

  async getAnalysisProgress(gameId: string): Promise<{
    gameId: string;
    status: string;
    progress: AnalysisProgress;
  }> {
    return this.request(`/games/analyze/${gameId}/progress`);
  }

  // Position Analysis APIs
  async analyzePosition(request: PositionAnalysisRequest): Promise<PositionAnalysisResult> {
    return this.request<PositionAnalysisResult>('/positions/analyze', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  // Opening Database APIs
  async searchOpenings(params: {
    eco?: string;
    fen?: string;
    moves?: string;
    name?: string;
  }): Promise<{ results: Opening[]; count: number }> {
    const searchParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value) searchParams.append(key, value);
    });

    const query = searchParams.toString();
    return this.request<{ results: Opening[]; count: number }>(
      `/openings/search${query ? `?${query}` : ''}`
    );
  }

  async getOpening(eco: string): Promise<Opening> {
    return this.request<Opening>(`/openings/${eco}`);
  }

  async getAllOpenings(): Promise<{ results: Opening[]; count: number }> {
    return this.request<{ results: Opening[]; count: number }>('/openings');
  }

  async getOpeningCategories(): Promise<{
    categories: Record<string, Opening[]>;
    total: number;
  }> {
    return this.request('/openings/categories');
  }

  // Player Statistics APIs
  async getPlayerStats(playerName: string): Promise<PlayerStatsResult> {
    return this.request<PlayerStatsResult>(`/stats/player/${encodeURIComponent(playerName)}`);
  }

  async getPlayerGames(playerName: string): Promise<{
    playerName: string;
    games: Array<{
      gameId: string;
      opponent: string;
      result: string;
      accuracy: number;
      date: string;
      opening: string;
      eco: string;
    }>;
    totalGames: number;
  }> {
    return this.request(`/stats/player/${encodeURIComponent(playerName)}/games`);
  }

  async getAllPlayers(): Promise<{ players: string[]; count: number }> {
    return this.request<{ players: string[]; count: number }>('/stats/players');
  }

  async getLeaderboard(limit?: number): Promise<{
    rankings: Array<{
      playerName: string;
      gamesAnalyzed: number;
      averageAccuracy: number;
      currentRating?: number;
    }>;
    count: number;
    limit: number;
  }> {
    const query = limit ? `?limit=${limit}` : '';
    return this.request(`/stats/leaderboard${query}`);
  }

  // Engine Configuration APIs
  async getEngineConfig(): Promise<{
    version: string;
    features: string[];
    limits: {
      maxDepth: number;
      maxTime: number;
      maxNodes: number;
    };
    currentConfig: {
      threads: number;
      hash: number;
      contempt: number;
      analysisContempt: string;
    };
  }> {
    return this.request('/engine/config');
  }

  async updateEngineConfig(config: {
    threads?: number;
    hash?: number;
    contempt?: number;
    analysisContempt?: string;
  }): Promise<{
    message: string;
    config: typeof config;
  }> {
    return this.request('/engine/config', {
      method: 'POST',
      body: JSON.stringify(config),
    });
  }

  // Health Check APIs
  async getHealth(): Promise<{
    status: string;
    timestamp: string;
    services: {
      stockfish: string;
      cache: string;
      database: string;
    };
    uptime: number;
    version: string;
  }> {
    return this.request('/health');
  }

  async getSimpleHealth(): Promise<{
    status: string;
    timestamp: string;
  }> {
    return this.request('/health', { headers: {} }); // Use simple endpoint
  }
}

export const apiClient = new ApiClient();
export { ApiError }; 