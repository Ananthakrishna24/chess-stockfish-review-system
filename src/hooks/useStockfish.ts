'use client';

import { useState, useCallback, useEffect, useRef } from 'react';
import { EngineEvaluation, AnalysisProgress, StockfishConfig } from '@/types/analysis';
import { apiClient, ApiError } from '@/lib/api';

// Compatibility layer for the old Stockfish hook
// This maintains the same interface but uses the backend API
export function useStockfish(initialConfig?: Partial<StockfishConfig>) {
  const [isReady, setIsReady] = useState(true); // API is always "ready"
  const [isInitializing, setIsInitializing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentEvaluation, setCurrentEvaluation] = useState<EngineEvaluation | null>(null);
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [analysisProgress, setAnalysisProgress] = useState<AnalysisProgress>({
    currentMove: 0,
    totalMoves: 0,
    isAnalyzing: false,
    progress: 0
  });

  const configRef = useRef<StockfishConfig>({
    depth: 15,
    time: 1000,
    threads: 1,
    hash: 128,
    ...initialConfig
  });

  const updateConfig = useCallback((newConfig: Partial<StockfishConfig>) => {
    configRef.current = { ...configRef.current, ...newConfig };
    console.log('Engine config updated (API mode):', configRef.current);
  }, []);

  const analyzePosition = useCallback(async (
    fen: string, 
    depth?: number
  ): Promise<EngineEvaluation | null> => {
    setIsAnalyzing(true);
    setError(null);

    try {
      const result = await apiClient.analyzePosition({
        fen,
        depth: depth || configRef.current.depth,
        multiPv: 1,
        timeLimit: configRef.current.time
      });
      
      setCurrentEvaluation(result.evaluation);
      return result.evaluation;
    } catch (err) {
      const errorMessage = err instanceof ApiError ? err.message : 'Position analysis failed';
      setError(errorMessage);
      console.error('Position analysis error:', err);
      return null;
    } finally {
      setIsAnalyzing(false);
    }
  }, []);

  const analyzeGame = useCallback(async (
    positions: string[],
    onProgress?: (progress: AnalysisProgress) => void
  ): Promise<EngineEvaluation[]> => {
    // This method is deprecated in favor of the new API-based game analysis
    // But we keep it for backward compatibility
    console.warn('analyzeGame is deprecated. Use the new API-based game analysis instead.');
    
    setError(null);
    setAnalysisProgress({
      currentMove: 0,
      totalMoves: positions.length,
      isAnalyzing: true,
      progress: 0
    });

    const evaluations: EngineEvaluation[] = [];
    
    try {
      for (let i = 0; i < positions.length; i++) {
        const position = positions[i];
        const evaluation = await analyzePosition(position);
        
        if (evaluation) {
          evaluations.push(evaluation);
        }

        const progress = {
          currentMove: i + 1,
          totalMoves: positions.length,
          isAnalyzing: true,
          progress: ((i + 1) / positions.length) * 100
        };

        setAnalysisProgress(progress);
        onProgress?.(progress);

        // Small delay to prevent overwhelming the API
        await new Promise(resolve => setTimeout(resolve, 100));
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Game analysis failed';
      setError(errorMessage);
      console.error('Game analysis error:', err);
    } finally {
      setAnalysisProgress(prev => ({
        ...prev,
        isAnalyzing: false
      }));
    }

    return evaluations;
  }, [analyzePosition]);

  const stopAnalysis = useCallback(() => {
    setIsAnalyzing(false);
    setAnalysisProgress(prev => ({
      ...prev,
      isAnalyzing: false
    }));
    console.log('Analysis stopped (API mode)');
  }, []);

  const getBestMove = useCallback(async (fen: string): Promise<string | null> => {
    try {
      const evaluation = await analyzePosition(fen);
      return evaluation?.bestMove || null;
    } catch (err) {
      console.error('Best move analysis error:', err);
      return null;
    }
  }, [analyzePosition]);

  const classifyMove = useCallback((
    positionBefore: EngineEvaluation,
    positionAfter: EngineEvaluation,
    playedMove: string,
    bestMove: string,
    playerRating: number = 1500
  ) => {
    // Simple classification for backward compatibility
    // The real classification is now done on the server
    const beforeScore = positionBefore.score;
    const afterScore = -positionAfter.score; // Flip for current player
    const scoreDiff = afterScore - beforeScore;
    
    if (playedMove === bestMove) return 'best';
    if (scoreDiff >= 50) return 'excellent';
    if (scoreDiff >= 0) return 'good';
    if (scoreDiff >= -50) return 'inaccuracy';
    if (scoreDiff >= -150) return 'mistake';
    return 'blunder';
  }, []);

  const calculateAccuracy = useCallback((evaluations: EngineEvaluation[]): number => {
    if (evaluations.length === 0) return 0;
    
    // Simple accuracy calculation for backward compatibility
    let totalAccuracy = 0;
    for (let i = 1; i < evaluations.length; i++) {
      const scoreDiff = Math.abs(evaluations[i].score - evaluations[i-1].score);
      const moveAccuracy = Math.max(0, 100 - scoreDiff / 10);
      totalAccuracy += moveAccuracy;
    }
    
    return totalAccuracy / (evaluations.length - 1);
  }, []);

  // Mock engine object for backward compatibility
  const engine = {
    analyzeTacticalPatterns: () => ({
      patterns: [],
      isForcing: false,
      isTactical: false,
      threatLevel: 'low' as const,
      description: 'Tactical analysis not available in API mode'
    }),
    detectCriticalMoments: () => [],
    analyzeGamePhases: () => ({
      opening: 10,
      middlegame: 25,
      endgame: 40,
      openingAccuracy: 85,
      middlegameAccuracy: 80,
      endgameAccuracy: 90
    }),
    classifyMove: classifyMove,
    calculateAccuracy: calculateAccuracy
  };

  return {
    // State
    isReady,
    isInitializing,
    error,
    currentEvaluation,
    isAnalyzing,
    analysisProgress,
    engine,
    
    // Methods
    updateConfig,
    analyzePosition,
    analyzeGame,
    stopAnalysis,
    getBestMove,
    classifyMove,
    calculateAccuracy
  };
} 