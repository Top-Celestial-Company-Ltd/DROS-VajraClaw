$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$projectDir = "$scriptDir\..\vajraclaw_sdk\android_test"
$gradleZip = "E:\gradle-8.7-bin.zip"
$gradleDir = "E:\gradle-8.7"
$gradleBin = "$gradleDir\bin\gradle.bat"

if (-not (Test-Path $gradleBin)) {
    Write-Host "Downloading Gradle 8.7..."
    Invoke-WebRequest -Uri "https://services.gradle.org/distributions/gradle-8.7-bin.zip" -OutFile $gradleZip
    Write-Host "Extracting Gradle..."
    Expand-Archive -Path $gradleZip -DestinationPath "E:\" -Force
}

$env:JAVA_HOME = "E:\Android\Java\latest"
$env:ANDROID_HOME = "E:\Android\Sdk"

Write-Host "Building Android Test APK..."
Set-Location $projectDir
& $gradleBin assembleDebug

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ APK Built Successfully: $projectDir\app\build\outputs\apk\debug\app-debug.apk"
} else {
    Write-Error "❌ Build failed"
}
