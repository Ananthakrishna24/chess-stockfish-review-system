'use client';

import { useState, useCallback, useEffect, useRef } from 'react';
import { useChessGame } from './useChessGame';
import { GameAnalysis, MoveAnalysis, PlayerStatistics, EngineEvaluation, AnalysisProgress } from '@/types/analysis';
import { apiClient, ApiError } from '@/lib/api';

interface AnalysisOptions {
  depth?: number;
  timePerMove?: number;
  includeBookMoves?: boolean;
  includeTacticalAnalysis?: boolean;
}

export function useGameAnalysis() {
  const chessGame = useChessGame();
  
  const [gameAnalysis, setGameAnalysis] = useState<GameAnalysis | null>(null);
  const [isAnalyzingGame, setIsAnalyzingGame] = useState(false);
  const [analysisError, setAnalysisError] = useState<string | null>(null);
  const [analysisProgress, setAnalysisProgress] = useState<number>(0);
  const [currentGameId, setCurrentGameId] = useState<string | null>(null);
  
  // Polling ref for checking analysis progress
  const progressIntervalRef = useRef<NodeJS.Timeout | null>(null);
  const abortControllerRef = useRef<AbortController | null>(null);

  // Load analysis from localStorage on mount
  useEffect(() => {
    const savedAnalysis = localStorage.getItem('chess-analysis');
    const savedGameState = localStorage.getItem('chess-game-state');
    
    if (savedAnalysis && savedGameState) {
      try {
        const analysis = JSON.parse(savedAnalysis);
        const gameState = JSON.parse(savedGameState);
        
        // Verify the saved analysis matches current game
        if (gameState.pgn && chessGame.gameState?.pgn === gameState.pgn) {
          setGameAnalysis(analysis);
        }
      } catch (error) {
        console.error('Failed to load saved analysis:', error);
        localStorage.removeItem('chess-analysis');
        localStorage.removeItem('chess-game-state');
      }
    }
  }, [chessGame.gameState?.pgn]);

  // Save analysis to localStorage when it changes
  useEffect(() => {
    if (gameAnalysis && chessGame.gameState) {
      try {
        localStorage.setItem('chess-analysis', JSON.stringify(gameAnalysis));
        localStorage.setItem('chess-game-state', JSON.stringify({
          pgn: chessGame.gameState.pgn,
          timestamp: Date.now()
        }));
      } catch (error) {
        console.error('Failed to save analysis:', error);
      }
    }
  }, [gameAnalysis, chessGame.gameState]);

  // Clean up polling on unmount
  useEffect(() => {
    return () => {
      if (progressIntervalRef.current) {
        clearInterval(progressIntervalRef.current);
      }
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, []);

  const analyzeCurrentPosition = useCallback(async (): Promise<EngineEvaluation | null> => {
    if (!chessGame.currentPosition) return null;
    
    try {
      const result = await apiClient.analyzePosition({
        fen: chessGame.currentPosition,
        depth: 15,
        multiPv: 1,
        timeLimit: 5000
      });
      
      return result.evaluation;
    } catch (error) {
      console.error('Position analysis failed:', error);
      setAnalysisError(error instanceof ApiError ? error.message : 'Position analysis failed');
      return null;
    }
  }, [chessGame.currentPosition]);

  const pollAnalysisProgress = useCallback(async (gameId: string) => {
    try {
      const response = await apiClient.getAnalysisProgress(gameId);
      
      if (response.status === 'completed') {
        // Analysis is complete, get the results
        const analysisResult = await apiClient.getGameAnalysis(gameId);
        setGameAnalysis(analysisResult.analysis);
        setIsAnalyzingGame(false);
        setAnalysisProgress(100);
        setCurrentGameId(null);
        
        if (progressIntervalRef.current) {
          clearInterval(progressIntervalRef.current);
          progressIntervalRef.current = null;
        }
        
        return true; // Analysis complete
      } else if (response.status === 'failed') {
        throw new Error('Analysis failed on server');
      } else {
        // Still analyzing, update progress
        const progressPercent = response.progress.percentage || 
          (response.progress.currentMove / response.progress.totalMoves) * 100;
        setAnalysisProgress(progressPercent);
        return false; // Still analyzing
      }
    } catch (error) {
      console.error('Failed to check analysis progress:', error);
      setAnalysisError(error instanceof ApiError ? error.message : 'Failed to check analysis progress');
      setIsAnalyzingGame(false);
      setCurrentGameId(null);
      
      if (progressIntervalRef.current) {
        clearInterval(progressIntervalRef.current);
        progressIntervalRef.current = null;
      }
      
      return true; // Stop polling on error
    }
  }, []);

  const analyzeCompleteGame = useCallback(async (options?: AnalysisOptions & { pgn?: string }) => {
    // Use provided PGN or fall back to gameState
    const gameState = chessGame.gameState;
    const pgnToAnalyze = options?.pgn || gameState?.pgn;
    
    if (!pgnToAnalyze) {
      setAnalysisError('No game or PGN provided for analysis');
      return;
    }

    setIsAnalyzingGame(true);
    setAnalysisError(null);
    setAnalysisProgress(0);

    // Create abort controller for this analysis session
    abortControllerRef.current = new AbortController();

    try {
      // Extract player ratings from game info if available
      const playerRatings = gameState ? {
        white: gameState.gameInfo.whiteRating,
        black: gameState.gameInfo.blackRating
      } : undefined;

      console.log('Starting game analysis with PGN:', pgnToAnalyze.substring(0, 100) + '...');

      // Start game analysis via API
      const analysisResponse = await apiClient.analyzeGame({
        pgn: pgnToAnalyze,
        options: {
          depth: options?.depth || 15,
          timePerMove: options?.timePerMove || 1000,
          includeBookMoves: options?.includeBookMoves ?? true,
          includeTacticalAnalysis: options?.includeTacticalAnalysis ?? true,
          playerRatings
        }
      });

      console.log('Analysis started with gameId:', analysisResponse.gameId);
      setCurrentGameId(analysisResponse.gameId);

      // Start polling for progress
      progressIntervalRef.current = setInterval(async () => {
        const isComplete = await pollAnalysisProgress(analysisResponse.gameId);
        if (isComplete && progressIntervalRef.current) {
          clearInterval(progressIntervalRef.current);
          progressIntervalRef.current = null;
        }
      }, 1000); // Poll every second

    } catch (error) {
      console.error('Game analysis failed:', error);
      setAnalysisError(error instanceof ApiError ? error.message : 'Game analysis failed');
      setIsAnalyzingGame(false);
      setCurrentGameId(null);
    }
  }, [chessGame.gameState, pollAnalysisProgress]);

  const stopAnalysis = useCallback(() => {
    // Stop polling
    if (progressIntervalRef.current) {
      clearInterval(progressIntervalRef.current);
      progressIntervalRef.current = null;
    }
    
    // Abort any ongoing requests
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
      abortControllerRef.current = null;
    }
    
    setIsAnalyzingGame(false);
    setAnalysisProgress(0);
    setCurrentGameId(null);
    
    console.log('Analysis stopped by user');
  }, []);

  // Helper function to get best move for current position
  const getBestMove = useCallback(async (): Promise<string | null> => {
    const evaluation = await analyzeCurrentPosition();
    return evaluation?.bestMove || null;
  }, [analyzeCurrentPosition]);

  // Helper function to classify a move (kept for backward compatibility)
  const classifyMove = useCallback((
    positionBefore: EngineEvaluation,
    positionAfter: EngineEvaluation,
    playedMove: string,
    bestMove: string,
    playerRating: number = 1500
  ) => {
    // Since classification is now done on the server, this is mainly for legacy support
    // The actual classification will be provided by the API analysis results
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

  // Helper function to calculate accuracy (kept for backward compatibility)
  const calculateAccuracy = useCallback((evaluations: EngineEvaluation[]): number => {
    if (evaluations.length === 0) return 0;
    
    // Simple accuracy calculation based on evaluation consistency
    let totalAccuracy = 0;
    for (let i = 1; i < evaluations.length; i++) {
      const scoreDiff = Math.abs(evaluations[i].score - evaluations[i-1].score);
      const moveAccuracy = Math.max(0, 100 - scoreDiff / 10);
      totalAccuracy += moveAccuracy;
    }
    
    return totalAccuracy / (evaluations.length - 1);
  }, []);

  return {
    // Analysis state
    gameAnalysis,
    isAnalyzingGame,
    analysisError,
    analysisProgress,
    currentGameId,
    
    // Analysis functions
    analyzeCompleteGame,
    analyzeCurrentPosition,
    stopAnalysis,
    
    // Helper functions (for backward compatibility)
    getBestMove,
    classifyMove,
    calculateAccuracy,
    
    // Chess game state (delegated)
    ...chessGame
  };
}

// Export helper functions for backward compatibility
export const calculatePlayerStats = (playerMoves: MoveAnalysis[]): PlayerStatistics => {
  const stats: PlayerStatistics = {
    accuracy: 0,
    brilliant: 0,
    great: 0,
    best: 0,
    excellent: 0,
    good: 0,
    book: 0,
    inaccuracy: 0,
    mistake: 0,
    blunder: 0,
    miss: 0,
    moveCounts: {
      brilliant: 0,
      great: 0,
      best: 0,
      excellent: 0,
      good: 0,
      book: 0,
      inaccuracy: 0,
      mistake: 0,
      blunder: 0,
      miss: 0
    },
    tacticalMoves: 0,
    forcingMoves: 0,
    criticalMoments: 0
  };

  if (playerMoves.length === 0) return stats;

  // Count moves by classification
  playerMoves.forEach(move => {
    const classification = move.classification;
    stats[classification]++;
    stats.moveCounts[classification]++;
    
    if (move.tacticalAnalysis?.isTactical) {
      stats.tacticalMoves = (stats.tacticalMoves || 0) + 1;
    }
    if (move.tacticalAnalysis?.isForcing) {
      stats.forcingMoves = (stats.forcingMoves || 0) + 1;
    }
  });

  // Calculate accuracy based on move quality
  const totalMoves = playerMoves.length;
  const weightedScore = 
    stats.brilliant * 100 +
    stats.great * 95 +
    stats.best * 90 +
    stats.excellent * 85 +
    stats.good * 80 +
    stats.book * 85 +
    stats.inaccuracy * 70 +
    stats.mistake * 50 +
    stats.blunder * 20 +
    stats.miss * 30;

  stats.accuracy = totalMoves > 0 ? weightedScore / totalMoves : 0;

  return stats;
}; 