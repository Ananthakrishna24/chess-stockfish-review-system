'use client';

import React from 'react';
import Button from '@/components/ui/Button';
import { Progress } from '@/components/ui/Progress';
import { cn } from '@/lib/utils';
import { 
  SkipBack, 
  ChevronLeft, 
  ChevronRight, 
  SkipForward 
} from 'lucide-react';

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
    <div className={cn("space-y-4", className)}>
      {/* Move indicator */}
      <div className="text-center">
        <div className="text-sm text-muted-foreground">
          {currentMoveIndex === -1 ? (
            'Starting position'
          ) : (
            `Move ${currentMoveIndex + 1} of ${totalMoves}`
          )}
        </div>
      </div>

      {/* Navigation buttons */}
      <div className="flex justify-center gap-2">
        <Button
          variant="outline"
          size="icon"
          onClick={onGoToStart}
          disabled={isAtStart}
          title="Go to start (Home)"
        >
          <SkipBack className="h-4 w-4" />
        </Button>
        
        <Button
          variant="outline"
          size="icon"
          onClick={onGoBackward}
          disabled={!canGoBackward}
          title="Previous move (←)"
        >
          <ChevronLeft className="h-4 w-4" />
        </Button>
        
        <Button
          variant="outline"
          size="icon"
          onClick={onGoForward}
          disabled={!canGoForward}
          title="Next move (→)"
        >
          <ChevronRight className="h-4 w-4" />
        </Button>
        
        <Button
          variant="outline"
          size="icon"
          onClick={onGoToEnd}
          disabled={isAtEnd}
          title="Go to end (End)"
        >
          <SkipForward className="h-4 w-4" />
        </Button>
      </div>

      {/* Progress bar */}
      <Progress 
        value={totalMoves > 0 ? ((currentMoveIndex + 1) / totalMoves) * 100 : 0}
        className="h-2"
      />

      {/* Keyboard shortcuts info */}
      <div className="text-xs text-muted-foreground text-center">
        ← → Navigate • Home/End Jump to start/end
      </div>
    </div>
  );
} 