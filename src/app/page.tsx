'use client';

import React from 'react';
import ChessBoard from '@/components/chess/ChessBoard';
import MoveList from '@/components/chess/MoveList';
import GameControls from '@/components/chess/GameControls';
import { PlayerStats } from '@/components/analysis/PlayerStats';
import { EvaluationChart } from '@/components/analysis/EvaluationChart';
import { GameSummary } from '@/components/analysis/GameSummary';
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

  // Get move classification color and styling
  const getClassificationStyle = (classification: string) => {
    switch (classification.toLowerCase()) {
      case 'brilliant': 
        return { 
          bg: 'bg-gradient-to-r from-cyan-50 to-blue-50', 
          text: 'text-cyan-700', 
          border: 'border-cyan-200',
          icon: '✨' 
        };
      case 'great': 
        return { 
          bg: 'bg-gradient-to-r from-blue-50 to-indigo-50', 
          text: 'text-blue-700', 
          border: 'border-blue-200',
          icon: '👍' 
        };
      case 'best': 
        return { 
          bg: 'bg-gradient-to-r from-green-50 to-emerald-50', 
          text: 'text-green-700', 
          border: 'border-green-200',
          icon: '✓' 
        };
      case 'good': 
        return { 
          bg: 'bg-gradient-to-r from-green-50 to-lime-50', 
          text: 'text-green-600', 
          border: 'border-green-200',
          icon: '○' 
        };
      case 'inaccuracy': 
        return { 
          bg: 'bg-gradient-to-r from-yellow-50 to-amber-50', 
          text: 'text-yellow-700', 
          border: 'border-yellow-200',
          icon: '?!' 
        };
      case 'mistake': 
        return { 
          bg: 'bg-gradient-to-r from-orange-50 to-red-50', 
          text: 'text-orange-700', 
          border: 'border-orange-200',
          icon: '?' 
        };
      case 'blunder': 
        return { 
          bg: 'bg-gradient-to-r from-red-50 to-pink-50', 
          text: 'text-red-700', 
          border: 'border-red-200',
          icon: '??' 
        };
      case 'miss': 
        return { 
          bg: 'bg-gradient-to-r from-red-50 to-rose-50', 
          text: 'text-red-600', 
          border: 'border-red-200',
          icon: '!!' 
        };
      default: 
        return { 
          bg: 'bg-gray-50', 
          text: 'text-gray-600', 
          border: 'border-gray-200',
          icon: '○' 
        };
    }
  };

  const renderEngineStatus = () => {
    if (engineInitializing) {
      return (
        <div className="flex items-center gap-2 px-3 py-1.5 bg-blue-50 border border-blue-200 rounded-lg">
          <div className="w-2 h-2 bg-blue-500 rounded-full animate-pulse"></div>
          <span className="text-sm font-medium text-blue-700">Initializing engine...</span>
        </div>
      );
    }
    
    if (engineError) {
      return (
        <div className="flex items-center gap-2 px-3 py-1.5 bg-red-50 border border-red-200 rounded-lg">
          <div className="w-2 h-2 bg-red-500 rounded-full"></div>
          <span className="text-sm font-medium text-red-700">Engine error</span>
        </div>
      );
    }
    
    if (engineReady) {
      return (
        <div className="flex items-center gap-2 px-3 py-1.5 bg-green-50 border border-green-200 rounded-lg">
          <div className="w-2 h-2 bg-green-500 rounded-full"></div>
          <span className="text-sm font-medium text-green-700">Engine ready</span>
        </div>
      );
    }
    
    return null;
  };

  return (
    <div className="min-h-screen" style={{ backgroundColor: 'var(--background)' }}>
      {/* Professional Header */}
      <header className="sticky top-0 z-50 backdrop-blur-sm border-b" style={{ 
        backgroundColor: 'var(--surface)', 
        borderColor: 'var(--border-light)',
        boxShadow: 'var(--shadow-sm)'
      }}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center gap-6">
              <div className="flex items-center gap-3">
                <div className="text-2xl">♔</div>
                <div>
                  <h1 className="text-xl font-semibold" style={{ color: 'var(--text-primary)' }}>
                    Chess Game Review
                  </h1>
                  <p className="text-xs" style={{ color: 'var(--text-muted)' }}>
                    Professional analysis powered by Stockfish
                  </p>
                </div>
              </div>
              {renderEngineStatus()}
            </div>
            
            <div className="flex items-center gap-3">
              <Button
                variant="ghost"
                size="sm"
                onClick={() => {
                  setPgnInput(samplePGN);
                  handleLoadPGN(samplePGN);
                }}
                leftIcon={<span className="text-sm">🎮</span>}
              >
                Load Sample
              </Button>
              {gameState && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={resetGame}
                  leftIcon={<span className="text-sm">🔄</span>}
                >
                  Reset
                </Button>
              )}
            </div>
          </div>
        </div>
      </header>

      {/* Analysis Progress Bar */}
      {isAnalyzingGame && (
        <div className="border-b" style={{ 
          backgroundColor: 'var(--surface)', 
          borderColor: 'var(--border-light)' 
        }}>
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
            <Card variant="minimal" className="border-0">
              <CardContent size="sm">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
                    <div>
                      <div className="text-sm font-medium" style={{ color: 'var(--text-primary)' }}>
                        Analyzing Game - {analysisProgress.progress.toFixed(1)}%
                      </div>
                      <div className="text-xs" style={{ color: 'var(--text-secondary)' }}>
                        Move {analysisProgress.currentMove} of {analysisProgress.totalMoves}
                      </div>
                    </div>
                  </div>
                  <Button
                    variant="outline"
                    size="xs"
                    onClick={stopAnalysis}
                  >
                    Stop
                  </Button>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className="h-2 rounded-full transition-all duration-300"
                    style={{
                      width: `${analysisProgress.progress}%`,
                      backgroundColor: 'var(--chess-success)'
                    }}
                  />
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      )}

      {/* Main Content Layout */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-1 xl:grid-cols-4 gap-8">
          
          {/* Chess Board Section - Takes 2/3 of the width */}
          <div className="xl:col-span-3 space-y-6">
            
            {/* Chess Board Card */}
            <Card variant="elevated">
              <CardHeader noBorder>
                <div className="flex justify-between items-start">
                  <div>
                    <CardTitle size="lg">Game Board</CardTitle>
                    {gameState && (
                      <CardDescription>
                        {gameState.gameInfo.white} vs {gameState.gameInfo.black}
                        {gameState.gameInfo.event && ` • ${gameState.gameInfo.event}`}
                      </CardDescription>
                    )}
                  </div>
                  {currentEval && (
                    <div className="text-right">
                      <div 
                        className={`text-2xl font-bold mb-1 ${getScoreColor(currentEval.score)}`}
                      >
                        {convertScoreToString(currentEval.score, currentEval.mate)}
                      </div>
                      <div className="text-xs" style={{ color: 'var(--text-muted)' }}>
                        Depth {currentEval.depth}
                      </div>
                    </div>
                  )}
                </div>
              </CardHeader>
              <CardContent className="flex justify-center bg-gradient-to-br from-slate-50 to-slate-100">
                <div className="p-6">
                  <ChessBoard
                    position={currentPosition}
                    orientation="white"
                  />
                </div>
              </CardContent>
            </Card>

            {/* Move Analysis Card */}
            {currentMoveAnalysis && (
              <Card variant="elevated">
                <CardHeader>
                  <CardTitle>Move Analysis</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-4">
                        <div className="text-2xl font-bold font-mono" style={{ color: 'var(--text-primary)' }}>
                          {currentMoveAnalysis.san}
                        </div>
                        <div className="text-sm" style={{ color: 'var(--text-secondary)' }}>
                          {currentMoveAnalysis.move}
                        </div>
                      </div>
                      
                      {(() => {
                        const style = getClassificationStyle(currentMoveAnalysis.classification);
                        return (
                          <div className={`flex items-center gap-2 px-4 py-2 rounded-xl border ${style.bg} ${style.text} ${style.border}`}>
                            <span className="text-lg">{style.icon}</span>
                            <span className="font-medium capitalize">
                              {currentMoveAnalysis.classification}
                            </span>
                          </div>
                        );
                      })()}
                    </div>
                    
                    {currentMoveAnalysis.alternativeMoves && currentMoveAnalysis.alternativeMoves.length > 0 && (
                      <div className="pt-4 border-t" style={{ borderColor: 'var(--border-light)' }}>
                        <div className="text-sm font-medium mb-2" style={{ color: 'var(--text-primary)' }}>
                          Engine Recommendation:
                        </div>
                        <div className="font-mono text-base font-semibold" style={{ color: 'var(--chess-accent)' }}>
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
              <Card variant="elevated">
                <CardHeader>
                  <CardTitle size="lg">Import Your Game</CardTitle>
                  <CardDescription className="text-balance">
                    Paste your PGN game notation below to start professional analysis with Stockfish engine
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <Textarea
                      label="Game PGN"
                      placeholder="[Event &quot;Your Game&quot;]
[Site &quot;lichess.org&quot;]
[Date &quot;2024.01.15&quot;]
[White &quot;You&quot;]
[Black &quot;Opponent&quot;]
...

1. e4 e5 2. Nf3 Nc6..."
                      value={pgnInput}
                      onChange={(e) => setPgnInput(e.target.value)}
                      rows={10}
                      className="font-mono text-sm"
                    />
                    <div className="flex gap-3">
                      <Button
                        onClick={() => handleLoadPGN(pgnInput)}
                        isLoading={isLoading}
                        fullWidth
                        disabled={!engineReady || !pgnInput.trim()}
                        leftIcon={!isLoading && <span className="text-lg">🚀</span>}
                      >
                        {isLoading ? 'Loading Game...' : 'Start Analysis'}
                      </Button>
                      <Button
                        variant="outline"
                        onClick={() => setPgnInput('')}
                        disabled={!pgnInput.trim()}
                      >
                        Clear
                      </Button>
                    </div>
                    {error && (
                      <div className="p-3 bg-red-50 border border-red-200 rounded-lg">
                        <div className="text-sm text-red-700 font-medium">Error</div>
                        <div className="text-sm text-red-600 mt-1">{error}</div>
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>
            )}
          </div>

          {/* Right Sidebar - Analysis Panel */}
          <div className="xl:col-span-1 space-y-6">
            {/* Game Summary */}
            {gameAnalysis && gameState && (
              <GameSummary
                gameAnalysis={gameAnalysis}
                gameInfo={{
                  white: gameState.gameInfo.white,
                  black: gameState.gameInfo.black,
                  whiteRating: gameState.gameInfo.whiteRating?.toString(),
                  blackRating: gameState.gameInfo.blackRating?.toString(),
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
          </div>
        </div>

        {/* Performance Visualization Section - Full Width */}
        {gameAnalysis && gameState && (
          <div className="mt-12 space-y-8">
            {/* Evaluation Chart */}
            <Card variant="elevated">
              <CardHeader>
                <CardTitle size="lg">Game Analysis</CardTitle>
                <CardDescription>
                  Evaluation timeline showing position advantages and critical moments
                </CardDescription>
              </CardHeader>
              <CardContent size="lg">
                <EvaluationChart
                  evaluations={gameAnalysis.evaluationHistory}
                  currentMoveIndex={currentMoveIndex}
                  criticalMoments={gameAnalysis.criticalMoments}
                  onMoveClick={goToMove}
                />
              </CardContent>
            </Card>

            {/* Player Performance Comparison */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
              <PlayerStats
                playerName={gameState.gameInfo.white}
                playerRating={gameState.gameInfo.whiteRating?.toString()}
                statistics={gameAnalysis.whiteStats}
                color="white"
                isWinner={gameState.gameInfo.result === '1-0'}
              />
              <PlayerStats
                playerName={gameState.gameInfo.black}
                playerRating={gameState.gameInfo.blackRating?.toString()}
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
