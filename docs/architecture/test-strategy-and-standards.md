# Test Strategy and Standards

## Testing Philosophy
- **Approach:** Test-After Development with comprehensive coverage
- **Coverage Goals:** 80% for core packages, 60% for UI packages
- **Test Pyramid:** 70% unit tests, 20% integration tests, 10% manual TUI tests

## Test Types and Organization

### Unit Tests
- **Framework:** testing (Go standard library)
- **File Convention:** {filename}_test.go in same package
- **Location:** Alongside source files
- **Mocking Library:** testify/mock 1.8.4
- **Coverage Requirement:** 80% for business logic, 60% for UI

**AI Agent Requirements:**
- Generate tests for all public methods
- Cover edge cases and error conditions
- Follow AAA pattern (Arrange, Act, Assert)
- Mock all external dependencies

### Integration Tests
- **Scope:** End-to-end workflows, file system operations, event flow
- **Location:** test/integration/
- **Test Infrastructure:**
  - **File System:** Real temp directories via t.TempDir()
  - **Event Bus:** Real implementation with test subscribers
  - **Sessions:** Test fixtures in test/fixtures/sessions/

### End-to-End Tests
- **Framework:** Manual testing procedures documented
- **Scope:** TUI interactions, user workflows
- **Environment:** Local development environment
- **Test Data:** Scripted session generation

## Test Data Management
- **Strategy:** Fixtures for deterministic tests, generated data for property tests
- **Fixtures:** test/fixtures/ with JSON test data
- **Factories:** Test builders for complex objects
- **Cleanup:** Automatic via t.Cleanup() and t.TempDir()

## Continuous Testing
- **CI Integration:** All unit and integration tests run on every push
- **Performance Tests:** Benchmark critical paths (session loading, indexing)
- **Security Tests:** Static analysis with gosec in CI pipeline

## Testing Patterns

**Table-Driven Tests (Preferred):**
```go
func TestSessionManager_LoadSession(t *testing.T) {
    tests := []struct {
        name    string
        path    string
        want    *Session
        wantErr bool
    }{
        {"valid session", "valid.json", &Session{...}, false},
        {"missing file", "missing.json", nil, true},
        {"invalid json", "invalid.json", nil, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

**Mock Usage:**
```go
type MockEventBus struct {
    mock.Mock
}

func (m *MockEventBus) Publish(event Event) error {
    args := m.Called(event)
    return args.Error(0)
}
```

**TUI Testing Approach:**
Since Bubbletea TUIs are challenging to test automatically:
1. Test business logic separately from UI
2. Test model updates without rendering
3. Manual test scripts for visual verification
4. Screenshot-based regression tests for critical views

## Test Organization Rules

1. **One test file per source file**
2. **Test package matches source package**
3. **Integration tests isolated in test/integration/**
4. **Benchmarks suffixed with _bench_test.go**
5. **Test helpers in testutil package**

## Critical Testing Requirements

- **Never skip tests in CI:** All tests must pass for merge
- **Test error paths explicitly:** Every error return must be tested
- **Use subtests for related cases:** Better test output and filtering
- **Parallel tests where safe:** Mark with t.Parallel() for speed
- **Clean up resources:** Use defer or t.Cleanup()
- **Deterministic tests only:** No random data without seed

## Performance Benchmarks

Key operations requiring benchmarks:
- Session JSON parsing/serialization
- Document indexing and search
- Event bus message throughput
- File watcher responsiveness

Benchmark format:
```go
func BenchmarkSessionManager_LoadSession(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Benchmark code
    }
}
```
