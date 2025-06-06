'use client';

import React, { useState } from 'react';
import { Card, CardHeader, CardContent, CardTitle } from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import { GameAnalysis } from '@/types/analysis';
import { GameState } from '@/types/chess';
import { exportToPGN, exportToPNG, ExportOptions as ExportOptionsType, downloadFile, generateFilename } from '@/utils/exportUtils';

interface ExportOptionsProps {
  gameState: GameState;
  gameAnalysis?: GameAnalysis;
  currentPosition: string;
  className?: string;
}

export function ExportOptions({ gameState, gameAnalysis, currentPosition, className = '' }: ExportOptionsProps) {
  const [isExporting, setIsExporting] = useState(false);
  const [exportOptions, setExportOptions] = useState<ExportOptionsType>({
    includeAnalysis: true,
    includeComments: true,
    includeStatistics: true,
    format: 'annotated'
  });

  const handleExportPGN = async () => {
    setIsExporting(true);
    try {
      const pgn = exportToPGN(gameState, gameAnalysis, exportOptions);
      const filename = generateFilename(gameState.gameInfo, 'pgn');
      downloadFile(pgn, filename);
    } catch (error) {
      console.error('Error exporting PGN:', error);
      alert('Failed to export PGN');
    } finally {
      setIsExporting(false);
    }
  };

  const handleExportPNG = async () => {
    setIsExporting(true);
    try {
      const pngBlob = await exportToPNG(currentPosition, {
        size: 400,
        coordinates: true,
        orientation: 'white'
      });
      const filename = generateFilename(gameState.gameInfo, 'png');
      downloadFile(pngBlob, filename);
    } catch (error) {
      console.error('Error exporting PNG:', error);
      alert('Failed to export PNG');
    } finally {
      setIsExporting(false);
    }
  };

  return (
    <Card className={className}>
      <CardHeader>
        <CardTitle>Export Options</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {/* Export Format Options */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              PGN Format
            </label>
            <select
              value={exportOptions.format}
              onChange={(e) => setExportOptions({
                ...exportOptions,
                format: e.target.value as ExportOptionsType['format']
              })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
            >
              <option value="standard">Standard PGN</option>
              <option value="annotated">Annotated with Analysis</option>
              <option value="analysis_only">Analysis Comments Only</option>
            </select>
          </div>

          {/* Include Options */}
          <div className="space-y-2">
            <label className="flex items-center">
              <input
                type="checkbox"
                checked={exportOptions.includeAnalysis}
                onChange={(e) => setExportOptions({
                  ...exportOptions,
                  includeAnalysis: e.target.checked
                })}
                className="mr-2"
              />
              <span className="text-sm">Include Analysis</span>
            </label>
            
            <label className="flex items-center">
              <input
                type="checkbox"
                checked={exportOptions.includeComments}
                onChange={(e) => setExportOptions({
                  ...exportOptions,
                  includeComments: e.target.checked
                })}
                className="mr-2"
              />
              <span className="text-sm">Include Comments</span>
            </label>
            
            <label className="flex items-center">
              <input
                type="checkbox"
                checked={exportOptions.includeStatistics}
                onChange={(e) => setExportOptions({
                  ...exportOptions,
                  includeStatistics: e.target.checked
                })}
                className="mr-2"
              />
              <span className="text-sm">Include Statistics</span>
            </label>
          </div>

          {/* Export Buttons */}
          <div className="flex space-x-2">
            <Button
              onClick={handleExportPGN}
              variant="outline"
              size="sm"
              isLoading={isExporting}
              className="flex-1"
            >
              Export PGN
            </Button>
            
            <Button
              onClick={handleExportPNG}
              variant="outline"
              size="sm"
              isLoading={isExporting}
              className="flex-1"
            >
              Export PNG
            </Button>
          </div>

          {/* Statistics Display */}
          {gameAnalysis && (
            <div className="pt-2 border-t border-gray-200">
              <div className="text-xs text-gray-500 space-y-1">
                <div>Game Analysis Ready</div>
                <div>White Accuracy: {gameAnalysis.whiteStats.accuracy.toFixed(1)}%</div>
                <div>Black Accuracy: {gameAnalysis.blackStats.accuracy.toFixed(1)}%</div>
                <div>Total Moves: {gameAnalysis.moves.length}</div>
              </div>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
} 