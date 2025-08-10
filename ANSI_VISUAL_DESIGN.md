# TerminalMail - Advanced ANSI Visual Design System

## 🎨 Core Visual Philosophy

**No emojis. Pure ANSI art.** Every visual element uses gradients, animations, and advanced terminal capabilities to create the most beautiful command-line email client ever built.

## 📐 ANSI Color Palette & Gradients

### 256-Color Gradient System
```typescript
// True gradient generation using 256-color ANSI
class GradientEngine {
  // Generate smooth gradients between any two RGB colors
  createGradient(startRGB: RGB, endRGB: RGB, steps: number): string[] {
    const gradientCodes: string[] = [];
    
    for (let i = 0; i < steps; i++) {
      const ratio = i / (steps - 1);
      const r = Math.round(startRGB.r + (endRGB.r - startRGB.r) * ratio);
      const g = Math.round(startRGB.g + (endRGB.g - startRGB.g) * ratio);
      const b = Math.round(startRGB.b + (endRGB.b - startRGB.b) * ratio);
      
      // 256-color ANSI escape sequence
      gradientCodes.push(`\x1b[38;2;${r};${g};${b}m`);
    }
    
    return gradientCodes;
  }
}
```

### Category Color Gradients
```typescript
const CATEGORY_GRADIENTS = {
  // Financial: Gold to Green gradient
  financial: {
    start: { r: 255, g: 215, b: 0 },   // Gold
    end: { r: 0, g: 255, b: 127 },     // Spring Green
    char: '▰' // Block character for visual weight
  },
  
  // Security: Deep Red to Bright Red pulse
  security: {
    start: { r: 139, g: 0, b: 0 },     // Dark Red
    end: { r: 255, g: 69, b: 0 },      // Orange Red
    char: '▲', // Triangle for alert
    pulse: true // Animated pulsing
  },
  
  // Newsletter: Ocean gradient
  newsletter: {
    start: { r: 0, g: 119, b: 190 },   // Deep Blue
    end: { r: 0, g: 191, b: 255 },     // Deep Sky Blue
    char: '▪'
  },
  
  // Reminder: Sunset gradient
  reminder: {
    start: { r: 255, g: 94, b: 77 },   // Sunset Orange
    end: { r: 255, g: 206, b: 84 },    // Sunset Yellow
    char: '◉',
    blink: true // Subtle blink animation
  },
  
  // Receipt: Purple to Pink
  receipt: {
    start: { r: 128, g: 0, b: 128 },   // Purple
    end: { r: 255, g: 105, b: 180 },   // Hot Pink
    char: '◈'
  },
  
  // Marketing: Cyan wave
  marketing: {
    start: { r: 0, g: 255, b: 255 },   // Cyan
    end: { r: 64, g: 224, b: 208 },    // Turquoise
    char: '◊',
    wave: true // Wave animation
  },
  
  // Waiting: White to Yellow urgency
  waiting: {
    start: { r: 255, g: 255, b: 255 }, // White
    end: { r: 255, g: 255, b: 0 },     // Yellow
    char: '⬤',
    throb: true // Throbbing animation
  }
};
```

## 🎭 Animation System

### Frame-Based Animation Engine
```typescript
class ANSIAnimator {
  private frameRate = 60; // 60 FPS for smooth animations
  private animations: Map<string, Animation> = new Map();
  
  // Pulse animation for security alerts
  pulse(text: string, gradient: Gradient): string {
    const frame = this.getFrame();
    const intensity = Math.sin(frame * 0.1) * 0.5 + 0.5;
    return this.applyGradientWithIntensity(text, gradient, intensity);
  }
  
  // Wave animation for marketing emails
  wave(text: string, gradient: Gradient): string {
    const frame = this.getFrame();
    const chars = text.split('');
    
    return chars.map((char, i) => {
      const offset = Math.sin((frame * 0.1) + (i * 0.5)) * 0.5 + 0.5;
      return this.applyGradientAtPosition(char, gradient, offset);
    }).join('');
  }
  
  // Smooth fade-in for new emails
  fadeIn(text: string, duration: number = 500): string {
    const progress = this.getAnimationProgress('fadeIn', duration);
    const opacity = Math.floor(progress * 255);
    return `\x1b[38;2;${opacity};${opacity};${opacity}m${text}\x1b[0m`;
  }
  
  // Slide-in from right for notifications
  slideIn(text: string, width: number): string {
    const progress = this.getAnimationProgress('slideIn', 300);
    const offset = Math.floor((1 - progress) * width);
    return ' '.repeat(offset) + text;
  }
}
```

### Loading Animations
```typescript
// Beautiful gradient spinners
const SPINNERS = {
  dots: {
    frames: ['⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'],
    interval: 80
  },
  
  gradient_bar: {
    frames: [
      '▰▱▱▱▱▱▱▱',
      '▰▰▱▱▱▱▱▱',
      '▰▰▰▱▱▱▱▱',
      '▰▰▰▰▱▱▱▱',
      '▰▰▰▰▰▱▱▱',
      '▰▰▰▰▰▰▱▱',
      '▰▰▰▰▰▰▰▱',
      '▰▰▰▰▰▰▰▰',
    ],
    interval: 100,
    gradient: true // Apply gradient coloring
  },
  
  orbit: {
    frames: ['◐', '◓', '◑', '◒'],
    interval: 120,
    gradient: true
  }
};
```

## 📊 Email List Visual Design

### Gradient List Headers
```typescript
// Beautiful gradient header bar
const renderHeader = () => {
  const gradient = createGradient(
    { r: 100, g: 100, b: 255 },  // Blue
    { r: 255, g: 100, b: 255 },  // Magenta
    process.stdout.columns
  );
  
  const header = '│ CAT │ FROM                │ SUBJECT                              │ DATE     │';
  return applyGradientToText(header, gradient);
};
```

### Row Rendering with Visual Indicators
```typescript
const renderEmailRow = (email: Email, index: number, selected: boolean) => {
  const colors = {
    unread: '\x1b[1m',        // Bold
    selected: '\x1b[7m',      // Reverse video
    focused: '\x1b[48;2;40;40;60m', // Subtle blue background
  };
  
  // Category indicator with gradient
  const categoryGradient = CATEGORY_GRADIENTS[email.ai_category];
  const categoryIndicator = renderCategoryIndicator(categoryGradient);
  
  // Unread indicator - animated dot
  const unreadIndicator = email.unread 
    ? animatePulse('●', { r: 0, g: 255, b: 0 }) 
    : ' ';
  
  // Thread depth visualization
  const threadIndent = '│ '.repeat(email.threadDepth);
  const threadColor = `\x1b[38;2;${100 + email.threadDepth * 30};100;100m`;
  
  return `${unreadIndicator} ${categoryIndicator} ${threadColor}${threadIndent}${email.subject}\x1b[0m`;
};
```

## 🌊 Smooth Scrolling & Transitions

### Virtual Scrolling with Gradient Fade
```typescript
class SmoothScroller {
  // Fade out items at edges for depth perception
  renderWithEdgeFade(items: string[], viewport: Viewport): string[] {
    const fadeZone = 3; // Number of rows to fade
    
    return items.map((item, index) => {
      if (index < fadeZone) {
        // Fade in from top
        const opacity = (index + 1) / fadeZone;
        return this.applyOpacity(item, opacity);
      } else if (index > items.length - fadeZone) {
        // Fade out to bottom
        const opacity = (items.length - index) / fadeZone;
        return this.applyOpacity(item, opacity);
      }
      return item;
    });
  }
  
  // Smooth scroll animation
  animateScroll(from: number, to: number, duration: number = 200) {
    const frames = Math.ceil(duration / 16); // 60fps
    const delta = (to - from) / frames;
    
    for (let i = 0; i < frames; i++) {
      const position = from + (delta * this.easeInOutCubic(i / frames));
      this.renderAtPosition(position);
      this.sleep(16);
    }
  }
}
```

## 💫 Status Bar & Progress Indicators

### Multi-Layer Status Bar
```typescript
class StatusBar {
  render(state: AppState): string {
    // Layer 1: Gradient background
    const bgGradient = createGradient(
      { r: 20, g: 20, b: 40 },
      { r: 40, g: 20, b: 60 },
      process.stdout.columns
    );
    
    // Layer 2: Animated sync indicator
    const syncIndicator = state.syncing 
      ? this.animateSync() 
      : this.renderCheckmark();
    
    // Layer 3: Email counts with color coding
    const counts = [
      { label: 'Total', value: state.total, color: { r: 100, g: 100, b: 100 } },
      { label: 'Unread', value: state.unread, color: { r: 0, g: 255, b: 0 } },
      { label: 'Waiting', value: state.waiting, color: { r: 255, g: 255, b: 0 } }
    ];
    
    // Compose the status bar
    return this.compose([
      bgGradient,
      syncIndicator,
      ...counts.map(c => this.renderCount(c))
    ]);
  }
  
  private animateSync(): string {
    const frame = Date.now() % 1000 / 1000;
    const gradient = createRadialGradient(frame);
    return `[${gradient}SYNCING${reset}]`;
  }
}
```

### Progress Bars with Live Updates
```typescript
class ProgressBar {
  render(current: number, total: number, label: string): string {
    const width = 40;
    const progress = current / total;
    const filled = Math.floor(progress * width);
    
    // Create gradient from green to blue
    const gradient = createGradient(
      { r: 0, g: 255, b: 0 },
      { r: 0, g: 100, b: 255 },
      width
    );
    
    // Build the bar with gradient
    const bar = gradient.map((color, i) => {
      if (i < filled) {
        return `${color}█`;
      } else if (i === filled) {
        return `${color}▓`;
      } else {
        return '\x1b[38;2;40;40;40m░';
      }
    }).join('');
    
    // Add percentage with animated color
    const percentage = Math.floor(progress * 100);
    const percentColor = this.getPercentageColor(percentage);
    
    return `${label} ${bar}\x1b[0m ${percentColor}${percentage}%\x1b[0m (${current}/${total})`;
  }
}
```

## 🎯 Focus & Selection Visualization

### Multi-Level Selection System
```typescript
const SELECTION_STYLES = {
  // Cursor position (keyboard navigation)
  cursor: {
    background: { r: 60, g: 60, b: 100 },
    border: '▶', // Arrow indicator
    animation: 'breathe' // Subtle pulsing
  },
  
  // Single selection
  selected: {
    background: { r: 40, g: 80, b: 120 },
    border: '┃', // Vertical bar
    animation: null
  },
  
  // Multi-selection
  multiSelected: {
    background: { r: 80, g: 60, b: 100 },
    border: '║', // Double vertical bar
    animation: 'shimmer' // Subtle shimmer effect
  },
  
  // Hover (if mouse support enabled)
  hover: {
    background: { r: 50, g: 50, b: 70 },
    border: '│', // Thin vertical bar
    animation: null
  }
};
```

## 🌈 Email Reading View

### Gradient-Enhanced Text Display
```typescript
class EmailReader {
  renderEmail(email: Email): string {
    const output: string[] = [];
    
    // Header with gradient border
    const headerGradient = createGradient(
      { r: 100, g: 100, b: 255 },
      { r: 255, g: 100, b: 100 },
      process.stdout.columns
    );
    
    output.push(this.renderBorder('╔', '═', '╗', headerGradient));
    
    // Metadata with color coding
    output.push(this.renderField('From', email.from, { r: 0, g: 200, b: 255 }));
    output.push(this.renderField('To', email.to, { r: 0, g: 255, b: 200 }));
    output.push(this.renderField('Date', email.date, { r: 200, g: 200, b: 200 }));
    
    // Subject with emphasis animation
    output.push(this.renderSubject(email.subject));
    
    // Body with syntax highlighting for quotes
    output.push(this.renderBody(email.body));
    
    // Footer with gradient
    output.push(this.renderBorder('╚', '═', '╝', headerGradient));
    
    return output.join('\n');
  }
  
  private renderBody(body: string): string {
    // Highlight quoted text with gradients
    const lines = body.split('\n');
    return lines.map(line => {
      if (line.startsWith('>')) {
        // Quote level coloring
        const level = line.match(/^>+/)[0].length;
        const color = this.getQuoteColor(level);
        return `${color}${line}\x1b[0m`;
      }
      return line;
    }).join('\n');
  }
}
```

## 🔮 Notification System

### Non-Intrusive Gradient Notifications
```typescript
class NotificationSystem {
  show(message: string, type: 'success' | 'error' | 'info'): void {
    const styles = {
      success: {
        gradient: [{ r: 0, g: 255, b: 0 }, { r: 0, g: 200, b: 0 }],
        border: '✓',
        duration: 2000
      },
      error: {
        gradient: [{ r: 255, g: 0, b: 0 }, { r: 200, g: 0, b: 0 }],
        border: '✗',
        duration: 3000
      },
      info: {
        gradient: [{ r: 0, g: 100, b: 255 }, { r: 0, g: 200, b: 255 }],
        border: 'ℹ',
        duration: 2500
      }
    };
    
    const style = styles[type];
    
    // Slide in from right with fade
    this.animator.slideIn(message, {
      gradient: style.gradient,
      duration: 300,
      position: 'top-right'
    });
    
    // Auto-dismiss with fade out
    setTimeout(() => {
      this.animator.fadeOut(message, 300);
    }, style.duration);
  }
}
```

## 🎪 Command Palette Visualization

### Gradient Command Palette
```typescript
class CommandPalette {
  render(commands: Command[], filter: string): string {
    const gradient = createGradient(
      { r: 30, g: 30, b: 50 },
      { r: 50, g: 30, b: 70 },
      this.height
    );
    
    // Fuzzy match highlighting
    const filtered = this.fuzzyFilter(commands, filter);
    
    return filtered.map((cmd, i) => {
      const highlighted = this.highlightMatches(cmd.name, filter);
      const shortcut = this.renderShortcut(cmd.shortcut);
      const selected = i === this.selectedIndex;
      
      if (selected) {
        // Animated selection indicator
        const arrow = this.animateArrow();
        return `${arrow} ${highlighted} ${shortcut}`;
      }
      
      return `  ${highlighted} ${shortcut}`;
    }).join('\n');
  }
}
```

## 🚀 Performance Considerations

### Optimized Rendering Pipeline
```typescript
class RenderOptimizer {
  // Batch ANSI codes to minimize escape sequences
  batchRender(elements: RenderElement[]): string {
    const optimized: string[] = [];
    let currentStyle: Style | null = null;
    
    for (const element of elements) {
      if (!this.styleEquals(currentStyle, element.style)) {
        optimized.push(this.generateStyleCode(element.style));
        currentStyle = element.style;
      }
      optimized.push(element.content);
    }
    
    return optimized.join('') + '\x1b[0m'; // Reset at end
  }
  
  // Cache gradient calculations
  private gradientCache = new Map<string, string[]>();
  
  getCachedGradient(key: string, generator: () => string[]): string[] {
    if (!this.gradientCache.has(key)) {
      this.gradientCache.set(key, generator());
    }
    return this.gradientCache.get(key)!;
  }
}
```

## 🎨 Theme System

### User-Customizable Gradients
```typescript
const THEMES = {
  cyberpunk: {
    primary: [{ r: 255, g: 0, b: 255 }, { r: 0, g: 255, b: 255 }],
    secondary: [{ r: 255, g: 255, b: 0 }, { r: 255, g: 0, b: 0 }],
    background: [{ r: 20, g: 0, b: 40 }, { r: 0, g: 20, b: 40 }]
  },
  
  matrix: {
    primary: [{ r: 0, g: 255, b: 0 }, { r: 0, g: 128, b: 0 }],
    secondary: [{ r: 0, g: 200, b: 0 }, { r: 0, g: 100, b: 0 }],
    background: [{ r: 0, g: 10, b: 0 }, { r: 0, g: 20, b: 0 }]
  },
  
  sunset: {
    primary: [{ r: 255, g: 94, b: 77 }, { r: 255, g: 206, b: 84 }],
    secondary: [{ r: 255, g: 154, b: 0 }, { r: 255, g: 206, b: 84 }],
    background: [{ r: 40, g: 20, b: 20 }, { r: 60, g: 30, b: 20 }]
  }
};
```

## 🏁 Summary

This ANSI visual design system transforms the terminal into a rich, animated, gradient-filled interface that rivals modern GUI applications while maintaining the speed and efficiency of command-line tools. Every element uses smooth gradients, subtle animations, and advanced terminal capabilities to create an unprecedented visual experience in the terminal.

**Key Differentiators:**
- No emoji dependencies - pure ANSI art
- 60 FPS animations for smooth transitions
- True RGB gradients using 256-color mode
- Multi-layer rendering for depth
- Performance-optimized batch rendering
- Customizable theme system
- Non-intrusive notification system
- Beautiful, functional, and fast