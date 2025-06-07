'use client';

import { EngineEvaluation } from '@/types/analysis';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card';
import { convertScoreToString } from '@/utils/stockfish';

interface EvaluationChartProps {
  evaluations: EngineEvaluation[];
  currentMoveIndex: number;
  criticalMoments?: number[];
  onMoveClick?: (moveIndex: number) => void;
}

export function EvaluationChart({
  evaluations,
  currentMoveIndex,
  criticalMoments = [],
  onMoveClick
}: EvaluationChartProps) {
  if (evaluations.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Game Evaluation</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-center text-gray-500 py-8">
            No evaluation data available
          </div>
        </CardContent>
      </Card>
    );
  }

  // Calculate chart dimensions and scaling
  const chartWidth = 300; // Fit within reduced sidebar
  const chartHeight = 140;
  const margin = { top: 20, right: 40, bottom: 40, left: 40 };
  const innerWidth = chartWidth - margin.left - margin.right;
  const innerHeight = chartHeight - margin.top - margin.bottom;

  // Find min/max scores for scaling
  const scores = evaluations.map(evaluation => evaluation.mate ? (evaluation.mate > 0 ? 1000 : -1000) : evaluation.score);
  const minScore = Math.min(-500, Math.min(...scores));
  const maxScore = Math.max(500, Math.max(...scores));

  // Scale functions
  const xScale = (index: number) => (index / (evaluations.length - 1)) * innerWidth;
  const yScale = (score: number) => {
    const normalizedScore = Math.max(minScore, Math.min(maxScore, score));
    return innerHeight - ((normalizedScore - minScore) / (maxScore - minScore)) * innerHeight;
  };

  // Generate path data
  const pathData = evaluations.map((evaluation, index) => {
    const score = evaluation.mate 
      ? (evaluation.mate > 0 ? 1000 : -1000)
      : evaluation.score;
    const x = xScale(index);
    const y = yScale(score);
    return `${index === 0 ? 'M' : 'L'} ${x} ${y}`;
  }).join(' ');

  // Calculate zero line position
  const zeroY = yScale(0);

  // Get current evaluation
  const currentEval = evaluations[currentMoveIndex];
  const currentScore = currentEval?.mate 
    ? (currentEval.mate > 0 ? 1000 : -1000)
    : currentEval?.score || 0;

  return (
    <Card>
      <CardHeader>
        <div className="flex justify-between items-center">
          <CardTitle>Game Evaluation</CardTitle>
          <div className="text-sm">
            <span className="text-gray-500">Current: </span>
            <span className={`font-bold ${
              currentScore > 50 ? 'text-green-600' : 
              currentScore < -50 ? 'text-red-600' : 'text-gray-600'
            }`}>
              {currentEval ? convertScoreToString(currentEval.score, currentEval.mate) : '0.0'}
            </span>
          </div>
        </div>
      </CardHeader>
      
      <CardContent>
        <div className="relative">
          {/* Chart SVG */}
          <svg 
            width={chartWidth} 
            height={chartHeight}
            className="border rounded overflow-hidden bg-gray-50 w-full"
            viewBox={`0 0 ${chartWidth} ${chartHeight}`}
          >
            {/* Grid lines */}
            <defs>
              <pattern id="grid" width="20" height="20" patternUnits="userSpaceOnUse">
                <path d="M 20 0 L 0 0 0 20" fill="none" stroke="#e5e7eb" strokeWidth="1"/>
              </pattern>
            </defs>
            <rect width="100%" height="100%" fill="url(#grid)" />

            {/* Zero line */}
            <line
              x1={margin.left}
              y1={margin.top + zeroY}
              x2={margin.left + innerWidth}
              y2={margin.top + zeroY}
              stroke="#6b7280"
              strokeWidth="2"
              strokeDasharray="5,5"
            />

            {/* Critical moments indicators */}
            {criticalMoments.map((moveIndex, index) => (
              <line
                key={`critical-${moveIndex}-${index}`}
                x1={margin.left + xScale(moveIndex)}
                y1={margin.top}
                x2={margin.left + xScale(moveIndex)}
                y2={margin.top + innerHeight}
                stroke="#f59e0b"
                strokeWidth="2"
                opacity="0.7"
              />
            ))}

            {/* Evaluation line */}
            <path
              d={pathData}
              fill="none"
              stroke="#059669"
              strokeWidth="3"
              transform={`translate(${margin.left}, ${margin.top})`}
            />

            {/* Current move indicator */}
            <circle
              cx={margin.left + xScale(currentMoveIndex)}
              cy={margin.top + yScale(currentScore)}
              r="5"
              fill="#dc2626"
              stroke="white"
              strokeWidth="2"
            />

            {/* Data points (clickable) */}
            {evaluations.map((evaluation, index) => {
              const score = evaluation.mate 
                ? (evaluation.mate > 0 ? 1000 : -1000)
                : evaluation.score;
              return (
                <circle
                  key={index}
                  cx={margin.left + xScale(index)}
                  cy={margin.top + yScale(score)}
                  r="3"
                  fill={index === currentMoveIndex ? "#dc2626" : "#059669"}
                  className="cursor-pointer hover:r-4"
                  onClick={() => onMoveClick?.(index)}
                />
              );
            })}

            {/* Y-axis labels */}
            <text x="5" y={margin.top + 5} fontSize="10" fill="#6b7280">+5</text>
            <text x="5" y={margin.top + zeroY + 4} fontSize="10" fill="#6b7280">0</text>
            <text x="5" y={margin.top + innerHeight} fontSize="10" fill="#6b7280">-5</text>

            {/* X-axis labels */}
            <text x={margin.left} y={chartHeight - 5} fontSize="10" fill="#6b7280">1</text>
            <text 
              x={margin.left + innerWidth} 
              y={chartHeight - 5} 
              fontSize="10" 
              fill="#6b7280"
              textAnchor="end"
            >
              {evaluations.length}
            </text>
          </svg>

          {/* Legend */}
          <div className="mt-4 flex flex-wrap gap-4 text-sm">
            <div className="flex items-center space-x-2">
              <div className="w-3 h-0.5 bg-green-600"></div>
              <span>Evaluation</span>
            </div>
            <div className="flex items-center space-x-2">
              <div className="w-3 h-3 bg-red-600 rounded-full"></div>
              <span>Current Move</span>
            </div>
            {criticalMoments.length > 0 && (
              <div className="flex items-center space-x-2">
                <div className="w-0.5 h-3 bg-yellow-500"></div>
                <span>Critical Moment</span>
              </div>
            )}
            <div className="flex items-center space-x-2">
              <div className="w-3 h-0.5 bg-gray-400 border-dashed border"></div>
              <span>Equal Position</span>
            </div>
          </div>

          {/* Game phases indicator */}
          <div className="mt-3 flex justify-between text-xs text-gray-500">
            <span>Opening</span>
            <span>Middlegame</span>
            <span>Endgame</span>
          </div>

          {/* Evaluation summary */}
          <div className="mt-4 grid grid-cols-3 gap-4 text-sm">
            <div className="text-center">
              <div className="font-medium text-gray-900">
                {Math.max(...scores.filter(s => s > 0)).toFixed(0)}
              </div>
              <div className="text-gray-500">Max Advantage</div>
            </div>
            <div className="text-center">
              <div className="font-medium text-gray-900">
                {criticalMoments.length}
              </div>
              <div className="text-gray-500">Critical Moments</div>
            </div>
            <div className="text-center">
              <div className="font-medium text-gray-900">
                {scores.filter(s => Math.abs(s) < 50).length}
              </div>
              <div className="text-gray-500">Equal Positions</div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
} 