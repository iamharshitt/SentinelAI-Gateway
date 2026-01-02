# Path to Native Messaging manifest
$manifestPath = "D:\SentinelAI-Gateway\com.sentinelai.gateway.json"

# Registry keys for Chrome and Edge
$chromeKey = "HKCU:\Software\Google\Chrome\NativeMessagingHosts\com.sentinelai.gateway"
$edgeKey = "HKCU:\Software\Microsoft\Edge\NativeMessagingHosts\com.sentinelai.gateway"

# Create registry entries
New-Item -Path $chromeKey -Force | Out-Null
Set-ItemProperty -Path $chromeKey -Name "(Default)" -Value $manifestPath

New-Item -Path $edgeKey -Force | Out-Null
Set-ItemProperty -Path $edgeKey -Name "(Default)" -Value $manifestPath

Write-Host "Native Messaging manifest registered for Chrome and Edge."
