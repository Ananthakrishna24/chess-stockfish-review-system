'use client';

import React from 'react';
import ChessBoard from '@/components/chess/ChessBoard';
import PlayerInfoBar from '@/components/chess/PlayerInfoBar';
import GameControls from '@/components/chess/GameControls';
import { EvaluationBar } from '@/components/chess/EvaluationBar';
import { MoveClassificationIcon } from '@/components/chess/MoveClassificationIcon';
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
import { Play, Pause, RotateCcw, BarChart3, Star, ThumbsUp, X, AlertTriangle, Volume2, Search, SkipBack, ChevronLeft, ChevronRight, SkipForward } from 'lucide-react';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/Select";
import { cn } from '@/lib/utils';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { startReviewMode, exitReviewMode } from '@/store/reviewModeSlice';

// New Component for Analysis Panel
const AnalysisPanel = ({ 
  gameAnalysis, 
  gameState, 
  goToMove, 
  goToStart, 
  currentMoveIndex, 
  isAnalyzingGame, 
  analysisProgress, 
  stopAnalysis, 
  analyzeCompleteGame,
  goToEnd,
  goForward,
  goBackward,
  canGoBackward,
  canGoForward,
  isAtStart,
  isAtEnd,
  totalMoves
}) => {
  const dispatch = useAppDispatch();
  const { isReviewMode } = useAppSelector(state => state.reviewMode);
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
    return <MoveClassificationIcon classification={classification} size="sm" />;
  };

  const renderPlayerStats = (player, stats) => {
    if (!stats || !stats.moveCounts) {
      return (
        <div>
          <div className="flex justify-between items-center mb-2">
            <div className="flex items-center gap-2">
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
    <div className="w-[400px] h-[750px] bg-card flex flex-col border border-border rounded-lg">
      <div className="h-14 flex items-center justify-between px-3 border-b border-border">
        <div className="flex items-center gap-2">
          <h2 className="font-semibold text-lg">Game Review</h2>
          {isReviewMode && (
            <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
          )}
        </div>
        <div className="flex items-center gap-1">
          {isReviewMode ? (
            <>
              <Button 
                variant="ghost" 
                size="icon" 
                className="h-7 w-7" 
                onClick={goToStart}
                disabled={isAtStart}
                title="Go to start"
              >
                <SkipBack className="h-3 w-3" />
              </Button>
              <Button 
                variant="ghost" 
                size="icon" 
                className="h-7 w-7" 
                onClick={goBackward}
                disabled={!canGoBackward}
                title="Previous move"
              >
                <ChevronLeft className="h-3 w-3" />
              </Button>
              <span className="text-xs text-muted-foreground mx-1 min-w-[3rem] text-center">
                {currentMoveIndex === -1 ? 'Start' : `${currentMoveIndex + 1}/${totalMoves}`}
              </span>
              <Button 
                variant="ghost" 
                size="icon" 
                className="h-7 w-7" 
                onClick={goForward}
                disabled={!canGoForward}
                title="Next move"
              >
                <ChevronRight className="h-3 w-3" />
              </Button>
              <Button 
                variant="ghost" 
                size="icon" 
                className="h-7 w-7" 
                onClick={goToEnd}
                disabled={isAtEnd}
                title="Go to end"
              >
                <SkipForward className="h-3 w-3" />
              </Button>
            </>
          ) : (
            <>
              <Button variant="ghost" size="icon" className="h-8 w-8">
                <Volume2 className="h-4 w-4" />
              </Button>
              <Button variant="ghost" size="icon" className="h-8 w-8">
                <Search className="h-4 w-4" />
              </Button>
            </>
          )}
        </div>
      </div>
      {isReviewMode && (
        <div className="px-3 py-2 border-b border-border">
          <Progress 
            value={totalMoves > 0 ? ((currentMoveIndex + 1) / totalMoves) * 100 : 0}
            className="h-1"
          />
        </div>
      )}
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
        {isReviewMode ? (
          <Button 
            onClick={() => dispatch(exitReviewMode())}
            variant="outline"
            className="w-full h-10 text-base font-semibold"
          >
            Exit Review
          </Button>
        ) : (
          <Button 
            onClick={() => {
              goToStart();
              dispatch(startReviewMode());
            }}
            className="w-full h-10 text-base font-semibold bg-primary hover:bg-primary/90"
          >
            Start Review
          </Button>
        )}
      </div>
    </div>
  );
};

export default function Home() {
  const dispatch = useAppDispatch();
  const { isReviewMode } = useAppSelector(state => state.reviewMode);
  
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

  // Keyboard event handling for review mode
  React.useEffect(() => {
    if (!isReviewMode) return;

    const handleKeyDown = (event: KeyboardEvent) => {
      switch (event.key) {
        case 'ArrowLeft':
          event.preventDefault();
          if (currentMoveIndex > -1) goBackward();
          break;
        case 'ArrowRight':
          event.preventDefault();
          if (currentMoveIndex < gameState.moves.length - 1) goForward();
          break;
        case 'Home':
          event.preventDefault();
          goToStart();
          break;
        case 'End':
          event.preventDefault();
          goToEnd();
          break;
        case 'Escape':
          event.preventDefault();
          dispatch(exitReviewMode());
          break;
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [isReviewMode, currentMoveIndex, gameState?.moves.length, goBackward, goForward, goToStart, goToEnd, dispatch]);

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

  // Get current evaluation for the evaluation bar
  const getCurrentEvaluation = () => {
    if (!gameAnalysis) return 0;
    const evalHistory = gameAnalysis.evaluationHistory;
    if (!evalHistory || evalHistory.length === 0) return 0;
    
    // For starting position (currentMoveIndex === -1), use first evaluation
    if (currentMoveIndex < 0) {
      return evalHistory[0]?.score || 0;
    }
    
    // For moves, use the evaluation after the move
    const evalIndex = Math.min(currentMoveIndex + 1, evalHistory.length - 1);
    return evalHistory[evalIndex]?.score || 0;
  };

  const mainContent = gameState ? (
    <div className="flex-1 flex justify-center items-center overflow-hidden min-h-screen">
      {/* Chess Board Section with Evaluation Bar */}
      <div className="flex items-center gap-6">
        {/* Evaluation Bar */}
        {gameAnalysis && (
          <EvaluationBar 
            evaluation={getCurrentEvaluation()}
            className="h-[700px]" // Match board height
          />
        )}
        
        {/* Board and Player Info */}
        <div className="flex flex-col">
          <PlayerInfoBar
            playerName={gameState.gameInfo.white || 'Player 1'}
            playerRating={gameState.gameInfo.whiteRating}
            className="mb-4 w-[700px]"
          />
          <div className="w-[700px] h-[700px]">
            <ChessBoard
              position={currentPosition}
              orientation="white"
            />
          </div>
          <PlayerInfoBar
            playerName={gameState.gameInfo.black || 'Player 2'}
            playerRating={gameState.gameInfo.blackRating}
            className="mt-4 w-[700px]"
          />
        </div>
      </div>
      
      {/* Game Review Panel - positioned next to board with matching height */}
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
          goToEnd={goToEnd}
          goForward={goForward}
          goBackward={goBackward}
          canGoBackward={currentMoveIndex > -1}
          canGoForward={currentMoveIndex < gameState.moves.length - 1}
          isAtStart={currentMoveIndex === -1}
          isAtEnd={currentMoveIndex === gameState.moves.length - 1}
          totalMoves={gameState.moves.length}
        />
      )}
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

      </div>
    </div>
  );
}
