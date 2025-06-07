'use client';

import { useState, useCallback, useEffect, useRef } from 'react';
import { EngineEvaluation, AnalysisProgress, StockfishConfig } from '@/types/analysis';
import { StockfishEngine } from '@/utils/stockfish';

export function useStockfish(initialConfig?: Partial<StockfishConfig>) {
  const [engine, setEngine] = useState<StockfishEngine | null>(null);
  const [isReady, setIsReady] = useState(false);
  const [isInitializing, setIsInitializing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentEvaluation, setCurrentEvaluation] = useState<EngineEvaluation | null>(null);
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [analysisProgress, setAnalysisProgress] = useState<AnalysisProgress>({
    currentMove: 0,
    totalMoves: 0,
    isAnalyzing: false,
    progress: 0
  });

  const engineRef = useRef<StockfishEngine | null>(null);
  const configRef = useRef(initialConfig);
  const abortControllerRef = useRef<AbortController | null>(null);

  const updateConfig = useCallback((newConfig: Partial<StockfishConfig>) => {
    configRef.current = { ...configRef.current, ...newConfig };
    if (engineRef.current) {
      engineRef.current.setConfig(newConfig);
    }
  }, []);

  const initializeEngine = useCallback(async () => {
    if (engineRef.current || isInitializing) return;
    
    setIsInitializing(true);
    setError(null);
    
    try {
      const newEngine = new StockfishEngine(configRef.current);
      await newEngine.initialize();
      
      engineRef.current = newEngine;
      setEngine(newEngine);
      setIsReady(true);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to initialize Stockfish';
      setError(errorMessage);
      console.error('Stockfish initialization error:', err);
    } finally {
      setIsInitializing(false);
    }
  }, [isInitializing]);

  const analyzePosition = useCallback(async (
    fen: string, 
    depth?: number
  ): Promise<EngineEvaluation | null> => {
    if (!engineRef.current || !isReady) {
      await initializeEngine();
      if (!engineRef.current) return null;
    }

    setIsAnalyzing(true);
    setError(null);

    try {
      const evaluation = await engineRef.current.analyzePosition(fen, depth);
      setCurrentEvaluation(evaluation);
      return evaluation;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Analysis failed';
      setError(errorMessage);
      console.error('Analysis error:', err);
      return null;
    } finally {
      setIsAnalyzing(false);
    }
  }, [isReady, initializeEngine]);

  const analyzeGame = useCallback(async (
    positions: string[],
    onProgress?: (progress: AnalysisProgress) => void
  ): Promise<EngineEvaluation[]> => {
    if (!engineRef.current || !isReady) {
      await initializeEngine();
      if (!engineRef.current) return [];
    }

    setError(null);
    setAnalysisProgress({
      currentMove: 0,
      totalMoves: positions.length,
      isAnalyzing: true,
      progress: 0
    });

    const evaluations: EngineEvaluation[] = [];
    
    // Create abort controller for this analysis session
    abortControllerRef.current = new AbortController();
    
    try {
      for (let i = 0; i < positions.length; i++) {
        // Check if analysis was aborted
        if (abortControllerRef.current?.signal.aborted) {
          break;
        }

        const position = positions[i];
        const evaluation = await engineRef.current.analyzePosition(position);
        evaluations.push(evaluation);

        const progress = {
          currentMove: i + 1,
          totalMoves: positions.length,
          isAnalyzing: true,
          progress: ((i + 1) / positions.length) * 100
        };

        setAnalysisProgress(progress);
        onProgress?.(progress);

        // Small delay to prevent UI blocking
        await new Promise(resolve => setTimeout(resolve, 10));
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Game analysis failed';
      setError(errorMessage);
      console.error('Game analysis error:', err);
    } finally {
      setAnalysisProgress(prev => ({
        ...prev,
        isAnalyzing: false
      }));
      abortControllerRef.current = null;
    }

    return evaluations;
  }, [isReady, initializeEngine]);

  const stopAnalysis = useCallback(() => {
    if (engineRef.current) {
      engineRef.current.stop();
    }
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
    setIsAnalyzing(false);
    setAnalysisProgress(prev => ({
      ...prev,
      isAnalyzing: false
    }));
  }, []);

  const getBestMove = useCallback(async (fen: string): Promise<string | null> => {
    if (!engineRef.current || !isReady) {
      await initializeEngine();
      if (!engineRef.current) return null;
    }

    try {
      return await engineRef.current.findBestMove(fen);
    } catch (err) {
      console.error('Best move analysis error:', err);
      return null;
    }
  }, [isReady, initializeEngine]);

  const classifyMove = useCallback((
    positionBefore: EngineEvaluation,
    positionAfter: EngineEvaluation,
    playedMove: string,
    bestMove: string
  ) => {
    if (!engineRef.current) return 'good';
    
    return engineRef.current.classifyMove(
      positionBefore,
      positionAfter,
      playedMove,
      bestMove
    );
  }, []);

  const calculateAccuracy = useCallback((evaluations: EngineEvaluation[]): number => {
    if (!engineRef.current) return 0;
    return engineRef.current.calculateAccuracy(evaluations);
  }, []);

  // Auto-initialize on mount
  useEffect(() => {
    initializeEngine();
  }, [initializeEngine]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (engineRef.current) {
        engineRef.current.quit();
      }
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, []);

  return {
    // State
    isReady,
    isInitializing,
    isAnalyzing,
    error,
    currentEvaluation,
    analysisProgress,
    
    // Actions
    initializeEngine,
    analyzePosition,
    analyzeGame,
    stopAnalysis,
    getBestMove,
    classifyMove,
    calculateAccuracy,
    updateConfig,
    
    // Engine reference (for advanced usage)
    engine: engineRef.current
  };
} 