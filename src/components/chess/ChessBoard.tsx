'use client';

import React from 'react';
import { Chess } from 'chess.js';
import { BoardOrientation } from '@/types/chess';
import { parseSquareColor } from '@/utils/chess';
import ChessSquare from './ChessSquare';
import ChessPiece from './ChessPiece';
import { cn } from '@/lib/utils';

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
  const chess = new Chess(position);
  const board = chess.board();
  
  const files = orientation === 'white' ? ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'] : ['h', 'g', 'f', 'e', 'd', 'c', 'b', 'a'];
  const ranks = orientation === 'white' ? ['8', '7', '6', '5', '4', '3', '2', '1'] : ['1', '2', '3', '4', '5', '6', '7', '8'];

  const handleSquareClick = (square: string) => {
    onSquareClick?.(square);
  };

  return (
    <div className={cn("aspect-square w-full", className)}>
      <div className="grid grid-cols-8 grid-rows-8 w-full h-full">
        {ranks.map((rank, rankIndex) =>
          files.map((file, fileIndex) => {
            const square = file + rank;
            const squareColor = parseSquareColor(square);
            const piece = board[rankIndex]?.[fileIndex];

            // Determine if the coordinate should be shown
            const showRank = fileIndex === 0;
            const showFile = rankIndex === 7;

            return (
              <ChessSquare
                key={square}
                square={square}
                color={squareColor}
                isHighlighted={highlightedSquares.includes(square)}
                onClick={() => handleSquareClick(square)}
                className="relative"
              >
                {/* Add rank and file coordinates inside the squares */}
                {showRank && (
                  <span
                    className={cn(
                      "absolute top-0 left-1 text-xs font-bold pointer-events-none",
                      squareColor === 'light' ? 'text-board-dark' : 'text-board-light'
                    )}
                  >
                    {rank}
                  </span>
                )}
                {showFile && (
                  <span
                    className={cn(
                      "absolute bottom-0 right-1 text-xs font-bold pointer-events-none",
                       squareColor === 'light' ? 'text-board-dark' : 'text-board-light'
                    )}
                  >
                    {file}
                  </span>
                )}
                
                {piece && (
                  <ChessPiece
                    type={piece.type}
                    color={piece.color}
                  />
                )}
              </ChessSquare>
            );
          })
        )}
      </div>
    </div>
  );
} 