import React from 'react';
import { cn } from '@/lib/utils';

interface EvaluationBarProps {
  evaluation: number; // Centipawns from white's perspective
  orientation?: 'white' | 'black'; // Board orientation
  className?: string;
}

export function EvaluationBar({ evaluation, orientation = 'white', className }: EvaluationBarProps) {
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

  // Determine layout based on orientation
  const isFlipped = orientation === 'black';
  
  // When flipped, swap the visual layout to match board orientation
  const topPlayer = isFlipped ? 'white' : 'black';
  const bottomPlayer = isFlipped ? 'black' : 'white';
  
  const topEvaluation = isFlipped ? evaluation : -evaluation;
  const topPercentage = isFlipped ? displayWhite : displayBlack;
  const topBgColor = isFlipped ? 'bg-gray-100' : 'bg-gray-800';
  const topTextColor = isFlipped ? 'text-black' : 'text-white';
  const topBarColor = isFlipped ? 'bg-gray-100' : 'bg-gray-800';
  
  const bottomEvaluation = isFlipped ? -evaluation : evaluation;
  const bottomPercentage = isFlipped ? displayBlack : displayWhite;
  const bottomBgColor = isFlipped ? 'bg-gray-800' : 'bg-gray-100';
  const bottomTextColor = isFlipped ? 'text-white' : 'text-black';
  const bottomBarColor = isFlipped ? 'bg-gray-800' : 'bg-gray-100';

  return (
    <div className={cn("flex flex-col items-center", className)}>
      {/* Top player evaluation number */}
      <div className={cn("text-xs font-bold px-1 py-0.5 rounded-t mb-1 min-w-[32px] text-center", topBgColor, topTextColor)}>
        {formatEvaluation(topEvaluation)}
      </div>
      
      {/* Evaluation bar */}
      <div className={cn("w-8 flex-1 bg-background border border-border rounded-sm overflow-hidden relative")}>
        {/* Top section */}
        <div 
          className={cn("transition-all duration-300 ease-in-out", topBarColor)}
          style={{ height: `${topPercentage}%` }}
        />
        {/* Bottom section */}
        <div 
          className={cn("transition-all duration-300 ease-in-out", bottomBarColor)}
          style={{ height: `${bottomPercentage}%` }}
        />
      </div>
      
      {/* Bottom player evaluation number */}
      <div className={cn("text-xs font-bold px-1 py-0.5 rounded-b mt-1 min-w-[32px] text-center", bottomBgColor, bottomTextColor)}>
        {formatEvaluation(bottomEvaluation)}
      </div>
    </div>
  );
} 