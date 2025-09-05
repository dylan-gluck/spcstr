#!/bin/sh
set +e  # Don't exit on error - CRITICAL: Never block Claude Code

# Source common functions
. "/Users/dylan/Workspace/projects/spcstr/.spcstr/hooks/common.sh" 2>/dev/null || true

# Ensure log directory exists
ensure_log_dir

# Log file modification event
log_info "File-modified hook triggered"

# Get file path from environment if available
FILE_PATH="${CLAUDE_MODIFIED_FILE:-unknown}"
SESSION_ID="${CLAUDE_SESSION_ID:-unknown}"

# Track file modification
if [ "${SESSION_ID}" != "unknown" ] && [ "${FILE_PATH}" != "unknown" ]; then
    SESSION_FILE="/Users/dylan/Workspace/projects/spcstr/.spcstr/sessions/active/sess_${SESSION_ID}.json"
    if [ -f "${SESSION_FILE}" ]; then
        log_info "File modified: ${FILE_PATH} in session ${SESSION_ID}"
        # Could extend to track modified files in session data
    fi
fi

# Update cache if document was modified
if printf "%s" "${FILE_PATH}" | grep -q "\.md$"; then
    # Mark cache as stale for document updates
    CACHE_FILE="/Users/dylan/Workspace/projects/spcstr/.spcstr/cache/document-index.json"
    if [ -f "${CACHE_FILE}" ]; then
        touch "/Users/dylan/Workspace/projects/spcstr/.spcstr/cache/.stale" 2>/dev/null || true
        log_info "Marked document cache as stale due to: ${FILE_PATH}"
    fi
fi

# Always exit successfully
exit 0
