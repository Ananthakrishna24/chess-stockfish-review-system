'use client';

import { EngineEvaluation, DisplayEvaluation } from '@/types/analysis';

interface EvaluationChartProps {
  evaluations: EngineEvaluation[];
  displayEvaluations?: DisplayEvaluation[]; // Enhanced Lichess evaluation data
  currentMoveIndex: number;
  criticalMoments?: number[];
  onMoveClick?: (moveIndex: number) => void;
  moveClassifications?: string[];
}

export function EvaluationChart({
  evaluations,
  displayEvaluations,
  currentMoveIndex,
  criticalMoments = [],
  onMoveClick,
  moveClassifications = []
}: EvaluationChartProps) {
  if (evaluations.length === 0) {
    return (
      <div className="w-full h-32 bg-gray-100 rounded-lg flex items-center justify-center">
        <div className="text-gray-500 text-sm">No evaluation data</div>
      </div>
    );
  }

  const width = 400; // Match the full container width
  const height = 120;
  const margin = { top: 0, right: 0, bottom: 0, left: 0 }; // Remove all margins
  const chartWidth = width - margin.left - margin.right;
  const chartHeight = height - margin.top - margin.bottom;

  // Use display evaluations if available (Lichess algorithm), otherwise fall back to raw evaluations
  const useDisplayEvaluations = displayEvaluations && displayEvaluations.length === evaluations.length;

  // Scale functions - Updated for Lichess algorithm
  const xScale = (index: number) => (index / Math.max(1, evaluations.length - 1)) * chartWidth;
  const yScale = (score: number) => {
    // Use Lichess capping of ±1000 instead of ±800, and use display scores if available
    const clampedScore = Math.max(-1000, Math.min(1000, score));
    return chartHeight - ((clampedScore + 1000) / 2000) * chartHeight;
  };

  // Calculate zero line position
  const zeroY = yScale(0);

  // Transform evaluation data - Use Lichess display evaluations when available
  const points = evaluations.map((evaluation, index) => {
    let score: number;
    
    if (useDisplayEvaluations && displayEvaluations![index]) {
      // Use Lichess display score (already smoothed and capped)
      score = displayEvaluations![index].displayScore;
    } else {
      // Fallback to raw evaluation with mate handling
      score = evaluation.mate 
        ? (evaluation.mate > 0 ? 1000 : -1000)
        : evaluation.score;
    }
    
    return {
      x: xScale(index),
      y: yScale(score),
      score,
      index,
      classification: moveClassifications[index] || 'none',
      // Add Lichess data if available
      winProbability: useDisplayEvaluations ? displayEvaluations![index]?.winProbability : undefined,
      positionAssessment: useDisplayEvaluations ? displayEvaluations![index]?.positionAssessment : undefined,
      isStable: useDisplayEvaluations ? displayEvaluations![index]?.isStable : undefined
    };
  });

  // Create smooth curve using cubic bezier interpolation
  // Note: Lichess data is already smoothed, so we can use gentler smoothing
  const createSmoothPath = (points: any[]) => {
    if (points.length < 2) return '';
    
    let path = `M ${points[0].x} ${points[0].y}`;
    
    // Use less aggressive smoothing if we have Lichess data (already smoothed)
    const smoothingFactor = useDisplayEvaluations ? 0.1 : 0.15;
    
    for (let i = 1; i < points.length; i++) {
      const prev = points[i - 1];
      const current = points[i];
      const next = points[i + 1];
      
      if (i === 1) {
        // First curve segment
        const cp1x = prev.x + (current.x - prev.x) * 0.3;
        const cp1y = prev.y;
        const cp2x = current.x - (current.x - prev.x) * 0.3;
        const cp2y = current.y;
        path += ` C ${cp1x} ${cp1y}, ${cp2x} ${cp2y}, ${current.x} ${current.y}`;
      } else if (i === points.length - 1) {
        // Last curve segment
        const cp1x = prev.x + (current.x - prev.x) * 0.3;
        const cp1y = prev.y;
        const cp2x = current.x - (current.x - prev.x) * 0.3;
        const cp2y = current.y;
        path += ` C ${cp1x} ${cp1y}, ${cp2x} ${cp2y}, ${current.x} ${current.y}`;
      } else {
        // Middle segments with Lichess-aware smoothing
        const prevPoint = points[i - 2] || prev;
        const cp1x = prev.x + (current.x - prevPoint.x) * smoothingFactor;
        const cp1y = prev.y + (current.y - prevPoint.y) * smoothingFactor;
        const cp2x = current.x - (next.x - prev.x) * smoothingFactor;
        const cp2y = current.y - (next.y - prev.y) * smoothingFactor;
        path += ` C ${cp1x} ${cp1y}, ${cp2x} ${cp2y}, ${current.x} ${current.y}`;
      }
    }
    
    return path;
  };

  const smoothLinePath = createSmoothPath(points);

  // Create white area path (above evaluation line to top)
  const whiteAreaPath = `${smoothLinePath} L ${chartWidth} 0 L 0 0 Z`;

  // Create black area path (below evaluation line to bottom)  
  const blackAreaPath = `${smoothLinePath} L ${chartWidth} ${chartHeight} L 0 ${chartHeight} Z`;

  const getClassificationColor = (classification: string) => {
    switch (classification) {
      case 'brilliant': return '#1e40af';
      case 'great': return '#059669';
      case 'book': return '#7c3aed';
      case 'mistake': return '#f59e0b';
      case 'blunder': return '#dc2626';
      default: return 'transparent';
    }
  };

  return (
    <div className="w-full">
      <div className="relative rounded-lg overflow-hidden">
        <svg 
          width="100%" 
          height={height} 
          viewBox={`0 0 ${width} ${height}`}
          className="w-full"
          preserveAspectRatio="none"
        >
          <g transform={`translate(${margin.left}, ${margin.top})`}>
            {/* White area fill (above evaluation line) */}
            <path
              d={whiteAreaPath}
              fill="white"
              fillOpacity="0.3"
            />

            {/* Black area fill (below evaluation line) */}
            <path
              d={blackAreaPath}
              fill="#374151"
              fillOpacity="0.3"
            />

            {/* Smooth evaluation line */}
            <path
              d={smoothLinePath}
              fill="none"
              stroke="#6b7280"
              strokeWidth={2}
            />

            {/* Move classification dots */}
            {points.map((point) => {
              const showClassifications = ['brilliant', 'great', 'book', 'mistake', 'blunder'];
              if (!showClassifications.includes(point.classification)) return null;
              
              return (
                <circle
                  key={`classification-${point.index}`}
                  cx={point.x}
                  cy={point.y}
                  r={4}
                  fill={getClassificationColor(point.classification)}
                  stroke="white"
                  strokeWidth={1}
                  className="cursor-pointer"
                  onClick={() => onMoveClick?.(point.index)}
                />
              );
            })}

            {/* Current move indicator */}
            {currentMoveIndex >= 0 && currentMoveIndex < points.length && (
              <circle
                cx={points[currentMoveIndex].x}
                cy={points[currentMoveIndex].y}
                r={3}
                fill="#dc2626"
                stroke="white"
                strokeWidth={2}
              />
            )}

            {/* Invisible clickable areas for navigation */}
            {points.map((point) => (
              <circle
                key={`clickable-${point.index}`}
                cx={point.x}
                cy={point.y}
                r={8}
                fill="transparent"
                className="cursor-pointer"
                onClick={() => onMoveClick?.(point.index)}
              />
            ))}
          </g>
        </svg>
      </div>
    </div>
  );
} 