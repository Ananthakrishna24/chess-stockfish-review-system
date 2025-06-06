// Positional analysis and strategic evaluation
export interface PositionalFactors {
  pawnStructure: {
    score: number;
    issues: string[];
    strengths: string[];
  };
  pieceActivity: {
    score: number;
    activenesss: number;
    coordination: number;
    comments: string[];
  };
  kingSafety: {
    score: number;
    threats: string[];
    protection: string[];
  };
  spaceAdvantage: {
    score: number;
    controlledSquares: number;
    comments: string[];
  };
  weakSquares: {
    holes: string[];
    outposts: string[];
    comments: string[];
  };
}

export interface PositionalEvaluation {
  overallScore: number; // Positive favors white, negative favors black
  phase: 'opening' | 'middlegame' | 'endgame';
  factors: PositionalFactors;
  recommendations: {
    immediate: string[];
    strategic: string[];
    warnings: string[];
  };
  imbalances: string[];
  characterization: string;
}

// Analyze position from strategic perspective
export function analyzePosition(fen: string): PositionalEvaluation {
  // Parse FEN string
  const [position, activeColor, castling, enPassant, halfmove, fullmove] = fen.split(' ');
  
  // Determine game phase
  const phase = determineGamePhase(position);
  
  // Evaluate each positional factor
  const factors = evaluatePositionalFactors(position, activeColor);
  
  // Calculate overall positional score
  const overallScore = calculateOverallScore(factors);
  
  // Generate recommendations
  const recommendations = generateRecommendations(factors, phase, activeColor);
  
  // Identify imbalances
  const imbalances = identifyImbalances(factors, position);
  
  // Characterize the position
  const characterization = characterizePosition(factors, phase, overallScore);
  
  return {
    overallScore,
    phase,
    factors,
    recommendations,
    imbalances,
    characterization
  };
}

function determineGamePhase(position: string): 'opening' | 'middlegame' | 'endgame' {
  const pieces = position.replace(/\d/g, '').replace(/\//g, '');
  const totalPieces = pieces.length;
  
  // Count major pieces (queens and rooks)
  const majorPieces = (pieces.match(/[QRqr]/g) || []).length;
  
  if (totalPieces <= 12 || majorPieces <= 2) {
    return 'endgame';
  } else if (totalPieces <= 20) {
    return 'middlegame';
  } else {
    return 'opening';
  }
}

function evaluatePositionalFactors(position: string, activeColor: string): PositionalFactors {
  return {
    pawnStructure: analyzePawnStructure(position),
    pieceActivity: analyzePieceActivity(position, activeColor),
    kingSafety: analyzeKingSafety(position),
    spaceAdvantage: analyzeSpaceAdvantage(position),
    weakSquares: analyzeWeakSquares(position)
  };
}

function analyzePawnStructure(position: string): PositionalFactors['pawnStructure'] {
  const issues: string[] = [];
  const strengths: string[] = [];
  let score = 0;
  
  // Mock pawn structure analysis (in real implementation, would analyze FEN)
  const pawnCount = (position.match(/[Pp]/g) || []).length;
  
  // Simulate various pawn structure evaluations
  if (Math.random() < 0.3) {
    issues.push('Doubled pawns on f-file');
    score -= 0.2;
  }
  
  if (Math.random() < 0.2) {
    issues.push('Isolated pawn on d4');
    score -= 0.3;
  }
  
  if (Math.random() < 0.4) {
    strengths.push('Strong pawn center');
    score += 0.3;
  }
  
  if (Math.random() < 0.3) {
    strengths.push('Advanced passed pawn');
    score += 0.4;
  }
  
  if (Math.random() < 0.25) {
    issues.push('Backward pawn on c6');
    score -= 0.2;
  }
  
  return { score, issues, strengths };
}

function analyzePieceActivity(position: string, activeColor: string): PositionalFactors['pieceActivity'] {
  const comments: string[] = [];
  let activeness = Math.random() * 2 - 1; // -1 to 1
  let coordination = Math.random() * 2 - 1; // -1 to 1
  
  // Mock piece activity analysis
  if (activeness > 0.3) {
    comments.push('Pieces are well-placed and active');
  } else if (activeness < -0.3) {
    comments.push('Pieces lack activity and scope');
  }
  
  if (coordination > 0.3) {
    comments.push('Good piece coordination');
  } else if (coordination < -0.3) {
    comments.push('Poor piece coordination');
  }
  
  // Check for specific piece placements
  if (Math.random() < 0.3) {
    comments.push('Centralized knight on e5');
    activeness += 0.2;
  }
  
  if (Math.random() < 0.25) {
    comments.push('Bishop pair advantage');
    coordination += 0.3;
  }
  
  if (Math.random() < 0.2) {
    comments.push('Rook on open file');
    activeness += 0.3;
  }
  
  const score = (activeness + coordination) / 2;
  
  return { score, activenesss: activeness, coordination, comments };
}

function analyzeKingSafety(position: string): PositionalFactors['kingSafety'] {
  const threats: string[] = [];
  const protection: string[] = [];
  let score = 0;
  
  // Mock king safety analysis
  if (Math.random() < 0.2) {
    threats.push('Exposed king in center');
    score -= 0.5;
  }
  
  if (Math.random() < 0.3) {
    threats.push('Weakened kingside pawn shield');
    score -= 0.3;
  }
  
  if (Math.random() < 0.4) {
    protection.push('King safely castled');
    score += 0.3;
  }
  
  if (Math.random() < 0.25) {
    protection.push('Strong pawn shield');
    score += 0.2;
  }
  
  if (Math.random() < 0.2) {
    threats.push('Enemy pieces attacking kingside');
    score -= 0.4;
  }
  
  return { score, threats, protection };
}

function analyzeSpaceAdvantage(position: string): PositionalFactors['spaceAdvantage'] {
  const comments: string[] = [];
  const controlledSquares = Math.floor(Math.random() * 20) + 10; // 10-30 squares
  let score = (controlledSquares - 20) / 10; // Normalize around 20 squares
  
  if (score > 0.3) {
    comments.push('Significant space advantage');
  } else if (score < -0.3) {
    comments.push('Cramped position');
  } else {
    comments.push('Balanced space control');
  }
  
  if (Math.random() < 0.3) {
    comments.push('Control of center squares');
    score += 0.2;
  }
  
  return { score, controlledSquares, comments };
}

function analyzeWeakSquares(position: string): PositionalFactors['weakSquares'] {
  const holes: string[] = [];
  const outposts: string[] = [];
  const comments: string[] = [];
  
  // Mock weak square analysis
  const weakSquareChance = Math.random();
  
  if (weakSquareChance < 0.3) {
    holes.push('d6', 'f6');
    comments.push('Weak squares in enemy position');
  }
  
  if (weakSquareChance < 0.2) {
    outposts.push('e5');
    comments.push('Excellent outpost for knight');
  }
  
  if (weakSquareChance < 0.25) {
    holes.push('h7');
    comments.push('Potential mating attack target');
  }
  
  return { holes, outposts, comments };
}

function calculateOverallScore(factors: PositionalFactors): number {
  return (
    factors.pawnStructure.score * 0.25 +
    factors.pieceActivity.score * 0.30 +
    factors.kingSafety.score * 0.25 +
    factors.spaceAdvantage.score * 0.20
  );
}

function generateRecommendations(
  factors: PositionalFactors, 
  phase: string, 
  activeColor: string
): PositionalEvaluation['recommendations'] {
  const immediate: string[] = [];
  const strategic: string[] = [];
  const warnings: string[] = [];
  
  // Immediate recommendations based on factors
  if (factors.kingSafety.score < -0.3) {
    immediate.push('Improve king safety');
    immediate.push('Consider defensive moves');
  }
  
  if (factors.pieceActivity.score < -0.2) {
    immediate.push('Activate passive pieces');
    immediate.push('Look for piece coordination');
  }
  
  // Strategic recommendations based on game phase
  if (phase === 'opening') {
    strategic.push('Complete development');
    strategic.push('Fight for center control');
    strategic.push('Ensure king safety');
  } else if (phase === 'middlegame') {
    strategic.push('Create pawn breaks');
    strategic.push('Improve piece positions');
    strategic.push('Look for tactical opportunities');
  } else { // endgame
    strategic.push('Centralize the king');
    strategic.push('Create passed pawns');
    strategic.push('Activate all pieces');
  }
  
  // Warnings based on positional issues
  if (factors.pawnStructure.issues.length > 0) {
    warnings.push('Pawn structure weaknesses detected');
  }
  
  if (factors.kingSafety.threats.length > 0) {
    warnings.push('King safety concerns');
  }
  
  return { immediate, strategic, warnings };
}

function identifyImbalances(factors: PositionalFactors, position: string): string[] {
  const imbalances: string[] = [];
  
  // Material imbalances (simplified)
  const pieces = position.replace(/\d/g, '').replace(/\//g, '');
  const queens = (pieces.match(/[Qq]/g) || []).length;
  const rooks = (pieces.match(/[Rr]/g) || []).length;
  const bishops = (pieces.match(/[Bb]/g) || []).length;
  const knights = (pieces.match(/[Nn]/g) || []).length;
  
  if (bishops > knights) {
    imbalances.push('Bishop vs Knight imbalance');
  } else if (knights > bishops) {
    imbalances.push('Knight vs Bishop imbalance');
  }
  
  // Structural imbalances
  if (factors.pawnStructure.issues.includes('Doubled pawns')) {
    imbalances.push('Doubled pawn structure');
  }
  
  if (factors.spaceAdvantage.score > 0.3) {
    imbalances.push('Space advantage');
  } else if (factors.spaceAdvantage.score < -0.3) {
    imbalances.push('Space disadvantage');
  }
  
  return imbalances;
}

function characterizePosition(
  factors: PositionalFactors, 
  phase: string, 
  overallScore: number
): string {
  const characteristics: string[] = [];
  
  // Overall evaluation
  if (overallScore > 0.4) {
    characteristics.push('Clearly better');
  } else if (overallScore > 0.2) {
    characteristics.push('Slightly better');
  } else if (overallScore < -0.4) {
    characteristics.push('Clearly worse');
  } else if (overallScore < -0.2) {
    characteristics.push('Slightly worse');
  } else {
    characteristics.push('Balanced');
  }
  
  // Position type
  if (factors.pieceActivity.score > 0.3) {
    characteristics.push('dynamic');
  } else if (factors.pawnStructure.score > 0.3) {
    characteristics.push('positional');
  } else {
    characteristics.push('strategic');
  }
  
  // Phase-specific characterization
  characteristics.push(phase);
  
  return characteristics.join(' ') + ' position';
}

// Time management analysis
export interface TimeAnalysis {
  timeSpent: number;
  averageTimePerMove: number;
  criticalMoments: number[];
  timeDistribution: {
    opening: number;
    middlegame: number;
    endgame: number;
  };
  recommendations: string[];
}

export function analyzeTimeManagement(
  moveTimes: number[], 
  gamePhases: { opening: number; middlegame: number; endgame: number }
): TimeAnalysis {
  const totalTime = moveTimes.reduce((sum, time) => sum + time, 0);
  const averageTime = totalTime / moveTimes.length;
  
  // Find moves where significantly more time was spent
  const criticalMoments = moveTimes
    .map((time, index) => ({ time, index }))
    .filter(move => move.time > averageTime * 2)
    .map(move => move.index);
  
  // Calculate time distribution by phase
  const openingTime = moveTimes.slice(0, gamePhases.opening).reduce((sum, time) => sum + time, 0);
  const middlegameTime = moveTimes.slice(gamePhases.opening, gamePhases.middlegame).reduce((sum, time) => sum + time, 0);
  const endgameTime = moveTimes.slice(gamePhases.middlegame).reduce((sum, time) => sum + time, 0);
  
  const timeDistribution = {
    opening: openingTime,
    middlegame: middlegameTime,
    endgame: endgameTime
  };
  
  // Generate recommendations
  const recommendations: string[] = [];
  
  if (openingTime > totalTime * 0.4) {
    recommendations.push('Spent too much time in opening - study opening theory');
  }
  
  if (criticalMoments.length > moveTimes.length * 0.2) {
    recommendations.push('Consider time management - many long thinks');
  }
  
  if (averageTime > 120) { // 2 minutes per move
    recommendations.push('Play faster to avoid time pressure');
  } else if (averageTime < 30) { // 30 seconds per move
    recommendations.push('Consider spending more time on critical positions');
  }
  
  return {
    timeSpent: totalTime,
    averageTimePerMove: averageTime,
    criticalMoments,
    timeDistribution,
    recommendations
  };
} 