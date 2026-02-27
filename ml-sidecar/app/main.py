from fastapi import FastAPI
from pydantic import BaseModel
from .embeddings import get_embedding

class Text(BaseModel):
    text: str

app = FastAPI()

@app.get("/")
def read_root():
    return {"Hello": "World"}

@app.post("/embedding")
def get_embedding_endpoint(text: Text):
    return get_embedding(text.text)
