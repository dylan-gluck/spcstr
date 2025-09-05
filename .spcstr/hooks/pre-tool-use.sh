#!/bin/sh
# spcstr PreToolUse hook - tracks tool usage before execution
# Must exit 0 to not block Claude Code operations

# Source common functions
. "${CLAUDE_PROJECT_DIR}/.spcstr/hooks/common.sh" 2>/dev/null || true

# Parse JSON input from stdin
input="$(cat)"
session_id="$(echo "${input}" | grep -o '"session_id":"[^"]*' | cut -d'"' -f4)" 2>/dev/null || true
tool_name="$(echo "${input}" | grep -o '"tool_name":"[^"]*' | cut -d'"' -f4)" 2>/dev/null || true

# Session file path
session_dir="${CLAUDE_PROJECT_DIR}/.spcstr/sessions"
session_file="${session_dir}/sess_${session_id}.json"

# Log tool usage event (simplified for POSIX sh)
if [ -f "${session_file}" ]; then
    timestamp="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
    # Track tool about to be used
    log_info "Tool '${tool_name}' preparing to execute in session ${session_id}"
    
    # For file operations, extract the file path
    if echo "${tool_name}" | grep -qE "Edit|Write|MultiEdit|Read"; then
        file_path="$(echo "${input}" | grep -o '"file_path":"[^"]*' | cut -d'"' -f4)" 2>/dev/null || true
        [ -n "${file_path}" ] && log_info "File operation on: ${file_path}"
    fi
fi

# Always exit 0 to not block Claude Code
exit 0