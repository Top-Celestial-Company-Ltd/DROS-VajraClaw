import sys
from database import get_db_connection

def test_connection():
    try:
        conn = get_db_connection()
        print("Success! Connected to MariaDB (top_db)")
        
        # 測試查詢資料表清單，確保權限正常
        with conn.cursor() as cursor:
            cursor.execute("SHOW TABLES")
            tables = cursor.fetchall()
            print("\nTables found (first 5):")
            count = 0
            for row in tables:
                table_name = list(row.values())[0]
                print(f"  - {table_name}")
                count += 1
                if count >= 5:
                    break
            
        conn.close()
        sys.exit(0)
    except Exception as e:
        print(f"Failed to connect: {e}")
        sys.exit(1)

if __name__ == "__main__":
    test_connection()
