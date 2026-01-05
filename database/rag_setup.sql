-- RAG Setup for Course Materials
-- This schema handles the storage of text chunks for the 5 Course x 5 Unit structure.

-- 1. Create the table for storing text chunks (RAG)
CREATE TABLE IF NOT EXISTS COURSE_MATERIAL_CHUNK (
    chunk_id SERIAL PRIMARY KEY,
    course_id VARCHAR(10) NOT NULL,
    unit_no INTEGER NOT NULL,       -- Matches your 5 units logic (1-5)
    content_text TEXT NOT NULL,     -- The actual note content (paragraph size)
    embedding vector(768),          -- Vector embedding (matching your model)
    source_file VARCHAR(255),       -- e.g., 'Physics_Unit1.pdf'
    chunk_index INTEGER,            -- To maintain reading order (0, 1, 2...)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Link to your existing COURSE table
    CONSTRAINT fk_material_course FOREIGN KEY (course_id) REFERENCES COURSE(course_id)
);

-- Performance Note:
-- For your scale (25 units -> ~5,000 chunks), a specific vector index (IVFFlat/HNSW) 
-- is NOT strictly necessary because PostgreSQL can scan 5,000 rows instantly (< 50ms).
-- If you grow to > 100,000 chunks, uncomment the following:

-- CREATE INDEX ON COURSE_MATERIAL_CHUNK USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
