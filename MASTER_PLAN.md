# TerminalMail - Master Implementation Plan (Go Edition)
*Single Source of Truth - All Previous Documents Superseded*

## 🎯 Project Vision

Build a **blazing fast, beautiful terminal email client** in Go that ports the proven functionality of AI Mailbox to the command line, using Bubbletea for UI and Lipgloss for stunning ANSI gradients.

## ⚡ Quick Facts

- **Language**: Go 1.21+
- **Timeline**: 2-3 weeks to MVP
- **UI Framework**: Bubbletea (Elm-inspired TUI)
- **Styling**: Lipgloss (gradient support)
- **Database**: SQLite3 with FTS5 (exact schema from AI Mailbox)
- **State Management**: Single struct with 5-7 fields
- **Memory Target**: <10MB
- **Startup Target**: <20ms
- **Binary Size**: <10MB

## 🏗️ Technical Stack (FINAL)

```go
// go.mod
module github.com/yourusername/terminalmail

go 1.21

require (
    github.com/charmbracelet/bubbletea v0.24.0
    github.com/charmbracelet/lipgloss v0.9.1
    github.com/mattn/go-sqlite3 v1.14.17
    google.golang.org/api v0.150.0
    golang.org/x/oauth2 v0.13.0
    github.com/zalando/go-keyring v0.2.3
    github.com/muesli/termenv v0.15.2
    github.com/spf13/cobra v1.8.0
)
```

## 📊 Database Schema (Direct Port)

```sql
-- Exact schema from AI Mailbox - DO NOT CHANGE
CREATE TABLE emails (
  id TEXT PRIMARY KEY,
  subject TEXT,
  body_text TEXT,
  body_html TEXT,
  from_email TEXT,
  to_email TEXT,
  cc_email TEXT,
  date TEXT,
  labels TEXT,
  ai_summary TEXT,
  status TEXT,
  thread_id TEXT,
  ai_category TEXT
);

CREATE VIRTUAL TABLE email_search USING fts5(
  subject, body_text, from_email, to_email, labels,
  content='emails', tokenize='porter'
);

-- Indexes proven to work
CREATE INDEX idx_compound_inbox ON emails(labels, date);
CREATE INDEX idx_ai_category ON emails(ai_category);

-- Terminal-specific preferences
CREATE TABLE terminal_preferences (
  id INTEGER PRIMARY KEY CHECK (id = 1),
  use_emoji BOOLEAN DEFAULT 0,
  theme TEXT DEFAULT 'cyberpunk',
  list_density INTEGER DEFAULT 2
);
```

## 🔄 Gmail Sync (Proven Algorithm - Go Implementation)

```go
// DO NOT CHANGE - This 2-phase sync is perfect
type GmailSync struct {
    service *gmail.Service
    db      *EmailDB
}

func (g *GmailSync) PrioritySync() error {
    // Phase 1: Last 48 hours (immediate)
    recent, err := g.fetchEmails("after:2d")
    if err != nil {
        return err
    }
    
    if err := g.db.SaveEmails(recent); err != nil {
        return err
    }
    
    // Phase 2: Background backfill (non-blocking)
    go g.backfillOlderEmails()
    
    return nil
}
```

## 🎨 Visual Design System (Lipgloss)

```go
// styles.go - ANSI Gradient Categories
var CategoryStyles = map[string]lipgloss.Style{
    "financial": lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{Light: "#FFD700", Dark: "#00FF7F"}).
        Background(lipgloss.AdaptiveColor{Light: "#FFF8DC", Dark: "#1A1A1A"}),
    
    "security": lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{Light: "#8B0000", Dark: "#FF4500"}).
        Bold(true),
    
    "newsletter": lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{Light: "#0077BE", Dark: "#00BFFF"}),
        
    "reminder": lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{Light: "#FF5E00", Dark: "#FFD700"}),
        
    "receipt": lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{Light: "#800080", Dark: "#FF69B4"}),
        
    "marketing": lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{Light: "#00CED1", Dark: "#40E0D0"}),
        
    "waiting": lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{Light: "#FFD700", Dark: "#FFFF00"}).
        Blink(true),
}

// Gradient support
var HeaderGradient = lipgloss.NewStyle().
    Foreground(lipgloss.Color("86")).
    Background(lipgloss.Color("63")).
    Bold(true)
```

## 🏛️ Application Architecture

```go
// main.go - Complete application structure
package main

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

// Complete state in one struct (5-7 fields)
type Model struct {
    view      ViewType         // list | read | compose
    folder    string           // inbox | sent | archive  
    selection map[string]bool  // selected email IDs
    cursor    int              // current position
    search    string           // active search query
    emails    []Email          // current email list
    viewport  viewport.Model   // for scrolling
}

// View types
type ViewType int
const (
    ListView ViewType = iota
    ReadView
    ComposeView
)

// Initialize the model
func initialModel() Model {
    return Model{
        view:      ListView,
        folder:    "inbox",
        selection: make(map[string]bool),
        cursor:    0,
        emails:    []Email{},
    }
}

// Update handles all state changes
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "j", "down":
            m.cursor++
            if m.cursor >= len(m.emails) {
                m.cursor = len(m.emails) - 1
            }
        case "k", "up":
            m.cursor--
            if m.cursor < 0 {
                m.cursor = 0
            }
        case "enter":
            m.view = ReadView
        case "c":
            m.view = ComposeView
        case "q":
            return m, tea.Quit
        }
    }
    return m, nil
}

// View renders the UI
func (m Model) View() string {
    switch m.view {
    case ListView:
        return m.renderEmailList()
    case ReadView:
        return m.renderEmailReader()
    case ComposeView:
        return m.renderComposer()
    }
    return ""
}

func main() {
    p := tea.NewProgram(initialModel(), tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
}
```

## 🤖 AI Integration (Direct Port to Go)

```go
// grok_queue.go - Exact queue implementation from AI Mailbox
type GrokQueue struct {
    queue      []QueueItem
    processing bool
    delay      time.Duration
    mu         sync.Mutex
}

type QueueItem struct {
    task   func() (interface{}, error)
    result chan interface{}
    err    chan error
}

func NewGrokQueue() *GrokQueue {
    return &GrokQueue{
        delay: 200 * time.Millisecond, // Proven delay
    }
}

func (q *GrokQueue) Add(task func() (interface{}, error)) (interface{}, error) {
    item := QueueItem{
        task:   task,
        result: make(chan interface{}),
        err:    make(chan error),
    }
    
    q.mu.Lock()
    q.queue = append(q.queue, item)
    if !q.processing {
        go q.process()
    }
    q.mu.Unlock()
    
    select {
    case res := <-item.result:
        return res, nil
    case err := <-item.err:
        return nil, err
    }
}

func (q *GrokQueue) process() {
    q.processing = true
    defer func() { q.processing = false }()
    
    for len(q.queue) > 0 {
        q.mu.Lock()
        item := q.queue[0]
        q.queue = q.queue[1:]
        q.mu.Unlock()
        
        res, err := item.task()
        if err != nil {
            item.err <- err
        } else {
            item.result <- res
        }
        
        time.Sleep(q.delay)
    }
}
```

## 🛠️ Core Utility Functions

```go
// internal/utils/math.go - CRITICAL: Go 1.21+ doesn't have min/max in stdlib
func min(a, b int) int { 
    if a < b { return a }
    return b 
}

func max(a, b int) int { 
    if a > b { return a }
    return b 
}
```

## 📈 Detailed 14-Day Implementation Schedule with Pseudocode

### Week 1: Core Foundation

#### **Day 1-2: Project Setup & Database**

**Task 1: Initialize Go project**
```go
// Pseudocode for project initialization
func initializeProject() {
    // Create directory structure
    mkdir terminalmail/
    mkdir terminalmail/{cmd,internal,pkg}
    mkdir terminalmail/internal/{database,gmail,ui,ai}
    
    // Initialize Go module
    go mod init github.com/username/terminalmail
    
    // Add core dependencies
    go get github.com/charmbracelet/bubbletea
    go get github.com/charmbracelet/lipgloss
    go get github.com/mattn/go-sqlite3
    
    // Create main.go with basic Bubbletea app
    createFile("cmd/tmail/main.go", basicBubbleteaApp)
    
    // Verify it runs
    go run cmd/tmail/main.go
}
```

**Task 2: Port SQLite schema from AI Mailbox**
```go
// Pseudocode for database setup (FIXED for FTS5 support)
func setupDatabase() {
    // internal/database/schema.go
    const schema = `
        CREATE TABLE IF NOT EXISTS emails (
            id TEXT PRIMARY KEY,
            subject TEXT,
            body_text TEXT,
            body_html TEXT,
            from_email TEXT,
            to_email TEXT,
            cc_email TEXT,
            date TEXT,
            labels TEXT,
            ai_summary TEXT,
            status TEXT,
            thread_id TEXT,
            ai_category TEXT
        );
        
        CREATE VIRTUAL TABLE IF NOT EXISTS email_search 
        USING fts5(
            subject, body_text, from_email, to_email, labels,
            content='emails', tokenize='porter'
        );
        
        CREATE INDEX IF NOT EXISTS idx_compound_inbox 
        ON emails(labels, date);
    `
    
    // internal/database/db.go
    type EmailDB struct {
        conn *sql.DB
    }
    
    // Cross-platform home directory helper
    func getHomeDir() string {
        if runtime.GOOS == "windows" {
            return os.Getenv("USERPROFILE")
        }
        return os.Getenv("HOME")
    }
    
    func NewEmailDB() (*EmailDB, error) {
        // FIXED: Cross-platform home directory and FTS5 build requirement
        dbPath := filepath.Join(getHomeDir(), ".terminalmail", "emails.db")
        if err := os.MkdirAll(filepath.Dir(dbPath), 0700); err != nil {
            return nil, fmt.Errorf("failed to create db directory: %w", err)
        }
        
        // FIXED: No invalid DSN parameters - build with -tags sqlite_fts5
        db, err := sql.Open("sqlite3", dbPath)
        if err != nil {
            return nil, err
        }
        
        // FIXED: Error handling for critical operations
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
        
        // FIXED: Create schema first, then verify FTS5 
        if _, err := db.Exec(schema); err != nil {
            return nil, fmt.Errorf("schema creation failed: %w", err)
        }
        
        // FIXED: Verify FTS5 is available after schema creation
        if _, err := db.Exec("SELECT * FROM email_search LIMIT 0"); err != nil {
            return nil, fmt.Errorf("FTS5 not available. Build with: go build -tags sqlite_fts5")
        }
        
        return &EmailDB{conn: db}, nil
    }
}
```

**Task 3: Set up Lipgloss theme system**
```go
// Pseudocode for theme setup
func setupThemes() {
    // internal/ui/themes.go
    type Theme struct {
        Name       string
        Primary    lipgloss.Style
        Secondary  lipgloss.Style
        Success    lipgloss.Style
        Error      lipgloss.Style
        Categories map[string]lipgloss.Style
    }
    
    var CyberpunkTheme = Theme{
        Name: "cyberpunk",
        Primary: lipgloss.NewStyle().
            Foreground(lipgloss.Color("#FF00FF")).
            Background(lipgloss.Color("#1A1A2E")),
        Categories: map[string]lipgloss.Style{
            "financial": goldToGreenGradient(),
            "security":  redPulsingStyle(),
        },
    }
    
    func goldToGreenGradient() lipgloss.Style {
        // Use termenv for gradient calculation
        return lipgloss.NewStyle().
            Foreground(lipgloss.AdaptiveColor{
                Light: "#FFD700",
                Dark:  "#00FF7F",
            })
    }
}
```

#### **Day 3-4: Gmail Integration**

**Task 1: Port OAuth flow from AI Mailbox**
```go
// Pseudocode for OAuth setup
func setupGmailAuth() {
    // internal/gmail/auth.go
    type GmailAuth struct {
        config *oauth2.Config
        token  *oauth2.Token
    }
    
    func NewGmailAuth() *GmailAuth {
        config := &oauth2.Config{
            ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
            ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
            Endpoint:     google.Endpoint,
            RedirectURL:  "http://localhost:8080/callback",
            Scopes: []string{
                gmail.GmailReadonlyScope,
                gmail.GmailSendScope,
                gmail.GmailModifyScope,
            },
        }
        
        return &GmailAuth{config: config}
    }
    
    func (g *GmailAuth) Authenticate() error {
        // Check keyring for existing token
        tokenJSON, err := keyring.Get("terminalmail", "gmail_token")
        if err == nil {
            // Parse and use existing token
            if err := json.Unmarshal([]byte(tokenJSON), &g.token); err != nil {
                return fmt.Errorf("failed to parse stored token: %w", err)
            }
            return nil
        }
        
        // Start local server for OAuth callback with port retry
        server, err := startCallbackServer()
        if err != nil {
            return fmt.Errorf("failed to start callback server: %w", err)
        }
        defer server.Close()
        
        // Generate auth URL
        authURL := g.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
        
        // Open browser (cross-platform)
        if err := openBrowser(authURL); err != nil {
            return fmt.Errorf("failed to open browser: %w", err)
        }
        
        // Wait for callback
        code := <-server.codeChan
        
        // Exchange for token
        token, err := g.config.Exchange(context.Background(), code)
        if err != nil {
            return fmt.Errorf("failed to exchange code for token: %w", err)
        }
        
        // Save to keyring with fallback
        if err := saveToken(token); err != nil {
            return fmt.Errorf("failed to save token: %w", err)
        }
        
        g.token = token
        return nil
    }
    
    // Cross-platform browser opening
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
            return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
        }
        
        return exec.Command(cmd, args...).Start()
    }
    
    // Token storage with keyring fallback
    func saveToken(token *oauth2.Token) error {
        tokenJSON, err := json.Marshal(token)
        if err != nil {
            return fmt.Errorf("failed to marshal token: %w", err)
        }
        
        err = keyring.Set("terminalmail", "gmail_token", string(tokenJSON))
        if err != nil {
            // Fallback to encrypted file
            return saveTokenToFile(string(tokenJSON))
        }
        return nil
    }
    
    // Port retry callback server to handle conflicts
    func startCallbackServer() (*CallbackServer, error) {
        // Try ports 8080-8090 to avoid conflicts
        for port := 8080; port <= 8090; port++ {
            addr := fmt.Sprintf("localhost:%d", port)
            listener, err := net.Listen("tcp", addr)
            if err == nil {
                server := &CallbackServer{
                    listener: listener,
                    port:     port,
                    codeChan: make(chan string, 1),
                }
                go server.serve()
                return server, nil
            }
        }
        return nil, fmt.Errorf("no available ports in range 8080-8090")
    }
    
    func saveTokenToFile(token string) error {
        // Save to cross-platform config directory with secure permissions
        configDir := filepath.Join(getHomeDir(), ".config", "terminalmail")
        if err := os.MkdirAll(configDir, 0700); err != nil {
            return fmt.Errorf("failed to create config dir: %w", err)
        }
        
        tokenPath := filepath.Join(configDir, "token")
        if err := os.WriteFile(tokenPath, []byte(token), 0600); err != nil {
            return fmt.Errorf("failed to write token file: %w", err)
        }
        
        return nil
    }
}
```

**Task 2: Implement 2-phase priority sync**
```go
// Pseudocode for sync implementation
func implementSync() {
    // internal/gmail/sync.go
    type Syncer struct {
        service   *gmail.Service
        db        *EmailDB
        batchSize int
    }
    
    func (s *Syncer) PrioritySync() error {
        // Phase 1: Recent emails (48 hours)
        recentQuery := fmt.Sprintf("after:%d", 
            time.Now().Add(-48*time.Hour).Unix())
        
        recentEmails := s.fetchBatch(recentQuery, 100)
        s.db.SaveEmails(recentEmails)
        
        // Update UI with recent emails
        s.sendUpdate(SyncPhase1Complete{Count: len(recentEmails)})
        
        // Phase 2: Background backfill (non-blocking)
        go func() {
            olderQuery := fmt.Sprintf("after:%d before:%d",
                time.Now().Add(-365*24*time.Hour).Unix(),
                time.Now().Add(-48*time.Hour).Unix())
            
            pageToken := ""
            for {
                batch := s.fetchBatch(olderQuery, 100, pageToken)
                if len(batch) == 0 {
                    break
                }
                s.db.SaveEmails(batch)
                s.sendUpdate(SyncProgress{Count: len(batch)})
            }
            s.sendUpdate(SyncComplete{})
        }()
        
        return nil
    }
    
    func (s *Syncer) fetchBatch(query string, size int, token ...string) []Email {
        call := s.service.Users.Messages.List("me").
            Q(query).
            MaxResults(int64(size))
        
        if len(token) > 0 {
            call.PageToken(token[0])
        }
        
        response, _ := call.Do()
        
        emails := []Email{}
        for _, msg := range response.Messages {
            fullMsg, _ := s.service.Users.Messages.Get("me", msg.Id).
                Format("full").Do()
            emails = append(emails, s.parseGmailMessage(fullMsg))
        }
        
        return emails
    }
}
```

**Task 3: Add progress indicators with Bubbletea**
```go
// Pseudocode for progress display
func addProgressIndicators() {
    // internal/ui/progress.go
    type progressMsg struct {
        current int
        total   int
        phase   string
    }
    
    type Model struct {
        progress  progress.Model
        syncing   bool
        phase     string
    }
    
    func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        switch msg := msg.(type) {
        case progressMsg:
            m.syncing = true
            m.phase = msg.phase
            cmd := m.progress.SetPercent(float64(msg.current) / float64(msg.total))
            return m, cmd
            
        case SyncComplete:
            m.syncing = false
            return m, nil
        }
        return m, nil
    }
    
    func (m Model) View() string {
        if m.syncing {
            return lipgloss.JoinVertical(
                lipgloss.Left,
                HeaderGradient.Render("Syncing: "+m.phase),
                m.progress.View(),
            )
        }
        return m.renderEmailList()
    }
}
```

#### **Day 5-6: Email List View**

**Task 1: Build Bubbletea component for email list**
```go
// Pseudocode for email list component
func buildEmailList() {
    // internal/ui/email_list.go
    type EmailListModel struct {
        emails    []Email
        cursor    int
        selected  map[string]bool
        viewport  viewport.Model
        width     int
        height    int
    }
    
    func NewEmailList(emails []Email) EmailListModel {
        vp := viewport.New(80, 20)
        return EmailListModel{
            emails:   emails,
            cursor:   0,
            selected: make(map[string]bool),
            viewport: vp,
        }
    }
    
    func (m *EmailListModel) Update(msg tea.Msg) (EmailListModel, tea.Cmd) {
        switch msg := msg.(type) {
        case tea.WindowSizeMsg:
            m.width = msg.Width
            m.height = msg.Height - 4 // Reserve for header/footer
            m.viewport.Width = msg.Width
            m.viewport.Height = m.height
            
        case tea.KeyMsg:
            switch msg.String() {
            case "j", "down":
                m.cursor = min(m.cursor+1, len(m.emails)-1)
                m.updateViewport()
            case "k", "up":
                m.cursor = max(m.cursor-1, 0)
                m.updateViewport()
            case "space":
                // Toggle selection with nil map safety
                if m.selected == nil {
                    m.selected = make(map[string]bool)
                }
                if m.cursor >= 0 && m.cursor < len(m.emails) {
                    id := m.emails[m.cursor].ID
                    m.selected[id] = !m.selected[id]
                }
            case "a":
                // Archive selected
                return m, archiveEmails(m.getSelectedIDs())
            }
        }
        
        var cmd tea.Cmd
        m.viewport, cmd = m.viewport.Update(msg)
        return m, cmd
    }
    
    func (m EmailListModel) View() string {
        // Build email list with ANSI colors
        var lines []string
        
        for i, email := range m.emails {
            line := m.renderEmailLine(email, i)
            lines = append(lines, line)
        }
        
        m.viewport.SetContent(strings.Join(lines, "\n"))
        
        return lipgloss.JoinVertical(
            lipgloss.Left,
            m.renderHeader(),
            m.viewport.View(),
            m.renderFooter(),
        )
    }
    
    func (m EmailListModel) renderEmailLine(email Email, index int) string {
        // Category indicator
        categoryStyle := CategoryStyles[email.Category]
        indicator := categoryStyle.Render("█")
        
        // Selection state
        selected := m.selected[email.ID]
        cursor := index == m.cursor
        
        style := lipgloss.NewStyle()
        if cursor {
            style = style.Background(lipgloss.Color("240"))
        }
        if selected {
            style = style.Foreground(lipgloss.Color("220"))
        }
        if !email.Read {
            style = style.Bold(true)
        }
        
        // Format line
        from := truncate(email.From, 20)
        subject := truncate(email.Subject, 40)
        date := formatDate(email.Date)
        
        line := fmt.Sprintf("%s %-20s %-40s %s",
            indicator, from, subject, date)
        
        return style.Render(line)
    }
}
```

**Task 2: Implement vim-style navigation**
```go
// Pseudocode for vim navigation
func implementVimNavigation() {
    // internal/ui/keybindings.go
    type KeyMap struct {
        Up       key.Binding
        Down     key.Binding
        PageUp   key.Binding
        PageDown key.Binding
        Home     key.Binding
        End      key.Binding
        Select   key.Binding
        Open     key.Binding
        Archive  key.Binding
        Delete   key.Binding
        Compose  key.Binding
        Reply    key.Binding
        Search   key.Binding
        Quit     key.Binding
    }
    
    var DefaultKeyMap = KeyMap{
        Up:       key.NewBinding(key.WithKeys("k", "up")),
        Down:     key.NewBinding(key.WithKeys("j", "down")),
        PageUp:   key.NewBinding(key.WithKeys("ctrl+u")),
        PageDown: key.NewBinding(key.WithKeys("ctrl+d")),
        Home:     key.NewBinding(key.WithKeys("g")),
        End:      key.NewBinding(key.WithKeys("G")),
        Select:   key.NewBinding(key.WithKeys("space", "x")),
        Open:     key.NewBinding(key.WithKeys("enter", "o")),
        Archive:  key.NewBinding(key.WithKeys("e")),
        Delete:   key.NewBinding(key.WithKeys("d")),
        Compose:  key.NewBinding(key.WithKeys("c")),
        Reply:    key.NewBinding(key.WithKeys("r")),
        Search:   key.NewBinding(key.WithKeys("/")),
        Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c")),
    }
    
    func handleVimMotion(m Model, key string) Model {
        switch key {
        case "5j":
            // Move down 5
            m.cursor = min(m.cursor+5, len(m.emails)-1)
        case "5k":
            // Move up 5
            m.cursor = max(m.cursor-5, 0)
        case "gg":
            // Go to top
            m.cursor = 0
        case "G":
            // Go to bottom
            m.cursor = len(m.emails) - 1
        case "zz":
            // Center current line
            m.centerCursor()
        }
        return m
    }
}
```

**Task 3: Add gradient category indicators**
```go
// Pseudocode for gradient indicators
func addGradientIndicators() {
    // internal/ui/gradients.go
    func RenderCategoryGradient(category string) string {
        colors := map[string][]string{
            "financial": {"#FFD700", "#90EE90", "#00FF7F"},
            "security":  {"#8B0000", "#DC143C", "#FF4500"},
            "newsletter": {"#000080", "#0000FF", "#87CEEB"},
        }
        
        gradientColors := colors[category]
        if gradientColors == nil {
            gradientColors = []string{"#808080", "#A9A9A9", "#D3D3D3"}
        }
        
        // Create gradient block
        block := ""
        for i, color := range gradientColors {
            style := lipgloss.NewStyle().
                Foreground(lipgloss.Color(color))
            block += style.Render("█")
        }
        
        return block
    }
    
    func AnimateCategoryIndicator(category string, frame int) string {
        if category == "waiting" {
            // Pulsing animation
            intensity := (math.Sin(float64(frame)*0.1) + 1) / 2
            color := lerpColor("#FFFFFF", "#FFFF00", intensity)
            return lipgloss.NewStyle().
                Foreground(lipgloss.Color(color)).
                Render("●")
        }
        
        if category == "security" {
            // Blinking animation
            if frame%20 < 10 {
                return lipgloss.NewStyle().
                    Foreground(lipgloss.Color("#FF0000")).
                    Bold(true).
                    Render("▲")
            }
            return " "
        }
        
        return RenderCategoryGradient(category)
    }
}
```

#### **Day 7: Email Reading**

**Task 1: Port email reader logic**
```go
// Pseudocode for email reader
func implementEmailReader() {
    // internal/ui/email_reader.go
    type EmailReader struct {
        email    Email
        viewport viewport.Model
        width    int
        height   int
    }
    
    func NewEmailReader(email Email) EmailReader {
        vp := viewport.New(80, 20)
        content := formatEmailContent(email)
        vp.SetContent(content)
        
        return EmailReader{
            email:    email,
            viewport: vp,
        }
    }
    
    // Quote color palette for email threads
    var quoteColors = []string{"86", "208", "196", "118", "201", "141"}
    
    func formatEmailContent(email Email) string {
        var sb strings.Builder
        
        // Header with gradient border
        border := lipgloss.NewStyle().
            Border(lipgloss.RoundedBorder()).
            BorderForeground(lipgloss.Color("86")).
            Padding(1, 2)
        
        // Metadata
        header := fmt.Sprintf(
            "From: %s\nTo: %s\nDate: %s\nSubject: %s",
            email.From, email.To, email.Date, email.Subject,
        )
        
        sb.WriteString(border.Render(header))
        sb.WriteString("\n\n")
        
        // Body with quote detection
        lines := strings.Split(email.BodyText, "\n")
        for _, line := range lines {
            if strings.HasPrefix(line, ">") {
                // Quote level coloring
                level := strings.Count(line, ">")
                color := quoteColors[min(level-1, len(quoteColors)-1)]
                style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
                sb.WriteString(style.Render(line))
            } else {
                sb.WriteString(line)
            }
            sb.WriteString("\n")
        }
        
        return sb.String()
    }
    
    func (r EmailReader) Update(msg tea.Msg) (EmailReader, tea.Cmd) {
        switch msg := msg.(type) {
        case tea.KeyMsg:
            switch msg.String() {
            case "r":
                return r, composeReply(r.email)
            case "f":
                return r, composeForward(r.email)
            case "e":
                return r, archiveEmail(r.email.ID)
            case "d":
                return r, deleteEmail(r.email.ID)
            case "q", "escape":
                return r, returnToList()
            }
        }
        
        var cmd tea.Cmd
        r.viewport, cmd = r.viewport.Update(msg)
        return r, cmd
    }
    
    func (r EmailReader) View() string {
        header := HeaderGradient.Render("Reading Email")
        footer := lipgloss.NewStyle().
            Foreground(lipgloss.Color("241")).
            Render("r: reply | f: forward | e: archive | d: delete | q: back")
        
        return lipgloss.JoinVertical(
            lipgloss.Left,
            header,
            r.viewport.View(),
            footer,
        )
    }
}
```

**Task 2: Implement markdown rendering**
```go
// Pseudocode for markdown support
func implementMarkdownRendering() {
    // internal/ui/markdown.go
    func RenderMarkdown(content string) string {
        // Convert markdown to terminal-friendly format
        lines := strings.Split(content, "\n")
        var rendered []string
        
        for _, line := range lines {
            // Headers
            if strings.HasPrefix(line, "# ") {
                style := lipgloss.NewStyle().
                    Bold(true).
                    Underline(true).
                    Foreground(lipgloss.Color("86"))
                rendered = append(rendered, style.Render(line[2:]))
                continue
            }
            
            // Bold
            line = replaceBold(line)
            
            // Italic
            line = replaceItalic(line)
            
            // Links
            line = replaceLinks(line)
            
            // Code blocks
            if strings.HasPrefix(line, "```") {
                // Handle code block
                style := lipgloss.NewStyle().
                    Background(lipgloss.Color("236")).
                    Foreground(lipgloss.Color("229"))
                // ... accumulate code lines
            }
            
            rendered = append(rendered, line)
        }
        
        return strings.Join(rendered, "\n")
    }
    
    func replaceBold(line string) string {
        re := regexp.MustCompile(`\*\*(.*?)\*\*`)
        return re.ReplaceAllStringFunc(line, func(match string) string {
            content := match[2 : len(match)-2]
            return lipgloss.NewStyle().Bold(true).Render(content)
        })
    }
}
```

**Task 3: Add thread support**
```go
// Pseudocode for thread display
func addThreadSupport() {
    // internal/ui/thread.go
    type ThreadView struct {
        thread    []Email
        collapsed map[string]bool
        cursor    int
    }
    
    func (t ThreadView) View() string {
        var lines []string
        
        for i, email := range t.thread {
            indent := strings.Repeat("  ", email.ThreadDepth)
            
            // Thread line connector
            connector := "├─"
            if i == len(t.thread)-1 {
                connector = "└─"
            }
            
            // Collapsed indicator
            indicator := ""
            if t.collapsed[email.ID] {
                indicator = "▶"
            } else if email.HasReplies {
                indicator = "▼"
            }
            
            style := lipgloss.NewStyle()
            if i == t.cursor {
                style = style.Background(lipgloss.Color("240"))
            }
            
            line := fmt.Sprintf("%s%s %s %s - %s",
                indent, connector, indicator,
                email.From, truncate(email.Subject, 50))
            
            lines = append(lines, style.Render(line))
            
            // Show nested replies if not collapsed
            if !t.collapsed[email.ID] {
                // Recursively render replies
            }
        }
        
        return strings.Join(lines, "\n")
    }
    
    func (t ThreadView) toggleCollapse() ThreadView {
        email := t.thread[t.cursor]
        t.collapsed[email.ID] = !t.collapsed[email.ID]
        return t
    }
}
```

### Week 2: Features & Polish

#### **Day 8-9: Compose & Reply**

**Task 1: Integrate $EDITOR**
```go
// Pseudocode for editor integration
func integrateEditor() {
    // internal/compose/editor.go
    func OpenInEditor(initialContent string) (string, error) {
        // Get editor preference with cross-platform defaults
        editor := os.Getenv("EDITOR")
        if editor == "" {
            // Platform-specific defaults
            switch runtime.GOOS {
            case "windows":
                editor = "notepad.exe"
            default:
                editor = "vim"
            }
        }
        
        // Create temp file
        tmpfile, err := os.CreateTemp("", "tmail-*.md")
        if err != nil {
            return "", err
        }
        defer os.Remove(tmpfile.Name())
        
        // FIXED: Write initial content with error handling
        if _, err := tmpfile.WriteString(initialContent); err != nil {
            return "", fmt.Errorf("failed to write initial content: %w", err)
        }
        if err := tmpfile.Close(); err != nil {
            return "", fmt.Errorf("failed to close temp file: %w", err)
        }
        
        // Open editor
        cmd := exec.Command(editor, tmpfile.Name())
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        
        if err := cmd.Run(); err != nil {
            return "", err
        }
        
        // Read result
        content, err := os.ReadFile(tmpfile.Name())
        if err != nil {
            return "", err
        }
        
        return string(content), nil
    }
    
    // Message types for async editor
    type editorFinishedMsg struct {
        content string
        err     error
    }
    
    // Async editor command for Bubbletea
    func composeEmailCmd() tea.Cmd {
        return func() tea.Msg {
            content, err := ComposeEmail()
            return editorFinishedMsg{content: content.BodyText, err: err}
        }
    }
    
    func ComposeEmail() (Email, error) {
        template := `To: 
Cc: 
Subject: 

---
Type your message below this line (markdown supported)
---

`
        
        content, err := OpenInEditor(template)
        if err != nil {
            return Email{}, err
        }
        
        return parseEmailFromTemplate(content)
    }
    
    func parseEmailFromTemplate(content string) (Email, error) {
        lines := strings.Split(content, "\n")
        email := Email{}
        bodyStart := -1
        
        for i, line := range lines {
            if strings.HasPrefix(line, "---") {
                bodyStart = i + 1
                break
            }
            
            if strings.HasPrefix(line, "To: ") {
                email.To = strings.TrimPrefix(line, "To: ")
            } else if strings.HasPrefix(line, "Cc: ") {
                email.Cc = strings.TrimPrefix(line, "Cc: ")
            } else if strings.HasPrefix(line, "Subject: ") {
                email.Subject = strings.TrimPrefix(line, "Subject: ")
            }
        }
        
        if bodyStart > 0 && bodyStart < len(lines) {
            email.BodyText = strings.Join(lines[bodyStart:], "\n")
        }
        
        return email, nil
    }
}
```

**Task 2: Port send logic from AI Mailbox**
```go
// Pseudocode for email sending
func implementSendLogic() {
    // internal/gmail/send.go
    func SendEmail(service *gmail.Service, email Email) error {
        // Create RFC 2822 message
        var message strings.Builder
        message.WriteString(fmt.Sprintf("From: me\r\n"))
        message.WriteString(fmt.Sprintf("To: %s\r\n", email.To))
        if email.Cc != "" {
            message.WriteString(fmt.Sprintf("Cc: %s\r\n", email.Cc))
        }
        message.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))
        message.WriteString("\r\n")
        message.WriteString(email.BodyText)
        
        // Base64 encode
        encoded := base64.URLEncoding.EncodeToString([]byte(message.String()))
        
        // Send via Gmail API
        msg := &gmail.Message{
            Raw: encoded,
        }
        
        _, err := service.Users.Messages.Send("me", msg).Do()
        return err
    }
    
    func SendReply(service *gmail.Service, original Email, reply Email) error {
        // Add threading headers
        var message strings.Builder
        message.WriteString(fmt.Sprintf("From: me\r\n"))
        message.WriteString(fmt.Sprintf("To: %s\r\n", original.From))
        message.WriteString(fmt.Sprintf("Subject: Re: %s\r\n", original.Subject))
        message.WriteString(fmt.Sprintf("In-Reply-To: %s\r\n", original.ID))
        message.WriteString(fmt.Sprintf("References: %s\r\n", original.ID))
        message.WriteString("\r\n")
        
        // Quote original message
        message.WriteString(reply.BodyText)
        message.WriteString("\n\n")
        message.WriteString(fmt.Sprintf("On %s, %s wrote:\n", 
            original.Date, original.From))
        
        for _, line := range strings.Split(original.BodyText, "\n") {
            message.WriteString("> " + line + "\n")
        }
        
        // Send
        encoded := base64.URLEncoding.EncodeToString([]byte(message.String()))
        msg := &gmail.Message{
            Raw:      encoded,
            ThreadId: original.ThreadID,
        }
        
        _, err := service.Users.Messages.Send("me", msg).Do()
        return err
    }
}
```

**Task 3: Add reply threading**
```go
// Pseudocode for reply threading
func addReplyThreading() {
    // internal/compose/reply.go
    type ReplyComposer struct {
        original  Email
        reply     Email
        quoteLevl int
    }
    
    func NewReplyComposer(original Email) ReplyComposer {
        template := formatReplyTemplate(original)
        
        return ReplyComposer{
            original: original,
            reply: Email{
                To:       original.From,
                Subject:  "Re: " + stripRePrefix(original.Subject),
                ThreadID: original.ThreadID,
            },
        }
    }
    
    func formatReplyTemplate(original Email) string {
        var sb strings.Builder
        
        sb.WriteString(fmt.Sprintf("To: %s\n", original.From))
        sb.WriteString(fmt.Sprintf("Subject: Re: %s\n", 
            stripRePrefix(original.Subject)))
        sb.WriteString("\n---\n")
        sb.WriteString("Type your reply above the quoted text\n")
        sb.WriteString("---\n\n")
        
        // Add quoted original
        sb.WriteString(fmt.Sprintf("On %s, %s wrote:\n\n", 
            original.Date, original.From))
        
        scanner := bufio.NewScanner(strings.NewReader(original.BodyText))
        for scanner.Scan() {
            line := scanner.Text()
            if strings.HasPrefix(line, ">") {
                sb.WriteString(">" + line + "\n")
            } else {
                sb.WriteString("> " + line + "\n")
            }
        }
        
        return sb.String()
    }
    
    func stripRePrefix(subject string) string {
        re := regexp.MustCompile(`^(Re:\s*)+`)
        return re.ReplaceAllString(subject, "")
    }
}
```

#### **Day 10-11: Search & Filters**

**Task 1: Implement FTS5 search**
```go
// Pseudocode for FTS5 search
func implementFTS5Search() {
    // internal/database/search.go
    type SearchEngine struct {
        db *sql.DB
    }
    
    func (s *SearchEngine) Search(query string) ([]Email, error) {
        // FIXED: Parse query for special operators
        parsed := s.parseQuery(query)
        
        sqlQuery := `
            SELECT e.* 
            FROM emails e
            JOIN email_search es ON e.rowid = es.rowid
            WHERE 1=1
        `
        
        args := []interface{}{}
        
        // Full-text search
        if parsed.Text != "" {
            sqlQuery += " AND es.email_search MATCH ?"
            args = append(args, s.prepareFTSQuery(parsed.Text))
        }
        
        // From filter
        if parsed.From != "" {
            sqlQuery += " AND e.from_email LIKE ?"
            args = append(args, "%"+parsed.From+"%")
        }
        
        // Category filter
        if parsed.Category != "" {
            sqlQuery += " AND e.ai_category = ?"
            args = append(args, parsed.Category)
        }
        
        // Date range
        if !parsed.After.IsZero() {
            sqlQuery += " AND e.date > ?"
            args = append(args, parsed.After.Format(time.RFC3339))
        }
        
        sqlQuery += " ORDER BY e.date DESC LIMIT 100"
        
        rows, err := s.db.Query(sqlQuery, args...)
        if err != nil {
            return nil, err
        }
        defer rows.Close()
        
        return s.scanEmails(rows)
    }
    
    func (s *SearchEngine) parseQuery(query string) ParsedQuery {
        result := ParsedQuery{Text: query}
        
        // Extract special filters
        patterns := map[string]*regexp.Regexp{
            "from":     regexp.MustCompile(`from:(\S+)`),
            "category": regexp.MustCompile(`category:(\S+)`),
            "after":    regexp.MustCompile(`after:(\S+)`),
            "label":    regexp.MustCompile(`label:(\S+)`),
        }
        
        for key, re := range patterns {
            if matches := re.FindStringSubmatch(query); matches != nil {
                switch key {
                case "from":
                    result.From = matches[1]
                case "category":
                    result.Category = matches[1]
                case "after":
                    // Parse date
                    if t, err := parseRelativeDate(matches[1]); err == nil {
                        result.After = t
                    }
                }
                // Remove from text query
                query = re.ReplaceAllString(query, "")
            }
        }
        
        result.Text = strings.TrimSpace(query)
        return result
    }
    
    func (s *SearchEngine) prepareFTSQuery(text string) string {
        // Escape special FTS5 characters
        text = strings.ReplaceAll(text, `"`, ``)
        
        // Convert to FTS5 query
        terms := strings.Fields(text)
        ftsTerms := []string{}
        
        for _, term := range terms {
            ftsTerms = append(ftsTerms, fmt.Sprintf(`"%s"*`, term))
        }
        
        return strings.Join(ftsTerms, " AND ")
    }
}
```

**Task 2: Add filter commands**
```go
// Pseudocode for filter system
func implementFilters() {
    // internal/ui/filters.go
    type Filter struct {
        ID       string
        Name     string
        Query    string
        Color    lipgloss.Color
        Shortcut string
    }
    
    var DefaultFilters = []Filter{
        {
            ID:       "unread",
            Name:     "Unread",
            Query:    "is:unread",
            Color:    lipgloss.Color("220"),
            Shortcut: "u",
        },
        {
            ID:       "financial",
            Name:     "Financial",
            Query:    "category:financial",
            Color:    lipgloss.Color("226"),
            Shortcut: "f",
        },
        {
            ID:       "waiting",
            Name:     "Waiting",
            Query:    "category:waiting",
            Color:    lipgloss.Color("214"),
            Shortcut: "w",
        },
        {
            ID:       "today",
            Name:     "Today",
            Query:    "after:today",
            Color:    lipgloss.Color("87"),
            Shortcut: "t",
        },
    }
    
    type FilterBar struct {
        filters  []Filter
        active   string
        onSelect func(Filter)
    }
    
    func (f FilterBar) View() string {
        var items []string
        
        for _, filter := range f.filters {
            style := lipgloss.NewStyle().
                Foreground(filter.Color).
                Padding(0, 1)
            
            if filter.ID == f.active {
                style = style.
                    Background(lipgloss.Color("237")).
                    Bold(true)
            }
            
            label := fmt.Sprintf("[%s] %s", filter.Shortcut, filter.Name)
            items = append(items, style.Render(label))
        }
        
        return lipgloss.JoinHorizontal(lipgloss.Top, items...)
    }
    
    func (f FilterBar) HandleKey(key string) (FilterBar, tea.Cmd) {
        for _, filter := range f.filters {
            if key == filter.Shortcut {
                f.active = filter.ID
                return f, func() tea.Msg {
                    return FilterSelectedMsg{Filter: filter}
                }
            }
        }
        return f, nil
    }
}
```

**Task 3: Build command palette**
```go
// Pseudocode for command palette
func buildCommandPalette() {
    // internal/ui/command_palette.go
    type CommandPalette struct {
        visible  bool
        input    textinput.Model
        commands []Command
        filtered []Command
        cursor   int
    }
    
    type Command struct {
        Name        string
        Description string
        Shortcut    string
        Action      func() tea.Msg
    }
    
    func NewCommandPalette() CommandPalette {
        ti := textinput.New()
        ti.Placeholder = "Type a command..."
        ti.Focus()
        
        return CommandPalette{
            input: ti,
            commands: []Command{
                {Name: "Compose", Description: "Write new email", Shortcut: "c"},
                {Name: "Search", Description: "Search emails", Shortcut: "/"},
                {Name: "Archive", Description: "Archive selected", Shortcut: "e"},
                {Name: "Sync", Description: "Force sync", Shortcut: "S"},
                {Name: "Settings", Description: "Open settings", Shortcut: ","},
            },
        }
    }
    
    func (c CommandPalette) Update(msg tea.Msg) (CommandPalette, tea.Cmd) {
        switch msg := msg.(type) {
        case tea.KeyMsg:
            switch msg.String() {
            case "ctrl+k":
                c.visible = !c.visible
                if c.visible {
                    c.input.Focus()
                } else {
                    c.input.Blur()
                }
                
            case "enter":
                if c.visible && len(c.filtered) > 0 {
                    cmd := c.filtered[c.cursor]
                    c.visible = false
                    return c, cmd.Action
                }
                
            case "down":
                if c.cursor < len(c.filtered)-1 {
                    c.cursor++
                }
                
            case "up":
                if c.cursor > 0 {
                    c.cursor--
                }
                
            case "escape":
                c.visible = false
            }
        }
        
        // Update input
        var cmd tea.Cmd
        c.input, cmd = c.input.Update(msg)
        
        // Filter commands
        c.filterCommands()
        
        return c, cmd
    }
    
    func (c *CommandPalette) filterCommands() {
        query := strings.ToLower(c.input.Value())
        c.filtered = []Command{}
        
        for _, cmd := range c.commands {
            name := strings.ToLower(cmd.Name)
            desc := strings.ToLower(cmd.Description)
            
            if strings.Contains(name, query) || 
               strings.Contains(desc, query) {
                c.filtered = append(c.filtered, cmd)
            }
        }
        
        // Reset cursor if out of bounds
        if c.cursor >= len(c.filtered) {
            c.cursor = max(0, len(c.filtered)-1)
        }
    }
    
    func (c CommandPalette) View() string {
        if !c.visible {
            return ""
        }
        
        // Modal overlay
        width := 60
        height := min(len(c.filtered)+3, 15)
        
        modal := lipgloss.NewStyle().
            Border(lipgloss.RoundedBorder()).
            BorderForeground(lipgloss.Color("86")).
            Width(width).
            Height(height).
            Padding(1)
        
        var content strings.Builder
        content.WriteString(c.input.View())
        content.WriteString("\n\n")
        
        for i, cmd := range c.filtered {
            style := lipgloss.NewStyle()
            if i == c.cursor {
                style = style.Background(lipgloss.Color("237"))
            }
            
            line := fmt.Sprintf("%-20s %s [%s]",
                cmd.Name, cmd.Description, cmd.Shortcut)
            content.WriteString(style.Render(line))
            content.WriteString("\n")
        }
        
        return modal.Render(content.String())
    }
}
```

#### **Day 12: AI Classification**

**Task 1: Port Grok queue exactly**
```go
// Pseudocode for Grok API queue
func portGrokQueue() {
    // internal/ai/grok_queue.go
    type GrokQueue struct {
        queue      chan QueueItem
        processing bool
        delay      time.Duration
        client     *http.Client
        apiKey     string
    }
    
    type QueueItem struct {
        emails []Email
        result chan []EmailWithCategory
        err    chan error
    }
    
    type EmailWithCategory struct {
        Email
        Category string
    }
    
    func NewGrokQueue() (*GrokQueue, error) {
        apiKey := os.Getenv("GROK_API_KEY")
        if apiKey == "" {
            // Graceful degradation - AI features disabled
            return nil, fmt.Errorf("GROK_API_KEY not set, AI features disabled")
        }
        
        return &GrokQueue{
            queue:  make(chan QueueItem, 100),
            delay:  200 * time.Millisecond,
            client: &http.Client{Timeout: 10 * time.Second},
            apiKey: apiKey,
        }, nil
    }
    
    func (g *GrokQueue) Start() {
        go g.processLoop()
    }
    
    func (g *GrokQueue) processLoop() {
        for item := range g.queue {
            result, err := g.classifyBatch(item.emails)
            
            // FIXED: Prevent deadlock with timeout
            select {
            case item.result <- result:
                // Success
            case item.err <- err:
                // Error
            case <-time.After(5 * time.Second):
                // Timeout - prevent deadlock
                log.Printf("Queue item timeout - receiver may be blocked")
            }
            
            time.Sleep(g.delay)
        }
    }
    
    func (g *GrokQueue) ClassifyEmails(emails []Email) ([]EmailWithCategory, error) {
        item := QueueItem{
            emails: emails,
            result: make(chan []EmailWithCategory),
            err:    make(chan error),
        }
        
        g.queue <- item
        
        select {
        case result := <-item.result:
            return result, nil
        case err := <-item.err:
            return nil, err
        case <-time.After(30 * time.Second):
            return nil, fmt.Errorf("classification timeout")
        }
    }
    
    func (g *GrokQueue) classifyBatch(emails []Email) ([]EmailWithCategory, error) {
        // Build prompt
        prompt := g.buildBatchPrompt(emails)
        
        // Call Grok API
        requestBody := map[string]interface{}{
            "messages": []map[string]string{
                {"role": "user", "content": prompt},
            },
            "model":       "grok-3-mini-fast",
            "temperature": 0,
        }
        
        jsonBody, _ := json.Marshal(requestBody)
        
        req, _ := http.NewRequest("POST",
            "https://api.x.ai/v1/chat/completions",
            bytes.NewBuffer(jsonBody))
        
        req.Header.Set("Authorization", "Bearer "+g.apiKey)
        req.Header.Set("Content-Type", "application/json")
        
        resp, err := g.client.Do(req)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()
        
        // Parse response
        var apiResp GrokResponse
        json.NewDecoder(resp.Body).Decode(&apiResp)
        
        // Extract categories
        return g.parseCategories(emails, apiResp.Choices[0].Message.Content)
    }
    
    func (g *GrokQueue) buildBatchPrompt(emails []Email) string {
        var sb strings.Builder
        
        sb.WriteString("Classify these emails into categories:\n")
        sb.WriteString("Categories: financial, security, newsletter, ")
        sb.WriteString("reminder, receipt, marketing, waiting, other\n\n")
        
        for i, email := range emails {
            // Redact sensitive info
            body := g.redactSensitive(email.BodyText)
            if len(body) > 300 {
                body = body[:300]
            }
            
            sb.WriteString(fmt.Sprintf("Email %d:\n", i+1))
            sb.WriteString(fmt.Sprintf("From: %s\n", email.From))
            sb.WriteString(fmt.Sprintf("Subject: %s\n", email.Subject))
            sb.WriteString(fmt.Sprintf("Body: %s\n\n", body))
        }
        
        sb.WriteString("Return only the category names in order, one per line.")
        
        return sb.String()
    }
}
```

**Task 2: Add classification to sync**
```go
// Pseudocode for classification during sync
func addClassificationToSync() {
    // internal/gmail/sync_with_ai.go
    type SyncerWithAI struct {
        *Syncer
        classifier *GrokQueue
    }
    
    func (s *SyncerWithAI) SyncAndClassify() error {
        // Regular sync
        emails, err := s.PrioritySync()
        if err != nil {
            return err
        }
        
        // Filter unclassified
        var unclassified []Email
        for _, email := range emails {
            if email.Category == "" {
                unclassified = append(unclassified, email)
            }
        }
        
        if len(unclassified) == 0 {
            return nil
        }
        
        // Classify in batches (non-blocking)
        go s.classifyInBackground(unclassified)
        
        return nil
    }
    
    func (s *SyncerWithAI) classifyInBackground(emails []Email) {
        batchSize := 10
        
        for i := 0; i < len(emails); i += batchSize {
            end := min(i+batchSize, len(emails))
            batch := emails[i:end]
            
            classified, err := s.classifier.ClassifyEmails(batch)
            if err != nil {
                log.Printf("Classification error: %v", err)
                continue
            }
            
            // Update database
            for _, email := range classified {
                s.db.UpdateCategory(email.ID, email.Category)
            }
            
            // Notify UI
            s.sendUpdate(ClassificationUpdate{
                Processed: end,
                Total:     len(emails),
            })
        }
    }
}
```

**Task 3: Display categories with gradients**
```go
// Pseudocode for category display
func displayCategoriesWithGradients() {
    // internal/ui/category_display.go
    type CategoryIndicator struct {
        category  string
        animated  bool
        frame     int
    }
    
    func (c CategoryIndicator) Render() string {
        style := CategoryStyles[c.category]
        
        if c.animated {
            return c.renderAnimated()
        }
        
        // Static gradient block
        return style.Render("█")
    }
    
    func (c CategoryIndicator) renderAnimated() string {
        switch c.category {
        case "waiting":
            // Pulsing effect
            intensity := (math.Sin(float64(c.frame)*0.1) + 1) / 2
            color := lerpColor("#FFFFFF", "#FFFF00", intensity)
            return lipgloss.NewStyle().
                Foreground(lipgloss.Color(color)).
                Render("●")
                
        case "security":
            // Blinking alert
            if c.frame%20 < 10 {
                return lipgloss.NewStyle().
                    Foreground(lipgloss.Color("#FF0000")).
                    Bold(true).
                    Render("▲")
            }
            return " "
            
        case "financial":
            // Shimmer effect
            offset := c.frame % 3
            colors := []string{"#FFD700", "#90EE90", "#00FF7F"}
            return lipgloss.NewStyle().
                Foreground(lipgloss.Color(colors[offset])).
                Render("$")
                
        default:
            return CategoryStyles[c.category].Render("■")
        }
    }
    
    // Category stats bar
    type CategoryStats struct {
        counts map[string]int
    }
    
    func (cs CategoryStats) View() string {
        var items []string
        
        categories := []string{
            "financial", "security", "waiting", 
            "reminder", "newsletter", "receipt",
        }
        
        for _, cat := range categories {
            count := cs.counts[cat]
            if count == 0 {
                continue
            }
            
            style := CategoryStyles[cat]
            indicator := style.Render("●")
            text := fmt.Sprintf("%s %d", indicator, count)
            
            items = append(items, text)
        }
        
        return lipgloss.JoinHorizontal(
            lipgloss.Top,
            items...,
        )
    }
}
```

#### **Day 13-14: Testing & Release**

**Task 1: Cross-terminal testing**
```go
// Pseudocode for terminal compatibility testing
func implementTerminalTesting() {
    // internal/ui/terminal_test.go
    func TestTerminalCompatibility(t *testing.T) {
        terminals := []struct {
            name string
            env  map[string]string
        }{
            {"iTerm2", map[string]string{"TERM": "xterm-256color"}},
            {"Terminal.app", map[string]string{"TERM": "xterm-256color"}},
            {"Alacritty", map[string]string{"TERM": "alacritty"}},
            {"Windows Terminal", map[string]string{"TERM": "xterm-256color"}},
            {"Linux Console", map[string]string{"TERM": "linux"}},
        }
        
        for _, term := range terminals {
            t.Run(term.name, func(t *testing.T) {
                // Set environment
                for k, v := range term.env {
                    os.Setenv(k, v)
                }
                
                // Test color support
                if !termenv.HasTrueColor() && 
                   !termenv.Has256Colors() {
                    t.Logf("%s: Limited color support", term.name)
                }
                
                // Test gradient rendering
                gradient := lipgloss.NewStyle().
                    Foreground(lipgloss.AdaptiveColor{
                        Light: "#FFD700",
                        Dark:  "#00FF7F",
                    })
                
                output := gradient.Render("Test")
                if output == "" {
                    t.Errorf("%s: Gradient rendering failed", term.name)
                }
                
                // Test unicode support
                if !isUnicodeSupported() {
                    t.Logf("%s: Limited unicode support", term.name)
                }
            })
        }
    }
    
    func isUnicodeSupported() bool {
        // Check if terminal supports unicode
        term := os.Getenv("TERM")
        lang := os.Getenv("LANG")
        
        return strings.Contains(lang, "UTF-8") ||
               strings.Contains(lang, "utf8")
    }
}
```

**Task 2: Performance optimization**
```go
// Pseudocode for performance optimization
func optimizePerformance() {
    // internal/performance/profiling.go
    func ProfileStartup() {
        start := time.Now()
        
        // Measure each component
        dbStart := time.Now()
        db, _ := NewEmailDB("emails.db")
        dbTime := time.Since(dbStart)
        
        uiStart := time.Now()
        ui := NewUI()
        uiTime := time.Since(uiStart)
        
        totalTime := time.Since(start)
        
        log.Printf("Startup times: DB=%v, UI=%v, Total=%v",
            dbTime, uiTime, totalTime)
        
        if totalTime > 20*time.Millisecond {
            log.Printf("WARNING: Startup exceeds 20ms target")
        }
    }
    
    func OptimizeDatabase() {
        // Optimize SQLite settings
        db.Exec("PRAGMA synchronous = NORMAL")
        db.Exec("PRAGMA cache_size = 10000")
        db.Exec("PRAGMA temp_store = MEMORY")
        db.Exec("PRAGMA mmap_size = 30000000000")
        
        // Analyze tables for query optimization
        db.Exec("ANALYZE emails")
        db.Exec("ANALYZE email_search")
    }
    
    func BenchmarkOperations() {
        operations := []struct {
            name string
            fn   func()
        }{
            {"Search 100k emails", benchSearch},
            {"Render 1000 emails", benchRender},
            {"Classify 100 emails", benchClassify},
            {"Sync 1000 emails", benchSync},
        }
        
        for _, op := range operations {
            start := time.Now()
            op.fn()
            duration := time.Since(start)
            
            log.Printf("%s: %v", op.name, duration)
        }
    }
}
```

**Task 3: Documentation & release**
```go
// Pseudocode for release preparation
func prepareRelease() {
    // build.go
    func BuildReleases() {
        platforms := []struct {
            os   string
            arch string
        }{
            {"darwin", "amd64"},  // macOS Intel
            {"darwin", "arm64"},  // macOS M1
            {"linux", "amd64"},   // Linux x64
            {"linux", "arm64"},   // Linux ARM
            {"windows", "amd64"}, // Windows x64
        }
        
        version := getVersion()
        
        for _, platform := range platforms {
            // FIXED: Define output variable properly
            output := fmt.Sprintf("tmail-%s-%s-%s",
                version, platform.os, platform.arch)
            
            if platform.os == "windows" {
                output += ".exe"
            }
            
            cmd := exec.Command("go", "build",
                "-ldflags", "-s -w", // Strip debug info
                "-tags", "sqlite_fts5", // FIXED: Add required build tags
                "-o", output,
                "./cmd/tmail")
            
            cmd.Env = append(os.Environ(),
                fmt.Sprintf("GOOS=%s", platform.os),
                fmt.Sprintf("GOARCH=%s", platform.arch))
            
            if err := cmd.Run(); err != nil {
                log.Printf("Build failed for %s/%s: %v",
                    platform.os, platform.arch, err)
            }
            
            // Compress
            compressExecutable(output)
        }
    }
    
    func GenerateDocumentation() {
        // README.md
        readme := `# TerminalMail
        
A blazing fast terminal email client with AI-powered classification.

## Features
- ⚡ <20ms startup time
- 🎨 Beautiful ANSI gradients
- 🤖 AI email classification
- 📧 Full Gmail integration
- ⌨️ Vim-style navigation

## Installation

### macOS
brew tap yourusername/terminalmail
brew install tmail

### Linux
curl -L https://github.com/yourusername/terminalmail/releases/latest/download/tmail-linux-amd64.tar.gz | tar xz
sudo mv tmail /usr/local/bin/

### From Source
go install github.com/yourusername/terminalmail/cmd/tmail@latest

## Usage
tmail              # Launch UI
tmail compose      # Compose new email
tmail search query # Search emails

## Configuration
export GOOGLE_CLIENT_ID=your_client_id
export GOOGLE_CLIENT_SECRET=your_client_secret
export GROK_API_KEY=your_grok_key

## Keyboard Shortcuts
j/k     - Navigate up/down
Enter   - Open email
c       - Compose
r       - Reply
/       - Search
q       - Quit
`
        
        os.WriteFile("README.md", []byte(readme), 0644)
    }
}
```

## 🚀 CLI Interface

```bash
# Main commands (simple, memorable)
tmail                   # Open interactive TUI
tmail inbox            # View inbox
tmail compose          # Compose new email
tmail search "query"   # Search emails

# Configuration
tmail --emoji          # Enable emoji mode
tmail --theme cyber    # Set theme

# Power user features
tmail list --json      # JSON output for piping
tmail export           # Backup emails
tmail sync             # Force sync
```

## 💡 Key Technical Decisions

1. **Go over Node/TypeScript**: 5x faster, single binary, better for CLI
2. **Bubbletea over raw terminal**: Modern TUI framework with good docs
3. **Lipgloss for styling**: Beautiful gradients without complexity
4. **SQLite unchanged**: Proven to work at scale in AI Mailbox
5. **2-phase sync unchanged**: Users need recent emails immediately

## 📝 Development Guidelines

1. **Start with working AI Mailbox patterns** - Don't reinvent
2. **Test in multiple terminals daily** - Compatibility is critical
3. **Gradient everything** - It's our visual signature
4. **Keep state simple** - 5-7 fields maximum
5. **Ship in 2-3 weeks** - Go is slightly slower to develop but worth it

## 🎬 Final Notes

This plan represents the **complete pivot to Go** while preserving all proven patterns from AI Mailbox. The pseudocode for each task provides clear implementation guidance suitable for GitHub issues.

**No more React confusion. No more TypeScript. Just clean, fast Go.**

---

*This document supersedes ALL previous documentation. This is the final architecture.*