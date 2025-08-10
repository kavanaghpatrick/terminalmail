# TerminalMail - Project Guidelines & Orchestration

## 🧘 SIMPLICITY FIRST - THE PRIME DIRECTIVE

### When facing ANY problem, ask in order:
1. **"What can I DELETE to fix this?"**
2. **"What existing code already does this?"**
3. **"Can I fix this in 1 line?"**
4. **"Is it just a typo?"**

**Core Rule: "Every bug is an opportunity to DELETE code, not add it."**

### Red Flags of Over-Engineering:
- "Just in case" error handling
- New state variables
- Coordination logic between components
- Queues/flags/timers for race conditions
- Fix longer than 10 lines

## 🚀 DEVELOPMENT WORKFLOW (FOLLOW EXACTLY)

```
1. Review existing GitHub Issue (we have 32 ready!)
2. Validate with Grok + Gemini
3. Filter feedback (reject complexity)
4. Setup environment (Go modules, database)
5. Implement with Qwen Coder ONLY
6. Test implementation
7. Code review with Grok
8. Apply only critical fixes
9. Commit and close issue
```

**🚨 CRITICAL RULES 🚨**
- **NEVER write code directly** - Use Qwen for EVERYTHING
- **NO EXCEPTIONS** - Even 1-line changes need Qwen
- **COMPLETE ALL CRITERIA** - Keep using Qwen until 100% done
- **VIOLATION** = Writing any code yourself, even "to finish up"

**⚠️ TIMEOUT WORKAROUND:**
- Claude's Bash has 2-min default (can extend to 10 min with timeout parameter)
- opencode itself doesn't timeout (--timeout flag doesn't work)
- Break large tasks into <2min chunks
- If timeout: assess progress, continue with focused commands

**IMPORTANT Todo List Rules:**
- Follow workflow order EXACTLY when creating todos
- Environment prep (step 4) MUST come before implementation (step 5)
- Each todo maps to one workflow step

## 🎯 Simplicity Filter

**✅ ACCEPT:**
- Bug/security fixes
- Prevents common failures
- Essential missing functionality
- Minimal code improvements

**❌ REJECT:**
- "Nice to have" features
- Enterprise patterns for MVP
- "Future flexibility" abstractions
- Rare edge case handling
- Premature optimization

**Filter Examples:**
```
Grok: "Add retry with exponential backoff"
→ REJECT: Simple retry sufficient for MVP

Gemini: "Handle Gmail rate limits"  
→ ACCEPT: Prevents common failure, simple fix

Grok: "Add connection pooling"
→ REJECT: SQLite handles this, premature optimization
```

## 📋 Project Structure

### We Have 32 Pre-Created Issues!
- #1-4: Foundation & Database (Days 1-2)
- #5-8: Gmail Integration (Days 3-4)
- #9-12: Email List UI (Days 5-6)
- #13-16: Email Reading (Day 7)
- #17-20: Editor Integration (Days 8-9)
- #21-24: Search & Filters (Days 10-11)
- #25-28: AI Classification (Day 12)
- #29-32: Testing & Release (Days 13-14)

Each issue contains: Go implementation, ASCII mockups, acceptance criteria, testing checklist.

### Tech Stack:
- Go 1.21+ with Bubbletea TUI
- SQLite with FTS5
- Gmail API
- OpenAI (Phase 2 only)

## 🔧 Quick Reference

### Environment Setup:
```bash
# .env file
GH_TOKEN=xxx  # GitHub API access
GOOGLE_CLIENT_ID=xxx
GOOGLE_CLIENT_SECRET=xxx
GOOGLE_REDIRECT_URI=http://localhost:8080/callback
OPENAI_API_KEY=xxx  # For AI features (Phase 2)
GROK_API_KEY=xxx  # For code reviews
GEMINI_API_KEY=xxx  # For validation

source .env  # OR: export $(cat .env | xargs)
```

### Gmail API Scopes Required:
```go
var GmailScopes = []string{
    "https://www.googleapis.com/auth/gmail.readonly",  // Read emails
    "https://www.googleapis.com/auth/gmail.send",      // Send emails
    "https://www.googleapis.com/auth/gmail.modify",    // Archive/delete
    "https://www.googleapis.com/auth/gmail.compose",   // Compose drafts
}
```

### Issue Workflow:
```bash
# Start work
gh issue edit 1 --add-label "in-progress"

# Implement with Qwen (ALWAYS use --timeout 600 for 10 minutes)
opencode run -m "cerebras/qwen-3-coder-480b" --timeout 600 "Implement issue #1: $(gh issue view 1)"

# Continue if needed
opencode run --continue --timeout 600 "Complete [specific remaining items]"

# Test
go test ./...

# Close when done
gh issue close 1 --comment "✅ All criteria met"
```

### Validation:
```bash
# Grok (technical review)
curl -X POST https://api.x.ai/v1/chat/completions \
  -H "Authorization: Bearer $GROK_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"messages": [{"role": "user", "content": "Review: [code]"}], "model": "grok-4"}'

# Gemini (completeness)
gemini --prompt "Review: [issue content]"
```

## 🎨 Bubbletea Principles

- **State in Model** - Never in Views
- **Views are pure** - Only render Model
- **Updates immutable** - Return new Model
- **Commands for I/O** - Never block Update

```go
// GOOD
type Model struct {
    emails []Email
    selected int
}

// BAD - Mixed concerns
type Model struct {
    emails []Email
    db *sql.DB  // Keep separate
}
```

## 🚨 Division of Labor

**Claude handles:**
- Dependencies (`go get`, `go mod tidy`)
- Database setup
- Testing (`go test`)
- Git operations
- Reviews and validation

**Qwen handles:**
- ALL Go code writing
- ALL file modifications
- ALL debugging code

**Examples of VIOLATIONS:**
```bash
# WRONG: Claude finishes "simple" part
# RIGHT: opencode run --continue --timeout 600 "Complete remaining criteria"

# WRONG: Claude adds debug logs
# RIGHT: opencode run --continue --timeout 600 "Add debug logs for [issue]"

# WRONG: Claude fixes typo
# RIGHT: opencode run --continue --timeout 600 "Fix typo in [file]"
```

## 📊 Progress Tracking

```bash
# Current status
gh issue list --label "in-progress"

# Progress report
echo "Completed: $(gh issue list --state closed | wc -l)/32"
echo "In Progress: $(gh issue list --label 'in-progress' | wc -l)"
echo "Remaining: $(gh issue list --state open | wc -l)"
```

## 🔄 Issue States

**Only close when:**
- ✅ Code implemented and tested
- ✅ Tests pass
- ✅ Committed to branch
- ✅ No blocking issues

**Never close when:**
- Tests failing
- Partial implementation
- Known bugs exist

## 🛠️ Debugging

```bash
# Go testing
go test -v ./...
go test -cover ./...
go test -bench=. ./...

# Bubbletea debug
# Add to main.go: tea.WithDebug()

# Database
sqlite3 ~/.terminalmail/emails.db
.tables
SELECT COUNT(*) FROM emails;
```

## 📝 Code Review Filter

After Grok reviews:

**Apply fixes for:**
- Crashes/security issues
- Data loss prevention
- Critical bugs

**Reject suggestions for:**
- Comprehensive logging
- Unit test suites (important but not blocking)
- "Enterprise" patterns
- Performance optimization without measured problems

## 🎯 Success Criteria

**Phase 1 (Do First):**
- Full email client functionality
- Gmail sync
- <200ms operations
- Virtual scrolling for 10K+ emails

**Phase 2 (Only After Phase 1):**
- AI summaries
- Smart prioritization

**Out of Scope:**
- Complex sync resumption
- Multi-vendor abstractions
- Enterprise features
- Advanced monitoring

## 📚 Complete Implementation Example

```bash
# 1. Start work on issue
ISSUE_NUM=1
gh issue edit $ISSUE_NUM --add-label "in-progress"
gh issue comment $ISSUE_NUM --body "Starting implementation"

# 2. Create feature branch  
git checkout -b "feature/issue-$ISSUE_NUM-project-setup"

# 3. Validate with external AI
ISSUE=$(gh issue view $ISSUE_NUM --json body -q .body)
# Send to Grok for technical review
# Send to Gemini for completeness
# Apply simplicity filter to feedback

# 4. Implement with Qwen (ALWAYS use --timeout 600)
opencode run -m "cerebras/qwen-3-coder-480b" \
  --timeout 600 \
  "Implement issue #$ISSUE_NUM: $(gh issue view $ISSUE_NUM)"

# 5. Continue if not complete
opencode run --continue --timeout 600 "Complete remaining acceptance criteria"

# 6. Test
go test ./...
if [ $? -eq 0 ]; then
  gh issue comment $ISSUE_NUM --body "✅ Tests passing"
else  
  gh issue comment $ISSUE_NUM --body "❌ Tests failing"
  # Use Qwen to fix, not Claude!
  opencode run --continue --timeout 600 "Fix failing tests"
fi

# 7. Commit and close
git add .
git commit -m "Implement #$ISSUE_NUM: [description]"
git push origin feature/issue-$ISSUE_NUM-project-setup
gh issue close $ISSUE_NUM --comment "✅ All criteria met"
```

---

**Remember: Simplicity > Features. Working > Perfect. Delete > Add.**