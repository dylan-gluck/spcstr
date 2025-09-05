#!/bin/sh
# Common functions for spcstr hooks
# These functions are sourced by all hook scripts

# Set safe defaults
set +e  # Don't exit on error

# Logging functions with proper escaping
log_error() {
    log_file="${1:-${LOG_FILE}}"
    shift
    printf "[%s] ERROR: %s\n" "$(date '+%Y-%m-%d %H:%M:%S')" "$*" >> "${log_file}" 2>/dev/null || true
}

log_info() {
    log_file="${1:-${LOG_FILE}}"
    shift
    printf "[%s] INFO: %s\n" "$(date '+%Y-%m-%d %H:%M:%S')" "$*" >> "${log_file}" 2>/dev/null || true
}

log_debug() {
    log_file="${1:-${LOG_FILE}}"
    shift
    printf "[%s] DEBUG: %s\n" "$(date '+%Y-%m-%d %H:%M:%S')" "$*" >> "${log_file}" 2>/dev/null || true
}

# Ensure directory exists (with error handling)
ensure_dir() {
    dir_path="$1"
    if [ -z "${dir_path}" ]; then
        return 1
    fi
    
    if [ ! -d "${dir_path}" ]; then
        mkdir -p "${dir_path}" 2>/dev/null || return 1
    fi
    return 0
}

# Safe JSON field extraction (basic)
# Usage: json_get_field "$json_string" "fieldname"
json_get_field() {
    json="$1"
    field="$2"
    
    # Basic extraction using sed (POSIX compliant)
    printf "%s" "${json}" | sed -n "s/.*\"${field}\":\"\([^\"]*\)\".*/\1/p"
}

# Create simple JSON object
# Usage: json_create "key1" "value1" "key2" "value2" ...
json_create() {
    result="{"
    first=true
    
    while [ $# -ge 2 ]; do
        key="$1"
        value="$2"
        shift 2
        
        if [ "${first}" = true ]; then
            first=false
        else
            result="${result},"
        fi
        
        # Escape special characters in value
        escaped_value="$(printf "%s" "${value}" | sed 's/\\/\\\\/g; s/"/\\"/g')"
        result="${result}\"${key}\":\"${escaped_value}\""
    done
    
    result="${result}}"
    printf "%s" "${result}"
}

# Get current ISO timestamp
get_timestamp() {
    date -u '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null || date '+%Y-%m-%d %H:%M:%S'
}

# Archive file with timestamp
archive_file() {
    src_file="$1"
    archive_dir="$2"
    
    if [ ! -f "${src_file}" ]; then
        return 1
    fi
    
    ensure_dir "${archive_dir}" || return 1
    
    basename="$(basename "${src_file}")"
    timestamp="$(date '+%Y%m%d_%H%M%S')"
    archive_path="${archive_dir}/${basename}.${timestamp}"
    
    mv "${src_file}" "${archive_path}" 2>/dev/null || cp "${src_file}" "${archive_path}" 2>/dev/null
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Get file age in days
file_age_days() {
    file="$1"
    
    if [ ! -f "${file}" ]; then
        printf "0"
        return 1
    fi
    
    # Try to use stat command (varies by OS)
    if command_exists stat; then
        # Try GNU stat first
        mod_time="$(stat -c %Y "${file}" 2>/dev/null)"
        if [ -z "${mod_time}" ]; then
            # Try BSD/macOS stat
            mod_time="$(stat -f %m "${file}" 2>/dev/null)"
        fi
        
        if [ -n "${mod_time}" ]; then
            current_time="$(date +%s)"
            age_seconds=$((current_time - mod_time))
            age_days=$((age_seconds / 86400))
            printf "%d" "${age_days}"
            return 0
        fi
    fi
    
    # Fallback: can't determine age
    printf "0"
    return 1
}

# Atomic file write
# Usage: atomic_write "file_path" "content"
atomic_write() {
    file_path="$1"
    content="$2"
    
    temp_file="${file_path}.tmp.$$"
    
    # Write to temporary file
    printf "%s" "${content}" > "${temp_file}" 2>/dev/null || {
        rm -f "${temp_file}" 2>/dev/null
        return 1
    }
    
    # Atomic rename
    mv "${temp_file}" "${file_path}" 2>/dev/null || {
        rm -f "${temp_file}" 2>/dev/null
        return 1
    }
    
    return 0
}

# Export functions for use in hooks
export -f log_error log_info log_debug ensure_dir json_get_field json_create
export -f get_timestamp archive_file command_exists file_age_days atomic_write