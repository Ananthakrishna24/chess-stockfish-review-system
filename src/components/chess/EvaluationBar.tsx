import React from 'react';
import { cn } from '@/lib/utils';
import { DisplayEvaluation } from '@/types/analysis';

interface EvaluationBarProps {
  evaluation: number;
  displayEvaluation?: DisplayEvaluation;
  className?: string;
}

export function EvaluationBar({ evaluation, displayEvaluation, className }: EvaluationBarProps) {
  // Use enhanced display evaluation if available
  if (displayEvaluation) {
    // Convert evaluationBar (-1 to +1) to percentage (0 to 100)
    // evaluationBar: -1 = 0%, 0 = 50%, +1 = 100%
    const whitePercentage = Math.max(0, Math.min(100, 50 + (displayEvaluation.evaluationBar * 50)));
    const blackPercentage = 100 - whitePercentage;

    const formatDisplayEvaluation = (displayEval: DisplayEvaluation): string => {
      const score = displayEval.displayScore;
      if (Math.abs(score) > 2000) {
        const mateIn = Math.ceil((3000 - Math.abs(score)) / 100);
        return score > 0 ? `M${mateIn}` : `-M${mateIn}`;
      }
      const displayScore = (score / 100).toFixed(1);
      return score > 0 ? `+${displayScore}` : displayScore;
    };

    const evalText = formatDisplayEvaluation(displayEvaluation);

    return (
      <div className={cn("relative h-full w-8 flex flex-col justify-center items-center", className)}>
        <div 
          className="absolute top-0 w-full bg-gray-100 dark:bg-gray-200"
          style={{ height: `${blackPercentage}%`, transition: 'height 0.15s ease-out' }}
        />
        <div 
          className="absolute bottom-0 w-full bg-gray-800 dark:bg-gray-700"
          style={{ height: `${whitePercentage}%`, transition: 'height 0.15s ease-out' }}
        />
        
        <span className="relative text-white dark:text-gray-800 font-bold text-xs mix-blend-difference">
          {evalText}
        </span>
      </div>
    );
  }

  // Fallback to raw evaluation handling
  const VIZ_EVAL_CAP = 1000;
  const clampedEval = Math.max(-VIZ_EVAL_CAP, Math.min(VIZ_EVAL_CAP, evaluation));
  const whitePercentage = 50 + (clampedEval / VIZ_EVAL_CAP) * 50;

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