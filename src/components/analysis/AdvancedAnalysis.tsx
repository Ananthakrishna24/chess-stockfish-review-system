'use client';

import React, { useState, useEffect } from 'react';
import { Card, CardHeader, CardContent, CardTitle } from '@/components/ui/Card';
// import Button from '@/components/ui/Button';
import { GameAnalysis, TimeAnalysis } from '@/types/analysis';
import { GameState } from '@/types/chess';
import { analyzeTimeManagement } from '@/utils/positionalAnalysis';
import { getOpeningAnalysis } from '@/utils/openingBook';
import { analyzeEndgame } from '@/utils/endgameTablebase';
import { analyzePosition } from '@/utils/positionalAnalysis';

interface AdvancedAnalysisProps {
  gameState: GameState;
  gameAnalysis?: GameAnalysis;
  currentPosition: string;
  currentMoveIndex: number;
  className?: string;
}

export function AdvancedAnalysis({ 
  gameState, 
  gameAnalysis, 
  currentPosition, 
  currentMoveIndex, 
  className = '' 
}: AdvancedAnalysisProps) {
  const [activeTab, setActiveTab] = useState<'time' | 'positional' | 'opening' | 'endgame'>('time');
  const [timeAnalysis, setTimeAnalysis] = useState<TimeAnalysis | null>(null);
  const [positionalAnalysis, setPositionalAnalysis] = useState<any>(null); // eslint-disable-line @typescript-eslint/no-explicit-any
  const [openingInfo, setOpeningInfo] = useState<any>(null); // eslint-disable-line @typescript-eslint/no-explicit-any
  const [endgameInfo, setEndgameInfo] = useState<any>(null); // eslint-disable-line @typescript-eslint/no-explicit-any
  const [isAnalyzing, setIsAnalyzing] = useState(false);

  // Analyze time management when game analysis is available
  useEffect(() => {
    if (gameAnalysis && gameState) {
      // Mock time data - in a real implementation, you'd get this from the PGN or user input
      const mockMoveTimes = gameState.moves.map(() => Math.random() * 180 + 20); // 20-200 seconds
      const analysis = analyzeTimeManagement(mockMoveTimes, gameAnalysis.gamePhases);
      setTimeAnalysis(analysis);
    }
  }, [gameAnalysis, gameState]);

  // Analyze current position
  useEffect(() => {
    if (currentPosition) {
      setIsAnalyzing(true);
      try {
        const posAnalysis = analyzePosition(currentPosition);
        setPositionalAnalysis(posAnalysis);
      } catch (error) {
        console.error('Error analyzing position:', error);
      } finally {
        setIsAnalyzing(false);
      }
    }
  }, [currentPosition]);

  // Get opening information
  useEffect(() => {
    if (gameState && currentMoveIndex < 15) {
      const moves = gameState.moves.slice(0, currentMoveIndex + 1).map(m => m.san);
      const opening = getOpeningAnalysis(moves);
      setOpeningInfo(opening);
    }
  }, [gameState, currentMoveIndex]);

  // Query endgame tablebase for positions with few pieces
  useEffect(() => {
    if (currentPosition) {
      const pieceCount = currentPosition.split('').filter(c => /[a-zA-Z]/.test(c)).length;
      if (pieceCount <= 7) {
        try {
          const result = analyzeEndgame(currentPosition);
          setEndgameInfo(result);
        } catch (error) {
          console.error('Endgame analysis failed:', error);
          setEndgameInfo(null);
        }
      } else {
        setEndgameInfo(null);
      }
    }
  }, [currentPosition]);

  const formatTime = (seconds: number): string => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  const tabs = [
    { id: 'time', label: 'Time Management', icon: '⏱️' },
    { id: 'positional', label: 'Position Analysis', icon: '♟️' },
    { id: 'opening', label: 'Opening Book', icon: '📖' },
    { id: 'endgame', label: 'Endgame', icon: '👑' }
  ];

  return (
    <Card className={className}>
      <CardHeader>
        <CardTitle>Advanced Analysis</CardTitle>
        
        {/* Tab Navigation */}
        <div className="flex space-x-1 bg-gray-100 rounded-lg p-1">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id as any)}
              className={`flex-1 flex items-center justify-center px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                activeTab === tab.id
                  ? 'bg-white text-gray-900 shadow-sm'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              <span className="mr-1">{tab.icon}</span>
              <span className="hidden sm:inline">{tab.label}</span>
            </button>
          ))}
        </div>
      </CardHeader>

      <CardContent>
        {/* Time Management Tab */}
        {activeTab === 'time' && timeAnalysis && (
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="bg-gray-50 p-3 rounded-lg">
                <div className="text-sm font-medium text-gray-600">Total Time</div>
                <div className="text-lg font-bold">{formatTime(timeAnalysis.timeSpent)}</div>
              </div>
              <div className="bg-gray-50 p-3 rounded-lg">
                <div className="text-sm font-medium text-gray-600">Avg per Move</div>
                <div className="text-lg font-bold">{formatTime(timeAnalysis.averageTimePerMove)}</div>
              </div>
            </div>

            <div>
              <div className="text-sm font-medium text-gray-700 mb-2">Time Distribution</div>
              <div className="space-y-2">
                <div className="flex justify-between items-center">
                  <span className="text-sm">Opening</span>
                  <span className="text-sm font-medium">{formatTime(timeAnalysis.timeDistribution.opening)}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm">Middlegame</span>
                  <span className="text-sm font-medium">{formatTime(timeAnalysis.timeDistribution.middlegame)}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm">Endgame</span>
                  <span className="text-sm font-medium">{formatTime(timeAnalysis.timeDistribution.endgame)}</span>
                </div>
              </div>
            </div>

            {timeAnalysis.recommendations.length > 0 && (
              <div>
                <div className="text-sm font-medium text-gray-700 mb-2">Recommendations</div>
                <div className="space-y-1">
                  {timeAnalysis.recommendations.map((rec, index) => (
                    <div key={index} className="text-sm text-gray-600 flex items-start">
                      <span className="text-blue-500 mr-2">•</span>
                      {rec}
                    </div>
                  ))}
                </div>
              </div>
            )}

            {timeAnalysis.criticalMoments.length > 0 && (
              <div>
                <div className="text-sm font-medium text-gray-700 mb-2">
                  Long Thinks ({timeAnalysis.criticalMoments.length} moves)
                </div>
                <div className="text-xs text-gray-500">
                  Moves: {timeAnalysis.criticalMoments.map(m => m + 1).join(', ')}
                </div>
              </div>
            )}
          </div>
        )}

        {/* Positional Analysis Tab */}
        {activeTab === 'positional' && (
          <div className="space-y-4">
            {isAnalyzing ? (
              <div className="flex items-center justify-center py-8">
                <div className="text-sm text-gray-500">Analyzing position...</div>
              </div>
            ) : positionalAnalysis ? (
              <div className="space-y-4">
                <div className="bg-gray-50 p-3 rounded-lg">
                  <div className="text-sm font-medium text-gray-600">Position Type</div>
                  <div className="text-lg font-bold capitalize">{positionalAnalysis.characterization}</div>
                  <div className="text-sm text-gray-500">
                    Overall Score: {positionalAnalysis.overallScore > 0 ? '+' : ''}{positionalAnalysis.overallScore.toFixed(2)}
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-3">
                  <div className="bg-blue-50 p-3 rounded-lg">
                    <div className="text-xs font-medium text-blue-700">Pawn Structure</div>
                    <div className="text-sm font-bold text-blue-900">
                      {positionalAnalysis.factors.pawnStructure.score > 0 ? '+' : ''}{positionalAnalysis.factors.pawnStructure.score.toFixed(2)}
                    </div>
                  </div>
                  <div className="bg-green-50 p-3 rounded-lg">
                    <div className="text-xs font-medium text-green-700">Piece Activity</div>
                    <div className="text-sm font-bold text-green-900">
                      {positionalAnalysis.factors.pieceActivity.score > 0 ? '+' : ''}{positionalAnalysis.factors.pieceActivity.score.toFixed(2)}
                    </div>
                  </div>
                  <div className="bg-red-50 p-3 rounded-lg">
                    <div className="text-xs font-medium text-red-700">King Safety</div>
                    <div className="text-sm font-bold text-red-900">
                      {positionalAnalysis.factors.kingSafety.score > 0 ? '+' : ''}{positionalAnalysis.factors.kingSafety.score.toFixed(2)}
                    </div>
                  </div>
                  <div className="bg-purple-50 p-3 rounded-lg">
                    <div className="text-xs font-medium text-purple-700">Space Advantage</div>
                    <div className="text-sm font-bold text-purple-900">
                      {positionalAnalysis.factors.spaceAdvantage.score > 0 ? '+' : ''}{positionalAnalysis.factors.spaceAdvantage.score.toFixed(2)}
                    </div>
                  </div>
                </div>

                {positionalAnalysis.recommendations.immediate.length > 0 && (
                  <div>
                    <div className="text-sm font-medium text-gray-700 mb-2">Immediate Recommendations</div>
                    <div className="space-y-1">
                      {positionalAnalysis.recommendations.immediate.map((rec: string, index: number) => (
                        <div key={index} className="text-sm text-gray-600 flex items-start">
                          <span className="text-red-500 mr-2">•</span>
                          {rec}
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {positionalAnalysis.imbalances.length > 0 && (
                  <div>
                    <div className="text-sm font-medium text-gray-700 mb-2">Position Imbalances</div>
                    <div className="flex flex-wrap gap-1">
                      {positionalAnalysis.imbalances.map((imbalance: string, index: number) => (
                        <span key={index} className="px-2 py-1 bg-yellow-100 text-yellow-800 text-xs rounded-full">
                          {imbalance}
                        </span>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ) : (
              <div className="text-sm text-gray-500 text-center py-4">
                No positional analysis available
              </div>
            )}
          </div>
        )}

        {/* Opening Book Tab */}
        {activeTab === 'opening' && (
          <div className="space-y-4">
            {openingInfo ? (
              <div className="space-y-4">
                <div className="bg-gray-50 p-3 rounded-lg">
                  <div className="text-sm font-medium text-gray-600">Opening</div>
                  <div className="text-lg font-bold">{openingInfo.name}</div>
                  {openingInfo.eco && (
                    <div className="text-sm text-gray-500">ECO: {openingInfo.eco}</div>
                  )}
                </div>

                {openingInfo.moves && openingInfo.moves.length > 0 && (
                  <div>
                    <div className="text-sm font-medium text-gray-700 mb-2">Theory Moves</div>
                    <div className="space-y-2">
                      {openingInfo.moves.slice(0, 3).map((move: any, index: number) => (
                        <div key={index} className="flex justify-between items-center p-2 bg-gray-50 rounded">
                          <span className="font-mono text-sm">{move.san}</span>
                          <span className="text-xs text-gray-500">{move.frequency}% played</span>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {openingInfo.statistics && (
                  <div className="grid grid-cols-3 gap-2">
                    <div className="text-center p-2 bg-white rounded border">
                      <div className="text-xs text-gray-500">White Wins</div>
                      <div className="font-bold">{openingInfo.statistics.white}%</div>
                    </div>
                    <div className="text-center p-2 bg-gray-100 rounded border">
                      <div className="text-xs text-gray-500">Draws</div>
                      <div className="font-bold">{openingInfo.statistics.draws}%</div>
                    </div>
                    <div className="text-center p-2 bg-gray-900 text-white rounded border">
                      <div className="text-xs">Black Wins</div>
                      <div className="font-bold">{openingInfo.statistics.black}%</div>
                    </div>
                  </div>
                )}

                {openingInfo.ideas && openingInfo.ideas.length > 0 && (
                  <div>
                    <div className="text-sm font-medium text-gray-700 mb-2">Key Ideas</div>
                    <div className="space-y-1">
                      {openingInfo.ideas.map((idea: string, index: number) => (
                        <div key={index} className="text-sm text-gray-600 flex items-start">
                          <span className="text-green-500 mr-2">•</span>
                          {idea}
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ) : currentMoveIndex >= 15 ? (
              <div className="text-sm text-gray-500 text-center py-4">
                Opening phase completed
              </div>
            ) : (
              <div className="text-sm text-gray-500 text-center py-4">
                No opening information available
              </div>
            )}
          </div>
        )}

        {/* Endgame Tab */}
        {activeTab === 'endgame' && (
          <div className="space-y-4">
            {endgameInfo ? (
              <div className="space-y-4">
                                 <div className="bg-gray-50 p-3 rounded-lg">
                   <div className="text-sm font-medium text-gray-600">Endgame Analysis</div>
                   <div className="text-lg font-bold capitalize">
                     {endgameInfo.result === 'win' ? '🏆 Winning' : 
                      endgameInfo.result === 'loss' ? '💔 Losing' : 
                      endgameInfo.result === 'draw' ? '🤝 Draw' : 'Unknown'}
                   </div>
                   <div className="text-sm text-gray-500">
                     {endgameInfo.classification}
                   </div>
                   {endgameInfo.movesToMate && (
                     <div className="text-sm text-gray-500">
                       Mate in {endgameInfo.movesToMate} moves
                     </div>
                   )}
                 </div>

                 {endgameInfo.technique && endgameInfo.technique.length > 0 && (
                   <div>
                     <div className="text-sm font-medium text-gray-700 mb-2">Technique</div>
                     <div className="space-y-1">
                       {endgameInfo.technique.slice(0, 3).map((tech: string, index: number) => (
                         <div key={index} className="text-sm text-gray-600 flex items-start">
                           <span className="text-blue-500 mr-2">•</span>
                           {tech}
                         </div>
                       ))}
                     </div>
                   </div>
                 )}

                 {endgameInfo.winningMethod && (
                   <div>
                     <div className="text-sm font-medium text-gray-700 mb-2">Winning Method</div>
                     <div className="text-sm text-gray-600">
                       {endgameInfo.winningMethod}
                     </div>
                   </div>
                 )}

                                  <div className="bg-blue-50 p-3 rounded-lg">
                   <div className="text-sm text-blue-700">
                     💡 Difficulty: {endgameInfo.difficulty} • Theoretical endgame analysis
                   </div>
                 </div>
              </div>
            ) : (
              <div className="text-sm text-gray-500 text-center py-4">
                {currentPosition ? (
                  currentPosition.split('').filter(c => /[a-zA-Z]/.test(c)).length > 7 ? 
                    'Too many pieces for tablebase lookup' :
                    'Querying endgame tablebase...'
                ) : 'No position loaded'}
              </div>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  );
} 