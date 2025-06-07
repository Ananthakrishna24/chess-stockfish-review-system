import React from 'react';
import { cn } from '@/lib/utils';

interface EvaluationBarProps {
  evaluation: number; // Centipawns from white's perspective
  className?: string;
}

export function EvaluationBar({ evaluation, className }: EvaluationBarProps) {
  // Clamp evaluation to a practical range for visualization, e.g., -1000 to +1000 centipawns
  const VIZ_EVAL_CAP = 1000;
  const clampedEval = Math.max(-VIZ_EVAL_CAP, Math.min(VIZ_EVAL_CAP, evaluation));

  // Linear scaling: map evaluation from [-VIZ_EVAL_CAP, VIZ_EVAL_CAP] to [0, 100]
  const whitePercentage = 50 + (clampedEval / VIZ_EVAL_CAP) * 50;

  // Handle mate scores explicitly
  const isMate = Math.abs(evaluation) > 9000;
  let displayPercentage = whitePercentage;

  if (isMate) {
    displayPercentage = evaluation > 0 ? 100 : 0;
  }

  const whiteHeight = displayPercentage;
  const blackHeight = 100 - displayPercentage;

  const formatEvaluation = (evalValue: number): string => {
    if (Math.abs(evalValue) > 9000) {
      const mateIn = Math.ceil((10000 - Math.abs(evalValue)) / 2);
      return evalValue > 0 ? `M${mateIn}` : `-M${mateIn}`;
    }
    const displayScore = (evalValue / 100).toFixed(1);
    // Add a plus sign for positive evaluations
    return evalValue > 0 ? `+${displayScore}` : displayScore;
  };
  
  const evalText = formatEvaluation(evaluation);

  return (
    <div className={cn("relative h-full w-8 flex flex-col justify-center items-center", className)}>
      <div 
        className="absolute top-0 w-full bg-gray-100 dark:bg-gray-200"
        style={{ height: `${blackHeight}%`, transition: 'height 0.3s ease-in-out' }}
      />
      <div 
        className="absolute bottom-0 w-full bg-gray-800 dark:bg-gray-700"
        style={{ height: `${whiteHeight}%`, transition: 'height 0.3s ease-in-out' }}
      />
      <span className="relative text-white dark:text-gray-800 font-bold text-xs mix-blend-difference">
        {evalText}
      </span>
    </div>
  );
} 