export type MoveClassification = 
  | 'brilliant'
  | 'great'
  | 'best'
  | 'excellent'
  | 'good'
  | 'book'
  | 'inaccuracy'
  | 'mistake'
  | 'blunder'
  | 'miss';

export type TacticalPattern = 
  | 'fork'
  | 'pin'
  | 'skewer'
  | 'discovery'
  | 'double_attack'
  | 'deflection'
  | 'decoy'
  | 'sacrifice'
  | 'clearance'
  | 'interference'
  | 'zugzwang'
  | 'stalemate_trick'
  | 'back_rank'
  | 'smothered_mate'
  | 'none';

export interface TacticalAnalysis {
  patterns: TacticalPattern[];
  isForcing: boolean;
  isTactical: boolean;
  threatLevel: 'low' | 'medium' | 'high' | 'critical';
  description?: string;
}

export interface EngineEvaluation {
  score: number; // Centipawns from white's perspective
  depth: number;
  bestMove: string;
  principalVariation: string[];
  nodes: number;
  time: number;
  mate?: number; // Mate in X moves
}

export interface MoveAnalysis {
  move: string;
  san: string;
  evaluation: EngineEvaluation;
  classification: MoveClassification;
  tacticalAnalysis?: TacticalAnalysis;
  alternativeMoves?: {
    move: string;
    evaluation: EngineEvaluation;
  }[];
  comment?: string;
}

export interface PlayerStatistics {
  accuracy: number;
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
  // For UI compatibility
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
  // Tactical statistics
  tacticalMoves?: number;
  forcingMoves?: number;
  criticalMoments?: number;
}

export interface GameAnalysis {
  moves: MoveAnalysis[];
  whiteStats: PlayerStatistics;
  blackStats: PlayerStatistics;
  openingAnalysis?: {
    name: string;
    eco: string;
    accuracy: number;
  };
  middlegameAnalysis?: {
    accuracy: number;
    criticalMoments: number[];
  };
  endgameAnalysis?: {
    accuracy: number;
    technicalMoves: number;
  };
  gamePhases: {
    opening: number; // Move number where opening ends
    middlegame: number; // Move number where middlegame ends
    endgame: number; // Move number where endgame starts
  };
  // Enhanced analysis data for Phase 4
  criticalMoments: number[];
  evaluationHistory: EngineEvaluation[];
  phaseAnalysis: {
    openingAccuracy: number;
    middlegameAccuracy: number;
    endgameAccuracy: number;
  };
  gameResult?: {
    result: '1-0' | '0-1' | '1/2-1/2' | '*';
    termination: string;
    winningAdvantage?: number; // Max advantage achieved
  };
}

export interface StockfishConfig {
  depth: number;
  time: number; // Time limit in milliseconds
  threads: number;
  hash: number; // Hash table size in MB
}

export interface AnalysisProgress {
  currentMove: number;
  totalMoves: number;
  progress: number; // Percentage (0-100)
  estimatedTimeRemaining?: number; // Seconds
}

export interface CriticalPosition {
  moveNumber: number;
  beforeEval: number;
  afterEval: number;
  advantage: 'white' | 'black';
  description: string;
} 