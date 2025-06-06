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
      className={`text-3xl leading-none select-none pointer-events-none ${className}`}
      style={{
        textShadow: '1px 1px 2px rgba(0, 0, 0, 0.3)',
        filter: 'drop-shadow(0 1px 1px rgba(0, 0, 0, 0.2))'
      }}
    >
      {pieceSymbol}
    </div>
  );
} 