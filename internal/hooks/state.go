package hooks

import (
	"github.com/dylan-gluck/spcstr/internal/session"
)

// Expose session functions for backward compatibility
var LoadSessionState = session.LoadSessionState
var SaveSessionState = session.SaveSessionState
var generateSessionID = session.GenerateSessionID
