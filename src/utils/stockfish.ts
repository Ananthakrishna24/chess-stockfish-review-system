import { EngineEvaluation, MoveClassification } from '@/types/analysis';

// Helper functions for chess analysis
// The actual engine is now on the backend

export function convertScoreToString(score: number, mate?: number): string {
  if (mate !== undefined) {
    return mate > 0 ? `M${mate}` : `M${Math.abs(mate)}`;
  }
  
  const absScore = Math.abs(score);
  if (absScore < 10) {
    return (score / 100).toFixed(2);
  } else if (absScore < 100) {
    return (score / 100).toFixed(1);
  } else {
    return (score / 100).toFixed(0);
  }
}

export function getScoreColor(score: number): string {
  if (score > 100) return 'text-green-400';
  if (score > 50) return 'text-green-300';
  if (score > 0) return 'text-green-200';
  if (score === 0) return 'text-gray-400';
  if (score > -50) return 'text-red-200';
  if (score > -100) return 'text-red-300';
  return 'text-red-400';
}

// Legacy compatibility - these functions are now handled by the backend
export function classifyMove(
  positionBefore: EngineEvaluation,
  positionAfter: EngineEvaluation,
  playedMove: string,
  bestMove: string,
  playerRating: number = 1500
): MoveClassification {
  // Simple classification for backward compatibility
  const beforeScore = positionBefore.score;
  const afterScore = -positionAfter.score; // Flip for current player
  const scoreDiff = afterScore - beforeScore;
  
  if (playedMove === bestMove) return 'best';
  if (scoreDiff >= 50) return 'excellent';
  if (scoreDiff >= 0) return 'good';
  if (scoreDiff >= -50) return 'inaccuracy';
  if (scoreDiff >= -150) return 'mistake';
  return 'blunder';
}

export function calculateAccuracy(evaluations: EngineEvaluation[]): number {
  if (evaluations.length === 0) return 0;
  
  // Simple accuracy calculation for backward compatibility
  let totalAccuracy = 0;
  for (let i = 1; i < evaluations.length; i++) {
    const scoreDiff = Math.abs(evaluations[i].score - evaluations[i-1].score);
    const moveAccuracy = Math.max(0, 100 - scoreDiff / 10);
    totalAccuracy += moveAccuracy;
  }
  
  return totalAccuracy / (evaluations.length - 1);
}

// Deprecated - engine is now on the backend
export async function getStockfishEngine(): Promise<null> {
  console.warn('getStockfishEngine is deprecated. Use the API client instead.');
  return null;
} 