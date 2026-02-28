CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE memories (
    id UUID PRIMARY KEY,
    content TEXT NOT NULL,
    embedding VECTOR(384),
    confidence FLOAT NOT NULL,
    version INT NOT NULL
);
