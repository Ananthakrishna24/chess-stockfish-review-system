'use client';

import React, { useEffect } from 'react';
import Button from '@/components/ui/Button';
import { Progress } from '@/components/ui/Progress';
import { cn } from '@/lib/utils';
import { 
  SkipBack, 
  ChevronLeft, 
  ChevronRight, 
  SkipForward,
  Play,
  Pause,
  X,
  Settings
} from 'lucide-react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { 
  exitReviewMode, 
  toggleAutoPlay, 
  setCurrentReviewMove 
} from '@/store/reviewModeSlice';

interface ReviewModeControlsProps {
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

export default function ReviewModeControls({
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
}: ReviewModeControlsProps) {
  const dispatch = useAppDispatch();
  const { autoPlayMode, autoPlayInterval } = useAppSelector(state => state.reviewMode);

  const handleExitReview = () => {
    dispatch(exitReviewMode());
  };

  const handleToggleAutoPlay = () => {
    dispatch(toggleAutoPlay());
  };

  // Auto-play functionality
  useEffect(() => {
    if (autoPlayMode && canGoForward) {
      const timer = setTimeout(() => {
        onGoForward();
        dispatch(setCurrentReviewMove(currentMoveIndex + 1));
      }, autoPlayInterval);

      return () => clearTimeout(timer);
    } else if (autoPlayMode && isAtEnd) {
      // Auto-stop when reaching the end
      dispatch(toggleAutoPlay());
    }
  }, [autoPlayMode, canGoForward, isAtEnd, currentMoveIndex, autoPlayInterval, onGoForward, dispatch]);

  // Keyboard event handling
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      switch (event.key) {
        case 'ArrowLeft':
          event.preventDefault();
          if (canGoBackward) onGoBackward();
          break;
        case 'ArrowRight':
          event.preventDefault();
          if (canGoForward) onGoForward();
          break;
        case 'Home':
          event.preventDefault();
          onGoToStart();
          break;
        case 'End':
          event.preventDefault();
          onGoToEnd();
          break;
        case ' ':
          event.preventDefault();
          handleToggleAutoPlay();
          break;
        case 'Escape':
          event.preventDefault();
          handleExitReview();
          break;
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
     }, [canGoBackward, canGoForward, onGoBackward, onGoForward, onGoToStart, onGoToEnd, handleToggleAutoPlay, handleExitReview]);

  return (
    <div className={cn("bg-card border-t border-border p-4", className)}>
      {/* Review Mode Header */}
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
          <span className="text-sm font-medium text-green-600">Review Mode</span>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={handleToggleAutoPlay}
            className="flex items-center gap-1"
          >
            {autoPlayMode ? <Pause className="h-3 w-3" /> : <Play className="h-3 w-3" />}
            {autoPlayMode ? 'Pause' : 'Auto Play'}
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={handleExitReview}
            className="flex items-center gap-1"
          >
            <X className="h-3 w-3" />
            Exit
          </Button>
        </div>
      </div>

      {/* Move indicator */}
      <div className="text-center mb-3">
        <div className="text-sm text-muted-foreground">
          {currentMoveIndex === -1 ? (
            'Starting position'
          ) : (
            `Move ${currentMoveIndex + 1} of ${totalMoves}`
          )}
        </div>
      </div>

      {/* Navigation buttons */}
      <div className="flex justify-center gap-2 mb-3">
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
        className="h-2 mb-2"
      />

      {/* Auto-play status */}
      {autoPlayMode && (
        <div className="text-center">
          <div className="text-xs text-muted-foreground">
            Auto-playing at {autoPlayInterval / 1000}s intervals
          </div>
        </div>
      )}

      {/* Keyboard shortcuts info */}
      <div className="text-xs text-muted-foreground text-center mt-2">
        ← → Navigate • Space Auto-play • Esc Exit review
      </div>
    </div>
  );
} 