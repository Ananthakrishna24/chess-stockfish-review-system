import { EngineEvaluation, StockfishConfig, MoveClassification } from '@/types/analysis';

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

export class StockfishEngine {
  private isReady = false;
  private config: StockfishConfig = {
    depth: 15,
    time: 1000,
    threads: 1,
    hash: 128
  };

  constructor(config?: Partial<StockfishConfig>) {
    if (config) {
      this.config = { ...this.config, ...config };
    }
  }

  async initialize(): Promise<void> {
    // Only initialize on client side
    if (typeof window === 'undefined') {
      throw new Error('Stockfish can only be initialized on the client side');
    }

    // Mock initialization - simulate delay
    await new Promise(resolve => setTimeout(resolve, 1000));
    this.isReady = true;
    console.log('Mock Stockfish engine initialized');
  }

  async analyzePosition(fen: string, depth?: number): Promise<EngineEvaluation> {
    if (!this.isReady) {
      throw new Error('Stockfish engine not ready');
    }

    // Mock analysis - simulate analysis time
    const analysisTime = Math.random() * 500 + 200; // 200-700ms
    await new Promise(resolve => setTimeout(resolve, analysisTime));
    
    // Generate mock evaluation data
    const mockEvaluation: EngineEvaluation = {
      score: Math.floor(Math.random() * 400 - 200), // Score between -200 and +200
      depth: depth || this.config.depth,
      bestMove: this.generateMockMove(),
      principalVariation: [this.generateMockMove(), this.generateMockMove(), this.generateMockMove()],
      nodes: Math.floor(Math.random() * 1000000 + 50000),
      time: analysisTime,
      mate: Math.random() < 0.05 ? Math.floor(Math.random() * 10 + 1) : undefined
    };

    return mockEvaluation;
  }

  private generateMockMove(): string {
    const files = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'];
    const ranks = ['1', '2', '3', '4', '5', '6', '7', '8'];
    
    const fromFile = files[Math.floor(Math.random() * files.length)];
    const fromRank = ranks[Math.floor(Math.random() * ranks.length)];
    const toFile = files[Math.floor(Math.random() * files.length)];
    const toRank = ranks[Math.floor(Math.random() * ranks.length)];
    
    return `${fromFile}${fromRank}${toFile}${toRank}`;
  }

  async findBestMove(fen: string): Promise<string> {
    const evaluation = await this.analyzePosition(fen);
    return evaluation.bestMove;
  }

  classifyMove(
    positionBefore: EngineEvaluation,
    positionAfter: EngineEvaluation,
    playedMove: string,
    bestMove: string
  ): MoveClassification {
    const scoreDiff = Math.abs(positionBefore.score - positionAfter.score);
    const isPlayerTurn = positionBefore.score > 0; // Assuming white to move
    
    // Adjust score based on whose turn it is
    const adjustedScoreBefore = isPlayerTurn ? positionBefore.score : -positionBefore.score;
    const adjustedScoreAfter = isPlayerTurn ? -positionAfter.score : positionAfter.score;
    const evaluation = adjustedScoreAfter - adjustedScoreBefore;

    // Check if played move is the best move
    if (playedMove === bestMove) {
      if (evaluation > 200) return 'brilliant';
      if (evaluation > 100) return 'great';
      return 'best';
    }

    // Classify based on evaluation loss
    if (evaluation >= -50) return 'good';
    if (evaluation >= -100) return 'inaccuracy';
    if (evaluation >= -250) return 'mistake';
    if (evaluation >= -500) return 'blunder';
    
    return 'miss';
  }

  // New tactical pattern recognition methods
  analyzeTacticalPatterns(
    positionBefore: EngineEvaluation,
    positionAfter: EngineEvaluation,
    playedMove: string
  ): TacticalAnalysis {
    const patterns: TacticalPattern[] = [];
    let isForcing = false;
    let isTactical = false;
    let threatLevel: TacticalAnalysis['threatLevel'] = 'low';
    let description = '';

    const scoreDiff = Math.abs(positionBefore.score - positionAfter.score);
    const evaluationSwing = positionAfter.score - positionBefore.score;

    // Detect if move is tactical based on evaluation swing
    if (scoreDiff > 100) {
      isTactical = true;
      
      if (scoreDiff > 500) {
        threatLevel = 'critical';
        patterns.push('sacrifice'); // Large material exchanges often involve sacrifices
      } else if (scoreDiff > 300) {
        threatLevel = 'high';
        patterns.push('double_attack'); // Significant advantage often from double attacks
      } else if (scoreDiff > 150) {
        threatLevel = 'medium';
        patterns.push('fork'); // Moderate tactical gains often from forks
      }
    }

    // Detect forcing moves (checks, captures, threats)
    if (positionAfter.mate !== undefined) {
      isForcing = true;
      isTactical = true;
      threatLevel = 'critical';
      
      if (Math.abs(positionAfter.mate) <= 3) {
        patterns.push('back_rank');
        description = `Mate in ${Math.abs(positionAfter.mate)}`;
      } else if (Math.abs(positionAfter.mate) === 1) {
        patterns.push('smothered_mate');
        description = 'Forced mate in 1';
      }
    }

    // Mock tactical pattern detection based on move characteristics
    const moveString = playedMove.toLowerCase();
    
    // Simple heuristics for demonstration
    if (this.containsCapture(moveString)) {
      isTactical = true;
      patterns.push('double_attack');
    }
    
    if (this.isCheckMove(moveString)) {
      isForcing = true;
      isTactical = true;
      patterns.push('discovery');
    }

    // Detect sacrificial patterns
    if (evaluationSwing < -200 && positionAfter.score > positionBefore.score + 300) {
      patterns.push('sacrifice');
      description = 'Tactical sacrifice for positional advantage';
    }

    // Detect pins and skewers (mock detection)
    if (Math.random() < 0.1 && isTactical) {
      patterns.push(Math.random() < 0.5 ? 'pin' : 'skewer');
    }

    // Detect deflection and decoy patterns
    if (scoreDiff > 200 && !patterns.includes('sacrifice')) {
      patterns.push(Math.random() < 0.5 ? 'deflection' : 'decoy');
    }

    // If no specific patterns detected but move is tactical
    if (isTactical && patterns.length === 0) {
      patterns.push('none');
    }

    return {
      patterns: patterns.length > 0 ? patterns : ['none'],
      isForcing,
      isTactical,
      threatLevel,
      description
    };
  }

  private containsCapture(move: string): boolean {
    // Mock capture detection - in real implementation would check FEN
    return move.includes('x') || Math.random() < 0.2;
  }

  private isCheckMove(move: string): boolean {
    // Mock check detection - in real implementation would check if move gives check
    return move.includes('+') || Math.random() < 0.15;
  }

  detectCriticalMoments(evaluations: EngineEvaluation[]): number[] {
    const criticalMoments: number[] = [];
    
    for (let i = 1; i < evaluations.length; i++) {
      const prev = evaluations[i - 1];
      const curr = evaluations[i];
      
      // Large evaluation swings indicate critical moments
      const swing = Math.abs(curr.score - prev.score);
      
      if (swing > 200) {
        criticalMoments.push(i);
      }
      
      // Mate threats
      if (curr.mate !== undefined && Math.abs(curr.mate) <= 5) {
        criticalMoments.push(i);
      }
      
      // Evaluation crosses zero (advantage change)
      if ((prev.score > 50 && curr.score < -50) || (prev.score < -50 && curr.score > 50)) {
        criticalMoments.push(i);
      }
    }
    
    return criticalMoments;
  }

  analyzeGamePhases(moves: any[], evaluations: EngineEvaluation[]): {
    opening: number;
    middlegame: number;
    endgame: number;
    openingAccuracy: number;
    middlegameAccuracy: number;
    endgameAccuracy: number;
  } {
    // Simple heuristics for game phase detection
    const totalMoves = moves.length;
    
    // Opening typically ends around move 10-15
    const openingEnd = Math.min(Math.floor(totalMoves * 0.25), 15);
    
    // Endgame typically starts when few pieces remain (mock detection)
    const endgameStart = Math.max(Math.floor(totalMoves * 0.75), openingEnd + 10);
    
    // Calculate phase accuracies
    const openingEvals = evaluations.slice(0, openingEnd);
    const middlegameEvals = evaluations.slice(openingEnd, endgameStart);
    const endgameEvals = evaluations.slice(endgameStart);
    
    return {
      opening: openingEnd,
      middlegame: endgameStart,
      endgame: totalMoves,
      openingAccuracy: this.calculateAccuracy(openingEvals),
      middlegameAccuracy: this.calculateAccuracy(middlegameEvals),
      endgameAccuracy: this.calculateAccuracy(endgameEvals)
    };
  }

  calculateAccuracy(evaluations: EngineEvaluation[]): number {
    if (evaluations.length === 0) return 0;

    let totalLoss = 0;
    let moveCount = 0;

    for (let i = 1; i < evaluations.length; i++) {
      const prevEval = evaluations[i - 1];
      const currEval = evaluations[i];
      
      // Calculate evaluation loss (from perspective of player who moved)
      const isWhiteMove = i % 2 === 1;
      const scoreBefore = isWhiteMove ? prevEval.score : -prevEval.score;
      const scoreAfter = isWhiteMove ? -currEval.score : currEval.score;
      
      const loss = Math.max(0, scoreBefore - scoreAfter);
      totalLoss += loss;
      moveCount++;
    }

    const averageLoss = totalLoss / moveCount;
    
    // Convert to accuracy percentage
    // Formula inspired by chess.com's accuracy calculation
    const accuracy = Math.max(0, 100 - (averageLoss / 10));
    
    return Math.round(accuracy * 10) / 10;
  }

  stop(): void {
    console.log('Mock Stockfish analysis stopped');
  }

  quit(): void {
    this.isReady = false;
    console.log('Mock Stockfish engine quit');
  }
}

// Singleton instance for global use
let stockfishInstance: StockfishEngine | null = null;

export async function getStockfishEngine(): Promise<StockfishEngine> {
  if (!stockfishInstance) {
    stockfishInstance = new StockfishEngine();
    await stockfishInstance.initialize();
  }
  return stockfishInstance;
}

export function convertScoreToString(score: number, mate?: number): string {
  if (mate !== undefined) {
    return `M${mate}`;
  }
  
  const pawnValue = score / 100;
  const sign = pawnValue >= 0 ? '+' : '';
  
  return `${sign}${pawnValue.toFixed(1)}`;
}

export function getScoreColor(score: number): string {
  if (Math.abs(score) < 50) return 'text-gray-600';
  if (score > 0) return 'text-green-600';
  return 'text-red-600';
} 