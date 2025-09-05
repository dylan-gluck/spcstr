#!/bin/sh
# spcstr SessionEnd hook - finalizes session tracking
# Must exit 0 to not block Claude Code operations

# Source common functions
. "${CLAUDE_PROJECT_DIR}/.spcstr/hooks/common.sh" 2>/dev/null || true

# Parse JSON input from stdin
input="$(cat)"
session_id="$(echo "${input}" | grep -o '"session_id":"[^"]*' | cut -d'"' -f4)" 2>/dev/null || true
reason="$(echo "${input}" | grep -o '"reason":"[^"]*' | cut -d'"' -f4)" 2>/dev/null || true

# Session file path
session_dir="${CLAUDE_PROJECT_DIR}/.spcstr/sessions"
session_file="${session_dir}/sess_${session_id}.json"

# Mark session as ended if file exists
if [ -f "${session_file}" ]; then
    timestamp="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
    
    # Update status to ended (simplified for POSIX sh)
    temp_file="${session_file}.tmp"
    if sed "s/\"status\": \"active\"/\"status\": \"ended\"/g" "${session_file}" > "${temp_file}" 2>/dev/null; then
        mv "${temp_file}" "${session_file}" 2>/dev/null || true
    fi
    
    # Log end event
    log_info "Session ended: sess_${session_id} (reason: ${reason:-unknown})"
    
    # Archive session to history if configured
    history_dir="${session_dir}/history"
    if [ -d "${history_dir}" ] || mkdir -p "${history_dir}" 2>/dev/null; then
        archive_name="$(date '+%Y%m%d_%H%M%S')_sess_${session_id}.json"
        cp "${session_file}" "${history_dir}/${archive_name}" 2>/dev/null || true
    fi
fi

# Always exit 0 to not block Claude Code
exit 0