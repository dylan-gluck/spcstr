package hooks

import (
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplates_Valid(t *testing.T) {
	// Test data for template execution
	data := struct {
		ProjectPath string
		HooksPath   string
		LogFile     string
	}{
		ProjectPath: "/test/project",
		HooksPath:   "/test/project/hooks",
		LogFile:     "/test/project/logs/hook-errors.log",
	}
	
	templates := []struct {
		name     string
		template string
	}{
		{"commonFunctions", commonFunctions},
		{"preCommandTemplate", preCommandTemplate},
		{"postCommandTemplate", postCommandTemplate},
		{"fileModifiedTemplate", fileModifiedTemplate},
		{"sessionEndTemplate", sessionEndTemplate},
	}
	
	for _, tt := range templates {
		t.Run(tt.name, func(t *testing.T) {
			// Parse template
			tmpl, err := template.New(tt.name).Parse(tt.template)
			require.NoError(t, err, "Template %s should parse without error", tt.name)
			
			// Execute template
			var buf strings.Builder
			err = tmpl.Execute(&buf, data)
			assert.NoError(t, err, "Template %s should execute without error", tt.name)
			
			result := buf.String()
			
			// Verify common requirements
			assert.True(t, strings.HasPrefix(result, "#!/bin/sh\n"), 
				"Template %s should start with shebang", tt.name)
			
			assert.Contains(t, result, "set +e", 
				"Template %s should contain 'set +e'", tt.name)
			
			// All except common functions should exit 0
			if tt.name != "commonFunctions" {
				assert.Contains(t, result, "exit 0", 
					"Template %s should exit 0", tt.name)
			}
		})
	}
}

func TestPreCommandTemplate_Content(t *testing.T) {
	data := struct {
		ProjectPath string
		HooksPath   string
		LogFile     string
	}{
		ProjectPath: "/test/project",
		HooksPath:   "/test/project/hooks",
		LogFile:     "/test/project/logs/hook-errors.log",
	}
	
	tmpl, err := template.New("pre").Parse(preCommandTemplate)
	require.NoError(t, err)
	
	var buf strings.Builder
	err = tmpl.Execute(&buf, data)
	require.NoError(t, err)
	
	result := buf.String()
	
	// Check specific content
	assert.Contains(t, result, "Pre-command hook triggered")
	assert.Contains(t, result, "CLAUDE_SESSION_ID")
	assert.Contains(t, result, "CLAUDE_COMMAND")
	assert.Contains(t, result, "SESSION_DIR=")
	assert.Contains(t, result, "/test/project/sessions/active")
	assert.Contains(t, result, "sess_${SESSION_ID}.json")
	assert.Contains(t, result, "create_session_data")
	assert.Contains(t, result, "ensure_log_dir")
}

func TestPostCommandTemplate_Content(t *testing.T) {
	data := struct {
		ProjectPath string
		HooksPath   string
		LogFile     string
	}{
		ProjectPath: "/test/project",
		HooksPath:   "/test/project/hooks",
		LogFile:     "/test/project/logs/hook-errors.log",
	}
	
	tmpl, err := template.New("post").Parse(postCommandTemplate)
	require.NoError(t, err)
	
	var buf strings.Builder
	err = tmpl.Execute(&buf, data)
	require.NoError(t, err)
	
	result := buf.String()
	
	// Check specific content
	assert.Contains(t, result, "Post-command hook triggered")
	assert.Contains(t, result, "CLAUDE_SESSION_ID")
	assert.Contains(t, result, "CLAUDE_COMMAND")
	assert.Contains(t, result, "CLAUDE_EXIT_CODE")
	assert.Contains(t, result, "command_complete")
	assert.Contains(t, result, "/test/project/sessions/active/sess_${SESSION_ID}.json")
}

func TestFileModifiedTemplate_Content(t *testing.T) {
	data := struct {
		ProjectPath string
		HooksPath   string
		LogFile     string
	}{
		ProjectPath: "/test/project",
		HooksPath:   "/test/project/hooks",
		LogFile:     "/test/project/logs/hook-errors.log",
	}
	
	tmpl, err := template.New("file").Parse(fileModifiedTemplate)
	require.NoError(t, err)
	
	var buf strings.Builder
	err = tmpl.Execute(&buf, data)
	require.NoError(t, err)
	
	result := buf.String()
	
	// Check specific content
	assert.Contains(t, result, "File-modified hook triggered")
	assert.Contains(t, result, "CLAUDE_MODIFIED_FILE")
	assert.Contains(t, result, "CLAUDE_SESSION_ID")
	assert.Contains(t, result, "grep -q \"\\.md$\"")
	assert.Contains(t, result, "/test/project/cache/document-index.json")
	assert.Contains(t, result, "/test/project/cache/.stale")
}

func TestSessionEndTemplate_Content(t *testing.T) {
	data := struct {
		ProjectPath string
		HooksPath   string
		LogFile     string
	}{
		ProjectPath: "/test/project",
		HooksPath:   "/test/project/hooks",
		LogFile:     "/test/project/logs/hook-errors.log",
	}
	
	tmpl, err := template.New("session").Parse(sessionEndTemplate)
	require.NoError(t, err)
	
	var buf strings.Builder
	err = tmpl.Execute(&buf, data)
	require.NoError(t, err)
	
	result := buf.String()
	
	// Check specific content
	assert.Contains(t, result, "Session-end hook triggered")
	assert.Contains(t, result, "CLAUDE_SESSION_ID")
	assert.Contains(t, result, "ACTIVE_FILE=")
	assert.Contains(t, result, "ARCHIVE_DIR=")
	assert.Contains(t, result, "/test/project/sessions/active/sess_${SESSION_ID}.json")
	assert.Contains(t, result, "/test/project/sessions/archive")
	assert.Contains(t, result, "TIMESTAMP=")
	assert.Contains(t, result, "mv \"${ACTIVE_FILE}\" \"${ARCHIVE_FILE}\"")
}

func TestCommonFunctions_Content(t *testing.T) {
	data := struct {
		ProjectPath string
		HooksPath   string
		LogFile     string
	}{
		ProjectPath: "/test/project",
		HooksPath:   "/test/project/hooks",
		LogFile:     "/test/project/logs/hook-errors.log",
	}
	
	tmpl, err := template.New("common").Parse(commonFunctions)
	require.NoError(t, err)
	
	var buf strings.Builder
	err = tmpl.Execute(&buf, data)
	require.NoError(t, err)
	
	result := buf.String()
	
	// Check function definitions
	assert.Contains(t, result, "log_error()")
	assert.Contains(t, result, "log_info()")
	assert.Contains(t, result, "create_session_data()")
	assert.Contains(t, result, "ensure_log_dir()")
	
	// Check that log file path is properly substituted
	assert.Contains(t, result, "/test/project/logs/hook-errors.log")
	
	// Check JSON creation
	assert.Contains(t, result, `printf '{"sessionId":"%s","eventType":"%s","timestamp":"%s"}\n'`)
	
	// Check date formatting
	assert.Contains(t, result, "date -u '+%Y-%m-%dT%H:%M:%SZ'")
	assert.Contains(t, result, "date '+%Y-%m-%d %H:%M:%S'")
}

func TestTemplates_POSIXCompliance(t *testing.T) {
	// These patterns should NOT appear in POSIX-compliant scripts
	bashisms := []string{
		"[[",      // bash test
		"]]",      // bash test
		"function ", // bash function declaration
		"source ",  // bash source (should use .)
		"<<<",      // bash here-string
		"&>",       // bash redirect
		"==",       // bash string comparison in test
		"local ",   // bash local variables (in templates, not common.sh)
		"${BASH",   // bash-specific variables
		"$RANDOM",  // bash random
		"${!",      // bash indirect expansion
		"${#",      // bash array length (except for string length)
		"+=",       // bash append
		"((",       // bash arithmetic
		"))",       // bash arithmetic
	}
	
	templates := []struct {
		name     string
		template string
	}{
		{"preCommandTemplate", preCommandTemplate},
		{"postCommandTemplate", postCommandTemplate},
		{"fileModifiedTemplate", fileModifiedTemplate},
		{"sessionEndTemplate", sessionEndTemplate},
	}
	
	for _, tmpl := range templates {
		t.Run(tmpl.name, func(t *testing.T) {
			for _, pattern := range bashisms {
				assert.NotContains(t, tmpl.template, pattern, 
					"Template %s should not contain bashism: %s", tmpl.name, pattern)
			}
			
			// Check for POSIX-compliant constructs
			if strings.Contains(tmpl.template, "if ") {
				assert.Contains(t, tmpl.template, "if [", 
					"Template %s should use POSIX [ ] for tests", tmpl.name)
			}
			
			// Check for proper variable quoting
			assert.NotContains(t, tmpl.template, "$SESSION_ID\"", 
				"Variables should be properly quoted with curly braces")
		})
	}
}