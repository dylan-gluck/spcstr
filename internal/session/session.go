package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	models "github.com/dylan-gluck/spcstr/pkg/hooks"
)

var (
	stateMutex sync.RWMutex
	stateCache *models.SessionState
	cacheTime  time.Time
)

const cacheTTL = 100 * time.Millisecond

func LoadSessionState(projectRoot string) (*models.SessionState, error) {
	stateMutex.RLock()
	if stateCache != nil && time.Since(cacheTime) < cacheTTL {
		cachedState := *stateCache
		stateMutex.RUnlock()
		return &cachedState, nil
	}
	stateMutex.RUnlock()

	sessionFile := filepath.Join(projectRoot, ".spcstr", "session_state.json")

	data, err := os.ReadFile(sessionFile)
	if err != nil {
		if os.IsNotExist(err) {
			state := &models.SessionState{
				SessionID:     GenerateSessionID(),
				Source:        "startup",
				ProjectPath:   projectRoot,
				Timestamp:     models.CurrentTimestamp(),
				LastUpdate:    models.CurrentTimestamp(),
				Status:        "active",
				Agents:        []string{},
				AgentsHistory: []models.AgentHistoryEntry{},
				Files: models.FileOperations{
					New:    []string{},
					Edited: []string{},
					Read:   []string{},
				},
				ToolsUsed: make(map[string]int),
				Errors:    []models.ErrorEntry{},
				Modified:  false,
			}
			return state, nil
		}
		return nil, fmt.Errorf("failed to read session state: %w", err)
	}

	var state models.SessionState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse session state: %w", err)
	}

	if state.ToolsUsed == nil {
		state.ToolsUsed = make(map[string]int)
	}
	if state.Agents == nil {
		state.Agents = []string{}
	}
	if state.AgentsHistory == nil {
		state.AgentsHistory = []models.AgentHistoryEntry{}
	}
	if state.Files.New == nil {
		state.Files.New = []string{}
	}
	if state.Files.Edited == nil {
		state.Files.Edited = []string{}
	}
	if state.Files.Read == nil {
		state.Files.Read = []string{}
	}
	if state.Errors == nil {
		state.Errors = []models.ErrorEntry{}
	}

	stateMutex.Lock()
	stateCache = &state
	cacheTime = time.Now()
	stateMutex.Unlock()

	return &state, nil
}

func SaveSessionState(projectRoot string, state *models.SessionState) error {
	sessionDir := filepath.Join(projectRoot, ".spcstr")
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return fmt.Errorf("failed to create session directory: %w", err)
	}

	state.LastUpdate = models.CurrentTimestamp()
	state.Modified = true

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session state: %w", err)
	}

	sessionFile := filepath.Join(sessionDir, "session_state.json")
	tempFile := sessionFile + ".tmp"

	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := os.Rename(tempFile, sessionFile); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	stateMutex.Lock()
	stateCache = state
	cacheTime = time.Now()
	stateMutex.Unlock()

	return nil
}

func GenerateSessionID() string {
	return fmt.Sprintf("sess_%d", time.Now().UnixNano())
}
