# Security

**⚠️ MANDATORY:** These security requirements must be followed by all AI and human developers.

## Input Validation
- **Validation Library:** Go standard library with custom validators
- **Validation Location:** At CLI command entry points before processing
- **Required Rules:**
  - All external inputs MUST be validated
  - Validation at API boundary before processing
  - Whitelist approach preferred over blacklist

## Authentication & Authorization
- **Auth Method:** Local file system permissions only (no user auth needed)
- **Session Management:** File-based sessions with unique IDs
- **Required Patterns:**
  - Sessions identified by UUID (cryptographically secure)
  - No session sharing between projects

## Secrets Management
- **Development:** No secrets required for core functionality
- **Production:** No secrets required for core functionality
- **Code Requirements:**
  - NEVER hardcode secrets
  - Access via configuration service only
  - No secrets in logs or error messages

## API Security
- **Rate Limiting:** N/A - Local application only
- **CORS Policy:** N/A - No web interface
- **Security Headers:** N/A - Terminal application
- **HTTPS Enforcement:** N/A - No network communication

## Data Protection
- **Encryption at Rest:** Rely on OS file system encryption
- **Encryption in Transit:** N/A - No network transfer
- **PII Handling:** Never log file contents, only paths
- **Logging Restrictions:** No session content, no file contents, no user paths in logs

## Dependency Security
- **Scanning Tool:** dependabot via GitHub
- **Update Policy:** Monthly review of updates
- **Approval Process:** Security updates applied immediately, features reviewed

## Security Testing
- **SAST Tool:** gosec in CI pipeline
- **DAST Tool:** N/A - Not applicable for CLI tool
- **Penetration Testing:** N/A - Local tool only

## File System Security

**Critical Requirements:**
- **Directory Permissions:** .spcstr/ created with 0755 permissions
- **File Permissions:** Session files created with 0644 permissions
- **Path Traversal Prevention:** All paths sanitized with filepath.Clean()
- **Symlink Handling:** Do not follow symlinks outside project root
- **Temp File Security:** Use os.CreateTemp() with proper permissions

## Shell Script Security

**Hook Script Requirements:**
- **Input Sanitization:** Quote all variables in shell scripts
- **Command Injection Prevention:** Never use eval or backticks
- **Environment Variables:** Clear sensitive vars before execution
- **Script Permissions:** Hook scripts created with 0755 permissions

Example hook invocation:
```bash
# Claude settings.json configuration
"PreToolUse": "spcstr hook pre-tool-use --cwd=$CLAUDE_PROJECT_DIR"

# Hook receives JSON on stdin, returns exit code:
# 0 = success, continue
# 2 = blocking operation detected
```

## Terminal Security

**TUI Security Measures:**
- **Escape Sequence Sanitization:** Strip ANSI codes from user input
- **Buffer Overflow Prevention:** Limit input field lengths
- **Screen Scraping Protection:** No sensitive data displayed persistently
- **Clipboard Security:** No automatic clipboard operations

## Supply Chain Security

**Build Security:**
- **Reproducible Builds:** Vendor dependencies with go mod vendor
- **Build Environment:** Clean CI environment for releases
- **Binary Signing:** macOS code signing for notarization
- **Checksum Verification:** SHA256 checksums for all releases
