$ErrorActionPreference = "Stop"

$sdkPath = "E:\Android\Sdk"
$javaPath = "E:\Android\Java"
$javaZip = "E:\Android\openjdk.zip"
$javaLatest = "$javaPath\latest"
$cmdlineToolsZip = "$sdkPath\cmdline-tools.zip"
$cmdlineToolsDir = "$sdkPath\cmdline-tools"
$latestDir = "$cmdlineToolsDir\latest"

Write-Host "🚀 Starting VajraClaw Android & Java Environment Auto-Installer..."

# 0. Create base Android directory
if (-not (Test-Path "E:\Android")) {
    New-Item -ItemType Directory -Path "E:\Android" | Out-Null
}

# 1. Download & Setup OpenJDK if not present
if (-not (Test-Path $javaLatest)) {
    Write-Host "[1/6] Downloading Portable OpenJDK 17..."
    if (-not (Test-Path $javaPath)) {
        New-Item -ItemType Directory -Path $javaPath | Out-Null
    }
    curl.exe --retry 5 --retry-delay 3 -L -o $javaZip "https://github.com/adoptium/temurin17-binaries/releases/download/jdk-17.0.11%2B9/OpenJDK17U-jdk_x64_windows_hotspot_17.0.11_9.zip"
    
    Write-Host "  -> Extracting OpenJDK..."
    Expand-Archive -Path $javaZip -DestinationPath $javaPath -Force
    $extractedJdk = Get-ChildItem $javaPath | Where-Object { $_.Name -like "jdk-*" } | Select-Object -First 1
    Rename-Item -Path $extractedJdk.FullName -NewName "latest"
    Remove-Item -Path $javaZip -Force
} else {
    Write-Host "[1/6] OpenJDK 17 already configured. Skipping download."
}

# 2. Register JAVA_HOME and update PATH for current session and user env
$env:JAVA_HOME = $javaLatest
$env:PATH = "$javaLatest\bin;" + $env:PATH
[Environment]::SetEnvironmentVariable("JAVA_HOME", $javaLatest, "User")
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notmatch "Android\\Java\\latest") {
    $updatedPath = "$javaLatest\bin;$userPath"
    [Environment]::SetEnvironmentVariable("Path", $updatedPath, "User")
}
Write-Host "  -> JAVA_HOME registered successfully: $javaLatest"

# 3. Create SDK directories
Write-Host "[3/6] Configuring Android Sdk directory structure..."
if (-not (Test-Path $sdkPath)) { 
    New-Item -ItemType Directory -Path $sdkPath | Out-Null 
}
if (-not (Test-Path $cmdlineToolsDir)) { 
    New-Item -ItemType Directory -Path $cmdlineToolsDir | Out-Null 
}

# 4. Download Command Line Tools
if (-not (Test-Path $cmdlineToolsZip)) {
    Write-Host "[4/6] Downloading Android Command Line Tools..."
    curl.exe --retry 5 --retry-delay 3 -L -o $cmdlineToolsZip "https://dl.google.com/android/repository/commandlinetools-win-14742923_latest.zip"
} else {
    Write-Host "[4/6] Existing Zip found. Skipping download."
}

# 5. Extract cmdline-tools
if (-not (Test-Path $latestDir)) {
    Write-Host "[5/6] Extracting cmdline-tools..."
    Expand-Archive -Path $cmdlineToolsZip -DestinationPath $cmdlineToolsDir -Force
    Rename-Item -Path "$cmdlineToolsDir\cmdline-tools" -NewName "latest"
} else {
    Write-Host "[5/6] Directory already configured. Skipping extraction."
}

# 6. Accept Licenses and install NDK, build tools, etc.
Write-Host "[6/6] Accepting Android SDK licenses and downloading NDK & build tools..."
$sdkmanager = "$latestDir\bin\sdkmanager.bat"
$yes = "y`n" * 30
$yes | & $sdkmanager --licenses | Out-Null

Write-Host "  -> Installing platform-tools, platforms;android-34, build-tools;34.0.0, ndk;26.1.10909125..."
& $sdkmanager "platform-tools" "platforms;android-34" "build-tools;34.0.0" "ndk;26.1.10909125"

# Set ANDROID environment variables
[Environment]::SetEnvironmentVariable("ANDROID_HOME", $sdkPath, "User")
if (Test-Path "$sdkPath\ndk") {
    $ndkPath = Get-ChildItem "$sdkPath\ndk" | Sort-Object Name -Descending | Select-Object -First 1 | Select-Object -ExpandProperty FullName
    [Environment]::SetEnvironmentVariable("ANDROID_NDK_HOME", $ndkPath, "User")
}

$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
$newPaths = "$sdkPath\platform-tools;$latestDir\bin"
if ($userPath -notmatch "Android\\Sdk\\platform-tools") {
    $updatedPath = "$userPath;$newPaths"
    [Environment]::SetEnvironmentVariable("Path", $updatedPath, "User")
}

Write-Host "✅ VajraClaw SDK & Java Environment Setup Completed Successfully!"
