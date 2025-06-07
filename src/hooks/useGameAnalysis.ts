'use client';

import { useState, useCallback, useEffect } from 'react';
import { useChessGame } from './useChessGame';
import { useStockfish } from './useStockfish';
import { GameAnalysis, MoveAnalysis, PlayerStatistics, EngineEvaluation } from '@/types/analysis';
import { ChessGameManager } from '@/utils/chess';

interface AnalysisOptions {
  depth?: number;
}

export function useGameAnalysis() {
  const chessGame = useChessGame();
  const stockfish = useStockfish();
  
  const [gameAnalysis, setGameAnalysis] = useState<GameAnalysis | null>(null);
  const [isAnalyzingGame, setIsAnalyzingGame] = useState(false);
  const [analysisError, setAnalysisError] = useState<string | null>(null);

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

  const analyzeCurrentPosition = useCallback(async () => {
    if (!chessGame.currentPosition || !stockfish.isReady) return null;
    
    return await stockfish.analyzePosition(chessGame.currentPosition);
  }, [chessGame.currentPosition, stockfish]);

  const analyzeCompleteGame = useCallback(async (options?: AnalysisOptions) => {
    if (options?.depth) {
      stockfish.updateConfig({ depth: options.depth });
    }

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
          const isWhiteMove = i % 2 === 0;
          const playerRating = isWhiteMove 
            ? (chessGame.gameState.gameInfo.whiteRating || 1500)
            : (chessGame.gameState.gameInfo.blackRating || 1500);
          
          const classification = stockfish.classifyMove(
            positionBefore,
            positionAfter,
            move.from + move.to,
            bestMove,
            playerRating
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

      // Calculate accuracies using the evaluations from the actual move analysis
      // Since moveAnalyses already has the correct evaluations for each move,
      // we can extract them from there to ensure proper white/black separation
      const whiteEvaluations: EngineEvaluation[] = [];
      const blackEvaluations: EngineEvaluation[] = [];
      
      moveAnalyses.forEach((moveAnalysis, index) => {
        const isWhiteMove = index % 2 === 0;
        if (isWhiteMove) {
          whiteEvaluations.push(moveAnalysis.evaluation);
        } else {
          blackEvaluations.push(moveAnalysis.evaluation);
        }
      });
      
      whiteStats.accuracy = stockfish.calculateAccuracy(whiteEvaluations);
      blackStats.accuracy = stockfish.calculateAccuracy(blackEvaluations);

      // Debug logging
      console.log('White move count:', whiteEvaluations.length);
      console.log('Black move count:', blackEvaluations.length);
      console.log('Total moves:', moveAnalyses.length);
      console.log('White stats:', whiteStats);
      console.log('Black stats:', blackStats);

      // Detect critical moments and analyze game phases
      const criticalMoments = stockfish.engine?.detectCriticalMoments(evaluations) || [];
      const phaseAnalysis = stockfish.engine?.analyzeGamePhases(moves, evaluations) || {
        opening: Math.min(10, moves.length),
        middlegame: Math.min(25, moves.length), 
        endgame: moves.length,
        openingAccuracy: whiteStats.accuracy,
        middlegameAccuracy: whiteStats.accuracy,
        endgameAccuracy: whiteStats.accuracy
      };

      // Calculate tactical statistics
      whiteStats.tacticalMoves = 0;
      whiteStats.forcingMoves = 0;
      blackStats.tacticalMoves = 0;
      blackStats.forcingMoves = 0;

      // Analyze each move for tactical patterns
      for (let i = 0; i < moveAnalyses.length; i++) {
        const moveAnalysis = moveAnalyses[i];
        const positionBefore = evaluations[i];
        const positionAfter = evaluations[i + 1];
        
        if (positionBefore && positionAfter && stockfish.engine) {
          const tacticalAnalysis = stockfish.engine.analyzeTacticalPatterns(
            positionBefore,
            positionAfter,
            moveAnalysis.move
          );
          
          moveAnalysis.tacticalAnalysis = tacticalAnalysis;
          
          const isWhiteMove = i % 2 === 0;
          if (isWhiteMove) {
            if (tacticalAnalysis.isTactical) whiteStats.tacticalMoves!++;
            if (tacticalAnalysis.isForcing) whiteStats.forcingMoves!++;
          } else {
            if (tacticalAnalysis.isTactical) blackStats.tacticalMoves!++;
            if (tacticalAnalysis.isForcing) blackStats.forcingMoves!++;
          }
        }
      }

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
          opening: phaseAnalysis.opening,
          middlegame: phaseAnalysis.middlegame,
          endgame: phaseAnalysis.endgame
        },
        criticalMoments,
        evaluationHistory: evaluations,
        phaseAnalysis: {
          openingAccuracy: phaseAnalysis.openingAccuracy,
          middlegameAccuracy: phaseAnalysis.middlegameAccuracy,
          endgameAccuracy: phaseAnalysis.endgameAccuracy
        },
        gameResult: {
          result: chessGame.gameState.gameInfo.result as any || '*',
          termination: chessGame.gameState.gameInfo.termination || 'Unknown',
          winningAdvantage: Math.max(...evaluations.map(e => Math.abs(e.score)))
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
        miss: 0,
      }
    };

    playerMoves.forEach(move => {
      stats[move.classification]++;
      stats.moveCounts[move.classification]++;
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

  const resetGame = useCallback(() => {
    setGameAnalysis(null);
    setAnalysisError(null);
    setIsAnalyzingGame(false);
    localStorage.removeItem('chess-analysis');
    localStorage.removeItem('chess-game-state');
    chessGame.resetGame();
  }, [chessGame]);

  const loadGameAndAnalyze = useCallback(async (pgn: string, options?: AnalysisOptions) => {
    // Clear existing analysis when loading new game
    setGameAnalysis(null);
    localStorage.removeItem('chess-analysis');
    localStorage.removeItem('chess-game-state');
    
    chessGame.loadGame(pgn);
    // The useEffect will trigger analysis, but we need to set config first
    if (options?.depth) {
      stockfish.updateConfig({ depth: options.depth });
    }
  }, [chessGame, stockfish]);

  // Auto-analyze when game is loaded and engine is ready
  useEffect(() => {
    if (chessGame.gameState && stockfish.isReady && !gameAnalysis && !isAnalyzingGame) {
      console.log('Starting automatic game analysis...');
      setTimeout(() => {
        analyzeCompleteGame();
      }, 500); // Increased delay to ensure engine is fully ready
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
    resetGame,
    
    // Analysis data getters
    getMoveAnalysis,
    getCurrentMoveAnalysis,
    getPositionEvaluation,
    
    // Computed values
    whiteAccuracy: gameAnalysis?.whiteStats.accuracy || 0,
    blackAccuracy: gameAnalysis?.blackStats.accuracy || 0,
    currentMoveAnalysis: getCurrentMoveAnalysis(),
    
    // New function
    loadGame: loadGameAndAnalyze,
  };
} 