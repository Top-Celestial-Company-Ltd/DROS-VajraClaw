import os
import pymysql
import logging
from dotenv import load_dotenv

# 載入環境變數
load_dotenv()

DB_HOST = os.getenv("DB_HOST", "127.0.0.1")
DB_PORT = int(os.getenv("DB_PORT", 3308))
DB_USER = os.getenv("DB_USER", "topuser")
DB_PASS = os.getenv("DB_PASS", "")
DB_NAME = os.getenv("DB_NAME", "top_db")

logger = logging.getLogger("uvicorn.error")

def init_db():
    connection = get_db_connection()
    try:
        with connection.cursor() as cursor:
            cursor.execute("""
                CREATE TABLE IF NOT EXISTS Vajraclaw試用表T (
                    MachineID VARCHAR(255) PRIMARY KEY,
                    TrialStartDate DATETIME NOT NULL
                )
            """)
        connection.commit()
    except Exception as e:
        logger.error(f"資料庫初始化失敗: {e}")
    finally:
        connection.close()

def get_db_connection():
    """建立並回傳 MariaDB 連線。每次請求都建立新連線以確保無狀態。"""
    try:
        connection = pymysql.connect(
            host=DB_HOST,
            port=DB_PORT,
            user=DB_USER,
            password=DB_PASS,
            database=DB_NAME,
            cursorclass=pymysql.cursors.DictCursor
        )
        return connection
    except Exception as e:
        logger.error(f"資料庫連線失敗: {e}")
        raise e

def verify_license(license_key: str, machine_uuid: str) -> dict:
    """
    核對授權碼是否有效，並檢查是否綁定了正確的機器特徵碼。
    回傳字典包含: is_valid, status, expire_date 等資訊。
    """
    connection = get_db_connection()
    try:
        with connection.cursor() as cursor:
            if license_key == "TRIAL" or license_key == "":
                # 處理自動 30 天試用期
                cursor.execute("SELECT TrialStartDate FROM Vajraclaw試用表T WHERE MachineID = %s", (machine_uuid,))
                result = cursor.fetchone()
                
                if not result:
                    # 首次使用，自動註冊 30 天試用
                    cursor.execute("INSERT INTO Vajraclaw試用表T (MachineID, TrialStartDate) VALUES (%s, NOW())", (machine_uuid,))
                    connection.commit()
                    cursor.execute("SELECT TrialStartDate FROM Vajraclaw試用表T WHERE MachineID = %s", (machine_uuid,))
                    result = cursor.fetchone()
                
                start_date = result['TrialStartDate']
                # 取得 UNIX Timestamp 以計算 30 天後
                import datetime
                # 若 30 天已過，拒絕授權
                expires_date = start_date + datetime.timedelta(days=30)
                if datetime.datetime.now() > expires_date:
                    return {"is_valid": False, "reason": "Trial expired. Please purchase an Enterprise License."}
                
                return {
                    "is_valid": True,
                    "tier": "Trial",
                    "expires_at": int(expires_date.timestamp())
                }

            # 使用我們定義的 Vajraclaw授權表T 結構
            sql = "SELECT * FROM Vajraclaw授權表T WHERE 授權碼 = %s"
            cursor.execute(sql, (license_key,))
            result = cursor.fetchone()

            if not result:
                return {"is_valid": False, "reason": "License not found"}
            
            status = result.get('啟用狀態')
            if status != 'Active':
                return {"is_valid": False, "reason": f"License is {status}"}
            
            # 檢查機器特徵碼綁定
            bound_uuid = result.get('機器特徵碼')
            if bound_uuid and bound_uuid != machine_uuid:
                # 已經綁定其他機器
                return {"is_valid": False, "reason": "License bound to another machine"}
            elif not bound_uuid:
                # 首次使用，自動綁定
                bind_sql = "UPDATE Vajraclaw授權表T SET 機器特徵碼 = %s WHERE LicenseID = %s"
                cursor.execute(bind_sql, (machine_uuid, result['LicenseID']))
                connection.commit()
            
            # 這裡可以加入檢查 "到期日" 的邏輯
            # ...
            
            return {
                "is_valid": True, 
                "tier": result.get('方案層級'),
                "expires_at": str(result.get('到期日')) if result.get('到期日') else None
            }
    finally:
        connection.close()

def create_license(license_key: str, tier: str, customer_email: str, customer_name: str, expires_at: str = None) -> bool:
    """
    由 Webhook 呼叫，在資料庫中新建一筆授權碼。
    自動比對 '客戶資料表T'，若存在則綁定舊客戶，若不存在則新增客戶。
    """
    connection = get_db_connection()
    try:
        with connection.cursor() as cursor:
            # 1. 尋找既有客戶
            find_customer_sql = "SELECT CustomerID FROM 客戶資料表T WHERE 客戶Email = %s"
            cursor.execute(find_customer_sql, (customer_email,))
            customer_record = cursor.fetchone()
            
            customer_id = None
            if customer_record:
                # 老客戶
                customer_id = customer_record['CustomerID']
                logger.info(f"找到既有客戶: {customer_email} (ID: {customer_id})")
            else:
                # 新客戶，自動建檔
                insert_customer_sql = "INSERT INTO 客戶資料表T (客戶姓名, 客戶Email) VALUES (%s, %s)"
                cursor.execute(insert_customer_sql, (customer_name, customer_email))
                customer_id = cursor.lastrowid
                logger.info(f"建立新客戶: {customer_name} ({customer_email}), 分配 ID: {customer_id}")
                
            # 2. 寫入新的授權碼
            insert_license_sql = """
                INSERT INTO Vajraclaw授權表T 
                (CustomerID, 授權碼, 方案層級, 啟用狀態, 購買日期, 到期日) 
                VALUES (%s, %s, %s, 'Active', NOW(), %s)
            """
            cursor.execute(insert_license_sql, (customer_id, license_key, tier, expires_at))
            connection.commit()
            return True
    except Exception as e:
        logger.error(f"建立授權碼失敗: {e}")
        return False
    finally:
        connection.close()
