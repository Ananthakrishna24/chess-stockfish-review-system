export type PieceType = 'p' | 'r' | 'n' | 'b' | 'q' | 'k';
export type PieceColor = 'w' | 'b';

export interface ChessPiece {
  type: PieceType;
  color: PieceColor;
}

export interface ChessSquare {
  piece: ChessPiece | null;
  square: string; // e.g., 'e4'
}

export interface ChessMove {
  from: string;
  to: string;
  piece: PieceType;
  captured?: PieceType;
  promotion?: PieceType;
  san: string; // Standard Algebraic Notation
  fen: string; // Position after move
  moveNumber: number;
  color: PieceColor;
}

export interface GameInfo {
  white: string;
  black: string;
  whiteRating?: number;
  blackRating?: number;
  result: string; // '1-0', '0-1', '1/2-1/2', or '*'
  date?: string;
  event?: string;
  site?: string;
  opening?: string;
  eco?: string; // Encyclopedia of Chess Openings code
}

export interface GameState {
  moves: ChessMove[];
  currentMoveIndex: number;
  gameInfo: GameInfo;
  pgn: string;
  startingFen?: string;
}

export type BoardOrientation = 'white' | 'black';

export interface BoardProps {
  position: string; // FEN string
  orientation?: BoardOrientation;
  onMove?: (move: ChessMove) => void;
  highlightedSquares?: string[];
  draggable?: boolean;
} 