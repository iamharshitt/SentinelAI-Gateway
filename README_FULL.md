# SentinelAI Gateway

SentinelAI Gateway is a native-messaging bridge that enforces configurable regex-based content policies for browser text inputs. It consists of:

- a Chrome/Edge extension (MV3) in the `extension/` folder
- a Go native agent in `agent/` that loads policies and applies actions (block/warn/redact)
- a small tray launcher at `SentinelAITrayLauncher/` for Windows

This README explains how to build, install (Windows), and create a packaged release.

## Prerequisites

- Go 1.21+
- PowerShell (Windows) for packaging/install scripts
- zip utility (or use PowerShell `Compress-Archive`)

## Build

Build the Go agent and tray launcher:

```powershell
Set-Location -Path agent
go mod tidy
go build -o sentinelai.exe ./...

Set-Location -Path ..\SentinelAITrayLauncher
go build -o SentinelAITrayLauncher.exe
```

The project includes `build_full_mvp.ps1` which automates build and packaging steps for a full MVP release. Use it on Windows:

```powershell
.\build_full_mvp.ps1
```

## Install (Windows native messaging)

The extension communicates with the Go agent using the Chrome/Edge native messaging protocol. The repository contains `install_mvp.ps1` and `com.sentinelai.gateway.json` which register the native messaging host in the Windows registry.

Quick manual steps (PowerShell as current user):

```powershell
# Copy the agent binary to a stable path, e.g. C:\Program Files\SentinelAI\
# Ensure com.sentinelai.gateway.json references the exact binary path.
# Then register the native messaging host (the script does this):
.\install_mvp.ps1
```

The install script writes registry entries under:

- `HKCU:\Software\Google\Chrome\NativeMessagingHosts\com.sentinelai.gateway`
- `HKCU:\Software\Microsoft\Edge\NativeMessagingHosts\com.sentinelai.gateway`

Confirm registration by inspecting those keys with `Get-Item` in PowerShell.

## Running (developer)

Run the agent directly for development, pointing it at a policy file or using environment variable:

```powershell
# From project root
Set-Location -Path agent
./sentinelai.exe -policy ..\sentinel_policies.json
# or
$env:SENTINEL_POLICY = "D:\path\to\sentinel_policies.json"
./sentinelai.exe
```

There is a `test_client.go` at the repo root and `agent/test_client.go` which exercise the length-prefixed stdio protocol for manual testing.

## Packaging & Release Workflow

Suggested steps to produce a release artifact (MVP zip):

1. Build binaries (agent and tray launcher) as above.
2. Prepare extension files: ensure `extension/manifest.json` has the correct version and allowed origins.
3. Create a `dist/` folder and copy:
   - compiled `agent/sentinelai.exe`
   - `SentinelAITrayLauncher.exe`
   - `com.sentinelai.gateway.json`
   - `sentinel_policies.json` (default policies)
   - `extension/` folder contents
4. Zip `dist/` into `SentinelAI-Gateway-MVP.zip`.

PowerShell example:

```powershell
New-Item -ItemType Directory -Path dist -Force
Copy-Item -Path agent\sentinelai.exe -Destination dist\
Copy-Item -Path SentinelAITrayLauncher\SentinelAITrayLauncher.exe -Destination dist\
Copy-Item -Path com.sentinelai.gateway.json -Destination dist\
Copy-Item -Path sentinel_policies.json -Destination dist\
Copy-Item -Path extension\* -Destination dist\extension -Recurse
Compress-Archive -Path dist\* -DestinationPath SentinelAI-Gateway-MVP.zip -Force
```

The repository's `build_full_mvp.ps1` automates these steps; review and adapt it for your release channel.

## Tests

Run unit tests for the Go agent packages:

```powershell
cd agent
go test ./...
```

## Notes & Recommendations

- Consider updating the Go module path in `agent/go.mod` to a VCS-style path (e.g. `github.com/<you>/SentinelAI-Gateway`) before publishing.
- Make binary and policy paths configurable via flags or environment variables (the agent already supports `-policy` and `SENTINEL_POLICY`).
- Before packaging, verify `extension/manifest.json` scopes and permissions; narrow `content_scripts.matches` from "<all_urls>" if possible.
- Add CI (GitHub Actions) to run `go build` and `go test` on PRs.

## Troubleshooting

- If the browser extension cannot connect, confirm the native messaging host registry keys and that the JSON manifest points to the correct binary path.
- Check agent logs for JSON marshal/unmarshal errors or policy load errors.

---

If you'd like, I can add a CI workflow and a packaged GitHub Actions release job next.
