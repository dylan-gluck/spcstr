# Monitoring and Observability

## Monitoring Stack

- **Application Monitoring:** Built-in TUI error display and log file monitoring
- **Hook Execution Monitoring:** Comprehensive logging to `.spcstr/logs/` with success/failure tracking
- **Performance Monitoring:** Built-in timing for hook execution and TUI render performance
- **Health Checking:** Self-diagnostics via `spcstr config doctor` command (future enhancement)

## Key Metrics

**Hook Performance Metrics:**
- Hook execution time (target: <100ms per hook)
- Hook success/failure rates per session
- State file write latency and atomic operation success
- File watcher event processing delay

**TUI Performance Metrics:**
- View switch response time (target: <100ms)
- Document rendering time for large markdown files
- Memory usage during long-running TUI sessions
- File system event processing throughput

**System Health Metrics:**
- `.spcstr/` directory size and growth rate
- Log file rotation effectiveness
- State file corruption incidents
- Hook process spawn success rate

## Built-in Observability Features

```go
// Performance monitoring built into state operations
type PerformanceTracker struct {
    hookTimes    map[string][]time.Duration
    renderTimes  []time.Duration
    stateOpTimes []time.Duration
    mutex        sync.RWMutex
}

func (p *PerformanceTracker) TrackHookExecution(hookName string, duration time.Duration) {
    p.mutex.Lock()
    defer p.mutex.Unlock()
    
    if p.hookTimes[hookName] == nil {
        p.hookTimes[hookName] = make([]time.Duration, 0)
    }
    
    p.hookTimes[hookName] = append(p.hookTimes[hookName], duration)
    
    // Log slow operations
    if duration > time.Millisecond*100 {
        log.Printf("SLOW HOOK: %s took %v", hookName, duration)
    }
}

func (p *PerformanceTracker) GetStats() map[string]interface{} {
    p.mutex.RLock()
    defer p.mutex.RUnlock()
    
    stats := make(map[string]interface{})
    
    for hook, times := range p.hookTimes {
        if len(times) == 0 {
            continue
        }
        
        var total time.Duration
        for _, t := range times {
            total += t
        }
        
        avg := total / time.Duration(len(times))
        stats[hook] = map[string]interface{}{
            "count":   len(times),
            "average": avg.Milliseconds(),
            "total":   total.Milliseconds(),
        }
    }
    
    return stats
}
```

## Health Check System

```go
// Built-in diagnostics for troubleshooting
type HealthChecker struct {
    stateManager *state.Manager
    configPath   string
}

func (h *HealthChecker) RunDiagnostics() (*HealthReport, error) {
    report := &HealthReport{
        Timestamp: time.Now(),
        Version:   BuildVersion,
    }
    
    // Check .spcstr directory structure
    report.DirectoryStructure = h.checkDirectoryStructure()
    
    // Check state file integrity
    report.StateFileHealth = h.checkStateFiles()
    
    // Check log files
    report.LogFileHealth = h.checkLogFiles()
    
    // Check hook configuration
    report.HookConfiguration = h.checkHookConfiguration()
    
    // Performance summary
    report.Performance = h.getPerformanceMetrics()
    
    return report, nil
}

type HealthReport struct {
    Timestamp           time.Time              `json:"timestamp"`
    Version             string                 `json:"version"`
    DirectoryStructure  DirectoryHealthCheck   `json:"directory_structure"`
    StateFileHealth     StateFileHealthCheck   `json:"state_file_health"`
    LogFileHealth       LogFileHealthCheck     `json:"log_file_health"`
    HookConfiguration   HookConfigHealthCheck  `json:"hook_configuration"`
    Performance         PerformanceMetrics     `json:"performance"`
}
```

This comprehensive architecture document provides the complete technical foundation for building spcstr as a specialized CLI/TUI observability tool. The design emphasizes reliability through atomic operations, real-time responsiveness through file watching, and maintainability through clean Go architecture patterns. The single-binary approach with embedded hook functionality ensures consistent deployment while the filesystem-first design preserves user privacy and enables offline operation.