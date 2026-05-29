$ErrorActionPreference = "Stop"

$sdkPath = "E:\Android\Sdk"
$cmdlineToolsZip = "$sdkPath\cmdline-tools.zip"
$cmdlineToolsDir = "$sdkPath\cmdline-tools"
$latestDir = "$cmdlineToolsDir\latest"

Write-Host "🚀 開始 VajraClaw Android SDK 全自動化建置程序..."

# 1. Create directories
Write-Host "[1/5] 建立 E:\Android\Sdk 目錄結構..."
if (-not (Test-Path $sdkPath)) { 
    New-Item -ItemType Directory -Path $sdkPath | Out-Null 
}
if (-not (Test-Path $cmdlineToolsDir)) { 
    New-Item -ItemType Directory -Path $cmdlineToolsDir | Out-Null 
}

# 2. Download Command Line Tools
Write-Host "[2/5] 從 Google 伺服器下載 Android Command Line Tools..."
if (-not (Test-Path $cmdlineToolsZip)) {
    Write-Host "  -> 正在使用 curl.exe 進行穩定下載 (含自動重試)..."
    curl.exe --retry 5 --retry-delay 3 -L -o $cmdlineToolsZip "https://dl.google.com/android/repository/commandlinetools-win-14742923_latest.zip"
} else {
    Write-Host "  -> 發現已下載的 Zip，跳過下載。"
}

# 3. Extract and organize
Write-Host "[3/5] 解壓縮並配置目錄 (cmdline-tools/latest)..."
if (-not (Test-Path $latestDir)) {
    Expand-Archive -Path $cmdlineToolsZip -DestinationPath $cmdlineToolsDir -Force
    Rename-Item -Path "$cmdlineToolsDir\cmdline-tools" -NewName "latest"
} else {
    Write-Host "  -> 目錄已配置，跳過解壓縮。"
}

# 4. Accept Licenses and install NDK, build tools, etc.
Write-Host "[4/5] 接受授權條款並下載 NDK、Platform-Tools (這可能需要幾分鐘)..."
$sdkmanager = "$latestDir\bin\sdkmanager.bat"

# Pipe multiple 'y's into sdkmanager to accept all licenses
$yes = "y`n" * 30
Write-Host "  -> 正在寫入同意條款..."
$yes | & $sdkmanager --licenses | Out-Null

Write-Host "  -> 正在下載套件：platform-tools, platforms;android-34, build-tools;34.0.0, ndk;26.1.10909125"
& $sdkmanager "platform-tools" "platforms;android-34" "build-tools;34.0.0" "ndk;26.1.10909125"

# 5. Set Environment Variables for User
Write-Host "[5/5] 設定 Windows 環境變數..."
[Environment]::SetEnvironmentVariable("ANDROID_HOME", $sdkPath, "User")
Write-Host "  -> ANDROID_HOME = $sdkPath"

if (Test-Path "$sdkPath\ndk") {
    $ndkPath = Get-ChildItem "$sdkPath\ndk" | Sort-Object Name -Descending | Select-Object -First 1 | Select-Object -ExpandProperty FullName
    [Environment]::SetEnvironmentVariable("ANDROID_NDK_HOME", $ndkPath, "User")
    Write-Host "  -> ANDROID_NDK_HOME = $ndkPath"
} else {
    Write-Host "⚠️ 警告：未找到 NDK 目錄，可能是下載失敗。"
}

# Update PATH
$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
$newPaths = $sdkPath + '\platform-tools;' + $latestDir + '\bin'
if ($userPath -notmatch 'Android\\Sdk\\platform-tools') {
    $updatedPath = $userPath + ';' + $newPaths
    [Environment]::SetEnvironmentVariable('Path', $updatedPath, 'User')
    Write-Host "  -> PATH 已更新，加入 platform-tools 與 cmdline-tools。"
} else {
    Write-Host "  -> PATH 已經包含 Android 路徑，無需修改。"
}

Write-Host "✅ VajraClaw Android 開發環境建置完成！請在下一次開啟 PowerShell 時生效。"
