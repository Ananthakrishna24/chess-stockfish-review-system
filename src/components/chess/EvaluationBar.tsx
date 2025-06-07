import React from 'react';
import { cn } from '@/lib/utils';

interface EvaluationBarProps {
  evaluation: number; // Centipawns from white's perspective
  className?: string;
}

export function EvaluationBar({ evaluation, className }: EvaluationBarProps) {
  // Convert centipawn evaluation to percentage for visual display
  // Sigmoid function: 1 / (1 + 10^(-evaluation/400))
  const getEvaluationPercentage = (evaluation: number): number => {
    const normalized = 1 / (1 + Math.pow(10, -evaluation / 400));
    return Math.max(0, Math.min(100, normalized * 100));
  };

  const whitePercentage = getEvaluationPercentage(evaluation);
  const blackPercentage = 100 - whitePercentage;

  // Handle mate scores
  const isMate = Math.abs(evaluation) > 9000;
  let displayWhite = whitePercentage;
  let displayBlack = blackPercentage;

  if (isMate) {
    if (evaluation > 0) {
      displayWhite = 100;
      displayBlack = 0;
    } else {
      displayWhite = 0;
      displayBlack = 100;
    }
  }

  // Format evaluation for display
  const formatEvaluation = (evalValue: number): string => {
    if (Math.abs(evalValue) > 9000) {
      const mateIn = Math.ceil((10000 - Math.abs(evalValue)) / 2);
      return evalValue > 0 ? `M${mateIn}` : `M-${mateIn}`;
    }
    return (evalValue / 100).toFixed(1);
  };

  return (
    <div className={cn("flex flex-col items-center", className)}>
      {/* Black evaluation number (top) - from black's perspective */}
      <div className="text-xs font-bold text-white bg-gray-800 px-1 py-0.5 rounded-t mb-1 min-w-[32px] text-center">
        {formatEvaluation(-evaluation)}
      </div>
      
      {/* Evaluation bar */}
      <div className={cn("w-8 flex-1 bg-background border border-border rounded-sm overflow-hidden relative")}>
        {/* Black evaluation (top) */}
        <div 
          className="bg-gray-800 transition-all duration-300 ease-in-out"
          style={{ height: `${displayBlack}%` }}
        />
        {/* White evaluation (bottom) */}
        <div 
          className="bg-gray-100 transition-all duration-300 ease-in-out"
          style={{ height: `${displayWhite}%` }}
        />
      </div>
      
      {/* White evaluation number (bottom) - from white's perspective */}
      <div className="text-xs font-bold text-black bg-gray-100 px-1 py-0.5 rounded-b mt-1 min-w-[32px] text-center">
        {formatEvaluation(evaluation)}
      </div>
    </div>
  );
} 