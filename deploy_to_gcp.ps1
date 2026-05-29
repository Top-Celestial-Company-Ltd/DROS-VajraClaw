<#
.SYNOPSIS
Deploy VajraClaw Heartbeat Server to GCP Frontline Bunker via Tailscale SSH.

.DESCRIPTION
This script zips the heartbeat_server directory, copies it to the GCP VM over Tailscale,
unzips it, builds the Docker container, and restarts the service.

.PARAMETER GcpUser
The SSH username for the GCP VM (e.g., jimmychen666).
.PARAMETER GcpIp
The Tailscale IP or External IP of the GCP VM.
#>

param(
    [Parameter(Mandatory=$true)]
    [string]$GcpUser,
    
    [Parameter(Mandatory=$true)]
    [string]$GcpIp
)

$ZipPath = "$env:TEMP\heartbeat_server.zip"
$SourceDir = ".\heartbeat_server"
$RemoteDir = "/home/$GcpUser/heartbeat_server"

Write-Host "🚀 [1/4] Zipping files..." -ForegroundColor Cyan
if (Test-Path $ZipPath) { Remove-Item $ZipPath }
Compress-Archive -Path "$SourceDir\*" -DestinationPath $ZipPath -Force

Write-Host "🚀 [2/4] Uploading to GCP Bunker ($GcpIp)..." -ForegroundColor Cyan
scp $ZipPath "${GcpUser}@${GcpIp}:/tmp/heartbeat_server.zip"

Write-Host "🚀 [3/4] Unzipping and Building Docker on GCP..." -ForegroundColor Cyan
$RemoteCommands = @"
    mkdir -p $RemoteDir
    unzip -o /tmp/heartbeat_server.zip -d $RemoteDir
    cd $RemoteDir
    
    echo 'Building Docker image...'
    sudo docker build -t vajraclaw-heartbeat .
    
    echo 'Stopping old container (if exists)...'
    sudo docker stop vajraclaw-bunker || true
    sudo docker rm vajraclaw-bunker || true
    
    echo 'Starting new container...'
    sudo docker run -d --name vajraclaw-bunker --restart unless-stopped -p 80:8000 --env-file .env vajraclaw-heartbeat
    
    echo 'Cleaning up...'
    rm /tmp/heartbeat_server.zip
    
    echo 'Deployment Successful!'
"@

ssh "${GcpUser}@${GcpIp}" $RemoteCommands

Write-Host "✅ Deployment Completed!" -ForegroundColor Green
