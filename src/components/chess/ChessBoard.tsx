'use client';

import React from 'react';
import { Chess } from 'chess.js';
import { BoardOrientation } from '@/types/chess';
import { parseSquareColor } from '@/utils/chess';
import ChessSquare from './ChessSquare';
import ChessPiece from './ChessPiece';

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
  
  // Generate files and ranks based on orientation
  const files = orientation === 'white' ? ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'] : ['h', 'g', 'f', 'e', 'd', 'c', 'b', 'a'];
  const ranks = orientation === 'white' ? ['8', '7', '6', '5', '4', '3', '2', '1'] : ['1', '2', '3', '4', '5', '6', '7', '8'];

  const handleSquareClick = (square: string) => {
    onSquareClick?.(square);
  };

  return (
    <div className={`select-none ${className}`}>
      {/* Board container with coordinates */}
      <div className="relative inline-block">
        {/* Top coordinates */}
        <div className="flex mb-2">
          <div className="w-8 h-8"></div> {/* Corner space */}
          {files.map((file) => (
            <div
              key={file}
              className="w-14 h-8 flex items-center justify-center text-sm font-semibold"
              style={{ color: 'var(--text-secondary)' }}
            >
              {file.toUpperCase()}
            </div>
          ))}
        </div>

        {/* Board with side coordinates */}
        <div className="flex">
          {/* Left coordinates */}
          <div className="flex flex-col">
            {ranks.map((rank) => (
              <div
                key={rank}
                className="w-8 h-14 flex items-center justify-center text-sm font-semibold"
                style={{ color: 'var(--text-secondary)' }}
              >
                {rank}
              </div>
            ))}
          </div>

          {/* The actual chess board */}
          <div 
            className="grid grid-cols-8 rounded-lg overflow-hidden shadow-lg"
            style={{
              border: '3px solid var(--chess-board-border)',
              boxShadow: 'var(--shadow-lg)'
            }}
          >
            {ranks.map((rank) =>
              files.map((file) => {
                const square = file + rank;
                const squareColor = parseSquareColor(square);
                const piece = board[7 - parseInt(rank) + 1]?.[file.charCodeAt(0) - 97];
                const isHighlighted = highlightedSquares.includes(square);

                return (
                  <ChessSquare
                    key={square}
                    square={square}
                    color={squareColor}
                    isHighlighted={isHighlighted}
                    onClick={() => handleSquareClick(square)}
                  >
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

          {/* Right coordinates */}
          <div className="flex flex-col">
            {ranks.map((rank) => (
              <div
                key={`r-${rank}`}
                className="w-8 h-14 flex items-center justify-center text-sm font-semibold"
                style={{ color: 'var(--text-secondary)' }}
              >
                {rank}
              </div>
            ))}
          </div>
        </div>

        {/* Bottom coordinates */}
        <div className="flex mt-2">
          <div className="w-8 h-8"></div> {/* Corner space */}
          {files.map((file) => (
            <div
              key={`b-${file}`}
              className="w-14 h-8 flex items-center justify-center text-sm font-semibold"
              style={{ color: 'var(--text-secondary)' }}
            >
              {file.toUpperCase()}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
} 