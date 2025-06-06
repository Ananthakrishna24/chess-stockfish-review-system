'use client';

import { GameAnalysis } from '@/types/analysis';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card';

interface GameSummaryProps {
  gameAnalysis: GameAnalysis;
  gameInfo: {
    white: string;
    black: string;
    whiteRating?: string;
    blackRating?: string;
    result?: string;
    date?: string;
    event?: string;
    opening?: string;
    eco?: string;
  };
}

export function GameSummary({ gameAnalysis, gameInfo }: GameSummaryProps) {
  const totalMoves = gameAnalysis.moves.length;
  const gameLength = Math.ceil(totalMoves / 2);
  
  // Determine winner
  const getWinner = () => {
    if (!gameInfo.result) return null;
    if (gameInfo.result === '1-0') return 'white';
    if (gameInfo.result === '0-1') return 'black';
    return 'draw';
  };

  const winner = getWinner();
  
  // Calculate game statistics
  const whiteAccuracy = gameAnalysis.whiteStats.accuracy;
  const blackAccuracy = gameAnalysis.blackStats.accuracy;
  const averageAccuracy = (whiteAccuracy + blackAccuracy) / 2;
  
  // Game phase statistics
  const openingLength = gameAnalysis.gamePhases.opening;
  const middlegameLength = gameAnalysis.gamePhases.middlegame - gameAnalysis.gamePhases.opening;
  const endgameLength = totalMoves - gameAnalysis.gamePhases.middlegame;

  const getResultDisplay = () => {
    switch (gameInfo.result) {
      case '1-0': return '1-0';
      case '0-1': return '0-1';
      case '1/2-1/2': return '½-½';
      default: return '*';
    }
  };

  const getResultColor = () => {
    switch (winner) {
      case 'white': return 'text-green-600';
      case 'black': return 'text-red-600';
      default: return 'text-gray-600';
    }
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return 'Unknown';
    try {
      return new Date(dateString).toLocaleDateString();
    } catch {
      return dateString;
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Game Summary</CardTitle>
      </CardHeader>
      
      <CardContent className="space-y-6">
        {/* Game Information */}
        <div>
          <div className="text-sm font-medium text-gray-700 mb-3">Game Information</div>
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span className="text-gray-500">Event:</span>
              <div className="font-medium">{gameInfo.event || 'Unknown'}</div>
            </div>
            <div>
              <span className="text-gray-500">Date:</span>
              <div className="font-medium">{formatDate(gameInfo.date)}</div>
            </div>
            <div>
              <span className="text-gray-500">Result:</span>
              <div className={`font-bold text-lg ${getResultColor()}`}>
                {getResultDisplay()}
              </div>
            </div>
            <div>
              <span className="text-gray-500">Moves:</span>
              <div className="font-medium">{gameLength}</div>
            </div>
          </div>
        </div>

        {/* Opening Analysis */}
        {gameAnalysis.openingAnalysis && (
          <div className="border-t pt-4">
            <div className="text-sm font-medium text-gray-700 mb-3">Opening</div>
            <div className="space-y-2">
              <div className="flex justify-between items-center">
                <span className="text-sm text-gray-500">Name:</span>
                <span className="text-sm font-medium">
                  {gameAnalysis.openingAnalysis.name}
                </span>
              </div>
              {gameAnalysis.openingAnalysis.eco && (
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-500">ECO:</span>
                  <span className="text-sm font-medium">
                    {gameAnalysis.openingAnalysis.eco}
                  </span>
                </div>
              )}
              <div className="flex justify-between items-center">
                <span className="text-sm text-gray-500">Opening Accuracy:</span>
                <span className="text-sm font-medium">
                  {gameAnalysis.phaseAnalysis?.openingAccuracy.toFixed(1)}%
                </span>
              </div>
            </div>
          </div>
        )}

        {/* Game Phase Analysis */}
        <div className="border-t pt-4">
          <div className="text-sm font-medium text-gray-700 mb-3">Game Phases</div>
          <div className="space-y-3">
            <div className="flex justify-between items-center">
              <div className="flex items-center space-x-2">
                <div className="w-3 h-3 bg-blue-500 rounded-full"></div>
                <span className="text-sm">Opening</span>
              </div>
              <div className="text-right">
                <div className="text-sm font-medium">{openingLength} moves</div>
                <div className="text-xs text-gray-500">
                  {gameAnalysis.phaseAnalysis?.openingAccuracy.toFixed(1)}% accuracy
                </div>
              </div>
            </div>
            
            <div className="flex justify-between items-center">
              <div className="flex items-center space-x-2">
                <div className="w-3 h-3 bg-green-500 rounded-full"></div>
                <span className="text-sm">Middlegame</span>
              </div>
              <div className="text-right">
                <div className="text-sm font-medium">{middlegameLength} moves</div>
                <div className="text-xs text-gray-500">
                  {gameAnalysis.phaseAnalysis?.middlegameAccuracy.toFixed(1)}% accuracy
                </div>
              </div>
            </div>
            
            <div className="flex justify-between items-center">
              <div className="flex items-center space-x-2">
                <div className="w-3 h-3 bg-orange-500 rounded-full"></div>
                <span className="text-sm">Endgame</span>
              </div>
              <div className="text-right">
                <div className="text-sm font-medium">{endgameLength} moves</div>
                <div className="text-xs text-gray-500">
                  {gameAnalysis.phaseAnalysis?.endgameAccuracy.toFixed(1)}% accuracy
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Performance Overview */}
        <div className="border-t pt-4">
          <div className="text-sm font-medium text-gray-700 mb-3">Performance Overview</div>
          <div className="space-y-2">
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-500">Average Accuracy:</span>
              <span className="text-sm font-bold text-blue-600">
                {averageAccuracy.toFixed(1)}%
              </span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-500">Critical Moments:</span>
              <span className="text-sm font-medium">
                {gameAnalysis.criticalMoments?.length || 0}
              </span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-500">Total Blunders:</span>
              <span className="text-sm font-medium text-red-600">
                {gameAnalysis.whiteStats.blunder + gameAnalysis.blackStats.blunder}
              </span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-500">Brilliant Moves:</span>
              <span className="text-sm font-medium text-cyan-600">
                {gameAnalysis.whiteStats.brilliant + gameAnalysis.blackStats.brilliant}
              </span>
            </div>
          </div>
        </div>

        {/* Winner Analysis */}
        {winner && winner !== 'draw' && (
          <div className="border-t pt-4">
            <div className="text-sm font-medium text-gray-700 mb-3">Victory Analysis</div>
            <div className="space-y-2">
              <div className="flex justify-between items-center">
                <span className="text-sm text-gray-500">Winner:</span>
                <span className={`text-sm font-bold ${getResultColor()}`}>
                  {winner === 'white' ? gameInfo.white : gameInfo.black}
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-gray-500">Winner Accuracy:</span>
                <span className="text-sm font-medium">
                  {winner === 'white' ? whiteAccuracy.toFixed(1) : blackAccuracy.toFixed(1)}%
                </span>
              </div>
              {gameAnalysis.gameResult?.winningAdvantage && (
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-500">Max Advantage:</span>
                  <span className="text-sm font-medium">
                    +{(gameAnalysis.gameResult.winningAdvantage / 100).toFixed(1)}
                  </span>
                </div>
              )}
            </div>
          </div>
        )}

        {/* Game Quality Assessment */}
        <div className="border-t pt-4">
          <div className="text-sm font-medium text-gray-700 mb-3">Game Quality</div>
          <div className="flex justify-center">
            <div className="text-center">
              <div className={`text-2xl font-bold ${
                averageAccuracy >= 90 ? 'text-green-600' :
                averageAccuracy >= 85 ? 'text-blue-600' :
                averageAccuracy >= 80 ? 'text-yellow-600' :
                averageAccuracy >= 75 ? 'text-orange-600' : 'text-red-600'
              }`}>
                {averageAccuracy >= 90 ? 'Excellent' :
                 averageAccuracy >= 85 ? 'Very Good' :
                 averageAccuracy >= 80 ? 'Good' :
                 averageAccuracy >= 75 ? 'Fair' : 'Poor'}
              </div>
              <div className="text-sm text-gray-500">Overall Quality</div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
} 