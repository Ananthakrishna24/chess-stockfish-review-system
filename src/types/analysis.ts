export type MoveClassification = 
  | 'brilliant'
  | 'great'
  | 'best'
  | 'good'
  | 'inaccuracy'
  | 'mistake'
  | 'blunder'
  | 'miss';

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
  good: number;
  inaccuracy: number;
  mistake: number;
  blunder: number;
  miss: number;
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
}

export interface AnalysisProgress {
  currentMove: number;
  totalMoves: number;
  isAnalyzing: boolean;
  progress: number; // 0-100
}

export interface StockfishConfig {
  depth: number;
  time: number;
  threads: number;
  hash: number;
}

export interface CriticalPosition {
  moveNumber: number;
  beforeEval: number;
  afterEval: number;
  advantage: 'white' | 'black';
  description: string;
} 