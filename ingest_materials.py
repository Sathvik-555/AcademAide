import os
import psycopg2
from langchain_community.document_loaders import PyPDFLoader
from langchain_experimental.text_splitter import SemanticChunker
from langchain_community.embeddings import OllamaEmbeddings
import time

# --- Configuration ---
# Folder structure should be: ./materials/<CourseID>/<UnitNumber>/file.pdf
# Example: ./materials/CS101/1/intro.pdf
BASE_DIR = "./materials"

# DB Connection - UPDATE THESE VALUES
DB_HOST = "localhost"
DB_PORT = "5435"
DB_NAME = "academ_aide" # Or your actual DB name
DB_USER = "postgres"
DB_PASS = "postgres" # Update this!

# Model (Must match what you use in Go)
# Using nomic-embed-text as per your previous setup, or llama3.2 depending on your preference
MODEL_NAME = "nomic-embed-text" 

def get_db_connection():
    conn = psycopg2.connect(
        host=DB_HOST,
        port=DB_PORT,
        database=DB_NAME,
        user=DB_USER,
        password=DB_PASS
    )
    return conn

def setup_database():
    """Ensures the table exists (using float8[] instead of pgvector) and creates a similarity function."""
    conn = get_db_connection()
    cursor = conn.cursor()
    try:
        # 1. Create Cosine Similarity Function (for native array support)
        # This allows us to do vector search WITHOUT installing the pgvector extension
        cursor.execute("""
            CREATE OR REPLACE FUNCTION cosine_similarity(a float8[], b float8[])
            RETURNS float8 AS $$
            DECLARE
                dot_product float8 := 0;
                norm_a float8 := 0;
                norm_b float8 := 0;
                i int;
            BEGIN
                -- Assumes arrays are same length
                FOR i IN 1 .. array_upper(a, 1) LOOP
                    dot_product := dot_product + (a[i] * b[i]);
                    norm_a := norm_a + (a[i] * a[i]);
                    norm_b := norm_b + (b[i] * b[i]);
                END LOOP;
                
                IF norm_a = 0 OR norm_b = 0 THEN
                    RETURN 0;
                END IF;
                
                RETURN dot_product / (sqrt(norm_a) * sqrt(norm_b));
            END;
            $$ LANGUAGE plpgsql IMMUTABLE;
        """)

        # 2. Create Table
        # Changed embedding from 'vector(768)' to 'double precision[]' (float8[])
        cursor.execute("""
            CREATE TABLE IF NOT EXISTS COURSE_MATERIAL_CHUNK (
                chunk_id SERIAL PRIMARY KEY,
                course_id VARCHAR(10) NOT NULL,
                unit_no INTEGER NOT NULL,
                content_text TEXT NOT NULL,
                embedding double precision[], 
                source_file VARCHAR(255),
                chunk_index INTEGER,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                CONSTRAINT fk_material_course FOREIGN KEY (course_id) REFERENCES COURSE(course_id)
            );
        """)
        conn.commit()
        print("✅ Database setup complete (Native Array Mode).")
    except Exception as e:
        print(f"❌ Database setup failed: {e}")
        conn.rollback()
        raise e
    finally:
        cursor.close()
        conn.close()

def process_course_materials():
    # Ensure DB is ready
    try:
        setup_database()
    except:
        return # Exit if setup failed

    print(f"🚀 Starting Ingestion (Model: {MODEL_NAME})")
    
    # 1. Setup Embeddings
    try:
        from langchain_ollama import OllamaEmbeddings
    except ImportError:
        from langchain_community.embeddings import OllamaEmbeddings

    embeddings = OllamaEmbeddings(model=MODEL_NAME)
    text_splitter = SemanticChunker(embeddings) # Intelligent splitting

    conn = get_db_connection()
    cursor = conn.cursor()

    count = 0
    
    # 2. Walk through folders
    if not os.path.exists(BASE_DIR):
        print(f"❌ Error: '{BASE_DIR}' directory not found. Please create it and add files.")
        return

    for course_id in os.listdir(BASE_DIR):
        course_path = os.path.join(BASE_DIR, course_id)
        if not os.path.isdir(course_path): continue

        print(f"📂 Processing Course: {course_id}")
        
        for unit_no in os.listdir(course_path):
            unit_path = os.path.join(course_path, unit_no)
            if not os.path.isdir(unit_path): continue
            
            try:
                unit_int = int(unit_no)
            except ValueError:
                print(f"  ⚠️  Skipping non-numeric folder: {unit_no}")
                continue

            print(f"  Start Unit {unit_int}...")

            for filename in os.listdir(unit_path):
                if not filename.endswith(".pdf"): continue

                file_path = os.path.join(unit_path, filename)
                print(f"    📄 Reading {filename}...")

                # 3. Load PDF
                loader = PyPDFLoader(file_path)
                docs = loader.load()
                
                # 4. Split and Embed
                chunks = text_splitter.split_documents(docs)
                
                print(f"      -> Generated {len(chunks)} chunks. Inserting...")

                    for i, chunk in enumerate(chunks):
                    content = chunk.page_content.replace('\x00', '')
                    # Generate embedding for single chunk
                    vector = embeddings.embed_query(content)
                    
                    # 5. Insert into DB
                    try:
                        sql = """
                            INSERT INTO COURSE_MATERIAL_CHUNK 
                            (course_id, unit_no, content_text, embedding, source_file, chunk_index)
                            VALUES (%s, %s, %s, %s, %s, %s)
                        """
                        cursor.execute(sql, (course_id, unit_int, content, vector, filename, i))
                        count += 1
                    except psycopg2.errors.ForeignKeyViolation:
                        print(f"      ❌ Skipped chunk: Course '{course_id}' does not exist in COURSE table.")
                        conn.rollback() # Reset transaction
                        continue 
                    except Exception as e:
                        print(f"      ❌ Error inserting chunk: {e}")
                        conn.rollback()
                        continue
                
                conn.commit()

    cursor.close()
    conn.close()
    print(f"\n✅ Done! Inserted {count} total chunks.")

if __name__ == "__main__":
    process_course_materials()
