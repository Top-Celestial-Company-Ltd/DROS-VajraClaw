# VajraClaw Heartbeat Server — RUNBOOK

<!-- dros_component: dros-heartbeat-server -->
<!-- dros_status: Active -->
<!-- dros_description: 授權驗證心跳伺服器的維護手冊，含部署、重啟、Debug 完整 SOP -->

> **最後更新**：2026-07-29 by Antigravity Agent  
> **伺服器**：Google Cloud VM `dros-vajra-claw`  
> **維護帳號**：`jimmychen666`

---

## 一、服務概覽

| 項目 | 值 |
|:---|:---|
| 服務名稱 | `dros-heartbeat.service` |
| 監聽埠 | `0.0.0.0:8000` |
| 工作目錄 | `/home/jimmychen666/heartbeat_server/` |
| 執行身份 | `root` |
| Python 版本 | Python 3.10.12 (`/usr/bin/python3`) |
| Log 位置 | `/home/jimmychen666/heartbeat_server/uvicorn.log` |
| 資料庫 | MariaDB `top_db`（連線參數見 `.env`）|

---

## 二、日常維護指令

### 查看服務狀態
```bash
sudo systemctl status dros-heartbeat
```

### 啟動 / 停止 / 重啟
```bash
sudo systemctl start dros-heartbeat
sudo systemctl stop dros-heartbeat
sudo systemctl restart dros-heartbeat
```

### 查看即時 Log
```bash
tail -f /home/jimmychen666/heartbeat_server/uvicorn.log
# 或
sudo journalctl -u dros-heartbeat -f
```

### 健康檢查
```bash
curl -s http://localhost:8000/health
# 預期回傳：{"status":"ok","service":"VajraClaw Bunker"}
```

---

## 三、部署更新（從 Windows 本地更新程式碼）

> **注意**：VM 的 NAS 掛載路徑（`/opt/openclaw-dros/...`）**不一定可用**。
> 建議直接在 VM 上用 Python heredoc 打 patch，或透過以下方式同步。

### 方法 A：SCP 直接上傳（推薦）
```bash
# 在 Windows PowerShell 執行
scp "E:\vscode\AI知識庫\DROS-VajraClaw\heartbeat_server\main.py" jimmychen666@<VM_IP>:/home/jimmychen666/heartbeat_server/main.py
scp "E:\vscode\AI知識庫\DROS-VajraClaw\heartbeat_server\database.py" jimmychen666@<VM_IP>:/home/jimmychen666/heartbeat_server/database.py

# 上傳後重啟
ssh jimmychen666@<VM_IP> "sudo systemctl restart dros-heartbeat"
```

### 方法 B：直接在 VM 上編輯
```bash
sudo nano /home/jimmychen666/heartbeat_server/main.py
sudo systemctl restart dros-heartbeat
```

---

## 四、Python 套件重裝（環境損毀時）

> ⚠️ 歷史教訓：原本的 Python 3.12（`/usr/local/bin/python3.12`）與 uvicorn 執行檔於 2026-07 間被刪除，導致服務無法重啟。現已改用系統 Python 3.10。

```bash
sudo apt install -y python3-pip
sudo pip3 install uvicorn fastapi pyjwt pymysql python-dotenv python-dateutil
sudo systemctl restart dros-heartbeat
```

---

## 五、API 端點說明

| 端點 | 方法 | 說明 |
|:---|:---|:---|
| `/health` | GET | 服務健康檢查 |
| `/heartbeat` | POST | Agent 每日心跳驗證，回傳 JWT + concurrency 限制 |
| `/webhook/gumroad` | POST | Gumroad 購買成功 Webhook，自動建立授權碼 |
| `/webhook/lemonsqueezy` | POST | LemonSqueezy 購買成功 Webhook |

### `/heartbeat` 回傳範例
```json
{
  "status": "Active",
  "ephemeral_token": "eyJ...",
  "expires_in": 86400,
  "concurrency": 30,
  "tier": "Startup"
}
```

---

## 六、Tier ↔ Concurrency 對應表（程式碼黃金標準）

| Tier | Machine UUIDs | Concurrent Agents |
|:---|:---:|:---:|
| Trial | 1 | **2** |
| Hacker | 1 | **5** |
| Startup | 3 | **30**（每機 10）|
| Enterprise | 15 | **450**（每機 30）|
| Sovereign | 無限制 | **9999** |

> 此表對應 `main.py` 的 `TIER_CONCURRENCY` dict，**若商品規格變更必須同步修改**。

---

## 七、Gumroad Permalink → Tier 對應表

| Gumroad Permalink | Tier |
|:---|:---|
| `nebkzs` | Trial |
| `vajraclaw_hacker` | Hacker |
| `vajraclaw_startup` | Startup |
| `vajraclaw_enterprise` | Enterprise |

> 若 Gumroad 後台 permalink 有變更，必須同步更新 `main.py` 的 `PERMALINK_TIER` dict。

---

## 八、已知問題與歷史紀錄

| 日期 | 問題 | 解法 |
|:---|:---|:---|
| 2026-07-29 | 原 Python 3.12 執行檔消失，服務重啟失敗 | 重裝 pip + 套件，改用系統 Python 3.10，建立 systemd service |
| 2026-07-29 | `concurrency` 邏輯錯誤（非 Trial 一律給 100）| 修正為 TIER_CONCURRENCY mapping |
| 2026-07-29 | Gumroad/LemonSqueezy webhook 寫死 `tier="Startup"` | 加入 permalink 與 product_name 識別邏輯 |
| 2026-07-29 | Trial 到期訊息誤導客戶「請購買 Enterprise」| 修正為正確的定價頁連結 |

---

*VajraClaw Heartbeat Server RUNBOOK — DROS Agent 治理協議* 🛡️
