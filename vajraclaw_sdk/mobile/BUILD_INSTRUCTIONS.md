# VajraClaw Mobile SDK - Android AAR 編譯指南

這份文件記錄了如何將 Go 核心安全引擎 (`mobile_*.go`) 編譯為 Android 專用的 `.aar` 函式庫。

## ⚠️ 環境迷思澄清：不需要安裝 Windows GCC！

在 Windows 上編譯 Android AAR 時，**不需要**在系統中安裝 GCC 或 MinGW。
`gomobile bind` 底層會自動呼叫 **Android NDK (Native Development Kit)** 中內建的 Clang 交叉編譯器。只要確保電腦有安裝 Android NDK 即可（通常位於 `Android\Sdk\ndk`）。

---

## 🛠️ 必備環境與編譯步驟

### 1. 確認環境變數
確保您的系統或終端機中，設定了正確的 Android 變數，讓 `gomobile` 能夠找到 NDK：
```powershell
$env:ANDROID_HOME="您的 Android SDK 路徑，例如 E:\Android\Sdk"
$env:ANDROID_NDK_HOME="您的 Android NDK 路徑，例如 E:\Android\Sdk\ndk\26.1.10909125"
```

### 2. 安裝 Gomobile 工具鏈
確保您已經安裝了 Go (建議版本 >= 1.22)，並執行以下指令安裝 `gomobile`：
```bash
go install golang.org/x/mobile/cmd/gomobile@latest
go install golang.org/x/mobile/cmd/gobind@latest
gomobile init
```

### 3. 執行編譯 (Gomobile Bind)
進入目前的目錄 (`vajraclaw_sdk\mobile`)，然後執行以下指令將 Go 程式碼封裝為 `.aar`，並直接輸出到 Android 專案的 `libs` 資料夾：

```bash
gomobile bind -target=android -androidapi 24 -o ..\..\..\VajraClaw-Mobile\VajraClawDemoApp\app\libs\vajraclaw_release.aar .
```
*(註：`-androidapi 24` 是為了對齊 Android 專案的 `minSdk 24`。若未指定，部分新版 NDK 可能會因為不再支援預設的 API 16 而報錯)*

---

## 🔧 常見問題排除 (Troubleshooting)

### Q: 執行 `gomobile bind` 時，出現 `download go1.25 for windows/amd64: toolchain not available` 或 `go mod tidy failed`
**原因**：最新版的 `golang.org/x/mobile` 依賴尚未發布的 Go 1.25 toolchain，導致在 `go mod tidy` 階段無法下載而中斷。
**解決方案**：
1. **方案 A (推薦)**：將本機的 Go 版本升級到最新穩定版。
2. **方案 B**：手動降版 `x/mobile` 套件。編輯 `go.mod` 將 Go 版本改為您的本機版本 (例如 `go 1.22.0`)，然後執行：
   ```bash
   go get golang.org/x/mobile@v0.0.0-20231127183840-76ac6878050a
   ```
   即可避免觸發新版 toolchain 的下載機制。
