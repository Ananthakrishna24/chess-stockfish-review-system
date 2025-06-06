import React from 'react';

interface ChessSquareProps {
  square: string;
  color: 'light' | 'dark';
  isHighlighted?: boolean;
  onClick?: () => void;
  children?: React.ReactNode;
}

export default function ChessSquare({
  square,
  color,
  isHighlighted = false,
  onClick,
  children
}: ChessSquareProps) {
  const baseClasses = 'w-12 h-12 flex items-center justify-center relative cursor-pointer transition-colors';
  
  const colorClasses = {
    light: 'bg-amber-50 hover:bg-amber-100',
    dark: 'bg-green-600 hover:bg-green-700'
  };

  const highlightClasses = isHighlighted 
    ? 'ring-2 ring-yellow-400 ring-inset' 
    : '';

  return (
    <div
      className={`${baseClasses} ${colorClasses[color]} ${highlightClasses}`}
      onClick={onClick}
      role="button"
      tabIndex={0}
      aria-label={`Square ${square}`}
    >
      {children}
      {isHighlighted && (
        <div className="absolute inset-0 bg-yellow-400 bg-opacity-30 pointer-events-none" />
      )}
    </div>
  );
} 