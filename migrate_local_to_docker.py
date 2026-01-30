import psycopg2
from psycopg2.extras import execute_values

# Configuration
SOURCE_DB = {
    "host": "localhost",
    "port": "5432",
    "database": "academ_aide",
    "user": "postgres",
    "password": "password" # Trying 'password' first as it is common default, or will try sathvik555 if failed
}

DEST_DB = {
    "host": "localhost",
    "port": "5435",
    "database": "academ_aide",
    "user": "postgres",
    "password": "postgres"  # We know this is correct for Docker
}

def migrate_data():
    try:
        # 1. Connect to Source (Local)
        print("Connecting to Local DB (5432)...")
        try:
            src_conn = psycopg2.connect(**SOURCE_DB)
        except psycopg2.OperationalError:
            print("Auth failed with 'password', trying 'sathvik555'...")
            SOURCE_DB["password"] = "sathvik555"
            src_conn = psycopg2.connect(**SOURCE_DB)

        src_cur = src_conn.cursor()
        
        # 2. Fetch Data
        print("Fetching chunks from Local DB...")
        src_cur.execute("SELECT count(*) FROM COURSE_MATERIAL_CHUNK")
        count = src_cur.fetchone()[0]
        print(f"Found {count} chunks in Local DB.")
        
        if count == 0:
            print("Nothing to migrate.")
            return

        src_cur.execute("SELECT course_id, unit_no, content_text, embedding, source_file, chunk_index FROM COURSE_MATERIAL_CHUNK")
        rows = src_cur.fetchall()
        
        # 3. Connect to Destination (Docker)
        print("Connecting to Docker DB (5435)...")
        dest_conn = psycopg2.connect(**DEST_DB)
        dest_cur = dest_conn.cursor()
        
        # 4. Insert Data
        print("Inserting data into Docker DB...")
        insert_query = """
            INSERT INTO COURSE_MATERIAL_CHUNK (course_id, unit_no, content_text, embedding, source_file, chunk_index)
            VALUES %s
        """
        
        execute_values(dest_cur, insert_query, rows)
        
        dest_conn.commit()
        print(f"Successfully migrated {len(rows)} chunks!")
        
        src_cur.close()
        src_conn.close()
        dest_cur.close()
        dest_conn.close()

    except Exception as e:
        print(f"Migration Failed: {e}")

if __name__ == "__main__":
    migrate_data()
