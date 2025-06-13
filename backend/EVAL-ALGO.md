# Lichess Position Evaluation System Architecture

Lichess's position evaluation system represents a sophisticated chess analysis platform that transforms raw Stockfish engine output into intuitive visual feedback through carefully calibrated algorithms and distributed computing architecture. The system processes millions of positions daily, combining cloud-based pre-computed evaluations with real-time browser analysis to deliver comprehensive position assessment.

## Core evaluation algorithm and mathematical foundation

Lichess **does not use raw Stockfish centipawn values directly**. Instead, it transforms them using a statistically-derived winning percentage formula calibrated from real game data. The core conversion uses a logistic function:

```
Win% = 50 + 50 * (2 / (1 + exp(-0.00368208 * centipawns)) - 1)  
```

This coefficient (-0.00368208) was determined through logistic regression analysis of approximately 75,000 games from players rated 2300+ on Lichess. The formula addresses a fundamental problem with raw centipawn evaluations: **context independence**. Losing 300 centipawns in an equal position has dramatically different significance than losing 300 centipawns when already winning by 500 centipawns.

All evaluations are **capped at ±1000 centipawns** before processing to normalize extreme positions and prevent display anomalies. Forced mate positions are converted to these maximum values regardless of mate distance, ensuring consistent visual representation.

## Stockfish integration and engine communication 

Lichess employs a **dual-architecture approach** for Stockfish integration, combining distributed server analysis with client-side browser engines.

**Server-side distributed analysis** operates through "fishnet," a volunteer-based computing network written in Rust. The system uses standard UCI (Universal Chess Interface) protocol for engine communication:
- Fishnet clients acquire analysis jobs via HTTP API calls
- Each job processes positions using consistent single-threaded analysis 
- Minimum performance requirement: ~2 million nodes per 6 seconds
- Results stored in cloud database containing 15+ million evaluated positions

**Client-side browser analysis** runs Stockfish directly in users' browsers using WebAssembly for modern browsers and JavaScript fallback for older ones. This provides immediate feedback without server roundtrips, typically analyzing to depth 22-23 automatically before requiring manual continuation.

The system processes raw UCI output format:
```
info depth 18 score cp 24 nodes 1686023 nps 1670251 time 1004 pv e2e4 e7e5 g1f3
```

Extracting centipawn values, mate distances, principal variations, and performance metrics for conversion to display format.

## Visual representation and technical implementation

The evaluation bar and graph employ sophisticated algorithms to convert numerical evaluations into intuitive visual feedback. The **evaluation bar** maps winning percentages to visual position using a non-linear scale that emphasizes differences near equality more than extreme values.

Technical implementation uses TypeScript with virtual DOM rendering (snabbdom library) for efficient updates. Key components include:
- **CevalCtrl**: Client evaluation controller managing engine communication
- **AnalyseCtrl**: Main analysis controller coordinating display updates  
- **Evaluation view components**: Handle visual representation and user interaction

The **evaluation graph** utilizes Highcharts 4.2.5 for interactive SVG-based visualization, plotting evaluation changes over move sequences with click-to-navigate functionality and color-coded annotations for blunders (red), mistakes (orange), and inaccuracies (yellow).

## Processing algorithms and data transformation

Beyond basic centipawn conversion, Lichess applies several sophisticated processing algorithms:

**Accuracy calculation** uses an exponential formula that measures the quality of moves based on winning percentage changes:
```
Accuracy% = 103.1668 * exp(-0.04354 * (winPercentBefore - winPercentAfter)) - 3.1669
```

**Smoothing algorithms** implement windowing systems with standard deviation weighting for consistency across game analysis. Window sizes are dynamically calculated (typically game length divided by 10, clamped between 2-8 moves) to balance precision with stability.

**Performance optimizations** include evaluation caching, incremental updates, WebWorker isolation for non-blocking analysis, and memory management for extended analysis sessions.

## Key differences from raw Stockfish output

Lichess's evaluation display differs significantly from raw engine output in several critical ways:

1. **Statistical calibration**: Raw Stockfish centipawns represent relative position strength, while Lichess shows contextual winning chances based on actual game outcomes
2. **Value normalization**: Extreme evaluations are capped and converted to prevent display inconsistencies
3. **Human-centered scaling**: The visual bar emphasizes meaningful differences near equality rather than using linear centipawn scaling
4. **Accuracy integration**: Moves are evaluated not just on position quality but on practical winning chances

## Visual conversion and display algorithms

The conversion from engine values to visual representation follows a multi-stage process:

1. **Raw evaluation processing**: Centipawn extraction and mate score conversion
2. **Statistical transformation**: Application of the winning percentage formula  
3. **Visual mapping**: Non-linear positioning within the evaluation bar
4. **Color coding**: White advantage (top), black advantage (bottom) with gradient transitions
5. **Interactive features**: Hover displays, click navigation, and real-time updates

Mate scores receive special handling, displayed as "M5" notation while internally processed as extreme winning percentages to push the evaluation bar to maximum extent.

## Evolution and continuous improvement

The evaluation system has undergone significant evolution, including major upgrades to Stockfish 16.1 with NNUE neural network evaluation, complete rewrite of the fishnet system in Rust for improved reliability, and expansion of the cloud evaluation database from 7 million to over 15 million positions.

Recent improvements focus on **analysis quality consistency** through single-threaded server analysis and **performance optimization** through automatic depth limits and intelligent caching strategies. The system continues to evolve based on user feedback and advances in chess engine technology.

## Conclusion

Lichess's evaluation system transforms raw chess engine output into meaningful, contextual feedback through sophisticated mathematical modeling and distributed computing architecture. By converting centipawn values to empirically-calibrated winning percentages and applying careful visual design principles, the system provides users with intuitive position assessment that reflects real-world game outcomes rather than abstract engine calculations. This approach represents a significant advancement over traditional centipawn displays, offering chess players actionable insights grounded in statistical analysis of millions of actual games.