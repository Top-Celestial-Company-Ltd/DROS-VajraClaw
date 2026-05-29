$ErrorActionPreference = "Stop"

$sdkPath = "E:\Android\Sdk"
$javaPath = "E:\Android\Java\latest"
$sdkmanager = "$sdkPath\cmdline-tools\latest\bin\sdkmanager.bat"

# Set environmental variables dynamically for the execution scope
$env:JAVA_HOME = $javaPath
$env:PATH = "$javaPath\bin;" + $env:PATH

Write-Host "🚀 Installing NDK 25.2.9519653 for Go Mobile compatibility..."
& $sdkmanager "ndk;25.2.9519653"

Write-Host "✅ NDK 25 installed successfully!"
