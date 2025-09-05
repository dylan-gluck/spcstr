package hooks

// commonFunctions contains shared shell functions for all hooks
const commonFunctions = `#!/bin/sh
# Common functions for spcstr hooks
# All hooks must exit 0 to not block Claude Code
set +e  # Don't exit on error

# Safe logging function
log_error() {
    printf "[%s] ERROR: %s\n" "$(date '+%Y-%m-%d %H:%M:%S')" "$1" >> "{{.LogFile}}" 2>/dev/null || true
}

log_info() {
    printf "[%s] INFO: %s\n" "$(date '+%Y-%m-%d %H:%M:%S')" "$1" >> "{{.LogFile}}" 2>/dev/null || true
}

# Create session data with proper escaping
create_session_data() {
    session_id="$1"
    event_type="$2"
    timestamp="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
    
    # Create JSON with printf for safety
    printf '{"sessionId":"%s","eventType":"%s","timestamp":"%s"}\n' \
        "${session_id}" "${event_type}" "${timestamp}"
}

# Ensure log directory exists
ensure_log_dir() {
    log_dir="$(dirname "{{.LogFile}}")"
    [ -d "${log_dir}" ] || mkdir -p "${log_dir}" 2>/dev/null || true
}
`

// preCommandTemplate is the template for pre-command hook
const preCommandTemplate = `#!/bin/sh
set +e  # Don't exit on error - CRITICAL: Never block Claude Code

# Source common functions
. "{{.HooksPath}}/common.sh" 2>/dev/null || true

# Ensure log directory exists
ensure_log_dir

# Log pre-command event
log_info "Pre-command hook triggered"

# Get session ID from environment if available
SESSION_ID="${CLAUDE_SESSION_ID:-unknown}"
COMMAND="${CLAUDE_COMMAND:-unknown}"

# Create session directory if needed
SESSION_DIR="{{.ProjectPath}}/sessions/active"
[ -d "${SESSION_DIR}" ] || mkdir -p "${SESSION_DIR}" 2>/dev/null || true

# Record command start
if [ "${SESSION_ID}" != "unknown" ]; then
    SESSION_FILE="${SESSION_DIR}/sess_${SESSION_ID}.json"
    if [ -f "${SESSION_FILE}" ]; then
        log_info "Command starting: ${COMMAND} for session ${SESSION_ID}"
    else
        # Create new session file
        create_session_data "${SESSION_ID}" "command_start" > "${SESSION_FILE}" 2>/dev/null || \
            log_error "Failed to create session file"
    fi
fi

# Always exit successfully
exit 0
`

// postCommandTemplate is the template for post-command hook
const postCommandTemplate = `#!/bin/sh
set +e  # Don't exit on error - CRITICAL: Never block Claude Code

# Source common functions
. "{{.HooksPath}}/common.sh" 2>/dev/null || true

# Ensure log directory exists
ensure_log_dir

# Log post-command event
log_info "Post-command hook triggered"

# Get session ID from environment if available
SESSION_ID="${CLAUDE_SESSION_ID:-unknown}"
COMMAND="${CLAUDE_COMMAND:-unknown}"
EXIT_CODE="${CLAUDE_EXIT_CODE:-0}"

# Update session data
if [ "${SESSION_ID}" != "unknown" ]; then
    SESSION_FILE="{{.ProjectPath}}/sessions/active/sess_${SESSION_ID}.json"
    if [ -f "${SESSION_FILE}" ]; then
        log_info "Command completed: ${COMMAND} with exit code ${EXIT_CODE}"
        # Update timestamp (atomic write)
        TEMP_FILE="${SESSION_FILE}.tmp"
        create_session_data "${SESSION_ID}" "command_complete" > "${TEMP_FILE}" 2>/dev/null && \
            mv "${TEMP_FILE}" "${SESSION_FILE}" 2>/dev/null || \
            log_error "Failed to update session file"
    fi
fi

# Always exit successfully
exit 0
`

// fileModifiedTemplate is the template for file-modified hook
const fileModifiedTemplate = `#!/bin/sh
set +e  # Don't exit on error - CRITICAL: Never block Claude Code

# Source common functions
. "{{.HooksPath}}/common.sh" 2>/dev/null || true

# Ensure log directory exists
ensure_log_dir

# Log file modification event
log_info "File-modified hook triggered"

# Get file path from environment if available
FILE_PATH="${CLAUDE_MODIFIED_FILE:-unknown}"
SESSION_ID="${CLAUDE_SESSION_ID:-unknown}"

# Track file modification
if [ "${SESSION_ID}" != "unknown" ] && [ "${FILE_PATH}" != "unknown" ]; then
    SESSION_FILE="{{.ProjectPath}}/sessions/active/sess_${SESSION_ID}.json"
    if [ -f "${SESSION_FILE}" ]; then
        log_info "File modified: ${FILE_PATH} in session ${SESSION_ID}"
        # Could extend to track modified files in session data
    fi
fi

# Update cache if document was modified
if printf "%s" "${FILE_PATH}" | grep -q "\.md$"; then
    # Mark cache as stale for document updates
    CACHE_FILE="{{.ProjectPath}}/cache/document-index.json"
    if [ -f "${CACHE_FILE}" ]; then
        touch "{{.ProjectPath}}/cache/.stale" 2>/dev/null || true
        log_info "Marked document cache as stale due to: ${FILE_PATH}"
    fi
fi

# Always exit successfully
exit 0
`

// sessionEndTemplate is the template for session-end hook
const sessionEndTemplate = `#!/bin/sh
set +e  # Don't exit on error - CRITICAL: Never block Claude Code

# Source common functions
. "{{.HooksPath}}/common.sh" 2>/dev/null || true

# Ensure log directory exists
ensure_log_dir

# Log session end event
log_info "Session-end hook triggered"

# Get session ID from environment if available
SESSION_ID="${CLAUDE_SESSION_ID:-unknown}"

# Archive session if it exists
if [ "${SESSION_ID}" != "unknown" ]; then
    ACTIVE_FILE="{{.ProjectPath}}/sessions/active/sess_${SESSION_ID}.json"
    if [ -f "${ACTIVE_FILE}" ]; then
        # Move to archive with timestamp
        ARCHIVE_DIR="{{.ProjectPath}}/sessions/archive"
        [ -d "${ARCHIVE_DIR}" ] || mkdir -p "${ARCHIVE_DIR}" 2>/dev/null || true
        
        TIMESTAMP="$(date '+%Y%m%d_%H%M%S')"
        ARCHIVE_FILE="${ARCHIVE_DIR}/sess_${SESSION_ID}_${TIMESTAMP}.json"
        
        mv "${ACTIVE_FILE}" "${ARCHIVE_FILE}" 2>/dev/null && \
            log_info "Archived session ${SESSION_ID}" || \
            log_error "Failed to archive session ${SESSION_ID}"
    fi
fi

# Clean up old sessions based on retention policy
# This would normally check config for retention days
# For now, just log that cleanup would occur
log_info "Session cleanup check completed"

# Always exit successfully
exit 0
`