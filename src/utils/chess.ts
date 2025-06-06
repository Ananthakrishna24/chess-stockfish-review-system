import { Chess } from 'chess.js';
import { ChessMove, GameInfo, GameState, PieceType, PieceColor } from '@/types/chess';

export class ChessGameManager {
  private chess: Chess;
  private moveHistory: ChessMove[] = [];

  constructor(pgn?: string, fen?: string) {
    this.chess = new Chess();
    
    if (fen) {
      this.chess.load(fen);
    } else if (pgn) {
      this.loadPGN(pgn);
    }
  }

  loadPGN(pgn: string): GameState {
    try {
      this.chess.loadPgn(pgn);
      this.moveHistory = this.extractMoves();
      
      const gameInfo = this.extractGameInfo(pgn);
      
      return {
        moves: this.moveHistory,
        currentMoveIndex: this.moveHistory.length - 1,
        gameInfo,
        pgn,
        startingFen: 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1'
      };
    } catch (error) {
      throw new Error(`Invalid PGN: ${error}`);
    }
  }

  private extractMoves(): ChessMove[] {
    const tempChess = new Chess();
    const moves: ChessMove[] = [];
    const history = this.chess.history({ verbose: true });

    history.forEach((move, index) => {
      const chessMove: ChessMove = {
        from: move.from,
        to: move.to,
        piece: move.piece as PieceType,
        captured: move.captured as PieceType | undefined,
        promotion: move.promotion as PieceType | undefined,
        san: move.san,
        fen: move.after || tempChess.fen(),
        moveNumber: Math.ceil((index + 1) / 2),
        color: move.color as PieceColor
      };

      tempChess.move(move);
      moves.push(chessMove);
    });

    return moves;
  }

  private extractGameInfo(pgn: string): GameInfo {
    const headers = this.chess.header();
    
    return {
      white: headers.White || 'Unknown',
      black: headers.Black || 'Unknown',
      whiteRating: headers.WhiteElo ? parseInt(headers.WhiteElo) : undefined,
      blackRating: headers.BlackElo ? parseInt(headers.BlackElo) : undefined,
      result: headers.Result || '*',
      date: headers.Date,
      event: headers.Event,
      site: headers.Site,
      opening: headers.Opening,
      eco: headers.ECO
    };
  }

  getPosition(moveIndex?: number): string {
    if (moveIndex === undefined) {
      return this.chess.fen();
    }

    if (moveIndex < 0 || moveIndex >= this.moveHistory.length) {
      return 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1';
    }

    return this.moveHistory[moveIndex].fen;
  }

  makeMove(from: string, to: string, promotion?: PieceType): ChessMove | null {
    try {
      const move = this.chess.move({
        from,
        to,
        promotion
      });

      if (move) {
        const chessMove: ChessMove = {
          from: move.from,
          to: move.to,
          piece: move.piece as PieceType,
          captured: move.captured as PieceType | undefined,
          promotion: move.promotion as PieceType | undefined,
          san: move.san,
          fen: this.chess.fen(),
          moveNumber: Math.ceil(this.chess.moveNumber()),
          color: move.color as PieceColor
        };

        this.moveHistory.push(chessMove);
        return chessMove;
      }
    } catch (error) {
      console.error('Invalid move:', error);
    }

    return null;
  }

  isGameOver(): boolean {
    return this.chess.isGameOver();
  }

  isCheck(): boolean {
    return this.chess.inCheck();
  }

  isCheckmate(): boolean {
    return this.chess.isCheckmate();
  }

  isDraw(): boolean {
    return this.chess.isDraw();
  }

  getPossibleMoves(square?: string): string[] {
    if (square) {
      return this.chess.moves({ square, verbose: false });
    }
    return this.chess.moves();
  }

  getGameResult(): string {
    if (this.chess.isCheckmate()) {
      return this.chess.turn() === 'w' ? '0-1' : '1-0';
    } else if (this.chess.isDraw()) {
      return '1/2-1/2';
    }
    return '*';
  }

  reset(): void {
    this.chess.reset();
    this.moveHistory = [];
  }
}

export function parseSquareColor(square: string): 'light' | 'dark' {
  const file = square.charCodeAt(0) - 97; // a=0, b=1, etc.
  const rank = parseInt(square[1]) - 1;   // 1=0, 2=1, etc.
  
  return (file + rank) % 2 === 0 ? 'dark' : 'light';
}

export function getSquareCoordinates(square: string): { x: number; y: number } {
  const file = square.charCodeAt(0) - 97; // a=0, b=1, etc.
  const rank = parseInt(square[1]) - 1;   // 1=0, 2=1, etc.
  
  return { x: file, y: 7 - rank }; // Flip rank for display
}

export function coordinatesToSquare(x: number, y: number): string {
  const file = String.fromCharCode(97 + x); // 0=a, 1=b, etc.
  const rank = (8 - y).toString();         // Flip y for chess notation
  
  return file + rank;
}

export function isValidSquare(square: string): boolean {
  return /^[a-h][1-8]$/.test(square);
}

export function getPieceUnicode(piece: PieceType, color: PieceColor): string {
  const pieces = {
    w: {
      k: '♔', q: '♕', r: '♖', b: '♗', n: '♘', p: '♙'
    },
    b: {
      k: '♚', q: '♛', r: '♜', b: '♝', n: '♞', p: '♟'
    }
  };
  
  return pieces[color][piece];
} 