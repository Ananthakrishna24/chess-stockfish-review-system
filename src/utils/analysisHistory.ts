// Analysis history management
import { GameAnalysis } from '@/types/analysis';
import { ChessGameState } from '@/types/chess';

export interface AnalysisHistoryEntry {
  id: string;
  timestamp: number;
  gameState: ChessGameState;
  gameAnalysis: GameAnalysis;
  metadata: {
    analysisTime: number; // Time taken to analyze in seconds
    engineDepth: number;
    totalMoves: number;
    gameResult: string;
    playerNames: string;
  };
  tags: string[];
  notes?: string;
}

export interface AnalysisHistoryFilters {
  dateRange?: {
    start: Date;
    end: Date;
  };
  playerName?: string;
  opening?: string;
  result?: '1-0' | '0-1' | '1/2-1/2' | '*';
  minAccuracy?: number;
  maxAccuracy?: number;
  tags?: string[];
  sortBy?: 'date' | 'accuracy' | 'moves' | 'analysis_time';
  sortOrder?: 'asc' | 'desc';
}

export interface AnalysisStatistics {
  totalGames: number;
  averageAccuracy: {
    white: number;
    black: number;
    overall: number;
  };
  openingFrequency: Record<string, number>;
  resultDistribution: Record<string, number>;
  gamePhaseAccuracy: {
    opening: number;
    middlegame: number;
    endgame: number;
  };
  improvementTrend: {
    period: string;
    accuracy: number;
  }[];
  commonMistakes: {
    mistake: string;
    frequency: number;
  }[];
}

class AnalysisHistoryManager {
  private storageKey = 'chess-analysis-history';
  private maxEntries = 100; // Maximum number of entries to store

  // Save analysis to history
  saveAnalysis(
    gameState: ChessGameState,
    gameAnalysis: GameAnalysis,
    analysisTime: number,
    engineDepth: number,
    tags: string[] = [],
    notes?: string
  ): string {
    const id = this.generateId();
    const timestamp = Date.now();
    
    const entry: AnalysisHistoryEntry = {
      id,
      timestamp,
      gameState,
      gameAnalysis,
      metadata: {
        analysisTime,
        engineDepth,
        totalMoves: gameState.moves.length,
        gameResult: gameState.gameInfo.result || '*',
        playerNames: `${gameState.gameInfo.white} vs ${gameState.gameInfo.black}`
      },
      tags: [...tags],
      notes
    };

    const history = this.getHistory();
    history.unshift(entry); // Add to beginning

    // Limit history size
    if (history.length > this.maxEntries) {
      history.splice(this.maxEntries);
    }

    this.saveHistory(history);
    return id;
  }

  // Get all analysis history
  getHistory(): AnalysisHistoryEntry[] {
    try {
      const stored = localStorage.getItem(this.storageKey);
      return stored ? JSON.parse(stored) : [];
    } catch (error) {
      console.error('Failed to load analysis history:', error);
      return [];
    }
  }

  // Get analysis by ID
  getAnalysisById(id: string): AnalysisHistoryEntry | null {
    const history = this.getHistory();
    return history.find(entry => entry.id === id) || null;
  }

  // Filter and search history
  filterHistory(filters: AnalysisHistoryFilters): AnalysisHistoryEntry[] {
    let history = this.getHistory();

    // Apply filters
    if (filters.dateRange) {
      const startTime = filters.dateRange.start.getTime();
      const endTime = filters.dateRange.end.getTime();
      history = history.filter(entry => 
        entry.timestamp >= startTime && entry.timestamp <= endTime
      );
    }

    if (filters.playerName) {
      const searchName = filters.playerName.toLowerCase();
      history = history.filter(entry =>
        entry.gameState.gameInfo.white?.toLowerCase().includes(searchName) ||
        entry.gameState.gameInfo.black?.toLowerCase().includes(searchName)
      );
    }

    if (filters.opening) {
      const searchOpening = filters.opening.toLowerCase();
      history = history.filter(entry =>
        entry.gameAnalysis.openingAnalysis?.name.toLowerCase().includes(searchOpening) ||
        entry.gameAnalysis.openingAnalysis?.eco.toLowerCase().includes(searchOpening)
      );
    }

    if (filters.result) {
      history = history.filter(entry => entry.metadata.gameResult === filters.result);
    }

    if (filters.minAccuracy !== undefined || filters.maxAccuracy !== undefined) {
      history = history.filter(entry => {
        const avgAccuracy = (entry.gameAnalysis.whiteStats.accuracy + entry.gameAnalysis.blackStats.accuracy) / 2;
        if (filters.minAccuracy !== undefined && avgAccuracy < filters.minAccuracy) return false;
        if (filters.maxAccuracy !== undefined && avgAccuracy > filters.maxAccuracy) return false;
        return true;
      });
    }

    if (filters.tags && filters.tags.length > 0) {
      history = history.filter(entry =>
        filters.tags!.some(tag => entry.tags.includes(tag))
      );
    }

    // Sort results
    if (filters.sortBy) {
      history.sort((a, b) => {
        let valueA: number;
        let valueB: number;

        switch (filters.sortBy) {
          case 'date':
            valueA = a.timestamp;
            valueB = b.timestamp;
            break;
          case 'accuracy':
            valueA = (a.gameAnalysis.whiteStats.accuracy + a.gameAnalysis.blackStats.accuracy) / 2;
            valueB = (b.gameAnalysis.whiteStats.accuracy + b.gameAnalysis.blackStats.accuracy) / 2;
            break;
          case 'moves':
            valueA = a.metadata.totalMoves;
            valueB = b.metadata.totalMoves;
            break;
          case 'analysis_time':
            valueA = a.metadata.analysisTime;
            valueB = b.metadata.analysisTime;
            break;
          default:
            return 0;
        }

        if (filters.sortOrder === 'desc') {
          return valueB - valueA;
        } else {
          return valueA - valueB;
        }
      });
    }

    return history;
  }

  // Get analysis statistics
  getStatistics(): AnalysisStatistics {
    const history = this.getHistory();
    
    if (history.length === 0) {
      return this.getEmptyStatistics();
    }

    // Calculate averages
    let totalWhiteAccuracy = 0;
    let totalBlackAccuracy = 0;
    const openingFrequency: Record<string, number> = {};
    const resultDistribution: Record<string, number> = {};
    let totalOpeningAccuracy = 0;
    let totalMiddlegameAccuracy = 0;
    let totalEndgameAccuracy = 0;

    for (const entry of history) {
      totalWhiteAccuracy += entry.gameAnalysis.whiteStats.accuracy;
      totalBlackAccuracy += entry.gameAnalysis.blackStats.accuracy;

      // Opening frequency
      const opening = entry.gameAnalysis.openingAnalysis?.name || 'Unknown';
      openingFrequency[opening] = (openingFrequency[opening] || 0) + 1;

      // Result distribution
      const result = entry.metadata.gameResult;
      resultDistribution[result] = (resultDistribution[result] || 0) + 1;

      // Phase accuracy
      if (entry.gameAnalysis.phaseAnalysis) {
        totalOpeningAccuracy += entry.gameAnalysis.phaseAnalysis.openingAccuracy;
        totalMiddlegameAccuracy += entry.gameAnalysis.phaseAnalysis.middlegameAccuracy;
        totalEndgameAccuracy += entry.gameAnalysis.phaseAnalysis.endgameAccuracy;
      }
    }

    const count = history.length;
    const averageWhiteAccuracy = totalWhiteAccuracy / count;
    const averageBlackAccuracy = totalBlackAccuracy / count;
    const averageOverallAccuracy = (averageWhiteAccuracy + averageBlackAccuracy) / 2;

    // Generate improvement trend (last 10 games)
    const recentGames = history.slice(0, 10).reverse(); // Most recent first, then reverse for chronological order
    const improvementTrend = recentGames.map((entry, index) => ({
      period: `Game ${index + 1}`,
      accuracy: (entry.gameAnalysis.whiteStats.accuracy + entry.gameAnalysis.blackStats.accuracy) / 2
    }));

    // Identify common mistakes
    const mistakeTypes: Record<string, number> = {};
    for (const entry of history) {
      // Count blunders and mistakes
      mistakeTypes['blunder'] = (mistakeTypes['blunder'] || 0) + 
        entry.gameAnalysis.whiteStats.blunder + entry.gameAnalysis.blackStats.blunder;
      mistakeTypes['mistake'] = (mistakeTypes['mistake'] || 0) + 
        entry.gameAnalysis.whiteStats.mistake + entry.gameAnalysis.blackStats.mistake;
      mistakeTypes['inaccuracy'] = (mistakeTypes['inaccuracy'] || 0) + 
        entry.gameAnalysis.whiteStats.inaccuracy + entry.gameAnalysis.blackStats.inaccuracy;
    }

    const commonMistakes = Object.entries(mistakeTypes)
      .map(([mistake, frequency]) => ({ mistake, frequency }))
      .sort((a, b) => b.frequency - a.frequency);

    return {
      totalGames: count,
      averageAccuracy: {
        white: averageWhiteAccuracy,
        black: averageBlackAccuracy,
        overall: averageOverallAccuracy
      },
      openingFrequency,
      resultDistribution,
      gamePhaseAccuracy: {
        opening: totalOpeningAccuracy / count,
        middlegame: totalMiddlegameAccuracy / count,
        endgame: totalEndgameAccuracy / count
      },
      improvementTrend,
      commonMistakes
    };
  }

  // Delete analysis entry
  deleteAnalysis(id: string): boolean {
    const history = this.getHistory();
    const index = history.findIndex(entry => entry.id === id);
    
    if (index === -1) return false;
    
    history.splice(index, 1);
    this.saveHistory(history);
    return true;
  }

  // Update analysis entry
  updateAnalysis(id: string, updates: Partial<Pick<AnalysisHistoryEntry, 'tags' | 'notes'>>): boolean {
    const history = this.getHistory();
    const entry = history.find(e => e.id === id);
    
    if (!entry) return false;
    
    if (updates.tags !== undefined) entry.tags = updates.tags;
    if (updates.notes !== undefined) entry.notes = updates.notes;
    
    this.saveHistory(history);
    return true;
  }

  // Export history
  exportHistory(): string {
    const history = this.getHistory();
    return JSON.stringify(history, null, 2);
  }

  // Import history
  importHistory(data: string): boolean {
    try {
      const imported = JSON.parse(data) as AnalysisHistoryEntry[];
      const existing = this.getHistory();
      
      // Merge with existing, avoiding duplicates
      const merged = [...existing];
      for (const entry of imported) {
        if (!merged.find(e => e.id === entry.id)) {
          merged.push(entry);
        }
      }
      
      // Sort by timestamp and limit
      merged.sort((a, b) => b.timestamp - a.timestamp);
      if (merged.length > this.maxEntries) {
        merged.splice(this.maxEntries);
      }
      
      this.saveHistory(merged);
      return true;
    } catch (error) {
      console.error('Failed to import history:', error);
      return false;
    }
  }

  // Clear all history
  clearHistory(): void {
    localStorage.removeItem(this.storageKey);
  }

  // Get available tags
  getAvailableTags(): string[] {
    const history = this.getHistory();
    const tags = new Set<string>();
    
    for (const entry of history) {
      for (const tag of entry.tags) {
        tags.add(tag);
      }
    }
    
    return Array.from(tags).sort();
  }

  // Private methods
  private generateId(): string {
    return `analysis_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }

  private saveHistory(history: AnalysisHistoryEntry[]): void {
    try {
      localStorage.setItem(this.storageKey, JSON.stringify(history));
    } catch (error) {
      console.error('Failed to save analysis history:', error);
    }
  }

  private getEmptyStatistics(): AnalysisStatistics {
    return {
      totalGames: 0,
      averageAccuracy: { white: 0, black: 0, overall: 0 },
      openingFrequency: {},
      resultDistribution: {},
      gamePhaseAccuracy: { opening: 0, middlegame: 0, endgame: 0 },
      improvementTrend: [],
      commonMistakes: []
    };
  }
}

// Create singleton instance
export const analysisHistory = new AnalysisHistoryManager();

// Export types and manager
export { AnalysisHistoryManager }; 