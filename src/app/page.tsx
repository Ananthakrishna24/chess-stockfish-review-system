'use client';

import React from 'react';
import ChessBoard from '@/components/chess/ChessBoard';
import MoveList from '@/components/chess/MoveList';
import GameControls from '@/components/chess/GameControls';
import { PlayerStats } from '@/components/analysis/PlayerStats';
import { EvaluationChart } from '@/components/analysis/EvaluationChart';
import { GameSummary } from '@/components/analysis/GameSummary';
import { ExportOptions } from '@/components/analysis/ExportOptions';
import { AdvancedAnalysis } from '@/components/analysis/AdvancedAnalysis';
import { Card, CardHeader, CardContent, CardTitle, CardDescription } from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import { Textarea } from '@/components/ui/Input';
import { useGameAnalysis } from '@/hooks/useGameAnalysis';
import { convertScoreToString, getScoreColor } from '@/utils/stockfish';

export default function Home() {
  const {
    // Game state
    gameState,
    currentPosition,
    currentMoveIndex,
    isLoading,
    error,
    
    // Navigation
    goToMove,
    goToStart,
    goToEnd,
    goForward,
    goBackward,
    canGoForward,
    canGoBackward,
    isAtStart,
    isAtEnd,
    
    // Analysis
    gameAnalysis,
    isAnalyzingGame,
    analysisProgress,
    // whiteAccuracy,
    // blackAccuracy,
    currentPositionEvaluation,
    
    // Engine state
    engineReady,
    engineInitializing,
    engineError,
    
    // Actions
    loadGame,
    resetGame,
    stopAnalysis
  } = useGameAnalysis();

  const handleLoadPGN = async (pgn: string) => {
    if (!pgn.trim()) {
      alert('Please paste a PGN game first');
      return;
    }
    await loadGame(pgn);
  };

  const samplePGN = `[Event "Rated Blitz game"]
[Site "https://lichess.org/"]
[Date "2024.01.15"]
[White "Player1"]
[Black "Player2"]
[Result "1-0"]
[WhiteElo "1500"]
[BlackElo "1480"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3 O-O 9. h3 Nb8 10. d4 Nbd7 11. Nbd2 Bb7 12. Bc2 Re8 13. Nf1 Bf8 14. Ng3 g6 15. a4 c5 16. d5 Nc4 17. Ra2 c4 18. axb5 axb5 19. Nh4 Qc7 20. Nhf5 gxf5 21. Nxf5 Bg7 22. g3 1-0`;

  const [pgnInput, setPgnInput] = React.useState('');

  // Get current move analysis for display
  const currentMoveAnalysis = gameAnalysis?.moves[currentMoveIndex];
  const currentEval = currentPositionEvaluation || currentMoveAnalysis?.evaluation;

  // Get move classification color
  const getClassificationColor = (classification: string) => {
    switch (classification) {
      case 'brilliant': return 'text-cyan-600 bg-cyan-50';
      case 'great': return 'text-blue-600 bg-blue-50';
      case 'best': return 'text-green-600 bg-green-50';
      case 'good': return 'text-green-700 bg-green-50';
      case 'inaccuracy': return 'text-yellow-600 bg-yellow-50';
      case 'mistake': return 'text-orange-600 bg-orange-50';
      case 'blunder': return 'text-red-600 bg-red-50';
      case 'miss': return 'text-red-700 bg-red-50';
      default: return 'text-gray-600 bg-gray-50';
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center space-x-4">
              <h1 className="text-2xl font-bold text-gray-900">♔ Chess Game Review</h1>
              {engineInitializing && (
                <div className="text-sm text-blue-600">Initializing engine...</div>
              )}
              {engineReady && (
                <div className="text-sm text-green-600">✓ Engine ready</div>
              )}
              {engineError && (
                <div className="text-sm text-red-600">⚠ Engine error</div>
              )}
            </div>
            <div className="flex items-center space-x-4">
              <Button
                variant="outline"
                onClick={() => {
                  setPgnInput(samplePGN);
                  handleLoadPGN(samplePGN);
                }}
              >
                Load Sample Game
              </Button>
              {gameState && (
                <Button
                  variant="outline"
                  onClick={resetGame}
                >
                  Reset
                </Button>
              )}
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Analysis Progress */}
        {isAnalyzingGame && (
          <Card className="mb-6">
            <CardContent className="py-4">
              <div className="flex items-center justify-between">
                <div>
                  <div className="text-sm font-medium text-gray-900">
                    Analyzing Game... {analysisProgress.progress.toFixed(1)}%
                  </div>
                  <div className="text-sm text-gray-500">
                    Move {analysisProgress.currentMove} of {analysisProgress.totalMoves}
                  </div>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={stopAnalysis}
                >
                  Stop Analysis
                </Button>
              </div>
              <div className="mt-3 w-full bg-gray-200 rounded-full h-2">
                <div
                  className="bg-green-600 h-2 rounded-full transition-all duration-200"
                  style={{ width: `${analysisProgress.progress}%` }}
                />
              </div>
            </CardContent>
          </Card>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Chess Board Section */}
          <div className="lg:col-span-2 space-y-6">
            <Card>
              <CardHeader>
                <div className="flex justify-between items-center">
                  <CardTitle>Game Board</CardTitle>
                  {currentEval && (
                    <div className="text-right">
                      <div className={`text-lg font-bold ${getScoreColor(currentEval.score)}`}>
                        {convertScoreToString(currentEval.score, currentEval.mate)}
                      </div>
                      <div className="text-sm text-gray-500">
                        Depth {currentEval.depth}
                      </div>
                    </div>
                  )}
                </div>
              </CardHeader>
              <CardContent className="flex justify-center">
                <ChessBoard
                  position={currentPosition}
                  orientation="white"
                />
              </CardContent>
            </Card>

            {/* Move Analysis */}
            {currentMoveAnalysis && (
              <Card>
                <CardHeader>
                  <CardTitle>Move Analysis</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="text-lg font-semibold">{currentMoveAnalysis.san}</div>
                        <div className="text-sm text-gray-500">
                          {currentMoveAnalysis.move}
                        </div>
                      </div>
                      <div className={`px-3 py-1 rounded-full text-sm font-medium ${getClassificationColor(currentMoveAnalysis.classification)}`}>
                        {currentMoveAnalysis.classification.charAt(0).toUpperCase() + currentMoveAnalysis.classification.slice(1)}
                      </div>
                    </div>
                    
                    {currentMoveAnalysis.alternativeMoves && currentMoveAnalysis.alternativeMoves.length > 0 && (
                      <div>
                        <div className="text-sm font-medium text-gray-700 mb-2">Best Move:</div>
                        <div className="text-sm text-gray-600">
                          {currentMoveAnalysis.alternativeMoves[0].move}
                        </div>
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>
            )}

            {/* Game Controls */}
            {gameState && (
              <GameControls
                canGoBackward={canGoBackward}
                canGoForward={canGoForward}
                isAtStart={isAtStart}
                isAtEnd={isAtEnd}
                onGoToStart={goToStart}
                onGoBackward={goBackward}
                onGoForward={goForward}
                onGoToEnd={goToEnd}
                currentMoveIndex={currentMoveIndex}
                totalMoves={gameState.moves.length}
              />
            )}

            {/* Game Input Section */}
            {!gameState && (
              <Card>
                <CardHeader>
                  <CardTitle>Import Game</CardTitle>
                  <CardDescription>
                    Paste a PGN game to start analysis
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <Textarea
                      label="Paste PGN here"
                      placeholder="Paste your PGN game notation here..."
                      value={pgnInput}
                      onChange={(e) => setPgnInput(e.target.value)}
                      rows={8}
                      className="font-mono text-sm"
                    />
                    <div className="flex space-x-3">
                      <Button
                        onClick={() => handleLoadPGN(pgnInput)}
                        isLoading={isLoading}
                        className="flex-1"
                        disabled={!engineReady}
                      >
                        {isLoading ? 'Loading...' : 'Start Review'}
                      </Button>
                      <Button
                        variant="outline"
                        onClick={() => setPgnInput('')}
                      >
                        Clear
                      </Button>
                    </div>
                    {error && (
                      <div className="text-sm text-red-600">{error}</div>
                    )}
                  </div>
                </CardContent>
              </Card>
            )}
          </div>

          {/* Analysis Panel */}
          <div className="space-y-6">
            {/* Game Summary */}
            {gameAnalysis && gameState && (
              <GameSummary
                gameAnalysis={gameAnalysis}
                gameInfo={{
                  white: gameState.gameInfo.white,
                  black: gameState.gameInfo.black,
                  whiteRating: gameState.gameInfo.whiteRating,
                  blackRating: gameState.gameInfo.blackRating,
                  result: gameState.gameInfo.result,
                  date: gameState.gameInfo.date,
                  event: gameState.gameInfo.event,
                  opening: gameState.gameInfo.opening,
                  eco: gameState.gameInfo.eco
                }}
              />
            )}

            {/* Move List */}
            {gameState && (
              <MoveList
                moves={gameState.moves}
                currentMoveIndex={currentMoveIndex}
                onMoveClick={goToMove}
              />
            )}

            {/* Export Options */}
            {gameState && (
              <ExportOptions
                gameState={gameState}
                gameAnalysis={gameAnalysis}
                currentPosition={currentPosition}
              />
            )}

            {/* Advanced Analysis */}
            {gameState && (
              <AdvancedAnalysis
                gameState={gameState}
                gameAnalysis={gameAnalysis}
                currentPosition={currentPosition}
                currentMoveIndex={currentMoveIndex}
              />
            )}
          </div>
        </div>

        {/* Performance Visualization Section */}
        {gameAnalysis && gameState && (
          <div className="mt-8 space-y-8">
            {/* Evaluation Chart */}
            <EvaluationChart
              evaluations={gameAnalysis.evaluationHistory}
              currentMoveIndex={currentMoveIndex}
              criticalMoments={gameAnalysis.criticalMoments}
              onMoveClick={goToMove}
            />

            {/* Player Stats Row */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <PlayerStats
                playerName={gameState.gameInfo.white}
                playerRating={gameState.gameInfo.whiteRating}
                statistics={gameAnalysis.whiteStats}
                color="white"
                isWinner={gameState.gameInfo.result === '1-0'}
              />
              <PlayerStats
                playerName={gameState.gameInfo.black}
                playerRating={gameState.gameInfo.blackRating}
                statistics={gameAnalysis.blackStats}
                color="black"
                isWinner={gameState.gameInfo.result === '0-1'}
              />
            </div>
        </div>
        )}
      </main>
    </div>
  );
}
