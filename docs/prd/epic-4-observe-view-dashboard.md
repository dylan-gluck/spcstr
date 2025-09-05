# Epic 4: Observe View & Dashboard

Create the Observe View with session list and detailed dashboard. Implement the orchestration dashboard layout with organized sections for agents, tasks, files, and metrics. Enable session management capabilities.

## Story 4.1: Session List Implementation

As a user,
I want to see all active and completed sessions,
so that I can select and monitor any session.

### Acceptance Criteria
1: Sessions display with ID, status, and duration
2: Active sessions appear at top with visual indicator
3: List updates automatically as sessions change
4: Sessions can be filtered by status or date
5: Completed sessions show completion time

## Story 4.2: Orchestration Dashboard Layout

As a user,
I want a comprehensive dashboard showing session details,
so that I understand what's happening in each session.

### Acceptance Criteria
1: Dashboard sections for agents, tasks, files, tools
2: Agents show with current status and activity
3: Task progress displays as todo/in-progress/done
4: Tool usage shows with execution counts
5: Layout adjusts to terminal size intelligently

## Story 4.3: Real-time Dashboard Updates

As a user,
I want the dashboard to update in real-time,
so that I can monitor ongoing activity.

### Acceptance Criteria
1: Updates appear within 1 second of hook trigger
2: Animations indicate changing values
3: Historical data remains visible with timestamps
4: Performance remains smooth with rapid updates
5: User can pause updates to examine details

## Story 4.4: Session Management Controls

As a user,
I want to manage sessions from the TUI,
so that I can kill or resume sessions as needed.

### Acceptance Criteria
1: Ctrl+K kills selected session with confirmation
2: Ctrl+R resumes stopped session if possible
3: Session details can be exported to file
4: Clear error messages for invalid operations
5: Audit log tracks session management actions
