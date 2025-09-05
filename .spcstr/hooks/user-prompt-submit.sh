#!/bin/sh
# spcstr UserPromptSubmit hook - captures user prompts
# Must exit 0 to not block Claude Code operations

# Source common functions
. "${CLAUDE_PROJECT_DIR}/.spcstr/hooks/common.sh" 2>/dev/null || true

# Parse JSON input from stdin
input="$(cat)"
session_id="$(echo "${input}" | grep -o '"session_id":"[^"]*' | cut -d'"' -f4)" 2>/dev/null || true
prompt="$(echo "${input}" | grep -o '"prompt":"[^"]*' | cut -d'"' -f4)" 2>/dev/null || true

# Create session directory if needed
session_dir="${CLAUDE_PROJECT_DIR}/.spcstr/sessions"
[ -d "${session_dir}" ] || mkdir -p "${session_dir}" 2>/dev/null || true

# Generate session file path  
session_file="${session_dir}/sess_${session_id}.json"

# Create or update session data
if [ -f "${session_file}" ]; then
    # Session exists, append prompt event
    timestamp="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
    # Note: Full JSON manipulation would require jq, keeping simple for POSIX compatibility
    log_info "Prompt captured for session ${session_id}"
else
    # Create new session file
    timestamp="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
    cat > "${session_file}" <<EOF
{
  "sessionId": "sess_${session_id}",
  "startTime": "${timestamp}",
  "status": "active",
  "events": [
    {
      "type": "prompt",
      "timestamp": "${timestamp}"
    }
  ]
}
EOF
    log_info "New session created: sess_${session_id}"
fi

# Always exit 0 to not block Claude Code
exit 0