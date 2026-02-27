from sentence_transformers import SentenceTransformer

model = SentenceTransformer('all-MiniLM-L6-v2')

def get_embedding(text: str) -> list[float]:
    """
    Generates an embedding for the given text.
    """
    embedding = model.encode(text)
    return embedding.tolist()
