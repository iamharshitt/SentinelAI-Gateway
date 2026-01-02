# SentinelAI Gateway - AI Coding Agent Instructions

## Project Overview
SentinelAI Gateway is a browser security system that prevents sensitive data leakage in LLM prompts. It runs as a native messaging bridge between a Chrome/Edge extension and a Go agent, analyzing text input against security policies before sending to AI models.

## Architecture & Data Flow

### Three-Tier System
1. **Browser Extension** (`extension/`) - Chrome/Edge Manifest V3 extension
   - Injects content script into all web pages
   - Monitors textarea/input fields for user text
   - Communicates via native messaging protocol
   
2. **Native Messaging Agent** (`agent/main.go`) - Go binary (stdio-based)
   - Receives 4-byte length prefix + JSON message from extension
   - Loads security policies from `sentinel_policies.json`
   - Runs analysis via analyzer engine
   - Returns JSON response with same length-prefix protocol
   
3. **Policy & Analysis Engine** (`agent/analyzer/`, `agent/policy/`)
   - Loads policy rules from JSON configuration
   - Applies regex-based pattern matching
   - Three action types: block, warn, redact

### Data Flow Example
```
Browser Input → Content.js (onKeyup) → 
Background.js (native message) → 
Go Agent (stdin/stdout) → 
Analyzer.Analyze(prompt, policy) → 
JSON response → Browser UI feedback
```

## Critical Developer Workflows

### Build & Package (PowerShell)
```powershell
# Full MVP build: compiles Go agent, packages extension, creates zip
.\build_full_mvp.ps1

# Install native messaging registry entries (required for extension to communicate)
.\install_mvp.ps1
```

### Build Process Details
- `build_full_mvp.ps1` compiles `agent/main.go` → `sentinelai.exe`
- Copies extension files, tray launcher, and policies to `/dist`
- Creates `SentinelAI-Gateway-MVP.zip` for distribution
- The compiled binary path in `com.sentinelai.gateway.json` must match output location

### Development Iteration
1. Edit Go code in `agent/`
2. Run `build_full_mvp.ps1` to rebuild agent binary
3. Edit extension JS files directly (no build needed)
4. Reinstall registry entries if native messaging breaks: `.\install_mvp.ps1`

## Key Code Patterns & Integration Points

### Native Messaging Protocol (Critical!)
The Go agent and extension communicate via **length-prefixed JSON over stdio**:
- **Request**: `[4-byte uint32 LE length]` + `[JSON: {"prompt": "text"}]`
- **Response**: Same format with analysis result
- See [agent/main.go](../agent/main.go#L27-L34) for exact binary.LittleEndian usage
- See [extension/background.js](../extension/background.js#L8-L15) for extension side

### Policy Rule Structure
All rules in `sentinel_policies.json` have three matching actions:
```json
{
  "action": "block",     // Immediately reject
  "action": "warn",      // Allow but signal concern
  "action": "redact",    // Replace matched text with redaction_token
}
```
- Default policy action: "allow" (if no rules match)
- See [analyzer/engine.go](../agent/analyzer/engine.go#L10-L40) for action enforcement order

### Registry Integration (Windows Only)
- `install_mvp.ps1` adds native messaging manifest to:
  - `HKCU:\Software\Google\Chrome\NativeMessagingHosts\com.sentinelai.gateway`
  - `HKCU:\Software\Microsoft\Edge\NativeMessagingHosts\com.sentinelai.gateway`
- Path in manifest must match compiled binary location: [com.sentinelai.gateway.json](../com.sentinelai.gateway.json#L4)

### Extension Content Script Timing
- Analyzes text only after **1-second typing pause** to avoid performance issues
- Only monitors `TEXTAREA` and `INPUT` elements
- Applies red border on policy violation for UX feedback
- See [extension/content.js](../extension/content.js#L5-L25)

## Project-Specific Conventions

1. **Policy Path**: Hard-coded in Go agent at startup
   - Current: `D:/SentinelAI-Gateway/sentinel_policies.json`
   - Must be absolute path (forward slashes work on Windows with Go)

2. **Extension ID Registration**
   - `allowed_origins` in manifest must match actual extension ID
   - Currently whitelisted: `chrome-extension://iceenonoogfbpljkjfcpoplcmegphfch/`

3. **Module Name**: Go module is `sentinelai` (no repo prefix)
   - Imports use: `sentinelai/analyzer`, `sentinelai/policy`

4. **Analysis Result Contract**
   - Always return JSON with: `allowed` (bool), `action` (string), optional: `reason`, `severity`, `modified_prompt`
   - `allowed` is final output (true = proceed, false = block UI)
   - See [analyzer/structure.go](../agent/analyzer/structure.go#L3-L10)

## Common Modifications

### Adding a New Policy Rule
1. Add rule to `sentinel_policies.json` with regex patterns
2. Analyzer automatically picks it up (loaded at startup)
3. Test with `build_full_mvp.ps1` and reinstall registry

### Modifying Analysis Logic
- All matching happens in [analyzer/engine.go](../agent/analyzer/engine.go#L10-L45)
- Current: First matching rule with "block" returns immediately; "warn" and "redact" accumulate
- Note: Final `result.Allowed = true` at line 45 is called regardless; review logic before editing

### Adding Extension Permissions
- Update [extension/manifest.json](../extension/manifest.json) permissions and content_script matches
- Manifest V3: no background page, only service_worker

## External Dependencies
- **Go**: 1.21+ (minimal stdlib only, no external packages for core agent)
- **Tray Launcher**: Requires `github.com/getlantern/systray` Go package
  - Install with: `go get github.com/getlantern/systray`
  - Used only for system tray integration in [SentinelAITrayLauncher/tray.go](../SentinelAITrayLauncher/tray.go)
- **Chrome/Edge**: Native Messaging API, Manifest V3 support
- **Node.js**: None (pure JS extension)

## Testing & Debugging

### Go Agent Debugging
- Add `fmt.Println()` logs to see startup messages in console
- Agent crashes if policy file not found (use full absolute path)
- Check registry entries with: `Get-Item -Path "HKCU:\Software\Google\Chrome\NativeMessagingHosts\com.sentinelai.gateway"`

### Extension Debugging
- Open `chrome://extensions`, find SentinelAI, click "Details"
- View background worker logs: "Inspect views" > "service_worker"
- View content script logs: Open DevTools on any webpage (F12)

### Manual Protocol Testing
- Use Node.js or Python to send binary messages if needed
- Ensure exactly 4 bytes (uint32 LE) precedes every JSON message
