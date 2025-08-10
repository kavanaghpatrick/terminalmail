# TerminalMail Session Restart Prompt

## Copy this entire prompt to start a new session:

---

I need to continue working on the TerminalMail project - a terminal-based email client built with Go and Bubbletea. Here's the complete context:

## Project Overview
TerminalMail is a 14-day implementation project for a terminal email client with these key features:
- Terminal UI using Go + Bubbletea framework
- Gmail API integration for email sync
- SQLite with FTS5 for local storage and search
- Vim-style navigation
- Virtual scrolling for 10,000+ emails
- AI enhancement (Phase 2, after core client works)

## Current Status
We have created 32 comprehensive GitHub issues (#1-32) covering the entire implementation:
- Issues #1-4: Foundation & Database (Days 1-2)
- Issues #5-8: Gmail Integration (Days 3-4)
- Issues #9-12: Email List UI (Days 5-6)
- Issues #13-16: Email Reading (Day 7)
- Issues #17-20: Editor Integration (Days 8-9)
- Issues #21-24: Search & Filters (Days 10-11)
- Issues #25-28: AI Classification (Day 12)
- Issues #29-32: Testing & Release (Days 13-14)

Each issue contains complete Go implementations, ASCII mockups, acceptance criteria, and testing checklists.

## Critical Files to Read
Please read these files to understand the project and workflow:

1. `/Users/patrickkavanagh/terminalmail/CLAUDE.md` - Our project guidelines and orchestration rules
2. `/Users/patrickkavanagh/terminalmail/ARCHITECTURE.md` - System architecture and design decisions
3. Check the current state of GitHub issues: `gh issue list --state all`

## Working Rules (CRITICAL)
1. **NEVER write code directly** - Always use Qwen Coder via opencode
2. **Follow the 9-step workflow** in CLAUDE.md exactly
3. **Apply simplicity filter** - Reject enterprise patterns and complexity
4. **Use --timeout 600** with every opencode command
5. **Complete ALL acceptance criteria** before closing issues

## Next Steps
1. Review the open GitHub issues to see what needs implementation
2. Start with Issue #1 (Project Setup) if not complete
3. Follow the workflow in CLAUDE.md for each issue:
   - Validate with Grok/Gemini
   - Filter feedback through simplicity principles
   - Implement with Qwen Coder
   - Test and review
   - Commit and close

## Environment Check
Ensure these are set up:
```bash
# Check environment variables
cat .env  # Should have GOOGLE_CLIENT_ID, GROK_API_KEY, etc.

# Check GitHub CLI
gh auth status

# Check repository
git remote -v  # Should show github.com/kavanaghpatrick/terminalmail

# Check issues
gh issue list --state open | head -5
```

## Implementation Command Template
When starting work on an issue:
```bash
ISSUE_NUM=1  # Replace with actual issue number
gh issue view $ISSUE_NUM  # Review the issue
gh issue edit $ISSUE_NUM --add-label "in-progress"

# Implement with Qwen (ALWAYS use --timeout 600)
opencode run -m "cerebras/qwen-3-coder-480b" --timeout 600 \
  "Implement TerminalMail issue #$ISSUE_NUM: $(gh issue view $ISSUE_NUM)"
```

## Key Principles to Remember
- **Simplicity First**: Every bug is an opportunity to DELETE code, not add it
- **Phase 1 First**: Build working email client before ANY AI features
- **Use Existing Issues**: We have 32 detailed issues ready - don't create new ones
- **Test Everything**: Run `go test ./...` after every implementation
- **Qwen Only for Code**: Claude handles environment, Qwen handles ALL code

## Project Structure
```
terminalmail/
├── cmd/terminalmail/     # Main application entry point
├── pkg/
│   ├── models/           # Data models
│   ├── ui/              # Bubbletea UI components
│   ├── services/        # Gmail API, database services
│   ├── database/        # SQLite integration
│   └── config/          # Configuration management
├── CLAUDE.md            # Project guidelines (READ THIS)
├── ARCHITECTURE.md      # System design
└── .env                 # Environment variables
```

Please confirm you understand the context and are ready to continue development on TerminalMail following the strict workflow rules in CLAUDE.md.

---

## Additional Notes for Starting Fresh

If this is a completely new session and you need to understand what was previously done:

1. First run: `gh issue list --state closed` to see completed work
2. Then run: `gh issue list --state open` to see remaining work
3. Read CLAUDE.md thoroughly - it contains all workflow rules
4. Pick up from the next unopened issue in sequence

The most important rule: **NEVER write code yourself. Always use Qwen Coder with --timeout 600**