'use client';

import React from 'react';
import { ChessMove } from '@/types/chess';
import { cn } from '@/lib/utils';

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
  const scrollRef = React.useRef<HTMLDivElement>(null);

  // Auto-scroll to current move
  React.useEffect(() => {
    if (scrollRef.current) {
      const currentElement = scrollRef.current.querySelector(`[data-move-index="${currentMoveIndex}"]`);
      if (currentElement) {
        currentElement.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
      }
    }
  }, [currentMoveIndex]);

  if (moves.length === 0) {
    return (
      <div className={cn("flex items-center justify-center h-full text-muted-foreground", className)}>
        <div className="text-center space-y-2">
          <div className="text-4xl">♟️</div>
          <div>No moves to display</div>
        </div>
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
    <div className={cn("h-full flex flex-col", className)}>
      <div className="p-4 border-b border-border bg-muted/20">
        <h3 className="font-semibold text-sm">Game Moves</h3>
        <div className="text-xs text-muted-foreground mt-1">
          {moves.length} moves • Click to navigate
        </div>
      </div>
      
      <div 
        ref={scrollRef}
        className="flex-1 overflow-y-auto custom-scrollbar"
      >
        <div className="p-4 space-y-2">
          {movePairs.map((pair, pairIndex) => (
            <div key={pairIndex} className="flex items-center gap-3 py-1">
              {/* Move number */}
              <div className="w-8 text-muted-foreground font-medium text-sm flex-shrink-0">
                {pair.moveNumber}.
              </div>
              
              {/* White move */}
              <button
                data-move-index={pairIndex * 2}
                onClick={() => onMoveClick(pairIndex * 2)}
                className={cn(
                  "px-3 py-2 rounded-md hover:bg-accent transition-colors min-w-16 text-left font-mono text-sm flex-1",
                  currentMoveIndex === pairIndex * 2
                    ? 'bg-primary text-primary-foreground font-semibold'
                    : 'text-foreground hover:text-accent-foreground'
                )}
              >
                {pair.white.san}
              </button>
              
              {/* Black move */}
              {pair.black ? (
                <button
                  data-move-index={pairIndex * 2 + 1}
                  onClick={() => onMoveClick(pairIndex * 2 + 1)}
                  className={cn(
                    "px-3 py-2 rounded-md hover:bg-accent transition-colors min-w-16 text-left font-mono text-sm flex-1",
                    currentMoveIndex === pairIndex * 2 + 1
                      ? 'bg-primary text-primary-foreground font-semibold'
                      : 'text-foreground hover:text-accent-foreground'
                  )}
                >
                  {pair.black.san}
                </button>
              ) : (
                <div className="flex-1"></div>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
} 