import React from 'react';
import { cn } from '@/lib/utils';

interface ChessSquareProps {
  square: string;
  color: 'light' | 'dark';
  isHighlighted?: boolean;
  onClick?: () => void;
  children?: React.ReactNode;
  className?: string;
}

export default function ChessSquare({
  square,
  color,
  isHighlighted = false,
  onClick,
  children,
  className
}: ChessSquareProps) {
  return (
    <div
      className={cn(
        'w-full h-full flex items-center justify-center relative cursor-pointer',
        color === 'light' ? 'chess-square-light' : 'chess-square-dark',
        isHighlighted && 'chess-square-highlighted',
        className
      )}
      onClick={onClick}
      role="button"
      tabIndex={0}
      aria-label={`Square ${square}`}
    >
      {children}
      {isHighlighted && (
        <div className="absolute inset-0 bg-blue-400/20 pointer-events-none rounded-sm" />
      )}
    </div>
  );
} 