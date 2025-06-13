# Lichess Position Evaluation: Complete Algorithm Implementation

Lichess employs a sophisticated multi-layered evaluation system that converts Stockfish centipawn values into human-interpretable winning percentages and accuracy scores. **The core conversion uses an empirically-derived sigmoid function with a multiplier of -0.00368208, calibrated from analysis of over two million rated games rather than theoretical perfect play.**

## Core algorithm implementation in WinPercent.scala

The heart of Lichess's evaluation system resides in the **WinPercent.scala** file within `modules/common/src/main/scala/`. The primary conversion algorithm transforms raw Stockfish centipawn evaluations into winning percentages:

```scala
def fromCentiPawns(cp: Eval.Cp) = WinPercent { 
  50 + 50 * winningChances(cp.ceiled) 
}

private[analyse] def winningChances(cp: Eval.Cp) = {
  val MULTIPLIER = -0.00368208 // https://github.com/lichess-org/lila/pull/11148
  2 / (1 + Math.exp(MULTIPLIER * cp.value)) - 1
} atLeast -1 atMost +1
```

**The MULTIPLIER constant (-0.00368208) represents the key innovation**: derived from lichess player data rather than master games, it reflects amateur-level online play patterns. The sigmoid function `2 / (1 + e^(-0.00368208 * cp)) - 1` converts centipawns to winning chances bounded between -1 and +1, then scaled to 0-100% range.

## Comprehensive accuracy calculation system

Lichess implements a two-stage accuracy calculation system in **AccuracyPercent.scala**. Game accuracy uses sliding window analysis with volatility weighting:

```scala
def gameAccuracy(startColor: Color, cps: List[Eval.Cp]): Option[ByColor[AccuracyPercent]] = 
  val allWinPercents = (Eval.Cp.initial :: cps).map(WinPercent.fromCentiPawns)
  val windowSize = (cps.size / 10).atLeast(2).atMost(8)
  val windows = List.fill(windowSize.atMost(allWinPercentValues.size) - 2)(allWinPercentValues.take(windowSize))
  val weights = windows.map { xs => 
    Maths.standardDeviation(xs).orZero.atLeast(0.5).atMost(12) 
  }
```

Move-level accuracy conversion follows the formula: `Accuracy% = 103.1668 * exp(-0.04354 * (winPercentBefore - winPercentAfter)) - 3.1669`. This exponential decay function emphasizes larger evaluation drops, making blunders significantly more penalizing than minor inaccuracies.

## Edge cases and boundary handling mechanisms

**Mate score processing** represents a critical edge case. Lichess converts mate-in-N scores to finite centipawn equivalents using the formula `cp = 100*(21 - min(10, N))`, capping mate distance calculations at 10 moves. This ensures consistent processing while prioritizing shorter mates with higher centipawn values.

**Extreme centipawn value handling** implements systematic capping at **±1000 centipawns** for ACPL calculations. Values beyond these bounds are truncated to prevent statistical skewing, though raw engine output can exceed these limits. The sigmoid conversion naturally prevents overflow by asymptotically approaching 0% or 100% win probability for large values.

**Error handling and validation** follows defensive programming patterns throughout the codebase. All engine outputs undergo "careful validation" with malformed UCI responses filtered before processing. The system uses Optional types extensively (`Option[Int]` for evaluations) to handle missing or invalid data gracefully.

## Distributed analysis architecture and engine integration

Lichess operates a sophisticated distributed analysis system called **fishnet**, comprising Rust-based clients coordinating Stockfish analysis across volunteer computers. The architecture follows: `lila ↔ redis ↔ lila-fishnet ← HTTP ← fishnet-clients`.

**UCI protocol handling** processes standard Stockfish output formats:
- `score cp <centipawns>` for normal evaluations  
- `score mate <moves>` for forced mates
- Additional metadata including depth, nodes, time, and principal variations

**Browser versus server-side differences** create a hybrid system. Client-side analysis uses WebAssembly Stockfish (typically depth ~20) for real-time evaluation, while server-side fishnet provides deep analysis (depth 45-50+) using native Stockfish builds with NNUE neural networks.

## Performance optimization and caching strategies

**Cloud evaluation cache** stores over 14 million pre-computed positions accessible via `/api/cloud-eval` endpoint. Cache entries require analysis depth ≥23 ply and include evaluation, principal variation, and metadata. This system dramatically reduces computational load for common positions.

**Fishnet performance optimizations** include node-based analysis limits (2.25M nodes for NNUE, 4.05M for classical), asynchronous task management using Rust's tokio framework, and CPU feature detection for optimal Stockfish binary selection. The two-pass analysis system identifies key positions first, then focuses computational resources accordingly.

**Database optimization** employs MongoDB with Redis caching, storing 4.7+ billion games using periodic saves rather than per-move persistence. WiredTiger storage engine provides internal caching while Elasticsearch enables fast game querying.

## Algorithm evolution and version management

The evaluation system has undergone significant evolution. The current win percentage formula emerged from **pull request #11148** based on empirical analysis of rated game outcomes. Migration from centipawn loss (ACPL) to win%-based accuracy improved human comprehension while maintaining statistical rigor.

**Stockfish integration updates** track engine development, currently using Stockfish 16+ with NNUE architecture and SFNNv9 neural networks. Fishnet underwent complete rewrite from Python to Rust for improved performance and memory safety.

## System resilience and error recovery

**Fallback mechanisms** include multi-depth analysis progression, automatic batch reassignment for failed jobs, and client redundancy across the distributed network. Position-based time allocation prioritizes key moves while ECO opening positions are skipped to conserve computation.

**Boundary enforcement** implements multiple validation layers from engine output to user display. Memory safety benefits from WASM engine isolation preventing crashes, while network resilience provides automatic retry and reassignment for distributed analysis.

## Technical implementation specifics

**Key constants and functions** include:
- `MULTIPLIER = -0.00368208` (empirical conversion constant)
- `WinPercent.fromCentiPawns()` (main conversion function)  
- `gameAccuracy()` (comprehensive game analysis)
- `AccuracyPercent.fromWinPercents()` (move-by-move accuracy)

**Data structures** follow consistent patterns with evaluation objects containing optional centipawn and mate components, enabling robust handling of missing or invalid engine outputs.

This implementation represents a sophisticated balance of accuracy, performance, and scalability. The system successfully processes millions of daily games while maintaining evaluation quality through empirically-validated algorithms, comprehensive error handling, and distributed computational architecture calibrated specifically for the Lichess player base rather than theoretical perfect play.