$ErrorActionPreference = "Stop"

$sdkPath = "E:\Android\Sdk"
$javaPath = "E:\Android\Java\latest"
$sdkmanager = "$sdkPath\cmdline-tools\latest\bin\sdkmanager.bat"

# Set environmental variables dynamically for the execution scope
$env:JAVA_HOME = $javaPath
$env:PATH = "$javaPath\bin;" + $env:PATH

Write-Host "🚀 Installing NDK 21.4.7075529 (LTS r21e) for absolute Go Mobile API 16 compatibility..."
& $sdkmanager "ndk;21.4.7075529"

Write-Host "✅ NDK 21 installed successfully!"
