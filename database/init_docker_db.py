import psycopg2
from psycopg2.extensions import ISOLATION_LEVEL_AUTOCOMMIT
import os

# Config from .env logic
DB_HOST = "localhost"
DB_PORT = "5435"
DB_USER = "postgres"
DB_PASS = "postgres" # Default for docker usually, checking .env... user's .env said password=postgres
TARGET_DB = "academ_aide"

def read_sql_file(filepath):
    with open(filepath, 'r') as f:
        return f.read()

def init_db():
    # 1. Connect to default 'postgres' db to create new db
    print(f"Connecting to postgres on port {DB_PORT}...")
    try:
        conn = psycopg2.connect(
            host=DB_HOST,
            port=DB_PORT,
            user=DB_USER,
            password=DB_PASS,
            database="postgres"
        )
        conn.set_isolation_level(ISOLATION_LEVEL_AUTOCOMMIT)
        cur = conn.cursor()
        
        # Check if DB exists
        cur.execute(f"SELECT 1 FROM pg_catalog.pg_database WHERE datname = '{TARGET_DB}'")
        exists = cur.fetchone()
        
        if not exists:
            print(f"Creating database '{TARGET_DB}'...")
            cur.execute(f"CREATE DATABASE {TARGET_DB}")
        else:
            print(f"Database '{TARGET_DB}' already exists.")
        
        cur.close()
        conn.close()
    except Exception as e:
        print(f"Error creating DB: {e}")
        return

    # 2. Connect to new DB and run schemas
    print(f"Connecting to {TARGET_DB}...")
    try:
        conn = psycopg2.connect(
            host=DB_HOST,
            port=DB_PORT,
            user=DB_USER,
            password=DB_PASS,
            database=TARGET_DB
        )
        cur = conn.cursor()
        
        # List of files to run in order
        files = [
            "database/schema.sql",
            "database/rag_setup.sql",
            "database/01_add_wallet_auth.sql",
            "database/insert_real_data.sql"
        ]
        
        for fpath in files:
            if os.path.exists(fpath):
                print(f"Running {fpath}...")
                sql = read_sql_file(fpath)
                cur.execute(sql)
            else:
                print(f"Warning: File {fpath} not found.")
        
        conn.commit()
        print("✅ Database initialization complete!")
        cur.close()
        conn.close()
        
    except Exception as e:
        print(f"Error running SQL scripts: {e}")

if __name__ == "__main__":
    init_db()
