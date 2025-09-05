#!/bin/sh
set +e  # Don't exit on error - CRITICAL: Never block Claude Code

# Source common functions
. "/Users/dylan/Workspace/projects/spcstr/.spcstr/hooks/common.sh" 2>/dev/null || true

# Ensure log directory exists
ensure_log_dir

# Log session end event
log_info "Session-end hook triggered"

# Get session ID from environment if available
SESSION_ID="${CLAUDE_SESSION_ID:-unknown}"

# Archive session if it exists
if [ "${SESSION_ID}" != "unknown" ]; then
    ACTIVE_FILE="/Users/dylan/Workspace/projects/spcstr/.spcstr/sessions/active/sess_${SESSION_ID}.json"
    if [ -f "${ACTIVE_FILE}" ]; then
        # Move to archive with timestamp
        ARCHIVE_DIR="/Users/dylan/Workspace/projects/spcstr/.spcstr/sessions/archive"
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
