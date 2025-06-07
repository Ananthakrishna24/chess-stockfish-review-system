'use client';

import { PlayerStatistics } from '@/types/analysis';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card';

interface PlayerStatsProps {
  playerName: string;
  playerRating?: string;
  statistics: PlayerStatistics;
  color: 'white' | 'black';
  isWinner?: boolean;
}

export function PlayerStats({
  playerName,
  playerRating,
  statistics,
  color,
  isWinner = false
}: PlayerStatsProps) {
  const pieceSymbol = color === 'white' ? '♔' : '♚';
  const bgColor = color === 'white' ? 'bg-gray-100' : 'bg-gray-800';
  const textColor = color === 'white' ? 'text-gray-900' : 'text-white';

  const getAccuracyColor = (accuracy: number) => {
    if (accuracy >= 95) return 'text-green-600';
    if (accuracy >= 90) return 'text-blue-600';
    if (accuracy >= 85) return 'text-yellow-600';
    if (accuracy >= 80) return 'text-orange-600';
    return 'text-red-600';
  };

  const getAccuracyBadge = (accuracy: number) => {
    if (accuracy >= 95) return 'bg-green-100 text-green-800';
    if (accuracy >= 90) return 'bg-blue-100 text-blue-800';
    if (accuracy >= 85) return 'bg-yellow-100 text-yellow-800';
    if (accuracy >= 80) return 'bg-orange-100 text-orange-800';
    return 'bg-red-100 text-red-800';
  };

  const totalMoves = statistics.brilliant + statistics.great + statistics.best + 
                    statistics.excellent + statistics.good + (statistics.book || 0) + statistics.inaccuracy + statistics.mistake + 
                    statistics.blunder + statistics.miss;

  const positiveMovesPercentage = totalMoves > 0 
    ? ((statistics.brilliant + statistics.great + statistics.best + statistics.excellent + statistics.good + (statistics.book || 0)) / totalMoves * 100)
    : 0;

  return (
    <Card className="h-full">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className={`w-10 h-10 ${bgColor} rounded-full flex items-center justify-center ${textColor} text-xl`}>
              {pieceSymbol}
            </div>
            <div>
              <div className="flex items-center space-x-2">
                <CardTitle className="text-lg">{playerName}</CardTitle>
                {isWinner && (
                  <span className="text-xs bg-green-100 text-green-800 px-2 py-1 rounded-full font-medium">
                    Winner
                  </span>
                )}
              </div>
              <div className="text-sm text-gray-500">
                {playerRating || 'Unrated'}
              </div>
            </div>
          </div>
          <div className="text-right">
            <div className={`text-2xl font-bold ${getAccuracyColor(statistics.accuracy)}`}>
              {statistics.accuracy.toFixed(1)}%
            </div>
            <div className={`text-xs px-2 py-1 rounded-full font-medium ${getAccuracyBadge(statistics.accuracy)}`}>
              Accuracy
            </div>
          </div>
        </div>
      </CardHeader>
      
      <CardContent className="space-y-4">
        {/* Move Quality Breakdown */}
        <div>
          <div className="text-sm font-medium text-gray-700 mb-3">Move Quality</div>
          <div className="space-y-2">
            {[
              { type: 'brilliant', count: statistics.brilliant, color: 'bg-cyan-500', label: 'Brilliant' },
              { type: 'great', count: statistics.great, color: 'bg-blue-500', label: 'Great' },
              { type: 'best', count: statistics.best, color: 'bg-green-500', label: 'Best' },
              { type: 'excellent', count: statistics.excellent, color: 'bg-green-400', label: 'Excellent' },
              { type: 'good', count: statistics.good, color: 'bg-green-300', label: 'Good' },
              { type: 'book', count: statistics.book || 0, color: 'bg-purple-500', label: 'Book' },
              { type: 'inaccuracy', count: statistics.inaccuracy, color: 'bg-yellow-500', label: 'Inaccuracy' },
              { type: 'mistake', count: statistics.mistake, color: 'bg-orange-500', label: 'Mistake' },
              { type: 'blunder', count: statistics.blunder, color: 'bg-red-500', label: 'Blunder' },
              { type: 'miss', count: statistics.miss, color: 'bg-red-600', label: 'Miss' }
            ].filter(item => item.count > 0).map(({ type, count, color, label }) => (
              <div key={type} className="flex justify-between items-center">
                <div className="flex items-center space-x-2">
                  <div className={`w-3 h-3 ${color} rounded-full`}></div>
                  <span className="text-sm">{label}</span>
                </div>
                <div className="flex items-center space-x-2">
                  <span className="text-sm font-medium">{count}</span>
                  <span className="text-xs text-gray-500">
                    {totalMoves > 0 ? `${(count / totalMoves * 100).toFixed(0)}%` : '0%'}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Tactical Statistics */}
        {(statistics.tacticalMoves || statistics.forcingMoves || statistics.criticalMoments) && (
          <div className="border-t pt-4">
            <div className="text-sm font-medium text-gray-700 mb-3">Tactical Performance</div>
            <div className="space-y-2">
              {statistics.tacticalMoves !== undefined && (
                <div className="flex justify-between items-center">
                  <span className="text-sm">Tactical Moves</span>
                  <span className="text-sm font-medium">{statistics.tacticalMoves}</span>
                </div>
              )}
              {statistics.forcingMoves !== undefined && (
                <div className="flex justify-between items-center">
                  <span className="text-sm">Forcing Moves</span>
                  <span className="text-sm font-medium">{statistics.forcingMoves}</span>
                </div>
              )}
              {statistics.criticalMoments !== undefined && (
                <div className="flex justify-between items-center">
                  <span className="text-sm">Critical Moments</span>
                  <span className="text-sm font-medium">{statistics.criticalMoments}</span>
                </div>
              )}
            </div>
          </div>
        )}

        {/* Performance Summary */}
        <div className="border-t pt-4">
          <div className="text-sm font-medium text-gray-700 mb-3">Performance Summary</div>
          <div className="space-y-2">
            <div className="flex justify-between items-center">
              <span className="text-sm">Total Moves</span>
              <span className="text-sm font-medium">{totalMoves}</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm">Good Moves</span>
              <span className="text-sm font-medium text-green-600">
                {positiveMovesPercentage.toFixed(0)}%
              </span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm">Errors</span>
              <span className="text-sm font-medium text-red-600">
                {statistics.inaccuracy + statistics.mistake + statistics.blunder + statistics.miss}
              </span>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
} 