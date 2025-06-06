// Export utilities for PGN and PNG
import { GameAnalysis } from '@/types/analysis';
import { GameState } from '@/types/chess';

export interface ExportOptions {
  includeAnalysis: boolean;
  includeComments: boolean;
  includeStatistics: boolean;
  format: 'standard' | 'annotated' | 'analysis_only';
}

// Export game as annotated PGN with analysis
export function exportToPGN(
  gameState: GameState,
  gameAnalysis?: GameAnalysis,
  options: ExportOptions = {
    includeAnalysis: true,
    includeComments: true,
    includeStatistics: true,
    format: 'annotated'
  }
): string {
  const { gameInfo, moves } = gameState;
  
  // Build PGN header
  let pgn = '';
  
  // Standard PGN headers
  pgn += `[Event "${gameInfo.event || 'Game Analysis'}"]\\n`;
  pgn += `[Site "${gameInfo.site || 'Chess Review System'}"]\\n`;
  pgn += `[Date "${gameInfo.date || new Date().toISOString().split('T')[0]}"]\\n`;
  pgn += `[Round "${gameInfo.round || '1'}"]\\n`;
  pgn += `[White "${gameInfo.white || 'White'}"]\\n`;
  pgn += `[Black "${gameInfo.black || 'Black'}"]\\n`;
  pgn += `[Result "${gameInfo.result || '*'}"]\\n`;
  
  // Additional headers
  if (gameInfo.whiteRating) {
    pgn += `[WhiteElo "${gameInfo.whiteRating}"]\\n`;
  }
  if (gameInfo.blackRating) {
    pgn += `[BlackElo "${gameInfo.blackRating}"]\\n`;
  }
  if (gameInfo.eco) {
    pgn += `[ECO "${gameInfo.eco}"]\\n`;
  }
  if (gameInfo.opening) {
    pgn += `[Opening "${gameInfo.opening}"]\\n`;
  }
  
  // Analysis-specific headers
  if (gameAnalysis && options.includeStatistics) {
    pgn += `[WhiteAccuracy "${gameAnalysis.whiteStats.accuracy.toFixed(1)}%"]\\n`;
    pgn += `[BlackAccuracy "${gameAnalysis.blackStats.accuracy.toFixed(1)}%"]\\n`;
    pgn += `[AnalysisEngine "Chess Review System v1.0"]\\n`;
    pgn += `[AnalysisDate "${new Date().toISOString().split('T')[0]}"]\\n`;
  }
  
  pgn += '\\n';
  
  // Game moves
  const moveLines: string[] = [];
  let currentLine = '';
  
  for (let i = 0; i < moves.length; i++) {
    const move = moves[i];
    const moveNumber = Math.floor(i / 2) + 1;
    const isWhiteMove = i % 2 === 0;
    
    // Add move number for white moves
    if (isWhiteMove) {
      currentLine += `${moveNumber}. `;
    } else if (i === 0) {
      currentLine += `${moveNumber}... `;
    }
    
    // Add the move
    currentLine += move.san;
    
    // Add analysis annotations if requested
    if (gameAnalysis && options.includeAnalysis && i < gameAnalysis.moves.length) {
      const moveAnalysis = gameAnalysis.moves[i];
      
      // Add move classification symbols
      switch (moveAnalysis.classification) {
        case 'brilliant':
          currentLine += '!!';
          break;
        case 'great':
          currentLine += '!';
          break;
        case 'best':
          currentLine += '';
          break;
        case 'good':
          currentLine += '';
          break;
        case 'inaccuracy':
          currentLine += '?!';
          break;
        case 'mistake':
          currentLine += '?';
          break;
        case 'blunder':
          currentLine += '??';
          break;
        case 'miss':
          currentLine += '??';
          break;
      }
      
      // Add evaluation comment
      if (options.includeComments) {
        const evaluation = moveAnalysis.evaluation;
        let evalComment = '';
        
        if (evaluation.mate) {
          evalComment = `M${evaluation.mate}`;
        } else {
          const score = evaluation.score / 100;
          evalComment = score >= 0 ? `+${score.toFixed(2)}` : score.toFixed(2);
        }
        
        currentLine += ` {${evalComment}}`;
        
        // Add tactical pattern comment
        if (moveAnalysis.tacticalAnalysis?.isTactical && moveAnalysis.tacticalAnalysis.patterns.length > 0) {
          const patterns = moveAnalysis.tacticalAnalysis.patterns.filter(p => p !== 'none');
          if (patterns.length > 0) {
            currentLine += ` {${patterns.join(', ')}}`;
          }
        }
      }
    }
    
    currentLine += ' ';
    
    // Break line if too long
    if (currentLine.length > 70) {
      moveLines.push(currentLine.trim());
      currentLine = '';
    }
  }
  
  // Add remaining moves
  if (currentLine.trim()) {
    moveLines.push(currentLine.trim());
  }
  
  // Add result
  const lastLine = moveLines[moveLines.length - 1] || '';
  moveLines[moveLines.length - 1] = `${lastLine} ${gameInfo.result || '*'}`;
  
  pgn += moveLines.join('\\n') + '\\n';
  
  // Add analysis summary if requested
  if (gameAnalysis && options.includeStatistics && options.format === 'annotated') {
    pgn += '\\n';
    pgn += '{Game Analysis Summary:\\n';
    pgn += `White Accuracy: ${gameAnalysis.whiteStats.accuracy.toFixed(1)}%\\n`;
    pgn += `Black Accuracy: ${gameAnalysis.blackStats.accuracy.toFixed(1)}%\\n`;
    pgn += `Critical Moments: ${gameAnalysis.criticalMoments?.length || 0}\\n`;
    
    if (gameAnalysis.openingAnalysis) {
      pgn += `Opening: ${gameAnalysis.openingAnalysis.name}`;
      if (gameAnalysis.openingAnalysis.eco) {
        pgn += ` (${gameAnalysis.openingAnalysis.eco})`;
      }
      pgn += '\\n';
    }
    
    pgn += '}\\n';
  }
  
  return pgn;
}

// Export board position as PNG image
export async function exportToPNG(
  position: string,
  options: {
    size?: number;
    coordinates?: boolean;
    orientation?: 'white' | 'black';
    highlightSquares?: string[];
    arrows?: Array<{ from: string; to: string; color?: string }>;
  } = {}
): Promise<Blob> {
  const {
    size = 400,
    coordinates = true,
    orientation = 'white',
    highlightSquares = [],
    arrows = []
  } = options;
  
  // Create canvas
  const canvas = document.createElement('canvas');
  const ctx = canvas.getContext('2d');
  
  if (!ctx) {
    throw new Error('Could not create canvas context');
  }
  
  canvas.width = size;
  canvas.height = size;
  
  const squareSize = size / 8;
  
  // Chess piece Unicode symbols
  const pieceSymbols: Record<string, string> = {
    'K': '♔', 'Q': '♕', 'R': '♖', 'B': '♗', 'N': '♘', 'P': '♙',
    'k': '♚', 'q': '♛', 'r': '♜', 'b': '♝', 'n': '♞', 'p': '♟'
  };
  
  // Parse FEN position
  const board = parseFENPosition(position);
  
  // Draw board squares
  for (let rank = 0; rank < 8; rank++) {
    for (let file = 0; file < 8; file++) {
      const isLight = (rank + file) % 2 === 0;
      const x = orientation === 'white' ? file * squareSize : (7 - file) * squareSize;
      const y = orientation === 'white' ? (7 - rank) * squareSize : rank * squareSize;
      
      // Square color
      ctx.fillStyle = isLight ? '#F0D9B5' : '#B88762';
      ctx.fillRect(x, y, squareSize, squareSize);
      
      // Highlight squares if specified
      const square = String.fromCharCode(97 + file) + (rank + 1);
      if (highlightSquares.includes(square)) {
        ctx.fillStyle = 'rgba(255, 255, 0, 0.4)';
        ctx.fillRect(x, y, squareSize, squareSize);
      }
    }
  }
  
  // Draw coordinates if requested
  if (coordinates) {
    ctx.fillStyle = '#333';
    ctx.font = `${Math.floor(squareSize * 0.15)}px Arial`;
    
    // File labels (a-h)
    for (let file = 0; file < 8; file++) {
      const letter = String.fromCharCode(97 + file);
      const x = orientation === 'white' ? file * squareSize + 5 : (7 - file) * squareSize + 5;
      const y = size - 5;
      ctx.fillText(letter, x, y);
    }
    
    // Rank labels (1-8)
    for (let rank = 0; rank < 8; rank++) {
      const number = (rank + 1).toString();
      const x = 5;
      const y = orientation === 'white' ? (7 - rank) * squareSize + 15 : rank * squareSize + 15;
      ctx.fillText(number, x, y);
    }
  }
  
  // Draw pieces
  ctx.font = `${Math.floor(squareSize * 0.7)}px serif`;
  ctx.textAlign = 'center';
  ctx.textBaseline = 'middle';
  
  for (let rank = 0; rank < 8; rank++) {
    for (let file = 0; file < 8; file++) {
      const piece = board[rank][file];
      if (piece !== '') {
        const x = orientation === 'white' ? file * squareSize : (7 - file) * squareSize;
        const y = orientation === 'white' ? (7 - rank) * squareSize : rank * squareSize;
        
        ctx.fillStyle = piece === piece.toUpperCase() ? '#fff' : '#000';
        ctx.strokeStyle = piece === piece.toUpperCase() ? '#000' : '#fff';
        ctx.lineWidth = 1;
        
        const centerX = x + squareSize / 2;
        const centerY = y + squareSize / 2;
        
        ctx.fillText(pieceSymbols[piece], centerX, centerY);
        ctx.strokeText(pieceSymbols[piece], centerX, centerY);
      }
    }
  }
  
  // Draw arrows if specified
  if (arrows.length > 0) {
    ctx.strokeStyle = '#ff0000';
    ctx.lineWidth = 3;
    ctx.lineCap = 'round';
    
    for (const arrow of arrows) {
      const fromFile = arrow.from.charCodeAt(0) - 97;
      const fromRank = parseInt(arrow.from[1]) - 1;
      const toFile = arrow.to.charCodeAt(0) - 97;
      const toRank = parseInt(arrow.to[1]) - 1;
      
      const fromX = orientation === 'white' ? fromFile * squareSize + squareSize / 2 : (7 - fromFile) * squareSize + squareSize / 2;
      const fromY = orientation === 'white' ? (7 - fromRank) * squareSize + squareSize / 2 : fromRank * squareSize + squareSize / 2;
      const toX = orientation === 'white' ? toFile * squareSize + squareSize / 2 : (7 - toFile) * squareSize + squareSize / 2;
      const toY = orientation === 'white' ? (7 - toRank) * squareSize + squareSize / 2 : toRank * squareSize + squareSize / 2;
      
      ctx.beginPath();
      ctx.moveTo(fromX, fromY);
      ctx.lineTo(toX, toY);
      ctx.stroke();
      
      // Draw arrowhead
      const angle = Math.atan2(toY - fromY, toX - fromX);
      const headSize = 10;
      
      ctx.beginPath();
      ctx.moveTo(toX, toY);
      ctx.lineTo(
        toX - headSize * Math.cos(angle - Math.PI / 6),
        toY - headSize * Math.sin(angle - Math.PI / 6)
      );
      ctx.moveTo(toX, toY);
      ctx.lineTo(
        toX - headSize * Math.cos(angle + Math.PI / 6),
        toY - headSize * Math.sin(angle + Math.PI / 6)
      );
      ctx.stroke();
    }
  }
  
  // Convert canvas to blob
  return new Promise((resolve) => {
    canvas.toBlob((blob) => {
      resolve(blob!);
    }, 'image/png');
  });
}

function parseFENPosition(fen: string): string[][] {
  const position = fen.split(' ')[0];
  const ranks = position.split('/');
  const board: string[][] = [];
  
  for (const rank of ranks) {
    const row: string[] = [];
    for (const char of rank) {
      if (char >= '1' && char <= '8') {
        const emptySquares = parseInt(char);
        for (let i = 0; i < emptySquares; i++) {
          row.push('');
        }
      } else {
        row.push(char);
      }
    }
    board.push(row);
  }
  
  return board;
}

// Download file helper
export function downloadFile(content: string | Blob, filename: string): void {
  const blob = typeof content === 'string' 
    ? new Blob([content], { type: 'text/plain' })
    : content;
    
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}

// Generate filename based on game info
export function generateFilename(
  gameInfo: GameState['gameInfo'],
  type: 'pgn' | 'png'
): string {
  const date = gameInfo.date ? gameInfo.date.replace(/\./g, '-') : new Date().toISOString().split('T')[0];
  const white = (gameInfo.white || 'White').replace(/[^a-zA-Z0-9]/g, '');
  const black = (gameInfo.black || 'Black').replace(/[^a-zA-Z0-9]/g, '');
  
  return `${date}_${white}_vs_${black}.${type}`;
} 