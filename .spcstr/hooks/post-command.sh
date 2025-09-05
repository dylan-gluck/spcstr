#!/bin/sh
set +e  # Don't exit on error - CRITICAL: Never block Claude Code

# Source common functions
. "/Users/dylan/Workspace/projects/spcstr/.spcstr/hooks/common.sh" 2>/dev/null || true

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
    SESSION_FILE="/Users/dylan/Workspace/projects/spcstr/.spcstr/sessions/active/sess_${SESSION_ID}.json"
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
