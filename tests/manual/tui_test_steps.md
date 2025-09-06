# Manual TUI Testing Steps

## Prerequisites
1. Build the spcstr binary: `make build`
2. Ensure binary is in PATH or use `./bin/spcstr`

## Test Scenarios

### 1. Test Project Not Initialized
**Steps:**
1. Navigate to a directory without `.spcstr` folder
2. Run `spcstr`
3. **Expected:** 
   - TUI launches with "Project not initialized" message
   - Message suggests running `spcstr init`
   - Press `q` to quit works

### 2. Test Project Initialized
**Steps:**
1. Run `spcstr init` in a test directory
2. Run `spcstr`
3. **Expected:**
   - TUI launches successfully
   - Header shows "spcstr | Plan View" and "Session: active"
   - Footer shows keybinds: [p] Plan [o] Observe [q] Quit

### 3. Test Navigation
**Steps:**
1. Launch TUI in initialized project
2. Press `o` to switch to Observe view
3. Press `p` to switch back to Plan view
4. **Expected:**
   - View switches immediately (<100ms)
   - Header updates to show current view
   - Footer keybinds update based on view
   - No visual glitches during transition

### 4. Test Terminal Resize
**Steps:**
1. Launch TUI
2. Resize terminal window (make it smaller, then larger)
3. **Expected:**
   - Content reflows appropriately
   - Header and footer adjust to new width
   - No content is cut off or overlapping

### 5. Test Quit Functionality
**Steps:**
1. Launch TUI
2. Press `q` to quit
3. Try also with `Ctrl+C`
4. **Expected:**
   - TUI exits cleanly
   - Terminal is restored to normal state
   - No error messages

### 6. Test View-Specific Content
**Steps:**
1. Launch TUI and navigate to Plan view
2. Observe placeholder content
3. Switch to Observe view
4. Observe placeholder content
5. **Expected:**
   - Plan view shows "Document browser will be displayed here"
   - Observe view shows "Session monitoring dashboard will be displayed here"
   - Each view has distinct content

### 7. Test Performance
**Steps:**
1. Launch TUI
2. Rapidly switch between views (press `p` and `o` quickly)
3. **Expected:**
   - All transitions are smooth
   - No lag or stuttering
   - View state is maintained correctly

### 8. Test Different Terminal Sizes
**Steps:**
1. Test in minimum terminal (80x24)
2. Test in large terminal (200x60)
3. Test in narrow terminal (60x40)
4. **Expected:**
   - TUI adapts to all sizes
   - Content remains readable
   - Layout doesn't break

## Performance Verification
Run with logging to verify <100ms view switching:
```bash
spcstr 2>tui.log
# Check tui.log for any WARNING messages about view switch timing
```

## Known Issues to Check
- [ ] TUI should not flicker on startup
- [ ] Colors should be visible in all terminal types
- [ ] Unicode characters (borders) should render correctly
- [ ] No memory leaks during extended use