package uci

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// Engine represents a UCI chess engine
type Engine struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	scanner *bufio.Scanner
	mutex   sync.Mutex
	ready   bool
}

// EngineInfo contains engine identification information
type EngineInfo struct {
	Name    string
	Author  string
	Version string
	Options map[string]Option
}

// Option represents a UCI option
type Option struct {
	Name         string
	Type         string
	Default      string
	Min          int
	Max          int
	Var          []string
}

// SearchResult contains the result of a position search
type SearchResult struct {
	BestMove           string
	PonderMove         string
	Score              int
	ScoreType          string // "cp" for centipawns, "mate" for mate
	Depth              int
	SelDepth           int
	Nodes              int64
	NodesPerSecond     int64
	Time               int
	PrincipalVariation []string
	MultiPV            int
}

// NewEngine creates a new UCI engine instance
func NewEngine(binaryPath string) (*Engine, error) {
	cmd := exec.Command(binaryPath)
	
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %v", err)
	}
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start engine: %v", err)
	}
	
	engine := &Engine{
		cmd:     cmd,
		stdin:   stdin,
		stdout:  stdout,
		scanner: bufio.NewScanner(stdout),
	}
	
	return engine, nil
}

// Initialize sets up the engine for UCI communication
func (e *Engine) Initialize() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	
	// Send UCI command
	if err := e.sendCommand("uci"); err != nil {
		return err
	}
	
	// Wait for uciok
	for e.scanner.Scan() {
		line := strings.TrimSpace(e.scanner.Text())
		if line == "uciok" {
			break
		}
	}
	
	// Send isready command
	if err := e.sendCommand("isready"); err != nil {
		return err
	}
	
	// Wait for readyok
	for e.scanner.Scan() {
		line := strings.TrimSpace(e.scanner.Text())
		if line == "readyok" {
			e.ready = true
			break
		}
	}
	
	return nil
}

// SetOption sets a UCI option
func (e *Engine) SetOption(name, value string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	
	cmd := fmt.Sprintf("setoption name %s value %s", name, value)
	return e.sendCommand(cmd)
}

// NewGame prepares the engine for a new game
func (e *Engine) NewGame() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	
	return e.sendCommand("ucinewgame")
}

// SetPosition sets the current position
func (e *Engine) SetPosition(fen string, moves []string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	
	var cmd string
	if fen == "" || fen == "startpos" {
		cmd = "position startpos"
	} else {
		cmd = fmt.Sprintf("position fen %s", fen)
	}
	
	if len(moves) > 0 {
		cmd += " moves " + strings.Join(moves, " ")
	}
	
	return e.sendCommand(cmd)
}

// Search performs a search on the current position
func (e *Engine) Search(depth int, timeMs int, multiPV int) (*SearchResult, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	
	// Set MultiPV if specified
	if multiPV > 1 {
		if err := e.sendCommand(fmt.Sprintf("setoption name MultiPV value %d", multiPV)); err != nil {
			return nil, err
		}
	}
	
	// Build search command
	var searchCmd strings.Builder
	searchCmd.WriteString("go")
	
	if depth > 0 {
		searchCmd.WriteString(fmt.Sprintf(" depth %d", depth))
	}
	if timeMs > 0 {
		searchCmd.WriteString(fmt.Sprintf(" movetime %d", timeMs))
	}
	
	if err := e.sendCommand(searchCmd.String()); err != nil {
		return nil, err
	}
	
	result := &SearchResult{}
	var lastInfo map[string]interface{}
	
	// Read search output
	for e.scanner.Scan() {
		line := strings.TrimSpace(e.scanner.Text())
		
		if strings.HasPrefix(line, "info") {
			info := parseInfoLine(line)
			if info != nil {
				lastInfo = info
			}
		} else if strings.HasPrefix(line, "bestmove") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				result.BestMove = parts[1]
			}
			if len(parts) >= 4 && parts[2] == "ponder" {
				result.PonderMove = parts[3]
			}
			break
		}
	}
	
	// Convert last info to result
	if lastInfo != nil {
		if score, ok := lastInfo["score"]; ok {
			if scoreMap, ok := score.(map[string]interface{}); ok {
				if cp, ok := scoreMap["cp"]; ok {
					result.Score = cp.(int)
					result.ScoreType = "cp"
				} else if mate, ok := scoreMap["mate"]; ok {
					result.Score = mate.(int)
					result.ScoreType = "mate"
				}
			}
		}
		
		if depth, ok := lastInfo["depth"]; ok {
			result.Depth = depth.(int)
		}
		if seldepth, ok := lastInfo["seldepth"]; ok {
			result.SelDepth = seldepth.(int)
		}
		if nodes, ok := lastInfo["nodes"]; ok {
			result.Nodes = nodes.(int64)
		}
		if nps, ok := lastInfo["nps"]; ok {
			result.NodesPerSecond = nps.(int64)
		}
		if time, ok := lastInfo["time"]; ok {
			result.Time = time.(int)
		}
		if pv, ok := lastInfo["pv"]; ok {
			result.PrincipalVariation = pv.([]string)
		}
		if multipv, ok := lastInfo["multipv"]; ok {
			result.MultiPV = multipv.(int)
		}
	}
	
	return result, nil
}

// Stop stops the current search
func (e *Engine) Stop() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	
	return e.sendCommand("stop")
}

// Quit terminates the engine
func (e *Engine) Quit() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	
	if err := e.sendCommand("quit"); err != nil {
		return err
	}
	
	// Wait for process to exit
	return e.cmd.Wait()
}

// IsReady checks if the engine is ready
func (e *Engine) IsReady() bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.ready
}

// sendCommand sends a command to the engine (must be called with lock held)
func (e *Engine) sendCommand(cmd string) error {
	_, err := fmt.Fprintln(e.stdin, cmd)
	return err
}

// parseInfoLine parses a UCI info line
func parseInfoLine(line string) map[string]interface{} {
	parts := strings.Fields(line)
	if len(parts) < 2 || parts[0] != "info" {
		return nil
	}
	
	info := make(map[string]interface{})
	
	for i := 1; i < len(parts); i++ {
		switch parts[i] {
		case "depth":
			if i+1 < len(parts) {
				if val, err := strconv.Atoi(parts[i+1]); err == nil {
					info["depth"] = val
					i++
				}
			}
		case "seldepth":
			if i+1 < len(parts) {
				if val, err := strconv.Atoi(parts[i+1]); err == nil {
					info["seldepth"] = val
					i++
				}
			}
		case "time":
			if i+1 < len(parts) {
				if val, err := strconv.Atoi(parts[i+1]); err == nil {
					info["time"] = val
					i++
				}
			}
		case "nodes":
			if i+1 < len(parts) {
				if val, err := strconv.ParseInt(parts[i+1], 10, 64); err == nil {
					info["nodes"] = val
					i++
				}
			}
		case "nps":
			if i+1 < len(parts) {
				if val, err := strconv.ParseInt(parts[i+1], 10, 64); err == nil {
					info["nps"] = val
					i++
				}
			}
		case "multipv":
			if i+1 < len(parts) {
				if val, err := strconv.Atoi(parts[i+1]); err == nil {
					info["multipv"] = val
					i++
				}
			}
		case "score":
			if i+1 < len(parts) {
				scoreInfo := make(map[string]interface{})
				i++
				if parts[i] == "cp" && i+1 < len(parts) {
					if val, err := strconv.Atoi(parts[i+1]); err == nil {
						scoreInfo["cp"] = val
						i++
					}
				} else if parts[i] == "mate" && i+1 < len(parts) {
					if val, err := strconv.Atoi(parts[i+1]); err == nil {
						scoreInfo["mate"] = val
						i++
					}
				}
				info["score"] = scoreInfo
			}
		case "pv":
			// Principal variation - collect all remaining moves
			var pv []string
			for j := i + 1; j < len(parts); j++ {
				// Check if this is another UCI keyword
				if isUCIKeyword(parts[j]) {
					break
				}
				pv = append(pv, parts[j])
			}
			info["pv"] = pv
			i = len(parts) // End the loop as PV is typically the last item
		}
	}
	
	return info
}

// isUCIKeyword checks if a string is a UCI keyword
func isUCIKeyword(s string) bool {
	keywords := []string{"depth", "seldepth", "time", "nodes", "pv", "multipv", "score", "cp", "mate", "nps"}
	for _, keyword := range keywords {
		if s == keyword {
			return true
		}
	}
	return false
} 