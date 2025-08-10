# TerminalMail - Technical Implementation Guide

## 🏛️ System Architecture

### Component Diagram
```
┌──────────────────────────────────────────────────────────┐
│                     Terminal UI Layer                     │
│                    (Blessed Components)                   │
├──────────────────────────────────────────────────────────┤
│                    Application Core                       │
│         EmailManager │ SyncEngine │ AIProcessor          │
├──────────────────────────────────────────────────────────┤
│                     Data Access Layer                     │
│              SQLite │ FTS5 │ Cache Manager               │
├──────────────────────────────────────────────────────────┤
│                    External Services                      │
│          Gmail API │ Grok API │ Keytar Storage           │
└──────────────────────────────────────────────────────────┘
```

## 🎨 Terminal UI Implementation

### Core UI Framework Setup
```typescript
// src/ui/app.ts
import blessed from 'blessed';
import contrib from 'blessed-contrib';

export class TerminalMailApp {
  private screen: blessed.Widgets.Screen;
  private emailList: blessed.Widgets.ListTable;
  private statusBar: blessed.Widgets.Box;
  private commandBar: blessed.Widgets.Textbox;
  
  constructor() {
    this.screen = blessed.screen({
      smartCSR: true,
      title: 'TerminalMail',
      fullUnicode: true,
      dockBorders: true,
      warnings: true
    });
    
    this.setupLayout();
    this.bindKeyHandlers();
  }
  
  private setupLayout() {
    // Email list (main area - 80% height)
    this.emailList = blessed.listtable({
      top: 2,
      left: 0,
      width: '100%',
      height: '80%',
      border: 'line',
      align: 'left',
      tags: true,
      keys: true,
      vi: true,
      mouse: true,
      style: {
        border: { fg: 'cyan' },
        header: { fg: 'bright-cyan', bold: true },
        cell: { selected: { bg: 'blue' } }
      }
    });
    
    // Status bar (top)
    this.statusBar = blessed.box({
      top: 0,
      left: 0,
      width: '100%',
      height: 1,
      content: '{cyan-fg}📧 TerminalMail{/} │ {yellow-fg}Inbox (142){/} │ 🔄 Synced 2m ago',
      tags: true,
      style: {
        fg: 'white',
        bg: 'black'
      }
    });
    
    // Command bar (bottom)
    this.commandBar = blessed.textbox({
      bottom: 0,
      left: 0,
      width: '100%',
      height: 1,
      keys: true,
      mouse: true,
      inputOnFocus: false,
      style: {
        fg: 'white',
        bg: 'black'
      }
    });
    
    this.screen.append(this.emailList);
    this.screen.append(this.statusBar);
    this.screen.append(this.commandBar);
  }
}
```

### Email List Renderer
```typescript
// src/ui/components/email-list.ts
interface EmailRow {
  id: string;
  read: boolean;
  category?: string;
  from: string;
  subject: string;
  date: string;
}

class EmailListComponent {
  private categoryIcons = {
    financial: '💰',
    security: '🔒',
    newsletter: '📰',
    reminder: '⏰',
    receipt: '🧾',
    marketing: '📧',
    waiting: '⏳'
  };
  
  private categoryColors = {
    financial: '{green-fg}',
    security: '{red-fg}',
    newsletter: '{blue-fg}',
    reminder: '{yellow-fg}',
    receipt: '{magenta-fg}',
    marketing: '{cyan-fg}',
    waiting: '{white-fg}'
  };
  
  renderEmailList(emails: EmailRow[]): string[][] {
    const headers = ['', 'Tag', 'From', 'Subject', 'Date'];
    const rows = emails.map(email => [
      email.read ? ' ' : '●',
      this.renderCategory(email.category),
      this.truncate(email.from, 20),
      email.read ? email.subject : `{bold}${email.subject}{/bold}`,
      this.formatDate(email.date)
    ]);
    
    return [headers, ...rows];
  }
  
  private renderCategory(category?: string): string {
    if (!category) return '  ';
    const icon = this.categoryIcons[category] || '📧';
    const color = this.categoryColors[category] || '{white-fg}';
    return `${color}${icon}{/}`;
  }
}
```

## 🔑 Gmail OAuth CLI Flow

### OAuth Implementation
```typescript
// src/auth/gmail-oauth.ts
import { OAuth2Client } from 'google-auth-library';
import express from 'express';
import open from 'open';
import keytar from 'keytar';

export class GmailAuthManager {
  private oauth2Client: OAuth2Client;
  private readonly SCOPES = [
    'https://www.googleapis.com/auth/gmail.readonly',
    'https://www.googleapis.com/auth/gmail.send',
    'https://www.googleapis.com/auth/gmail.modify'
  ];
  
  async authenticate(): Promise<void> {
    // Check for existing token
    const existingToken = await keytar.getPassword('terminalmail', 'gmail_token');
    if (existingToken) {
      this.oauth2Client.setCredentials(JSON.parse(existingToken));
      return;
    }
    
    // Start local server for OAuth callback
    const server = express();
    const port = await this.findFreePort();
    
    server.get('/oauth/callback', async (req, res) => {
      const code = req.query.code as string;
      const { tokens } = await this.oauth2Client.getToken(code);
      
      // Store token securely
      await keytar.setPassword('terminalmail', 'gmail_token', JSON.stringify(tokens));
      
      res.send('Authentication successful! You can close this window.');
      process.exit(0);
    });
    
    const authUrl = this.oauth2Client.generateAuthUrl({
      access_type: 'offline',
      scope: this.SCOPES,
      redirect_uri: `http://localhost:${port}/oauth/callback`
    });
    
    console.log('Opening browser for authentication...');
    await open(authUrl);
    
    server.listen(port);
  }
}
```

## 📨 Email Composition with $EDITOR

### Editor Integration
```typescript
// src/compose/editor.ts
import { spawn } from 'child_process';
import fs from 'fs/promises';
import os from 'os';
import path from 'path';

export class EmailComposer {
  async composeWithEditor(draft?: EmailDraft): Promise<EmailContent> {
    const editor = process.env.EDITOR || 'vim';
    const tempFile = path.join(os.tmpdir(), `terminalmail-${Date.now()}.md`);
    
    // Create template
    const template = this.createTemplate(draft);
    await fs.writeFile(tempFile, template);
    
    // Open editor
    return new Promise((resolve, reject) => {
      const child = spawn(editor, [tempFile], {
        stdio: 'inherit',
        shell: true
      });
      
      child.on('exit', async (code) => {
        if (code === 0) {
          const content = await fs.readFile(tempFile, 'utf-8');
          const parsed = this.parseEmailContent(content);
          await fs.unlink(tempFile);
          resolve(parsed);
        } else {
          reject(new Error('Editor exited with error'));
        }
      });
    });
  }
  
  private createTemplate(draft?: EmailDraft): string {
    return `# Email Composition
# Lines starting with # will be ignored
# Use markdown for formatting

To: ${draft?.to || ''}
Cc: ${draft?.cc || ''}
Subject: ${draft?.subject || ''}

---

${draft?.body || ''}

---
# Save and exit to send, or exit without saving to cancel
`;
  }
  
  private parseEmailContent(content: string): EmailContent {
    const lines = content.split('\n').filter(l => !l.startsWith('#'));
    const headers: any = {};
    let bodyStart = 0;
    
    for (let i = 0; i < lines.length; i++) {
      if (lines[i] === '---') {
        bodyStart = i + 1;
        break;
      }
      const [key, ...valueParts] = lines[i].split(':');
      if (key && valueParts.length > 0) {
        headers[key.toLowerCase().trim()] = valueParts.join(':').trim();
      }
    }
    
    const body = lines.slice(bodyStart).join('\n').replace(/^---$/gm, '').trim();
    
    return {
      to: headers.to,
      cc: headers.cc,
      subject: headers.subject,
      body: this.markdownToHtml(body)
    };
  }
}
```

## 🤖 AI Classification Pipeline

### Background AI Processor
```typescript
// src/ai/classifier.ts
import { EventEmitter } from 'events';
import PQueue from 'p-queue';

export class AIEmailClassifier extends EventEmitter {
  private queue: PQueue;
  private grokClient: GrokClient;
  
  constructor() {
    super();
    this.queue = new PQueue({ 
      concurrency: 2,
      interval: 1000,
      intervalCap: 5  // Rate limiting
    });
  }
  
  async classifyEmails(emails: Email[]): Promise<void> {
    for (const email of emails) {
      this.queue.add(async () => {
        try {
          const category = await this.classifyEmail(email);
          await this.db.updateEmailCategory(email.id, category);
          
          this.emit('classified', { 
            emailId: email.id, 
            category,
            remaining: this.queue.size 
          });
        } catch (error) {
          this.emit('error', { emailId: email.id, error });
        }
      });
    }
  }
  
  private async classifyEmail(email: Email): Promise<string> {
    const prompt = `Classify this email into one category:
    - financial (payments, invoices, financial transactions)
    - security (security alerts, 2FA, password resets)
    - newsletter (newsletters, digests, subscriptions)
    - reminder (calendar, tasks, deadlines)
    - receipt (purchase confirmations, orders)
    - marketing (promotional, sales, ads)
    - waiting (requires response or action)
    - other
    
    Email:
    From: ${email.from}
    Subject: ${email.subject}
    Body: ${email.body.substring(0, 500)}
    
    Respond with only the category name.`;
    
    const response = await this.grokClient.complete(prompt);
    return response.trim().toLowerCase();
  }
}
```

## 🔍 Full-Text Search Implementation

### FTS5 Search Engine
```typescript
// src/search/fts-engine.ts
export class EmailSearchEngine {
  private db: Database;
  
  async search(query: string): Promise<Email[]> {
    // Parse advanced search syntax
    const parsed = this.parseQuery(query);
    
    let sql = `
      SELECT e.* 
      FROM emails e
      JOIN email_search es ON e.id = es.id
      WHERE 1=1
    `;
    
    const params: any[] = [];
    
    // Full-text search
    if (parsed.text) {
      sql += ` AND es.email_search MATCH ?`;
      params.push(this.prepareFTSQuery(parsed.text));
    }
    
    // Filters
    if (parsed.from) {
      sql += ` AND e.from_email LIKE ?`;
      params.push(`%${parsed.from}%`);
    }
    
    if (parsed.category) {
      sql += ` AND e.ai_category = ?`;
      params.push(parsed.category);
    }
    
    if (parsed.after) {
      sql += ` AND e.date > ?`;
      params.push(parsed.after);
    }
    
    sql += ` ORDER BY e.date DESC LIMIT 100`;
    
    return this.db.prepare(sql).all(...params);
  }
  
  private parseQuery(query: string): SearchQuery {
    const result: SearchQuery = { text: query };
    
    // Extract special filters
    const filters = {
      from: /from:(\S+)/i,
      category: /category:(\S+)/i,
      after: /after:(\S+)/i,
      label: /label:(\S+)/i
    };
    
    for (const [key, regex] of Object.entries(filters)) {
      const match = query.match(regex);
      if (match) {
        result[key] = match[1];
        query = query.replace(regex, '').trim();
      }
    }
    
    result.text = query;
    return result;
  }
  
  private prepareFTSQuery(text: string): string {
    // Escape special FTS5 characters
    return text
      .replace(/['"]/g, '')
      .split(/\s+/)
      .map(term => `"${term}"*`)
      .join(' AND ');
  }
}
```

## ⌨️ Command System

### Vim-Style Commands
```typescript
// src/commands/command-handler.ts
export class CommandHandler {
  private commands: Map<string, CommandFunction> = new Map();
  
  constructor(private app: TerminalMailApp) {
    this.registerCommands();
  }
  
  private registerCommands() {
    // Navigation
    this.register('j', () => this.app.moveCursor(1));
    this.register('k', () => this.app.moveCursor(-1));
    this.register('g', () => this.app.moveCursor('top'));
    this.register('G', () => this.app.moveCursor('bottom'));
    
    // Actions
    this.register('Enter', () => this.app.openEmail());
    this.register('c', () => this.app.compose());
    this.register('r', () => this.app.reply());
    this.register('a', () => this.app.archive());
    this.register('d', () => this.app.delete());
    
    // Search
    this.register('/', () => this.app.startSearch());
    this.register('n', () => this.app.nextSearchResult());
    this.register('N', () => this.app.prevSearchResult());
    
    // Filters
    this.register(':split', (args) => this.app.createSplit(args));
    this.register(':filter', (args) => this.app.applyFilter(args));
    
    // AI
    this.register(':classify', () => this.app.classifyUnprocessed());
    this.register(':summarize', () => this.app.summarizeEmail());
  }
  
  async execute(input: string): Promise<void> {
    const [cmd, ...args] = input.split(' ');
    const handler = this.commands.get(cmd);
    
    if (handler) {
      await handler(args.join(' '));
    } else {
      this.app.showError(`Unknown command: ${cmd}`);
    }
  }
}
```

## 🎨 ANSI Color System

### Theme Engine
```typescript
// src/ui/themes.ts
export interface Theme {
  name: string;
  colors: {
    // UI Elements
    background: string;
    foreground: string;
    selection: string;
    border: string;
    statusBar: string;
    
    // Email Categories
    financial: string;
    security: string;
    newsletter: string;
    reminder: string;
    receipt: string;
    marketing: string;
    waiting: string;
    
    // States
    unread: string;
    selected: string;
    error: string;
    success: string;
  };
}

export const themes = {
  default: {
    name: 'Default',
    colors: {
      background: '\x1b[40m',      // Black
      foreground: '\x1b[37m',      // White
      selection: '\x1b[44m',       // Blue bg
      border: '\x1b[36m',          // Cyan
      statusBar: '\x1b[100m',      // Bright black bg
      
      financial: '\x1b[92m',       // Bright green
      security: '\x1b[91m',        // Bright red
      newsletter: '\x1b[94m',      // Bright blue
      reminder: '\x1b[93m',        // Bright yellow
      receipt: '\x1b[95m',         // Bright magenta
      marketing: '\x1b[96m',       // Bright cyan
      waiting: '\x1b[97m',         // Bright white
      
      unread: '\x1b[1m',           // Bold
      selected: '\x1b[7m',         // Reverse
      error: '\x1b[31m',           // Red
      success: '\x1b[32m'          // Green
    }
  },
  
  dracula: {
    name: 'Dracula',
    colors: {
      background: '\x1b[48;2;40;42;54m',
      foreground: '\x1b[38;2;248;248;242m',
      selection: '\x1b[48;2;68;71;90m',
      border: '\x1b[38;2;139;233;253m',
      statusBar: '\x1b[48;2;68;71;90m',
      
      financial: '\x1b[38;2;80;250;123m',
      security: '\x1b[38;2;255;85;85m',
      newsletter: '\x1b[38;2;139;233;253m',
      reminder: '\x1b[38;2;241;250;140m',
      receipt: '\x1b[38;2;255;121;198m',
      marketing: '\x1b[38;2;189;147;249m',
      waiting: '\x1b[38;2;248;248;242m',
      
      unread: '\x1b[1m',
      selected: '\x1b[48;2;98;114;164m',
      error: '\x1b[38;2;255;85;85m',
      success: '\x1b[38;2;80;250;123m'
    }
  }
};

export class ThemeManager {
  private currentTheme: Theme;
  
  constructor(themeName: string = 'default') {
    this.currentTheme = themes[themeName] || themes.default;
  }
  
  apply(element: string, text: string): string {
    const color = this.currentTheme.colors[element];
    return `${color}${text}\x1b[0m`;
  }
  
  icon(category: string): string {
    const icons = {
      financial: '💰',
      security: '🔒',
      newsletter: '📰',
      reminder: '⏰',
      receipt: '🧾',
      marketing: '📧',
      waiting: '⏳',
      default: '📧'
    };
    
    const icon = icons[category] || icons.default;
    const color = this.currentTheme.colors[category] || this.currentTheme.colors.foreground;
    
    return `${color}${icon}\x1b[0m`;
  }
}
```

## 🚀 Performance Optimizations

### Virtual Scrolling
```typescript
// src/ui/virtual-scroll.ts
export class VirtualScrollManager {
  private viewport: number = 50;  // Visible rows
  private buffer: number = 10;    // Extra rows for smooth scrolling
  private offset: number = 0;
  private totalItems: number = 0;
  
  async getVisibleItems<T>(
    fetchFn: (offset: number, limit: number) => Promise<T[]>
  ): Promise<T[]> {
    const start = Math.max(0, this.offset - this.buffer);
    const limit = this.viewport + (this.buffer * 2);
    
    return await fetchFn(start, limit);
  }
  
  scrollTo(position: number): void {
    this.offset = Math.max(0, Math.min(position, this.totalItems - this.viewport));
  }
  
  pageDown(): void {
    this.scrollTo(this.offset + this.viewport);
  }
  
  pageUp(): void {
    this.scrollTo(this.offset - this.viewport);
  }
}
```

### Background Sync Worker
```typescript
// src/sync/background-worker.ts
import { Worker } from 'worker_threads';

export class SyncWorker {
  private worker: Worker;
  private syncInterval: NodeJS.Timer;
  
  start(): void {
    this.worker = new Worker('./sync-worker.js');
    
    this.worker.on('message', (msg) => {
      if (msg.type === 'sync-complete') {
        this.emit('sync-complete', msg.data);
      } else if (msg.type === 'new-emails') {
        this.showNotification(msg.count);
      }
    });
    
    // Auto-sync every minute
    this.syncInterval = setInterval(() => {
      this.worker.postMessage({ type: 'sync' });
    }, 60000);
    
    // Initial sync
    this.worker.postMessage({ type: 'sync' });
  }
  
  private showNotification(count: number): void {
    // Terminal bell for new emails
    process.stdout.write('\x07');
    
    // Update status bar
    this.updateStatus(`📬 ${count} new emails`);
  }
}
```

## 📦 Package Structure

### Complete Package.json
```json
{
  "name": "terminalmail",
  "version": "1.0.0",
  "description": "AI-powered terminal email client",
  "main": "dist/index.js",
  "bin": {
    "tmail": "./dist/cli.js"
  },
  "scripts": {
    "build": "tsc",
    "dev": "tsc --watch",
    "start": "node dist/cli.js",
    "test": "jest",
    "lint": "eslint src/**/*.ts"
  },
  "dependencies": {
    "blessed": "^0.1.81",
    "blessed-contrib": "^4.11.0",
    "better-sqlite3": "^9.0.0",
    "googleapis": "^126.0.0",
    "google-auth-library": "^9.0.0",
    "keytar": "^7.9.0",
    "commander": "^11.0.0",
    "marked": "^9.0.0",
    "marked-terminal": "^6.0.0",
    "open": "^9.0.0",
    "express": "^4.18.0",
    "p-queue": "^7.3.0",
    "dotenv": "^16.0.0"
  },
  "devDependencies": {
    "@types/blessed": "^0.1.22",
    "@types/node": "^20.0.0",
    "typescript": "^5.0.0",
    "eslint": "^8.0.0",
    "jest": "^29.0.0"
  }
}
```

## 🎯 CLI Interface

### Command Structure
```bash
# Main commands
tmail                        # Open interactive UI
tmail inbox                  # View inbox
tmail compose                # Compose new email
tmail search "query"         # Search emails
tmail sync                   # Force sync

# Email operations
tmail show <id>              # Display email
tmail reply <id>             # Reply to email
tmail forward <id>           # Forward email
tmail archive <id>           # Archive email
tmail delete <id>            # Delete email

# AI commands
tmail classify               # Classify unprocessed emails
tmail summarize <id>         # Get AI summary
tmail portfolio              # View portfolio tracking

# Filters and views
tmail filter financial       # Show financial emails
tmail filter unread          # Show unread emails
tmail filter from:john       # Show emails from John

# Export/Import
tmail export --format json   # Export emails as JSON
tmail export --format csv    # Export as CSV
tmail backup                 # Backup database

# Configuration
tmail config set editor vim  # Set preferred editor
tmail config set theme dracula # Set color theme
tmail auth gmail             # Setup Gmail authentication
```

## 🔒 Security Considerations

### Token Storage
```typescript
// Use system keychain for sensitive data
await keytar.setPassword('terminalmail', 'gmail_refresh_token', token);
await keytar.setPassword('terminalmail', 'grok_api_key', apiKey);

// Never store in plain text files
// Never log sensitive information
// Use environment variables for development only
```

### Database Encryption
```typescript
// Optional SQLCipher for encrypted database
import SQLite3 from 'better-sqlite3';

const db = new SQLite3('emails.db');
db.pragma('key = "your-encryption-key"');
db.pragma('cipher_compatibility = 4');
```

## 🚢 Distribution Strategy

### Installation Methods
```bash
# NPM (primary)
npm install -g terminalmail

# Homebrew (macOS)
brew tap terminalmail/tap
brew install terminalmail

# Direct binary
curl -L https://github.com/terminalmail/releases/latest/tmail-linux -o /usr/local/bin/tmail
chmod +x /usr/local/bin/tmail
```

### Platform Binaries
- Build with `pkg` or `nexe` for standalone executables
- Target: Linux, macOS, Windows (WSL)
- Size: ~50MB compressed

## 📈 Success Metrics & Monitoring

### Analytics (Privacy-Respecting)
```typescript
// Anonymous usage metrics only
const metrics = {
  dailyActiveUsers: true,      // Count unique days used
  featuresUsed: true,          // Which features are popular
  performanceMetrics: true,    // Startup time, sync speed
  errorReporting: false,       // Opt-in only
  emailContent: false,         // Never track
  personalInfo: false          // Never track
};
```

### Performance Targets
- **Startup**: < 100ms
- **Email list render**: < 50ms for 1000 emails
- **Search**: < 200ms for 100k emails
- **Sync**: < 5s for 100 new emails
- **Memory**: < 50MB baseline, < 100MB with 10k emails

---

This technical implementation guide provides the complete blueprint for building TerminalMail with production-ready code examples and architectural decisions.