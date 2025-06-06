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
  const baseClasses = 'w-14 h-14 flex items-center justify-center relative cursor-pointer transition-all duration-200 hover:brightness-110';
  
  const getSquareStyle = () => {
    if (color === 'light') {
      return {
        backgroundColor: 'var(--chess-light-square)',
      };
    } else {
      return {
        backgroundColor: 'var(--chess-dark-square)',
      };
    }
  };

  const highlightClasses = isHighlighted 
    ? 'ring-2 ring-inset' 
    : '';

  const highlightStyle = isHighlighted 
    ? { ringColor: 'var(--chess-accent)' }
    : {};

  return (
    <div
      className={`${baseClasses} ${highlightClasses}`}
      style={{
        ...getSquareStyle(),
        ...highlightStyle
      }}
      onClick={onClick}
      role="button"
      tabIndex={0}
      aria-label={`Square ${square}`}
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          onClick?.();
        }
      }}
    >
      {children}
      {isHighlighted && (
        <div 
          className="absolute inset-0 pointer-events-none rounded-sm"
          style={{
            backgroundColor: 'var(--chess-accent)',
            opacity: 0.2
          }}
        />
      )}
    </div>
  );
} 