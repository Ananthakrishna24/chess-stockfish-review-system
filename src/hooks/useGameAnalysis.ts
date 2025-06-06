'use client';

import { useState, useCallback, useEffect } from 'react';
import { useChessGame } from './useChessGame';
import { useStockfish } from './useStockfish';
import { GameAnalysis, MoveAnalysis, PlayerStatistics, EngineEvaluation } from '@/types/analysis';
import { ChessGameManager } from '@/utils/chess';

export function useGameAnalysis() {
  const chessGame = useChessGame();
  const stockfish = useStockfish({ depth: 15, time: 1000 });
  
  const [gameAnalysis, setGameAnalysis] = useState<GameAnalysis | null>(null);
  const [isAnalyzingGame, setIsAnalyzingGame] = useState(false);
  const [analysisError, setAnalysisError] = useState<string | null>(null);

  const analyzeCurrentPosition = useCallback(async () => {
    if (!chessGame.currentPosition || !stockfish.isReady) return null;
    
    return await stockfish.analyzePosition(chessGame.currentPosition);
  }, [chessGame.currentPosition, stockfish]);

  const analyzeCompleteGame = useCallback(async () => {
    if (!chessGame.gameState || !stockfish.isReady) {
      setAnalysisError('Game or engine not ready');
      return;
    }

    setIsAnalyzingGame(true);
    setAnalysisError(null);

    try {
      const { moves } = chessGame.gameState;
      const gameManager = new ChessGameManager();
      
      // Get all positions in the game
      const positions: string[] = [];
      positions.push('rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1'); // Starting position
      
      // Load the game and get positions after each move
      gameManager.loadPGN(chessGame.gameState.pgn);
      for (let i = 0; i < moves.length; i++) {
        const position = gameManager.getPosition(i);
        positions.push(position);
      }

      // Analyze all positions
      const evaluations = await stockfish.analyzeGame(positions, (progress) => {
        // Progress callback could be used to update UI
        console.log(`Analysis progress: ${progress.progress.toFixed(1)}%`);
      });

      if (evaluations.length === 0) {
        throw new Error('Analysis failed - no evaluations received');
      }

      // Process move analysis
      const moveAnalyses: MoveAnalysis[] = [];
      
      for (let i = 0; i < moves.length; i++) {
        const move = moves[i];
        const positionBefore = evaluations[i];
        const positionAfter = evaluations[i + 1];
        
        if (positionBefore && positionAfter) {
          const bestMove = positionBefore.bestMove;
          const classification = stockfish.classifyMove(
            positionBefore,
            positionAfter,
            move.from + move.to,
            bestMove
          );

          const moveAnalysis: MoveAnalysis = {
            move: move.from + move.to,
            san: move.san,
            evaluation: positionAfter,
            classification,
            alternativeMoves: [{
              move: bestMove,
              evaluation: positionBefore
            }]
          };

          moveAnalyses.push(moveAnalysis);
        }
      }

      // Calculate player statistics
      const whiteStats = calculatePlayerStats(moveAnalyses.filter((_, i) => i % 2 === 0));
      const blackStats = calculatePlayerStats(moveAnalyses.filter((_, i) => i % 2 === 1));

      // Calculate accuracies
      const whiteEvaluations = evaluations.filter((_, i) => i % 2 === 0);
      const blackEvaluations = evaluations.filter((_, i) => i % 2 === 1);
      
      whiteStats.accuracy = stockfish.calculateAccuracy(whiteEvaluations);
      blackStats.accuracy = stockfish.calculateAccuracy(blackEvaluations);

      // Create complete game analysis
      const analysis: GameAnalysis = {
        moves: moveAnalyses,
        whiteStats,
        blackStats,
        openingAnalysis: {
          name: chessGame.gameState.gameInfo.opening || 'Unknown',
          eco: chessGame.gameState.gameInfo.eco || '',
          accuracy: Math.max(whiteStats.accuracy, blackStats.accuracy)
        },
        gamePhases: {
          opening: Math.min(10, moves.length),
          middlegame: Math.min(25, moves.length),
          endgame: Math.max(25, moves.length)
        }
      };

      setGameAnalysis(analysis);
      
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Analysis failed';
      setAnalysisError(errorMessage);
      console.error('Game analysis error:', error);
    } finally {
      setIsAnalyzingGame(false);
    }
  }, [chessGame.gameState, stockfish]);

  const calculatePlayerStats = (playerMoves: MoveAnalysis[]): PlayerStatistics => {
    const stats: PlayerStatistics = {
      accuracy: 0,
      brilliant: 0,
      great: 0,
      best: 0,
      good: 0,
      inaccuracy: 0,
      mistake: 0,
      blunder: 0,
      miss: 0
    };

    playerMoves.forEach(move => {
      stats[move.classification]++;
    });

    return stats;
  };

  const getMoveAnalysis = useCallback((moveIndex: number): MoveAnalysis | null => {
    if (!gameAnalysis || moveIndex < 0 || moveIndex >= gameAnalysis.moves.length) {
      return null;
    }
    return gameAnalysis.moves[moveIndex];
  }, [gameAnalysis]);

  const getCurrentMoveAnalysis = useCallback((): MoveAnalysis | null => {
    return getMoveAnalysis(chessGame.currentMoveIndex);
  }, [getMoveAnalysis, chessGame.currentMoveIndex]);

  const getPositionEvaluation = useCallback((moveIndex: number): EngineEvaluation | null => {
    const moveAnalysis = getMoveAnalysis(moveIndex);
    return moveAnalysis?.evaluation || null;
  }, [getMoveAnalysis]);

  const stopAnalysis = useCallback(() => {
    stockfish.stopAnalysis();
    setIsAnalyzingGame(false);
  }, [stockfish]);

  // Auto-analyze when game is loaded and engine is ready
  useEffect(() => {
    if (chessGame.gameState && stockfish.isReady && !gameAnalysis && !isAnalyzingGame) {
      // Small delay to ensure UI is ready
      setTimeout(() => {
        analyzeCompleteGame();
      }, 500);
    }
  }, [chessGame.gameState, stockfish.isReady, gameAnalysis, isAnalyzingGame, analyzeCompleteGame]);

  return {
    // Chess game state
    ...chessGame,
    
    // Stockfish state
    engineReady: stockfish.isReady,
    engineInitializing: stockfish.isInitializing,
    engineError: stockfish.error,
    
    // Analysis state
    gameAnalysis,
    isAnalyzingGame,
    analysisError,
    analysisProgress: stockfish.analysisProgress,
    
    // Current position analysis
    currentPositionEvaluation: stockfish.currentEvaluation,
    isAnalyzingPosition: stockfish.isAnalyzing,
    
    // Analysis actions
    analyzeCompleteGame,
    analyzeCurrentPosition,
    stopAnalysis,
    
    // Analysis data getters
    getMoveAnalysis,
    getCurrentMoveAnalysis,
    getPositionEvaluation,
    
    // Computed values
    whiteAccuracy: gameAnalysis?.whiteStats.accuracy || 0,
    blackAccuracy: gameAnalysis?.blackStats.accuracy || 0,
    currentMoveAnalysis: getCurrentMoveAnalysis(),
  };
} 