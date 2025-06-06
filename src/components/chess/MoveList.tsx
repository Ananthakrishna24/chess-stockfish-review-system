'use client';

import React from 'react';
import { ChessMove } from '@/types/chess';

interface MoveListProps {
  moves: ChessMove[];
  currentMoveIndex: number;
  onMoveClick: (moveIndex: number) => void;
  className?: string;
}

export default function MoveList({
  moves,
  currentMoveIndex,
  onMoveClick,
  className = ''
}: MoveListProps) {
  if (moves.length === 0) {
    return (
      <div className={`text-center text-gray-500 py-8 ${className}`}>
        No moves to display
      </div>
    );
  }

  // Group moves by pairs (white and black)
  const movePairs: { white: ChessMove; black?: ChessMove; moveNumber: number }[] = [];
  
  for (let i = 0; i < moves.length; i += 2) {
    const whiteMove = moves[i];
    const blackMove = moves[i + 1];
    
    if (whiteMove && whiteMove.color === 'w') {
      movePairs.push({
        white: whiteMove,
        black: blackMove && blackMove.color === 'b' ? blackMove : undefined,
        moveNumber: whiteMove.moveNumber
      });
    }
  }

  return (
    <div className={`bg-white border border-gray-200 rounded-lg ${className}`}>
      <div className="px-4 py-3 border-b border-gray-200">
        <h3 className="text-sm font-semibold text-gray-900">Game Moves</h3>
      </div>
      
      <div className="max-h-96 overflow-y-auto">
        <div className="p-2 space-y-1">
          {movePairs.map((pair, pairIndex) => (
            <div key={pairIndex} className="flex items-center space-x-2 text-sm">
              {/* Move number */}
              <div className="w-8 text-gray-500 font-medium">
                {pair.moveNumber}.
              </div>
              
              {/* White move */}
              <button
                onClick={() => onMoveClick(pairIndex * 2)}
                className={`px-2 py-1 rounded hover:bg-gray-100 transition-colors min-w-16 text-left ${
                  currentMoveIndex === pairIndex * 2
                    ? 'bg-green-100 text-green-800 font-semibold'
                    : 'text-gray-700 hover:text-gray-900'
                }`}
              >
                {pair.white.san}
              </button>
              
              {/* Black move */}
              {pair.black && (
                <button
                  onClick={() => onMoveClick(pairIndex * 2 + 1)}
                  className={`px-2 py-1 rounded hover:bg-gray-100 transition-colors min-w-16 text-left ${
                    currentMoveIndex === pairIndex * 2 + 1
                      ? 'bg-green-100 text-green-800 font-semibold'
                      : 'text-gray-700 hover:text-gray-900'
                  }`}
                >
                  {pair.black.san}
                </button>
              )}
            </div>
          ))}
        </div>
      </div>
      
      {/* Navigation controls */}
      <div className="px-4 py-3 border-t border-gray-200 bg-gray-50">
        <div className="text-xs text-gray-500 text-center">
          Use arrow keys or click moves to navigate
        </div>
      </div>
    </div>
  );
} 