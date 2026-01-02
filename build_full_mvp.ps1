# Build Go agent
Write-Host "Building SentinelAI agent..."
cd "D:\SentinelAI-Gateway\agent"
go build -o sentinelai.exe main.go

# Create dist folder
$distPath = "D:\SentinelAI-Gateway\dist"
if (!(Test-Path $distPath)) { New-Item -ItemType Directory -Path $distPath }

# Copy files to dist
Copy-Item -Path "D:\SentinelAI-Gateway\agent\sentinelai.exe" -Destination $distPath -Force
Copy-Item -Path "D:\SentinelAI-Gateway\SentinelAITrayLauncher\*" -Destination $distPath\SentinelAITrayLauncher -Recurse -Force
Copy-Item -Path "D:\SentinelAI-Gateway\extension\*" -Destination $distPath\extension -Recurse -Force
Copy-Item -Path "D:\SentinelAI-Gateway\sentinel_policies.json" -Destination $distPath -Force
Copy-Item -Path "D:\SentinelAI-Gateway\com.sentinelai.gateway.json" -Destination $distPath -Force

# Zip the MVP
$zipPath = "D:\SentinelAI-Gateway\SentinelAI-Gateway-MVP.zip"
if (Test-Path $zipPath) { Remove-Item $zipPath }
Compress-Archive -Path $distPath\* -DestinationPath $zipPath

Write-Host "MVP build complete: $zipPath"
