'use client';

import React from 'react';
import Button from '@/components/ui/Button';

interface GameControlsProps {
  canGoBackward: boolean;
  canGoForward: boolean;
  isAtStart: boolean;
  isAtEnd: boolean;
  onGoToStart: () => void;
  onGoBackward: () => void;
  onGoForward: () => void;
  onGoToEnd: () => void;
  currentMoveIndex: number;
  totalMoves: number;
  className?: string;
}

export default function GameControls({
  canGoBackward,
  canGoForward,
  isAtStart,
  isAtEnd,
  onGoToStart,
  onGoBackward,
  onGoForward,
  onGoToEnd,
  currentMoveIndex,
  totalMoves,
  className = ''
}: GameControlsProps) {
  return (
    <div className={`bg-white border border-gray-200 rounded-lg p-4 ${className}`}>
      <div className="flex flex-col space-y-4">
        {/* Move indicator */}
        <div className="text-center">
          <div className="text-sm text-gray-600">
            Move {currentMoveIndex + 1} of {totalMoves}
            {currentMoveIndex === -1 && ' (Starting position)'}
          </div>
        </div>

        {/* Navigation buttons */}
        <div className="flex justify-center space-x-2">
          <Button
            variant="outline"
            size="sm"
            onClick={onGoToStart}
            disabled={isAtStart}
            className="w-12"
            title="Go to start (Home)"
          >
            ⏮
          </Button>
          
          <Button
            variant="outline"
            size="sm"
            onClick={onGoBackward}
            disabled={!canGoBackward}
            className="w-12"
            title="Previous move (←)"
          >
            ◀
          </Button>
          
          <Button
            variant="outline"
            size="sm"
            onClick={onGoForward}
            disabled={!canGoForward}
            className="w-12"
            title="Next move (→)"
          >
            ▶
          </Button>
          
          <Button
            variant="outline"
            size="sm"
            onClick={onGoToEnd}
            disabled={isAtEnd}
            className="w-12"
            title="Go to end (End)"
          >
            ⏭
          </Button>
        </div>

        {/* Progress bar */}
        <div className="w-full bg-gray-200 rounded-full h-2">
          <div
            className="bg-green-600 h-2 rounded-full transition-all duration-200"
            style={{
              width: totalMoves > 0 ? `${((currentMoveIndex + 1) / totalMoves) * 100}%` : '0%'
            }}
          />
        </div>

        {/* Keyboard shortcuts info */}
        <div className="text-xs text-gray-500 text-center">
          ← → Navigate • Home/End Jump to start/end
        </div>
      </div>
    </div>
  );
} 