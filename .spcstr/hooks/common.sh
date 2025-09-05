#!/bin/sh
# Common functions for spcstr hooks
# All hooks must exit 0 to not block Claude Code
set +e  # Don't exit on error

# Safe logging function
log_error() {
    printf "[%s] ERROR: %s\n" "$(date '+%Y-%m-%d %H:%M:%S')" "$1" >> "<no value>" 2>/dev/null || true
}

log_info() {
    printf "[%s] INFO: %s\n" "$(date '+%Y-%m-%d %H:%M:%S')" "$1" >> "<no value>" 2>/dev/null || true
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
    log_dir="$(dirname "<no value>")"
    [ -d "${log_dir}" ] || mkdir -p "${log_dir}" 2>/dev/null || true
}
