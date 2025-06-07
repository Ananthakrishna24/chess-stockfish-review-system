'use client';

import React from 'react';
import ChessBoard from '@/components/chess/ChessBoard';
import PlayerInfoBar from '@/components/chess/PlayerInfoBar';
import GameControls from '@/components/chess/GameControls';
import { PlayerStats } from '@/components/analysis/PlayerStats';
import { EvaluationChart } from '@/components/analysis/EvaluationChart';
import { GameSummary } from '@/components/analysis/GameSummary';
import { Card, CardContent } from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import { Badge } from '@/components/ui/Badge';
import { Progress } from '@/components/ui/Progress';
import { Textarea } from '@/components/ui/Input';
import { useGameAnalysis } from '@/hooks/useGameAnalysis';
import { convertScoreToString, getScoreColor } from '@/utils/stockfish';
import { Play, Pause, RotateCcw, BarChart3, Star, ThumbsUp, X, AlertTriangle } from 'lucide-react';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/Select";
import { cn } from '@/lib/utils';

// New Component for Analysis Panel
const AnalysisPanel = ({ gameAnalysis, gameState, goToMove, goToStart, currentMoveIndex, isAnalyzingGame, analysisProgress, stopAnalysis, analyzeCompleteGame }) => {
  if (!gameAnalysis || !gameState) {
    return (
      <div className="w-[360px] bg-card flex flex-col border-l border-border">
        <div className="h-14 flex items-center justify-between px-3 border-b border-border">
          <h2 className="font-semibold text-lg">Game Review</h2>
        </div>
                 <div className="flex-1 flex items-center justify-center p-6">
           {isAnalyzingGame ? (
             <div className="text-center space-y-3">
               <div className="w-8 h-8 border-2 border-primary border-t-transparent rounded-full animate-spin mx-auto"></div>
               <div className="text-sm text-muted-foreground">Analyzing game...</div>
               <div className="text-xs text-muted-foreground mt-2">
                 {analysisProgress ? `${Math.round(analysisProgress)}%` : 'Starting...'}
               </div>
               <Button 
                 onClick={() => {
                   console.log('Stopping analysis...');
                   stopAnalysis();
                 }}
                 variant="outline" 
                 size="sm" 
                 className="mt-3"
               >
                 Cancel Analysis
               </Button>
             </div>
           ) : (
             <div className="text-center text-muted-foreground">
               <div className="text-4xl mb-2">🔍</div>
               <div>Load a game to start analysis</div>
               {gameState && (
                 <Button 
                   onClick={() => {
                     console.log('Manual restart analysis...');
                     analyzeCompleteGame();
                   }}
                   variant="outline" 
                   size="sm" 
                   className="mt-3"
                 >
                   Restart Analysis
                 </Button>
               )}
             </div>
           )}
         </div>
      </div>
    );
  }

  const getMoveIcon = (classification) => {
    switch (classification) {
      case 'brilliant': return '!!';
      case 'great': return '!';
      case 'best': return <Star className="w-4 h-4 text-green-400" />;
      case 'good': return <ThumbsUp className="w-4 h-4 text-yellow-400" />;
      case 'inaccuracy': return '?';
      case 'mistake': return '??';
      case 'blunder': return '⁉️';
      case 'miss': return <X className="w-4 h-4 text-red-500" />;
      default: return null;
    }
  };

  const renderPlayerStats = (player, stats) => {
    if (!stats || !stats.moveCounts) {
      return (
        <div>
          <div className="flex justify-between items-center mb-2">
            <div className="flex items-center gap-2">
               <div className="w-6 h-6 rounded bg-muted flex items-center justify-center">
                 <span className="text-sm">♟️</span>
               </div>
               <span className="font-semibold">{player.name}</span>
               {player.rating && <span className="text-sm text-muted-foreground">({player.rating})</span>}
             </div>
             <Badge variant="secondary" className="font-bold text-lg">...%</Badge>
           </div>
           <div className="text-sm text-muted-foreground">Analyzing...</div>
        </div>
      )
    }

    return (
      <div>
        <div className="flex justify-between items-center mb-2">
          <div className="flex items-center gap-2">
            <div className="w-6 h-6 rounded bg-muted flex items-center justify-center">
              <span className="text-sm">♟️</span>
            </div>
            <span className="font-semibold text-foreground">{player.name}</span>
            {player.rating && <span className="text-sm text-muted-foreground">({player.rating})</span>}
          </div>
          <Badge variant="secondary" className="font-bold text-lg">{stats.accuracy}%</Badge>
        </div>
        <div className="space-y-1 text-sm">
          {Object.entries(stats.moveCounts).map(([key, value], index) => value > 0 && (
            <div key={`${key}-${index}`} className="flex justify-between items-center">
              <span className="capitalize text-muted-foreground">{key}</span>
              <div className="flex items-center gap-2">
                <span className="font-semibold text-foreground">{value}</span>
                <span className="w-6 text-center">{getMoveIcon(key)}</span>
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  };

  return (
    <div className="w-[360px] bg-card flex flex-col border-l border-border">
      <div className="h-14 flex items-center justify-between px-3 border-b border-border">
        <h2 className="font-semibold text-lg">Game Review</h2>
        <div className="flex items-center gap-2">
          {/* Add sound/search icons here if needed */}
        </div>
      </div>
      <div className="flex-1 overflow-y-auto custom-scrollbar p-3 space-y-4">
        <EvaluationChart
          evaluations={gameAnalysis.evaluationHistory}
          currentMoveIndex={currentMoveIndex}
          criticalMoments={gameAnalysis.criticalMoments}
          onMoveClick={goToMove}
        />
        {renderPlayerStats({ name: gameState.gameInfo.white, rating: gameState.gameInfo.whiteRating }, gameAnalysis.whiteStats)}
        {renderPlayerStats({ name: gameState.gameInfo.black, rating: gameState.gameInfo.blackRating }, gameAnalysis.blackStats)}

        {/* Move List */}
        <Card>
          <CardContent className="p-0">
            <div className="p-3 border-b border-border">
              <h3 className="font-semibold text-sm">Game Moves</h3>
              <div className="text-xs text-muted-foreground mt-1">
                {gameState.moves.length} moves • Click to navigate
              </div>
            </div>
            <div className="max-h-48 overflow-y-auto custom-scrollbar p-3">
              {gameState.moves.length > 0 ? (
                <div className="space-y-1">
                  {Array.from({ length: Math.ceil(gameState.moves.length / 2) }, (_, pairIndex) => {
                    const whiteMove = gameState.moves[pairIndex * 2];
                    const blackMove = gameState.moves[pairIndex * 2 + 1];
                    
                    return (
                      <div key={pairIndex} className="flex items-center gap-2 text-sm">
                        <div className="w-6 text-muted-foreground font-medium flex-shrink-0">
                          {pairIndex + 1}.
                        </div>
                                                 <button
                           onClick={() => goToMove(pairIndex * 2)}
                           className={cn(
                             "px-2 py-1 rounded hover:bg-accent transition-colors font-mono text-xs flex-1 text-left",
                             currentMoveIndex === pairIndex * 2
                               ? 'bg-primary text-primary-foreground font-semibold'
                               : 'hover:text-accent-foreground'
                           )}
                         >
                           {whiteMove?.san}
                         </button>
                         {blackMove && (
                           <button
                             onClick={() => goToMove(pairIndex * 2 + 1)}
                             className={cn(
                               "px-2 py-1 rounded hover:bg-accent transition-colors font-mono text-xs flex-1 text-left",
                               currentMoveIndex === pairIndex * 2 + 1
                                 ? 'bg-primary text-primary-foreground font-semibold'
                                 : 'hover:text-accent-foreground'
                             )}
                           >
                             {blackMove.san}
                           </button>
                         )}
                      </div>
                    );
                  })}
                </div>
              ) : (
                <div className="text-center text-muted-foreground py-4">
                  No moves available
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
      <div className="p-3 border-t border-border">
        <Button 
          onClick={() => {
            goToStart();
            // Could add review mode state here if needed
          }}
          className="w-full h-10 text-base font-semibold bg-primary hover:bg-primary/90"
        >
          Start Review
        </Button>
      </div>
    </div>
  );
};

export default function Home() {
  const {
    gameState,
    currentPosition,
    currentMoveIndex,
    isLoading,
    error,
    goToMove,
    goToStart,
    goToEnd,
    goForward,
    goBackward,
    gameAnalysis,
    isAnalyzingGame,
    analysisProgress,
    loadGame,
    resetGame,
    stopAnalysis,
    analyzeCompleteGame
  } = useGameAnalysis();
  
  const [pgnInput, setPgnInput] = React.useState('');
  const [depth, setDepth] = React.useState(4); // Default depth for quick testing

  const handleLoadPGN = async (pgn: string) => {
    if (!pgn.trim()) return;
    await loadGame(pgn, { depth });
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

  const mainContent = gameState ? (
    <div className="flex-1 flex flex-col justify-center items-center p-6">
      <div className="w-full h-full flex flex-col justify-center items-center max-w-[min(calc(100vh-10rem),calc(100vw-400px))]">
        <PlayerInfoBar
          playerName={gameState.gameInfo.white || 'Player 1'}
          playerRating={gameState.gameInfo.whiteRating}
        />
        <div className="my-6 w-full">
          <ChessBoard
            position={currentPosition}
            orientation="white"
          />
        </div>
        <PlayerInfoBar
          playerName={gameState.gameInfo.black || 'Player 2'}
          playerRating={gameState.gameInfo.blackRating}
        />
      </div>
    </div>
  ) : (
    <div className="flex-1 flex items-center justify-center p-8">
      <div className="w-full max-w-lg space-y-6">
        <div className="text-center">
          <h1 className="text-3xl font-bold">Chess Game Review</h1>
          <p className="text-muted-foreground mt-2">
            Paste a PGN to start a professional analysis.
          </p>
        </div>
        <Card>
          <CardContent className="p-6">
            <Textarea
              placeholder="Paste PGN here..."
              value={pgnInput}
              onChange={(e) => setPgnInput(e.target.value)}
              rows={10}
              className="font-mono"
            />
            <div className="mt-4 pt-4 border-t border-border">
              <label htmlFor="depth" className="text-sm font-medium text-muted-foreground">
                Engine Depth
              </label>
              <Select value={String(depth)} onValueChange={(value) => setDepth(Number(value))}>
                <SelectTrigger id="depth" className="w-full mt-2">
                  <SelectValue placeholder="Select engine depth" />
                </SelectTrigger>
                <SelectContent>
                  {[4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24].map(d => (
                    <SelectItem key={d} value={String(d)}>Depth {d}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <p className="text-xs text-muted-foreground mt-2">
                Higher depth means stronger analysis, but it will take longer.
              </p>
            </div>
            <div className="flex gap-3 mt-4">
              <Button onClick={() => handleLoadPGN(pgnInput)} isLoading={isLoading} className="flex-1">
                Start Review
              </Button>
              <Button onClick={() => handleLoadPGN(samplePGN)} variant="secondary">
                Load Sample
              </Button>
            </div>
            {gameAnalysis && (
              <div className="mt-4 pt-4 border-t border-border">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Previous analysis restored</span>
                  <Button onClick={resetGame} variant="outline" size="sm">
                    Clear & New Game
                  </Button>
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );

  return (
    <div className="h-screen bg-background text-foreground flex overflow-hidden">
      <div className="flex-1 flex flex-col">
        {/* Optional top bar can go here */}
        <main className="flex-1 flex">
          {mainContent}
        </main>
        {gameState && (
          <footer className="h-16 flex items-center justify-center border-t border-border">
            <GameControls
              onGoToStart={goToStart}
              onGoBackward={goBackward}
              onGoForward={goForward}
              onGoToEnd={goToEnd}
              canGoBackward={currentMoveIndex > -1}
              canGoForward={currentMoveIndex < gameState.moves.length -1}
              isAtStart={currentMoveIndex === -1}
              isAtEnd={currentMoveIndex === gameState.moves.length -1}
              currentMoveIndex={currentMoveIndex}
              totalMoves={gameState.moves.length}
            />
          </footer>
        )}
      </div>
      {gameAnalysis && (
        <AnalysisPanel 
          gameAnalysis={gameAnalysis}
          gameState={gameState}
          goToMove={goToMove}
          goToStart={goToStart}
          currentMoveIndex={currentMoveIndex}
          isAnalyzingGame={isAnalyzingGame}
          analysisProgress={analysisProgress}
          stopAnalysis={stopAnalysis}
          analyzeCompleteGame={analyzeCompleteGame}
        />
      )}
    </div>
  );
}
