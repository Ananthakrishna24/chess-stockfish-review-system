import React from 'react';
import { PieceType, PieceColor } from '@/types/chess';
import { getPieceUnicode } from '@/utils/chess';

interface ChessPieceProps {
  type: PieceType;
  color: PieceColor;
  className?: string;
}

export default function ChessPiece({ type, color, className = '' }: ChessPieceProps) {
  const pieceSymbol = getPieceUnicode(type, color);
  
  return (
    <div
      className={`text-[8vmin] leading-none select-none pointer-events-none ${className}`}
      style={{
        textShadow: '1px 1px 3px rgba(0,0,0,0.5)'
      }}
    >
      {pieceSymbol}
    </div>
  );
} 