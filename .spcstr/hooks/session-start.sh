#!/bin/sh
# spcstr SessionStart hook - initializes session tracking
# Must exit 0 to not block Claude Code operations

# Source common functions  
. "${CLAUDE_PROJECT_DIR}/.spcstr/hooks/common.sh" 2>/dev/null || true

# Parse JSON input from stdin
input="$(cat)"
session_id="$(echo "${input}" | grep -o '"session_id":"[^"]*' | cut -d'"' -f4)" 2>/dev/null || true
source="$(echo "${input}" | grep -o '"source":"[^"]*' | cut -d'"' -f4)" 2>/dev/null || true

# Create session directory if needed
session_dir="${CLAUDE_PROJECT_DIR}/.spcstr/sessions"
[ -d "${session_dir}" ] || mkdir -p "${session_dir}" 2>/dev/null || true

# Generate session file path
session_file="${session_dir}/sess_${session_id}.json"

# Create session initialization data
timestamp="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
cat > "${session_file}" <<EOF
{
  "sessionId": "sess_${session_id}",
  "startTime": "${timestamp}",
  "source": "${source}",
  "status": "active",
  "agents": [],
  "tasks": [],
  "files": [],
  "tools": [],
  "errors": [],
  "events": [
    {
      "type": "session_start",
      "timestamp": "${timestamp}",
      "source": "${source}"
    }
  ]
}
EOF

log_info "Session initialized: sess_${session_id} (source: ${source})"

# Always exit 0 to not block Claude Code
exit 0