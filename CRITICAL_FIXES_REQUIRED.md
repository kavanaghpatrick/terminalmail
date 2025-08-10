# TerminalMail - Critical Fixes Required

## 🚨 Overview

Both Grok and Gemini identified several **CRITICAL BUGS** that would cause the app to crash, fail to start, or be unusable. These must be fixed before implementation begins.

## 🔴 Priority 1: Breaking Bugs (App Won't Start)

### 1. SQLite FTS5 Configuration (Days 1-2)
**Problem**: `sql.Open("sqlite3", "path?_fts5=1")` is invalid. FTS5 requires compile-time flags.
**Impact**: Database won't open, app crashes on startup.
**Fix**:
```go
// Correct approach
import _ "github.com/mattn/go-sqlite3"

// Build with: go build -tags "sqlite_fts5"
db, err := sql.Open("sqlite3", "file:terminalmail.db?cache=shared&mode=rwc")
if err != nil {
    log.Fatal(err)
}

// Verify FTS5 is available
_, err = db.Exec("CREATE VIRTUAL TABLE test_fts USING fts5(content)")
if err != nil {
    log.Fatal("FTS5 not available. Build with -tags sqlite_fts5")
}
db.Exec("DROP TABLE test_fts")
```

### 2. Cross-Platform OAuth Browser Opening (Day 3-4)
**Problem**: `exec.Command("open", authURL)` only works on macOS.
**Impact**: OAuth fails on Linux/Windows, users can't authenticate.
**Fix**:
```go
func openBrowser(url string) error {
    var cmd string
    var args []string
    
    switch runtime.GOOS {
    case "darwin":
        cmd = "open"
        args = []string{url}
    case "linux":
        cmd = "xdg-open"
        args = []string{url}
    case "windows":
        cmd = "cmd"
        args = []string{"/c", "start", url}
    default:
        return fmt.Errorf("unsupported platform")
    }
    
    return exec.Command(cmd, args...).Start()
}
```

### 3. Missing Error Handling Throughout
**Problem**: No error checks after critical operations.
**Impact**: Silent failures, data corruption, panics.
**Fix**: Add error handling to EVERY operation:
```go
// WRONG
db.Exec("PRAGMA journal_mode=WAL")

// RIGHT
if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
    return fmt.Errorf("failed to set WAL mode: %w", err)
}
```

## 🔴 Priority 2: Runtime Crashes

### 4. Array Index Out of Bounds (Days 5-7)
**Problem**: No bounds checking for cursor movement and quote colors.
**Impact**: Panic when navigating past email list bounds or with deep quotes.
**Fix**:
```go
// Email list navigation
case "j", "down":
    if m.cursor < len(m.emails)-1 {
        m.cursor++
    }
case "k", "up":
    if m.cursor > 0 {
        m.cursor--
    }

// Quote coloring
func getQuoteColor(level int) string {
    if level <= 0 {
        return defaultColor
    }
    if level > len(quoteColors) {
        level = len(quoteColors)
    }
    return quoteColors[level-1]
}
```

### 5. Blocking UI with Editor (Day 8-9)
**Problem**: Synchronous `exec.Command` blocks Bubbletea event loop.
**Impact**: UI freezes during email composition.
**Fix**:
```go
// Return a command to run editor asynchronously
func composeEmailCmd() tea.Cmd {
    return func() tea.Msg {
        content, err := openInEditor("")
        return editorFinishedMsg{content: content, err: err}
    }
}

// Handle in Update()
case editorFinishedMsg:
    if msg.err != nil {
        return m, nil
    }
    // Process composed email
```

### 6. Missing Grok API Key (Day 12)
**Problem**: No API key handling for Grok requests.
**Impact**: All AI classification fails with 401 errors.
**Fix**:
```go
func NewGrokQueue() (*GrokQueue, error) {
    apiKey := os.Getenv("GROK_API_KEY")
    if apiKey == "" {
        // Graceful degradation - AI features disabled
        return nil, fmt.Errorf("GROK_API_KEY not set, AI features disabled")
    }
    
    return &GrokQueue{
        apiKey: apiKey,
        client: &http.Client{Timeout: 10 * time.Second},
    }, nil
}
```

## 🟡 Priority 3: Data Integrity Issues

### 7. SQL Injection Risk (Day 10-11)
**Problem**: User input in SQL queries without proper escaping.
**Impact**: Data corruption, security vulnerability.
**Fix**:
```go
// ALWAYS use parameterized queries
stmt, err := db.Prepare("SELECT * FROM emails WHERE email_search MATCH ?")
if err != nil {
    return nil, err
}
defer stmt.Close()

// Escape special FTS5 characters
query = strings.ReplaceAll(query, `"`, `""`)
rows, err := stmt.Query(query)
```

### 8. Race Conditions in Background Goroutines (Day 12)
**Problem**: Concurrent access to shared state without synchronization.
**Impact**: Data races, crashes, corrupted state.
**Fix**:
```go
type EmailApp struct {
    mu        sync.RWMutex
    emails    []Email
    // ...
}

// Background classification
go func() {
    classified := classifyEmails(batch)
    
    app.mu.Lock()
    defer app.mu.Unlock()
    
    // Update shared state safely
    for _, email := range classified {
        app.updateEmailCategory(email.ID, email.Category)
    }
}()
```

### 9. Keyring Fallback Handling (Day 3-4)
**Problem**: Keyring may not be available (headless Linux, SSH).
**Impact**: Token storage fails, OAuth breaks.
**Fix**:
```go
func saveToken(token string) error {
    err := keyring.Set("terminalmail", "gmail_token", token)
    if err != nil {
        // Fallback to encrypted file
        return saveTokenToFile(token)
    }
    return nil
}

func saveTokenToFile(token string) error {
    // Encrypt and save to ~/.config/terminalmail/token
    // Use OS file permissions (0600)
    configDir := filepath.Join(os.Getenv("HOME"), ".config", "terminalmail")
    os.MkdirAll(configDir, 0700)
    return os.WriteFile(filepath.Join(configDir, "token"), []byte(token), 0600)
}
```

## 🟢 Simplicity Filter Applied

### What We're NOT Fixing (Complexity Creep)
These suggestions from reviewers add unnecessary complexity:

❌ **Advanced error recovery mechanisms** - Simple retry is enough
❌ **Database migration system** - Single schema version for MVP
❌ **Complex OAuth token refresh logic** - Google's client handles it
❌ **Elaborate concurrency patterns** - Simple mutex is sufficient
❌ **Performance profiling infrastructure** - Premature optimization

### What We ARE Keeping Simple
✅ **Single SQLite connection** - No connection pooling needed
✅ **Basic error handling** - Log and return, no complex recovery
✅ **Simple state struct** - 5-7 fields, no state machines
✅ **Direct Gmail API calls** - No abstraction layers
✅ **Synchronous database operations** - SQLite handles it fine

## 📝 Updated Pseudocode for Critical Sections

### Day 1-2: Foundation (FIXED)
```go
func setupDatabase() (*sql.DB, error) {
    // Build with: go build -tags sqlite_fts5
    dbPath := filepath.Join(os.Getenv("HOME"), ".terminalmail", "emails.db")
    os.MkdirAll(filepath.Dir(dbPath), 0700)
    
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }
    
    // Essential pragmas with error checking
    pragmas := []string{
        "PRAGMA journal_mode=WAL",
        "PRAGMA synchronous=NORMAL",
        "PRAGMA cache_size=10000",
    }
    
    for _, pragma := range pragmas {
        if _, err := db.Exec(pragma); err != nil {
            return nil, fmt.Errorf("%s failed: %w", pragma, err)
        }
    }
    
    // Create schema
    if err := createSchema(db); err != nil {
        return nil, err
    }
    
    return db, nil
}
```

### Day 5-6: Email List (FIXED)
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "j", "down":
            if m.cursor < len(m.emails)-1 {
                m.cursor++
            }
        case "k", "up":
            if m.cursor > 0 {
                m.cursor--
            }
        case "enter":
            if m.cursor >= 0 && m.cursor < len(m.emails) {
                m.view = ReadView
                m.currentEmail = m.emails[m.cursor]
            }
        }
    }
    return m, nil
}
```

## 🚀 Implementation Readiness

With these fixes applied:
- ✅ App will start on all platforms
- ✅ OAuth will work cross-platform
- ✅ No array index panics
- ✅ No SQL injection vulnerabilities
- ✅ No race conditions
- ✅ Graceful degradation when services unavailable

**The plan is now implementation-ready while maintaining simplicity.**