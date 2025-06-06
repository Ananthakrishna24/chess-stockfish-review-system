'use client';

import { useState, useCallback, useEffect } from 'react';
import { GameState } from '@/types/chess';
import { ChessGameManager } from '@/utils/chess';

export function useChessGame() {
  const [gameState, setGameState] = useState<GameState | null>(null);
  const [currentMoveIndex, setCurrentMoveIndex] = useState(-1);
  const [gameManager, setGameManager] = useState<ChessGameManager | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadGame = useCallback(async (pgn: string) => {
    setIsLoading(true);
    setError(null);
    
    try {
      const manager = new ChessGameManager();
      const state = manager.loadPGN(pgn);
      
      setGameManager(manager);
      setGameState(state);
      setCurrentMoveIndex(-1); // Start at initial position
      
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load game');
      console.error('Failed to load PGN:', err);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const goToMove = useCallback((moveIndex: number) => {
    if (!gameState) return;
    
    const maxIndex = gameState.moves.length - 1;
    const newIndex = Math.max(-1, Math.min(moveIndex, maxIndex));
    setCurrentMoveIndex(newIndex);
  }, [gameState]);

  const goToStart = useCallback(() => {
    setCurrentMoveIndex(-1);
  }, []);

  const goToEnd = useCallback(() => {
    if (gameState) {
      setCurrentMoveIndex(gameState.moves.length - 1);
    }
  }, [gameState]);

  const goForward = useCallback(() => {
    if (gameState && currentMoveIndex < gameState.moves.length - 1) {
      setCurrentMoveIndex(currentMoveIndex + 1);
    }
  }, [gameState, currentMoveIndex]);

  const goBackward = useCallback(() => {
    if (currentMoveIndex > -1) {
      setCurrentMoveIndex(currentMoveIndex - 1);
    }
  }, [currentMoveIndex]);

  const getCurrentPosition = useCallback(() => {
    if (!gameManager) {
      return 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1';
    }
    return gameManager.getPosition(currentMoveIndex);
  }, [gameManager, currentMoveIndex]);

  const getCurrentMove = useCallback(() => {
    if (!gameState || currentMoveIndex < 0 || currentMoveIndex >= gameState.moves.length) {
      return null;
    }
    return gameState.moves[currentMoveIndex];
  }, [gameState, currentMoveIndex]);

  const getMovesUpToCurrent = useCallback(() => {
    if (!gameState || currentMoveIndex < 0) {
      return [];
    }
    return gameState.moves.slice(0, currentMoveIndex + 1);
  }, [gameState, currentMoveIndex]);

  // Keyboard navigation
  useEffect(() => {
    const handleKeyPress = (event: KeyboardEvent) => {
      if (!gameState) return;
      
      switch (event.key) {
        case 'ArrowLeft':
          event.preventDefault();
          goBackward();
          break;
        case 'ArrowRight':
          event.preventDefault();
          goForward();
          break;
        case 'Home':
          event.preventDefault();
          goToStart();
          break;
        case 'End':
          event.preventDefault();
          goToEnd();
          break;
      }
    };

    window.addEventListener('keydown', handleKeyPress);
    return () => window.removeEventListener('keydown', handleKeyPress);
  }, [gameState, goBackward, goForward, goToStart, goToEnd]);

  const resetGame = useCallback(() => {
    setGameState(null);
    setGameManager(null);
    setCurrentMoveIndex(-1);
    setError(null);
  }, []);

  return {
    // State
    gameState,
    currentMoveIndex,
    isLoading,
    error,
    
    // Computed values
    currentPosition: getCurrentPosition(),
    currentMove: getCurrentMove(),
    movesUpToCurrent: getMovesUpToCurrent(),
    
    // Actions
    loadGame,
    goToMove,
    goToStart,
    goToEnd,
    goForward,
    goBackward,
    resetGame,
    
    // Status
    canGoForward: gameState ? currentMoveIndex < gameState.moves.length - 1 : false,
    canGoBackward: currentMoveIndex > -1,
    isAtStart: currentMoveIndex === -1,
    isAtEnd: gameState ? currentMoveIndex === gameState.moves.length - 1 : false
  };
} 