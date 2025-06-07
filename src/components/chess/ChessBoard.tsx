'use client';

import React, { useEffect, useRef } from 'react';
import { BoardOrientation } from '@/types/chess';
import { cn } from '@/lib/utils';

// Import ChessBoard.js
declare global {
  interface Window {
    Chessboard: any;
  }
}

interface ChessBoardProps {
  position: string; // FEN string
  orientation?: BoardOrientation;
  highlightedSquares?: string[];
  onSquareClick?: (square: string) => void;
  className?: string;
}

export default function ChessBoard({
  position,
  orientation = 'white',
  highlightedSquares = [],
  onSquareClick,
  className = ''
}: ChessBoardProps) {
  const boardRef = useRef<HTMLDivElement>(null);
  const chessboardRef = useRef<any>(null);

  useEffect(() => {
    // Dynamically load ChessBoard.js
    const loadChessBoard = async () => {
      // Load CSS
      if (!document.getElementById('chessboard-css')) {
        const css = document.createElement('link');
        css.id = 'chessboard-css';
        css.rel = 'stylesheet';
        css.href = 'https://unpkg.com/@chrisoakman/chessboardjs@1.0.0/dist/chessboard-1.0.0.min.css';
        document.head.appendChild(css);
      }

      // Load jQuery if not already loaded
      if (!window.jQuery) {
        const jquery = document.createElement('script');
        jquery.src = 'https://code.jquery.com/jquery-3.5.1.min.js';
        document.head.appendChild(jquery);
        await new Promise(resolve => jquery.onload = resolve);
      }

      // Load ChessBoard.js if not already loaded
      if (!window.Chessboard) {
        const chessboard = document.createElement('script');
        chessboard.src = 'https://unpkg.com/@chrisoakman/chessboardjs@1.0.0/dist/chessboard-1.0.0.min.js';
        document.head.appendChild(chessboard);
        await new Promise(resolve => chessboard.onload = resolve);
      }

      // Initialize the board
      if (boardRef.current && window.Chessboard && !chessboardRef.current) {
        chessboardRef.current = window.Chessboard(boardRef.current, {
          position: position,
          orientation: orientation,
          showNotation: true,
          pieceTheme: 'https://chessboardjs.com/img/chesspieces/wikipedia/{piece}.png',
          sparePieces: false,
          draggable: false,
          dropOffBoard: 'snapback',
          moveSpeed: 'fast',
          snapbackSpeed: 500,
          snapSpeed: 100,
          onSquareClick: onSquareClick
        });
      }
    };

    loadChessBoard();

    return () => {
      // Cleanup when component unmounts
      if (chessboardRef.current && chessboardRef.current.destroy) {
        chessboardRef.current.destroy();
        chessboardRef.current = null;
      }
    };
  }, []);

  // Update position when it changes
  useEffect(() => {
    if (chessboardRef.current && chessboardRef.current.position) {
      // Use position() method to update the board position immediately
      chessboardRef.current.position(position, false); // false = no animation
    }
  }, [position]);

  // Update orientation when it changes
  useEffect(() => {
    if (chessboardRef.current && chessboardRef.current.orientation) {
      chessboardRef.current.orientation(orientation);
    }
  }, [orientation]);

  return (
    <div className={cn("aspect-square", className)}>
      <div 
        ref={boardRef}
        style={{ width: '100%', height: '100%' }}
      />
    </div>
  );
} 