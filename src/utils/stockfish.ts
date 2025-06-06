import { EngineEvaluation, StockfishConfig, MoveClassification } from '@/types/analysis';

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