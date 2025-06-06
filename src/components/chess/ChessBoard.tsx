'use client';

import React from 'react';
import { Chess } from 'chess.js';
import { BoardOrientation, ChessPiece as ChessPieceType } from '@/types/chess';
import { parseSquareColor, getPieceUnicode } from '@/utils/chess';
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
        <div className="flex mb-1">
          <div className="w-6 h-6"></div> {/* Corner space */}
          {files.map((file) => (
            <div
              key={file}
              className="w-12 h-6 flex items-center justify-center text-sm font-medium text-gray-600"
            >
              {file}
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
                className="w-6 h-12 flex items-center justify-center text-sm font-medium text-gray-600"
              >
                {rank}
              </div>
            ))}
          </div>

          {/* The actual chess board */}
          <div className="grid grid-cols-8 border-2 border-gray-800 rounded-lg overflow-hidden">
            {ranks.map((rank, rankIndex) =>
              files.map((file, fileIndex) => {
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
                className="w-6 h-12 flex items-center justify-center text-sm font-medium text-gray-600"
              >
                {rank}
              </div>
            ))}
          </div>
        </div>

        {/* Bottom coordinates */}
        <div className="flex mt-1">
          <div className="w-6 h-6"></div> {/* Corner space */}
          {files.map((file) => (
            <div
              key={`b-${file}`}
              className="w-12 h-6 flex items-center justify-center text-sm font-medium text-gray-600"
            >
              {file}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
} 