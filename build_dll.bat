@echo off
echo ====================================================
echo VajraClaw Core Recompiler (C-Shared DLL)
echo ====================================================
echo This script requires GCC (MinGW-w64) to be installed 
echo and added to your PATH.
echo.

go env -w CGO_ENABLED=1
echo Compiling vajra_claw.dll...
go build -buildmode=c-shared -o vajra_claw.dll core\vajra_claw.go

if %ERRORLEVEL% equ 0 (
    echo [SUCCESS] vajra_claw.dll successfully built!
    echo You can now run the Python LangChain integrations.
) else (
    echo [ERROR] Compilation failed. Please ensure GCC is installed.
)
pause
