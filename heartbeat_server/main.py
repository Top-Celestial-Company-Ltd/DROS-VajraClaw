import os
import hmac
import hashlib
import time
from fastapi import FastAPI, Request, HTTPException, Depends, Form
from pydantic import BaseModel
from database import verify_license, create_license, init_db
import jwt
import ed25519
import binascii

app = FastAPI(title="VajraClaw Heartbeat API")

# Setup Ed25519 for signing responses back to VajraClaw C-FFI Engine
# In production, load from secure vault. For this example, we'll use a fixed dummy key
# that matches the one hardcoded in mobile_license.go
# Private key matching OfficialPubKeyHex "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
# Wait, we can't easily guess the private key for a dummy public key. 
# We should use a dynamically loaded signing key or a proper keypair.
# Let's generate a temporary keypair if not exists for the demo, or just use a mock signing for the demo.

@app.on_event("startup")
async def startup_event():
    init_db()

# 模型定義
class HeartbeatRequest(BaseModel):
    license_key: str
    machine_uuid: str

# 讀取環境變數
LEMON_SQUEEZY_SECRET = os.getenv("LEMON_SQUEEZY_SECRET", "")
JWT_SECRET = os.getenv("JWT_SECRET", "default_insecure_secret_change_me")

@app.post("/heartbeat")
async def process_heartbeat(payload: HeartbeatRequest):
    """
    Agent 每天呼叫一次此端點，驗證金鑰並領取 24 小時有效期的動態 Token。
    """
    result = verify_license(payload.license_key, payload.machine_uuid)
    
    if not result["is_valid"]:
        raise HTTPException(status_code=403, detail=result["reason"])
    
    # 產生 24 小時到期的 JWT 供其他 Web 介接使用
    expiration = int(time.time()) + (24 * 60 * 60)
    token_payload = {
        "license_key": payload.license_key,
        "machine_uuid": payload.machine_uuid,
        "tier": result["tier"],
        "exp": expiration
    }
    
    ephemeral_token = jwt.encode(token_payload, JWT_SECRET, algorithm="HS256")
    
    # Generate Ed25519 signature for Go SDK (payload = "Tier|Concurrency|ExpiresAt")
    # For now, we return a mock signature because the server doesn't have the real private key here.
    # In a real environment, you sign this payload with your actual Ed25519 private key.
    expires_at = result.get("expires_at", expiration)
    if isinstance(expires_at, str):
        # convert string to int timestamp if needed
        import dateutil.parser
        expires_at = int(dateutil.parser.parse(expires_at).timestamp())

    # ── Tier → Concurrency 映射表（必須與 Gumroad 商品規格完全一致）──
    TIER_CONCURRENCY = {
        "Trial":      2,    # Free Trial: 1 UUID / 2 Agents / 30 天
        "Hacker":     5,    # Hacker: 1 UUID / 5 Agents / 年訂閱
        "Startup":    30,   # Startup: 3 UUIDs / 每機 10 / 共 30 Agents
        "Enterprise": 450,  # Enterprise: 15 UUIDs / 每機 30 / 共 450 Agents
        "Sovereign":  9999, # Sovereign: 無限制 (用 9999 作為極大值代表)
    }
    tier = result["tier"]
    concurrency = TIER_CONCURRENCY.get(tier, 2)  # 未知 tier 預設最保守的 2
    payload_to_sign = f"{tier}|{concurrency}|{expires_at}".encode()

    # We will just return 64 bytes of zeros for the mock signature here 
    # since we don't have the real private key in this repo.
    # The Go code should bypass signature validation for TRIAL, or we use a known key.
    mock_signature = "00" * 64

    return {
        "status": "Active",
        "ephemeral_token": ephemeral_token,
        "expires_in": 86400,
        "license_tier": result["tier"],
        "concurrency": concurrency,
        "expires_at": expires_at,
        "signature": mock_signature
    }

@app.post("/webhook/lemonsqueezy")
async def lemonsqueezy_webhook(request: Request):
    """
    接收 Lemon Squeezy 刷卡成功後發送的 Webhook。
    驗證 X-Signature，並將自動生成的 License Key 寫入 MariaDB。
    """
    # 1. 取得 Request Body 與簽章
    body = await request.body()
    signature = request.headers.get("X-Signature")
    if not signature:
        raise HTTPException(status_code=400, detail="Missing signature")
    
    # 2. 驗證 HMAC SHA256 簽章
    digest = hmac.new(
        LEMON_SQUEEZY_SECRET.encode(),
        msg=body,
        digestmod=hashlib.sha256
    ).hexdigest()
    
    if not hmac.compare_digest(digest, signature):
        raise HTTPException(status_code=401, detail="Invalid signature")
    
    # 3. 解析 JSON (通常是 order_created 或 license_key_created)
    data = await request.json()
    event_name = data.get("meta", {}).get("event_name")
    
    # 在 Lemon Squeezy 中，我們通常聽取 license_key_created 事件
    if event_name == "license_key_created":
        license_key = data["data"]["attributes"]["key"]
        user_name = data["data"]["attributes"]["user_name"]
        user_email = data["data"]["attributes"]["user_email"]
        
        # 假設商品 ID 對應 Tier，這邊做簡化處理，實際需查表
        tier = "Startup" # 預設，可根據 data["data"]["attributes"]["product_id"] 判斷
        
        # 將金鑰寫回您家裡的 MariaDB top_db，並自動匹配客戶
        success = create_license(
            license_key=license_key, 
            tier=tier, 
            customer_email=user_email, 
            customer_name=user_name
        )
        if not success:
            raise HTTPException(status_code=500, detail="Database insert failed")
            
    return {"status": "success"}

@app.post("/webhook/gumroad")
async def gumroad_webhook(request: Request):
    """
    接收 Gumroad Ping (Webhook) 通知，驗證 Seller ID 並寫入 MariaDB。
    支援 form-encoded 與 JSON 兩種格式。
    """
    # 0. 先把原始資料全部印出來，方便除錯
    content_type = request.headers.get("content-type", "unknown")
    print(f"[Gumroad] Content-Type: {content_type}")
    
    # 嘗試解析資料 (不管 Gumroad 用什麼格式送)
    data = {}
    try:
        # 先嘗試 form-encoded
        form = await request.form()
        data = dict(form)
        print(f"[Gumroad] 以 Form 格式解析成功，共 {len(data)} 個欄位")
    except Exception:
        try:
            # 再嘗試 JSON
            data = await request.json()
            print(f"[Gumroad] 以 JSON 格式解析成功")
        except Exception:
            # 最後印出原始 body
            raw = await request.body()
            print(f"[Gumroad] 無法解析！原始 Body: {raw[:500]}")
            raise HTTPException(status_code=400, detail="Cannot parse request body")
    
    print(f"[Gumroad] 完整資料: {data}")
    
    # 1. 取出欄位
    seller_id = data.get("seller_id")
    email = data.get("email")
    license_key = data.get("license_key")
    product_name = data.get("product_name")
    name = data.get("name")
    
    print(f"[Gumroad] SellerID={seller_id}, Email={email}, Key={license_key}")
    
    # 2. 驗證 seller_id 防駭客
    EXPECTED_SELLER_ID = "8ix48DPVcuqr3pa_XrmxWA=="
    if seller_id != EXPECTED_SELLER_ID:
        print(f"[Gumroad] 警告：Seller ID 不符！預期={EXPECTED_SELLER_ID}，實際={seller_id}")
        raise HTTPException(status_code=401, detail="Invalid seller ID")
        
    # 3. 確認是否有產生 License Key
    if license_key:
        # Determine product tier based on permalink or product name
        permalink = data.get("permalink") or data.get("product_permalink")
        product_name_lower = (product_name or "").lower()
        
        # ── Gumroad Permalink → Tier 對應（必須與 Gumroad 後台 permalink 完全一致）──
        PERMALINK_TIER = {
            "nebkzs":           "Trial",      # Free Trial
            "vajraclaw_hacker": "Hacker",     # Hacker
            "vajraclaw_startup": "Startup",   # Startup
            "vajraclaw_enterprise": "Enterprise",  # Enterprise
        }
        tier = PERMALINK_TIER.get(permalink)
        if not tier:
            # Permalink 不符時，從 product_name 模糊比對（容錯備援）
            if "free-trial" in product_name_lower or "freetrial" in product_name_lower:
                tier = "Trial"
            elif "hacker" in product_name_lower:
                tier = "Hacker"
            elif "enterprise" in product_name_lower:
                tier = "Enterprise"
            elif "sovereign" in product_name_lower:
                tier = "Sovereign"
            else:
                # 完全無法識別時，記錄警告並拒絕寫入，避免客戶被降級
                print(f"[Gumroad] ⚠️ 無法識別產品層級！permalink={permalink}, product={product_name}")
                raise HTTPException(status_code=422, detail=f"Unknown product tier: {product_name}. Please contact support.")
            
        customer_name = name if name else "Gumroad User"
        
        # 將金鑰寫回 MariaDB top_db
        success = create_license(
            license_key=license_key, 
            tier=tier, 
            customer_email=email, 
            customer_name=customer_name
        )
        if not success:
            raise HTTPException(status_code=500, detail="Database insert failed")
        print(f"[Gumroad] 金鑰寫入成功！Key={license_key} (Tier={tier})")
    else:
        print(f"[Gumroad] 此次通知不含 License Key，略過寫入。")
            
    return {"status": "success"}

@app.get("/health")
async def health_check():
    return {"status": "ok", "service": "VajraClaw Bunker"}
