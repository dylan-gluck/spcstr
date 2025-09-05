#!/bin/sh
# spcstr PostToolUse hook - tracks tool results after execution
# Must exit 0 to not block Claude Code operations

# Source common functions
. "${CLAUDE_PROJECT_DIR}/.spcstr/hooks/common.sh" 2>/dev/null || true

# Parse JSON input from stdin
input="$(cat)"
session_id="$(echo "${input}" | grep -o '"session_id":"[^"]*' | cut -d'"' -f4)" 2>/dev/null || true
tool_name="$(echo "${input}" | grep -o '"tool_name":"[^"]*' | cut -d'"' -f4)" 2>/dev/null || true

# Check for tool response success
success="$(echo "${input}" | grep -o '"success":[^,}]*' | cut -d':' -f2 | tr -d ' ')" 2>/dev/null || true

# Session file path
session_dir="${CLAUDE_PROJECT_DIR}/.spcstr/sessions"
session_file="${session_dir}/sess_${session_id}.json"

# Log tool completion event
if [ -f "${session_file}" ]; then
    timestamp="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
    
    # Log based on success status
    if [ "${success}" = "true" ]; then
        log_info "Tool '${tool_name}' completed successfully in session ${session_id}"
    elif [ "${success}" = "false" ]; then
        log_error "Tool '${tool_name}' failed in session ${session_id}"
    else
        log_info "Tool '${tool_name}' completed in session ${session_id}"
    fi
    
    # For file operations, track the modified file
    if echo "${tool_name}" | grep -qE "Edit|Write|MultiEdit"; then
        file_path="$(echo "${input}" | grep -o '"file_path":"[^"]*' | cut -d'"' -f4)" 2>/dev/null || true
        if [ -n "${file_path}" ]; then
            log_info "File modified: ${file_path}"
            
            # Trigger file watcher notification (non-blocking)
            {
                echo "FILE_MODIFIED:${file_path}" >> "${session_dir}/.events" 2>/dev/null || true
            } &
        fi
    fi
fi

# Always exit 0 to not block Claude Code
exit 0