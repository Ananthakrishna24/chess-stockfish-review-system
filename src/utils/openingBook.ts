// Opening book database with ECO codes and variations
export interface OpeningVariation {
  name: string;
  eco: string;
  moves: string[];
  characteristics: string[];
  popularity: 'common' | 'uncommon' | 'rare';
  difficulty: 'beginner' | 'intermediate' | 'advanced' | 'master';
  themes: string[];
}

export interface OpeningData {
  name: string;
  eco: string;
  moves: string[];
  variations: OpeningVariation[];
  characteristics: string[];
  theory: {
    mainIdeas: string[];
    typicalPlans: string[];
    commonMistakes: string[];
  };
}

// Comprehensive opening database
export const OPENING_DATABASE: OpeningData[] = [
  {
    name: "Ruy Lopez",
    eco: "C60-C99",
    moves: ["e4", "e5", "Nf3", "Nc6", "Bb5"],
    variations: [
      {
        name: "Berlin Defense",
        eco: "C65",
        moves: ["e4", "e5", "Nf3", "Nc6", "Bb5", "Nf6"],
        characteristics: ["Solid", "Drawish", "Endgame-oriented"],
        popularity: "common",
        difficulty: "advanced",
        themes: ["solid_defense", "endgame_technique"]
      },
      {
        name: "Morphy Defense",
        eco: "C77-C78",
        moves: ["e4", "e5", "Nf3", "Nc6", "Bb5", "a6", "Ba4", "Nf6"],
        characteristics: ["Classical", "Balanced", "Strategic"],
        popularity: "common",
        difficulty: "intermediate",
        themes: ["center_control", "piece_development"]
      }
    ],
    characteristics: ["King-side attack", "Center control", "Long castling"],
    theory: {
      mainIdeas: ["Control the center", "Develop pieces quickly", "Castle king to safety"],
      typicalPlans: ["Kingside attack", "Central breakthrough", "Queenside expansion"],
      commonMistakes: ["Moving the same piece twice", "Neglecting king safety", "Premature attacks"]
    }
  },
  {
    name: "Italian Game",
    eco: "C50-C54",
    moves: ["e4", "e5", "Nf3", "Nc6", "Bc4"],
    variations: [
      {
        name: "Italian Game, Classical",
        eco: "C53",
        moves: ["e4", "e5", "Nf3", "Nc6", "Bc4", "Be7"],
        characteristics: ["Positional", "Strategic"],
        popularity: "common",
        difficulty: "beginner",
        themes: ["piece_development", "center_control"]
      }
    ],
    characteristics: ["Quick development", "Central control", "Tactical possibilities"],
    theory: {
      mainIdeas: ["Rapid piece development", "Control the center", "Attack weak points"],
      typicalPlans: ["Kingside attack", "Central advance", "Piece coordination"],
      commonMistakes: ["Rushing the attack", "Ignoring development", "Weakening king position"]
    }
  },
  {
    name: "Queen's Gambit",
    eco: "D06-D69",
    moves: ["d4", "d5", "c4"],
    variations: [
      {
        name: "Queen's Gambit Declined",
        eco: "D30-D69",
        moves: ["d4", "d5", "c4", "e6"],
        characteristics: ["Solid", "Positional", "Strategic"],
        popularity: "common",
        difficulty: "intermediate",
        themes: ["center_control", "positional_play"]
      },
      {
        name: "Queen's Gambit Accepted",
        eco: "D20-D29",
        moves: ["d4", "d5", "c4", "dxc4"],
        characteristics: ["Dynamic", "Active", "Tactical"],
        popularity: "common",
        difficulty: "intermediate",
        themes: ["piece_activity", "tactical_complications"]
      }
    ],
    characteristics: ["Central control", "Positional pressure", "Long-term advantage"],
    theory: {
      mainIdeas: ["Control the center", "Create long-term pressure", "Develop harmoniously"],
      typicalPlans: ["Minority attack", "Central breakthrough", "Kingside expansion"],
      commonMistakes: ["Passive piece play", "Ignoring central tension", "Premature piece trades"]
    }
  },
  {
    name: "Sicilian Defense",
    eco: "B20-B99",
    moves: ["e4", "c5"],
    variations: [
      {
        name: "Sicilian Dragon",
        eco: "B70-B79",
        moves: ["e4", "c5", "Nf3", "d6", "d4", "cxd4", "Nxd4", "Nf6", "Nc3", "g6"],
        characteristics: ["Sharp", "Tactical", "Double-edged"],
        popularity: "uncommon",
        difficulty: "advanced",
        themes: ["king_safety", "tactical_complications", "opposite_castling"]
      },
      {
        name: "Sicilian Najdorf",
        eco: "B90-B99",
        moves: ["e4", "c5", "Nf3", "d6", "d4", "cxd4", "Nxd4", "Nf6", "Nc3", "a6"],
        characteristics: ["Flexible", "Dynamic", "Complex"],
        popularity: "common",
        difficulty: "advanced",
        themes: ["flexibility", "counterplay", "central_control"]
      }
    ],
    characteristics: ["Counterattacking", "Asymmetrical", "Fighting spirit"],
    theory: {
      mainIdeas: ["Create imbalances", "Generate counterplay", "Fight for the initiative"],
      typicalPlans: ["Queenside counterplay", "Central breaks", "Kingside attacks"],
      commonMistakes: ["Playing too passively", "Neglecting king safety", "Rushing counterattacks"]
    }
  },
  {
    name: "King's Indian Defense",
    eco: "E60-E99",
    moves: ["d4", "Nf6", "c4", "g6", "Nc3", "Bg7"],
    variations: [
      {
        name: "King's Indian, Classical",
        eco: "E90-E99",
        moves: ["d4", "Nf6", "c4", "g6", "Nc3", "Bg7", "e4", "d6", "Nf3", "O-O", "Be2"],
        characteristics: ["Strategic", "Complex", "Long-term"],
        popularity: "common",
        difficulty: "advanced",
        themes: ["kingside_attack", "central_control", "piece_coordination"]
      }
    ],
    characteristics: ["Fianchetto", "Kingside attack", "Strategic complexity"],
    theory: {
      mainIdeas: ["Fianchetto the bishop", "Create kingside attacking chances", "Control dark squares"],
      typicalPlans: ["Kingside pawn storm", "Central breakthrough", "Piece sacrifices"],
      commonMistakes: ["Premature kingside attacks", "Neglecting the center", "Poor piece coordination"]
    }
  },
  {
    name: "French Defense",
    eco: "C00-C19",
    moves: ["e4", "e6"],
    variations: [
      {
        name: "French, Winawer",
        eco: "C15-C19",
        moves: ["e4", "e6", "d4", "d5", "Nc3", "Bb4"],
        characteristics: ["Sharp", "Tactical", "Imbalanced"],
        popularity: "uncommon",
        difficulty: "advanced",
        themes: ["pawn_structure", "tactical_complications", "piece_activity"]
      }
    ],
    characteristics: ["Solid structure", "Strategic battles", "Counterplay"],
    theory: {
      mainIdeas: ["Challenge white's center", "Create counterplay", "Use pawn structure advantages"],
      typicalPlans: ["Central breaks", "Queenside play", "Kingside counterattacks"],
      commonMistakes: ["Passive bishop play", "Neglecting piece activity", "Poor pawn structure decisions"]
    }
  }
];

// Position classification based on opening moves
export function identifyOpening(moves: string[]): OpeningData | null {
  if (moves.length < 2) return null;

  // Normalize moves (remove move numbers, annotations)
  const cleanMoves = moves.map(move => 
    move.replace(/^\d+\.+\s*/, '').replace(/[+#?!]*$/, '').trim()
  );

  // Find the best matching opening
  let bestMatch: OpeningData | null = null;
  let maxMatchLength = 0;

  for (const opening of OPENING_DATABASE) {
    const matchLength = getMatchingMoveCount(cleanMoves, opening.moves);
    
    if (matchLength >= 2 && matchLength > maxMatchLength) {
      maxMatchLength = matchLength;
      bestMatch = opening;
    }

    // Check variations for better matches
    for (const variation of opening.variations) {
      const variationMatchLength = getMatchingMoveCount(cleanMoves, variation.moves);
      if (variationMatchLength > matchLength && variationMatchLength > maxMatchLength) {
        maxMatchLength = variationMatchLength;
        bestMatch = {
          ...opening,
          name: variation.name,
          eco: variation.eco,
          moves: variation.moves,
          characteristics: [...opening.characteristics, ...variation.characteristics]
        };
      }
    }
  }

  return bestMatch;
}

function getMatchingMoveCount(gameMoves: string[], openingMoves: string[]): number {
  let count = 0;
  const minLength = Math.min(gameMoves.length, openingMoves.length);
  
  for (let i = 0; i < minLength; i++) {
    if (normalizeMove(gameMoves[i]) === normalizeMove(openingMoves[i])) {
      count++;
    } else {
      break;
    }
  }
  
  return count;
}

function normalizeMove(move: string): string {
  return move.replace(/[+#?!]*$/, '').toLowerCase().trim();
}

// Get opening statistics and recommendations
export function getOpeningAnalysis(moves: string[], playerLevel: 'beginner' | 'intermediate' | 'advanced' | 'master' = 'intermediate') {
  const opening = identifyOpening(moves);
  
  if (!opening) {
    return {
      identified: false,
      name: 'Unknown Opening',
      eco: '',
      analysis: {
        strength: 'neutral',
        recommendations: ['Focus on piece development', 'Control the center', 'Ensure king safety'],
        continuation: []
      }
    };
  }

  // Analyze opening choice based on player level
  const isAppropriate = isOpeningAppropriate(opening, playerLevel);
  const recommendations = generateOpeningRecommendations(opening, moves, playerLevel);
  
  return {
    identified: true,
    name: opening.name,
    eco: opening.eco,
    characteristics: opening.characteristics,
    theory: opening.theory,
    analysis: {
      strength: isAppropriate ? 'good' : 'questionable',
      recommendations,
      continuation: getOpeningContinuation(opening, moves),
      playerLevel: getOpeningDifficulty(opening),
      popularity: getOpeningPopularity(opening)
    }
  };
}

function isOpeningAppropriate(opening: OpeningData, playerLevel: string): boolean {
  const difficulty = getOpeningDifficulty(opening);
  
  switch (playerLevel) {
    case 'beginner':
      return difficulty === 'beginner' || difficulty === 'intermediate';
    case 'intermediate':
      return difficulty !== 'master';
    case 'advanced':
    case 'master':
      return true;
    default:
      return true;
  }
}

function getOpeningDifficulty(opening: OpeningData): string {
  // Extract difficulty from variations or default
  const difficulties = opening.variations.map(v => v.difficulty);
  if (difficulties.includes('master')) return 'master';
  if (difficulties.includes('advanced')) return 'advanced';
  if (difficulties.includes('intermediate')) return 'intermediate';
  return 'beginner';
}

function getOpeningPopularity(opening: OpeningData): string {
  const popularities = opening.variations.map(v => v.popularity);
  if (popularities.includes('common')) return 'common';
  if (popularities.includes('uncommon')) return 'uncommon';
  return 'rare';
}

function generateOpeningRecommendations(opening: OpeningData, moves: string[], playerLevel: string): string[] {
  const recommendations: string[] = [];
  
  // Add theory-based recommendations
  recommendations.push(...opening.theory.mainIdeas.slice(0, 2));
  
  // Add level-specific advice
  if (playerLevel === 'beginner') {
    recommendations.push('Focus on basic opening principles');
    recommendations.push('Complete development before attacking');
  } else if (playerLevel === 'advanced' || playerLevel === 'master') {
    recommendations.push(...opening.theory.typicalPlans.slice(0, 2));
  }
  
  return recommendations;
}

function getOpeningContinuation(opening: OpeningData, moves: string[]): string[] {
  // Generate likely continuation moves based on opening theory
  const movesPlayed = moves.length;
  const theoreticalMoves = opening.moves;
  
  if (movesPlayed < theoreticalMoves.length) {
    return theoreticalMoves.slice(movesPlayed, movesPlayed + 3);
  }
  
  // Return general developing moves if past theory
  return ['Develop pieces', 'Castle king', 'Control center'];
} 