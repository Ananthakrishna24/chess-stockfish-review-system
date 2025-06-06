// Endgame tablebase and theoretical analysis
export interface EndgameEvaluation {
  result: 'win' | 'draw' | 'loss' | 'unknown';
  movesToMate?: number;
  difficulty: 'trivial' | 'easy' | 'medium' | 'hard' | 'very_hard';
  technique: string[];
  classification: string;
  winningMethod?: string;
}

export interface EndgamePosition {
  material: {
    white: string[];
    black: string[];
  };
  classification: string;
  theoretical: EndgameEvaluation;
}

// Common endgame patterns and their evaluations
const ENDGAME_DATABASE: Record<string, EndgameEvaluation> = {
  // King and Pawn endgames
  'KP-K': {
    result: 'win',
    difficulty: 'easy',
    technique: ['Opposition', 'King activity', 'Pawn promotion'],
    classification: 'King and Pawn vs King',
    winningMethod: 'Promote the pawn with king support'
  },
  'KPP-KP': {
    result: 'win',
    difficulty: 'medium',
    technique: ['Pawn breakthroughs', 'Outside passed pawn', 'King activity'],
    classification: 'Pawn endgame',
    winningMethod: 'Create passed pawns and promote'
  },

  // Rook endgames
  'KR-KR': {
    result: 'draw',
    difficulty: 'hard',
    technique: ['Rook activity', 'King shelter', 'Perpetual check'],
    classification: 'Rook vs Rook',
    winningMethod: 'Generally drawn with accurate play'
  },
  'KRP-KR': {
    result: 'win',
    difficulty: 'hard',
    technique: ['Cut off the king', 'Lucena position', 'Philidor position'],
    classification: 'Rook and Pawn vs Rook',
    winningMethod: 'Use rook to support pawn promotion'
  },

  // Queen endgames
  'KQ-KQ': {
    result: 'draw',
    difficulty: 'very_hard',
    technique: ['Perpetual check', 'Stalemate tricks', 'Queen activity'],
    classification: 'Queen vs Queen',
    winningMethod: 'Extremely difficult to win without material advantage'
  },
  'KQ-KR': {
    result: 'win',
    movesToMate: 10,
    difficulty: 'medium',
    technique: ['Centralize queen', 'Avoid stalemate', 'Coordinate with king'],
    classification: 'Queen vs Rook',
    winningMethod: 'Force mate with queen and king coordination'
  },

  // Minor piece endgames
  'KN-K': {
    result: 'draw',
    difficulty: 'trivial',
    technique: ['Insufficient material'],
    classification: 'Knight vs King',
    winningMethod: 'Theoretical draw - insufficient material'
  },
  'KB-K': {
    result: 'draw',
    difficulty: 'trivial',
    technique: ['Insufficient material'],
    classification: 'Bishop vs King',
    winningMethod: 'Theoretical draw - insufficient material'
  },
  'KBB-K': {
    result: 'win',
    movesToMate: 19,
    difficulty: 'hard',
    technique: ['Bishop coordination', 'Corner mate', 'King centralization'],
    classification: 'Two Bishops vs King',
    winningMethod: 'Force mate by driving king to corner'
  },
  'KBN-K': {
    result: 'win',
    movesToMate: 33,
    difficulty: 'very_hard',
    technique: ['Force king to corner', 'Bishop and knight coordination', 'Precise technique'],
    classification: 'Bishop and Knight vs King',
    winningMethod: 'Force mate in corner of bishop\'s color'
  },

  // Complex endgames
  'KRR-KR': {
    result: 'win',
    difficulty: 'hard',
    technique: ['Rook coordination', 'Back rank mate', 'Cut off enemy king'],
    classification: 'Two Rooks vs Rook',
    winningMethod: 'Coordinate rooks for decisive attack'
  },
  'KQR-KQ': {
    result: 'win',
    difficulty: 'hard',
    technique: ['Piece coordination', 'Avoid perpetual check', 'King safety'],
    classification: 'Queen and Rook vs Queen',
    winningMethod: 'Use material advantage carefully'
  }
};

// Analyze endgame position
export function analyzeEndgame(fen: string): EndgameEvaluation | null {
  const materialSignature = getMaterialSignature(fen);
  
  if (!materialSignature) return null;
  
  const evaluation = ENDGAME_DATABASE[materialSignature];
  
  if (evaluation) {
    return {
      ...evaluation,
      // Add position-specific adjustments
      ...getPositionSpecificEvaluation(fen, materialSignature)
    };
  }
  
  // For unknown endgames, provide general guidance
  return analyzeUnknownEndgame(materialSignature);
}

function getMaterialSignature(fen: string): string | null {
  try {
    const position = fen.split(' ')[0];
    const pieces = position.replace(/\d/g, '').replace(/\//g, '');
    
    const white: string[] = [];
    const black: string[] = [];
    
    for (const piece of pieces) {
      if (piece === piece.toUpperCase()) {
        white.push(piece);
      } else {
        black.push(piece.toUpperCase());
      }
    }
    
    // Sort pieces for consistent signature
    white.sort();
    black.sort();
    
    const whiteStr = white.join('');
    const blackStr = black.join('');
    
    return `${whiteStr}-${blackStr}`;
  } catch {
    return null;
  }
}

function getPositionSpecificEvaluation(fen: string, materialSignature: string): Partial<EndgameEvaluation> {
  // Add specific position analysis based on material and position
  const position = fen.split(' ')[0];
  const activeColor = fen.split(' ')[1];
  
  // For pawn endgames, check for opposition and key squares
  if (materialSignature.includes('P')) {
    return analyzePawnEndgame(fen);
  }
  
  // For rook endgames, check for cut-off patterns
  if (materialSignature.includes('R')) {
    return analyzeRookEndgame(fen);
  }
  
  return {};
}

function analyzePawnEndgame(fen: string): Partial<EndgameEvaluation> {
  // Simplified pawn endgame analysis
  const techniques: string[] = [];
  
  // Check for passed pawns
  if (hasPassedPawn(fen)) {
    techniques.push('Support passed pawn');
    techniques.push('King activity crucial');
  }
  
  // Check for opposition potential
  techniques.push('Fight for opposition');
  techniques.push('Control key squares');
  
  return {
    technique: techniques,
    difficulty: 'medium'
  };
}

function analyzeRookEndgame(fen: string): Partial<EndgameEvaluation> {
  const techniques: string[] = [];
  
  techniques.push('Activate the rook');
  techniques.push('King safety important');
  
  if (hasPassedPawn(fen)) {
    techniques.push('Support passed pawn');
    techniques.push('Cut off enemy king');
  }
  
  return {
    technique: techniques,
    difficulty: 'hard'
  };
}

function hasPassedPawn(fen: string): boolean {
  // Simplified check for passed pawns (mock implementation)
  return Math.random() < 0.3; // 30% chance for demonstration
}

function analyzeUnknownEndgame(materialSignature: string): EndgameEvaluation {
  const whiteMaterial = materialSignature.split('-')[0];
  const blackMaterial = materialSignature.split('-')[1];
  
  // Calculate material balance
  const materialValue = calculateMaterialValue(whiteMaterial) - calculateMaterialValue(blackMaterial);
  
  let result: EndgameEvaluation['result'] = 'unknown';
  let difficulty: EndgameEvaluation['difficulty'] = 'medium';
  const technique: string[] = [];
  
  if (Math.abs(materialValue) >= 5) { // Significant material advantage
    result = materialValue > 0 ? 'win' : 'loss';
    difficulty = 'medium';
    technique.push('Convert material advantage');
    technique.push('Avoid unnecessary complications');
  } else if (Math.abs(materialValue) >= 3) { // Minor material advantage
    result = materialValue > 0 ? 'win' : 'loss';
    difficulty = 'hard';
    technique.push('Precise technique required');
    technique.push('Centralize pieces');
  } else {
    result = 'draw';
    difficulty = 'medium';
    technique.push('Active piece play');
    technique.push('Look for tactical opportunities');
  }
  
  return {
    result,
    difficulty,
    technique,
    classification: `Complex endgame: ${materialSignature}`,
    winningMethod: result === 'win' ? 'Use material/positional advantage' : 'Hold the balance'
  };
}

function calculateMaterialValue(material: string): number {
  const values: Record<string, number> = {
    'Q': 9, 'R': 5, 'B': 3, 'N': 3, 'P': 1, 'K': 0
  };
  
  let total = 0;
  for (const piece of material) {
    total += values[piece] || 0;
  }
  return total;
}

// Check if position is in endgame phase
export function isEndgame(fen: string): boolean {
  const materialSignature = getMaterialSignature(fen);
  if (!materialSignature) return false;
  
  const allMaterial = materialSignature.replace('-', '');
  const totalPieces = allMaterial.replace(/K/g, '').length; // Exclude kings
  
  // Consider it endgame if <= 6 pieces (excluding kings)
  return totalPieces <= 6;
}

// Get endgame learning resources
export function getEndgameResources(materialSignature: string): {
  studyMaterial: string[];
  keyPositions: string[];
  practiceRecommendations: string[];
} {
  const baseSignature = normalizeSignature(materialSignature);
  
  const resources = {
    studyMaterial: [
      'Study basic endgame principles',
      'Practice king and pawn vs king positions',
      'Learn opposition and triangulation'
    ],
    keyPositions: [
      'Lucena position',
      'Philidor position',
      'Basic checkmate patterns'
    ],
    practiceRecommendations: [
      'Solve endgame puzzles daily',
      'Play endgame positions against computer',
      'Study master endgames'
    ]
  };
  
  // Add specific resources based on material
  if (baseSignature.includes('R')) {
    resources.studyMaterial.push('Rook endgame fundamentals');
    resources.keyPositions.push('Rook vs pawn positions');
  }
  
  if (baseSignature.includes('Q')) {
    resources.studyMaterial.push('Queen endgame technique');
    resources.keyPositions.push('Queen vs pawn positions');
  }
  
  if (baseSignature.includes('B') || baseSignature.includes('N')) {
    resources.studyMaterial.push('Minor piece endgames');
    resources.keyPositions.push('Bishop vs knight endgames');
  }
  
  return resources;
}

function normalizeSignature(signature: string): string {
  return signature.replace(/K/g, '').toUpperCase();
} 