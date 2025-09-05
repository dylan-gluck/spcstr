package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Configuration
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid default config",
			cfg:     DefaultConfig(),
			wantErr: false,
		},
		{
			name:    "nil config",
			cfg:     nil,
			wantErr: true,
			errMsg:  "configuration is nil",
		},
		{
			name: "empty version",
			cfg: &Configuration{
				Version: "",
				Scope:   "project",
				Paths:   DefaultConfig().Paths,
				UI:      DefaultConfig().UI,
				Session: DefaultConfig().Session,
			},
			wantErr: true,
			errMsg:  "version is required",
		},
		{
			name: "invalid version format",
			cfg: &Configuration{
				Version: "1.0",
				Scope:   "project",
				Paths:   DefaultConfig().Paths,
				UI:      DefaultConfig().UI,
				Session: DefaultConfig().Session,
			},
			wantErr: true,
			errMsg:  "invalid version format",
		},
		{
			name: "invalid scope",
			cfg: &Configuration{
				Version: "1.0.0",
				Scope:   "invalid",
				Paths:   DefaultConfig().Paths,
				UI:      DefaultConfig().UI,
				Session: DefaultConfig().Session,
			},
			wantErr: true,
			errMsg:  "scope must be either 'project' or 'global'",
		},
		{
			name: "invalid theme",
			cfg: &Configuration{
				Version: "1.0.0",
				Scope:   "project",
				Paths:   DefaultConfig().Paths,
				UI: UIConfig{
					Theme:           "invalid",
					RefreshInterval: 500,
				},
				Session: DefaultConfig().Session,
			},
			wantErr: true,
			errMsg:  "invalid theme",
		},
		{
			name: "refresh interval too low",
			cfg: &Configuration{
				Version: "1.0.0",
				Scope:   "project",
				Paths:   DefaultConfig().Paths,
				UI: UIConfig{
					Theme:           "dark",
					RefreshInterval: 50,
				},
				Session: DefaultConfig().Session,
			},
			wantErr: true,
			errMsg:  "refresh interval must be at least 100ms",
		},
		{
			name: "refresh interval too high",
			cfg: &Configuration{
				Version: "1.0.0",
				Scope:   "project",
				Paths:   DefaultConfig().Paths,
				UI: UIConfig{
					Theme:           "dark",
					RefreshInterval: 20000,
				},
				Session: DefaultConfig().Session,
			},
			wantErr: true,
			errMsg:  "refresh interval must be at most 10000ms",
		},
		{
			name: "negative retention days",
			cfg: &Configuration{
				Version: "1.0.0",
				Scope:   "project",
				Paths:   DefaultConfig().Paths,
				UI:      DefaultConfig().UI,
				Session: SessionConfig{
					RetentionDays:    -1,
					AutoArchive:      true,
					MaxActiveSession: 10,
				},
			},
			wantErr: true,
			errMsg:  "retention days cannot be negative",
		},
		{
			name: "retention days too high",
			cfg: &Configuration{
				Version: "1.0.0",
				Scope:   "project",
				Paths:   DefaultConfig().Paths,
				UI:      DefaultConfig().UI,
				Session: SessionConfig{
					RetentionDays:    400,
					AutoArchive:      true,
					MaxActiveSession: 10,
				},
			},
			wantErr: true,
			errMsg:  "retention days cannot exceed 365",
		},
		{
			name: "max active sessions too low",
			cfg: &Configuration{
				Version: "1.0.0",
				Scope:   "project",
				Paths:   DefaultConfig().Paths,
				UI:      DefaultConfig().UI,
				Session: SessionConfig{
					RetentionDays:    30,
					AutoArchive:      true,
					MaxActiveSession: 0,
				},
			},
			wantErr: true,
			errMsg:  "max active sessions must be at least 1",
		},
		{
			name: "max active sessions too high",
			cfg: &Configuration{
				Version: "1.0.0",
				Scope:   "project",
				Paths:   DefaultConfig().Paths,
				UI:      DefaultConfig().UI,
				Session: SessionConfig{
					RetentionDays:    30,
					AutoArchive:      true,
					MaxActiveSession: 200,
				},
			},
			wantErr: true,
			errMsg:  "max active sessions cannot exceed 100",
		},
		{
			name: "empty hooks path",
			cfg: &Configuration{
				Version: "1.0.0",
				Scope:   "project",
				Paths: PathConfig{
					Hooks:    "",
					Sessions: ".spcstr/sessions",
					Docs:     []string{"docs/"},
				},
				UI:      DefaultConfig().UI,
				Session: DefaultConfig().Session,
			},
			wantErr: true,
			errMsg:  "hooks path is required",
		},
		{
			name: "path with parent directory traversal",
			cfg: &Configuration{
				Version: "1.0.0",
				Scope:   "project",
				Paths: PathConfig{
					Hooks:    "../hooks",
					Sessions: ".spcstr/sessions",
					Docs:     []string{"docs/"},
				},
				UI:      DefaultConfig().UI,
				Session: DefaultConfig().Session,
			},
			wantErr: true,
			errMsg:  "invalid hooks path",
		},
		{
			name: "no docs paths",
			cfg: &Configuration{
				Version: "1.0.0",
				Scope:   "project",
				Paths: PathConfig{
					Hooks:    ".spcstr/hooks",
					Sessions: ".spcstr/sessions",
					Docs:     []string{},
				},
				UI:      DefaultConfig().UI,
				Session: DefaultConfig().Session,
			},
			wantErr: true,
			errMsg:  "at least one docs path is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.cfg)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsValidVersion(t *testing.T) {
	tests := []struct {
		version string
		valid   bool
	}{
		{"1.0.0", true},
		{"0.1.0", true},
		{"10.20.30", true},
		{"1.0", false},
		{"1", false},
		{"1.0.0.0", false},
		{"", false},
		{"1.a.0", false},
		{"1.0.0-beta", false},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			result := isValidVersion(tt.version)
			assert.Equal(t, tt.valid, result)
		})
	}
}

func TestIsValidPath(t *testing.T) {
	tests := []struct {
		path  string
		valid bool
	}{
		{".spcstr/hooks", true},
		{"hooks", true},
		{"/absolute/path", true},
		{"relative/path", true},
		{"path/", true},
		{"../parent", false},
		{"path/../other", false},
		{"", false},
		{".", true},
		{"./subdir", true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := isValidPath(tt.path)
			assert.Equal(t, tt.valid, result)
		})
	}
}

func TestIsValidConfigPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid path",
			path:    ".spcstr/config.json",
			wantErr: false,
		},
		{
			name:    "valid absolute path",
			path:    "/home/user/.spcstr/config.json",
			wantErr: false,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
			errMsg:  "config path cannot be empty",
		},
		{
			name:    "missing extension",
			path:    ".spcstr/config",
			wantErr: true,
			errMsg:  "config file must have .json extension",
		},
		{
			name:    "wrong extension",
			path:    ".spcstr/config.yaml",
			wantErr: true,
			errMsg:  "config file must have .json extension",
		},
		{
			name:    "parent directory traversal",
			path:    "../config.json",
			wantErr: true,
			errMsg:  "config path cannot contain parent directory references",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidConfigPath(tt.path)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		slice    []string
		item     string
		expected bool
	}{
		{[]string{"a", "b", "c"}, "b", true},
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{}, "a", false},
		{[]string{"test"}, "test", true},
		{[]string{"test"}, "Test", false}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.item, func(t *testing.T) {
			result := contains(tt.slice, tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}
