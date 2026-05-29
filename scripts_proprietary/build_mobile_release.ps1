$ErrorActionPreference = "Stop"

# Detect script root dynamically to bypass any Chinese character/encoding issues in hardcoded paths
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$baseProjectDir = Resolve-Path "$scriptDir\.."

# Set up local build environment paths
$sdkPath = "E:\Android\Sdk"
$ndkPath = "$sdkPath\ndk\21.4.7075529"
$javaPath = "E:\Android\Java\latest"
$goPath = "E:\dev_tools\go"

Write-Host "Starting VajraClaw Mobile SDK RELEASE compilation pipeline..."
Write-Host "==========================================================="
Write-Host "[SECURITY] Engaging Anti-Reverse Engineering Protections..."
Write-Host "  -> Stripping DWARF debug symbols (-s)"
Write-Host "  -> Stripping DWARF symbol tables (-w)"
Write-Host "  -> Erasing local path identities (-trimpath)"
Write-Host "==========================================================="

# 1. Inject build variables into active process context
$env:ANDROID_HOME = $sdkPath
$env:ANDROID_NDK_HOME = $ndkPath
$env:JAVA_HOME = $javaPath
$env:GOTOOLCHAIN = "go1.25.10"
$userProfile = $env:USERPROFILE
$env:PATH = "$javaPath\bin;$sdkPath\platform-tools;$goPath\bin;$userProfile\go\bin;" + $env:PATH

# 2. Locate gomobile binary
$gomobilePaths = @(
    "$userProfile\go\bin\gomobile.exe",
    "$goPath\bin\gomobile.exe",
    "gomobile"
)

$gomobile = $null
foreach ($path in $gomobilePaths) {
    if (Get-Command $path -ErrorAction SilentlyContinue) {
        $gomobile = $path
        break
    }
}

if (-not $gomobile) {
    Write-Error "Error: gomobile.exe was not found in standard paths!"
}
Write-Host "  -> Using gomobile compiler: $gomobile"

# 3. Ensure gomobile is fully initialized
Write-Host "  -> Verifying compiler initialization..."
& $gomobile init

# 4. Compile Android AAR library package
Write-Host "Compiling Android AAR dynamic package (ARM64 / ARMv7)..."
$targetAar = "$baseProjectDir\vajraclaw_sdk\android\vajraclaw_release.aar"
$mobilePkg = "$baseProjectDir\vajraclaw_sdk\mobile"

# Navigate dynamically
Set-Location $mobilePkg

# Execute the release build with aggressive stripping
& $gomobile bind -target=android -ldflags="-s -w" -trimpath -o $targetAar -v

if (Test-Path $targetAar) {
    Write-Host "✅ Android SDK RELEASE compilation succeeded: $targetAar"
} else {
    Write-Error "❌ Compilation failed to produce RELEASE AAR file!"
}

Write-Host "🎉 VajraClaw Mobile SDK Release Build Pipeline Completed Successfully!"
