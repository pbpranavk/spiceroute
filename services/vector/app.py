from fastapi import FastAPI
from pydantic import BaseModel
import numpy as np
from typing import List
from sklearn.neighbors import NearestNeighbors

app = FastAPI()
embeddings = {}  # dish_id -> np.array
index = None
ids = []

class EmbedRequest(BaseModel):
    id: str
    text: str

class QueryRequest(BaseModel):
    text: str
    k: int = 20

@app.post("/embed")
def embed(req: EmbedRequest):
    vec = fake_embed(req.text)
    embeddings[req.id] = vec
    rebuild()
    return {"ok": True}

@app.post("/query")
def query(req: QueryRequest):
    q = fake_embed(req.text).reshape(1, -1)
    distances, idxs = index.kneighbors(q, n_neighbors=min(req.k, len(ids)))
    return {"ids": [ids[i] for i in idxs[0].tolist()]}

def fake_embed(text: str):
    return np.random.rand(384)

def rebuild():
    global index, ids
    ids = list(embeddings.keys())
    X = np.stack(list(embeddings.values()))
    index = NearestNeighbors(metric="cosine").fit(X)
