#!/bin/sh
set +e  # Don't exit on error - CRITICAL: Never block Claude Code

# Source common functions
. "/Users/dylan/Workspace/projects/spcstr/.spcstr/hooks/common.sh" 2>/dev/null || true

# Ensure log directory exists
ensure_log_dir

# Log pre-command event
log_info "Pre-command hook triggered"

# Get session ID from environment if available
SESSION_ID="${CLAUDE_SESSION_ID:-unknown}"
COMMAND="${CLAUDE_COMMAND:-unknown}"

# Create session directory if needed
SESSION_DIR="/Users/dylan/Workspace/projects/spcstr/.spcstr/sessions/active"
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
