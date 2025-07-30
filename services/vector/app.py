import os
import logging
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import numpy as np
from typing import List, Optional
from sklearn.neighbors import NearestNeighbors
import openai
from openai import OpenAI
import asyncio
import json

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(title="SpiceRoute Vector Service", version="1.0.0")

# Initialize OpenAI client
client = None
try:
    api_key = os.getenv("OPENAI_API_KEY")
    if api_key:
        client = OpenAI(api_key=api_key)
        logger.info("OpenAI client initialized successfully")
    else:
        logger.warning("OPENAI_API_KEY not found, using fallback embeddings")
except Exception as e:
    logger.error(f"Failed to initialize OpenAI client: {e}")

# In-memory storage (in production, use a proper vector database like Pinecone or Weaviate)
embeddings = {}  # dish_id -> np.array
index = None
ids = []

class EmbedRequest(BaseModel):
    id: str
    text: str
    metadata: Optional[dict] = None

class QueryRequest(BaseModel):
    text: str
    k: int = 20
    threshold: Optional[float] = 0.7

class EmbedResponse(BaseModel):
    id: str
    success: bool
    message: str
    embedding_size: Optional[int] = None

class QueryResponse(BaseModel):
    results: List[dict]
    query_text: str
    total_found: int

@app.get("/health")
def health_check():
    return {
        "status": "healthy",
        "embeddings_count": len(embeddings),
        "openai_available": client is not None
    }

@app.post("/embed", response_model=EmbedResponse)
async def embed(req: EmbedRequest):
    try:
        # Generate embedding using OpenAI
        if client:
            embedding = await generate_openai_embedding(req.text)
        else:
            # Fallback to a more sophisticated fake embedding
            embedding = generate_fallback_embedding(req.text)
        
        embeddings[req.id] = embedding
        
        # Store metadata if provided
        if req.metadata:
            # In a real implementation, you'd store this in a database
            logger.info(f"Stored metadata for {req.id}: {req.metadata}")
        
        # Rebuild the index
        rebuild_index()
        
        logger.info(f"Successfully embedded text for ID: {req.id}")
        return EmbedResponse(
            id=req.id,
            success=True,
            message="Embedding created successfully",
            embedding_size=len(embedding)
        )
        
    except Exception as e:
        logger.error(f"Error creating embedding for {req.id}: {e}")
        raise HTTPException(status_code=500, detail=f"Failed to create embedding: {str(e)}")

@app.post("/query", response_model=QueryResponse)
async def query(req: QueryRequest):
    try:
        if not embeddings:
            return QueryResponse(
                results=[],
                query_text=req.text,
                total_found=0
            )
        
        # Generate query embedding
        if client:
            query_embedding = await generate_openai_embedding(req.text)
        else:
            query_embedding = generate_fallback_embedding(req.text)
        
        query_embedding = query_embedding.reshape(1, -1)
        
        # Perform similarity search
        distances, indices = index.kneighbors(
            query_embedding, 
            n_neighbors=min(req.k, len(ids))
        )
        
        # Format results
        results = []
        for i, (distance, idx) in enumerate(zip(distances[0], indices[0])):
            similarity = 1 - distance  # Convert distance to similarity
            if req.threshold and similarity < req.threshold:
                continue
                
            results.append({
                "id": ids[idx],
                "similarity": float(similarity),
                "rank": i + 1
            })
        
        logger.info(f"Query '{req.text}' returned {len(results)} results")
        return QueryResponse(
            results=results,
            query_text=req.text,
            total_found=len(results)
        )
        
    except Exception as e:
        logger.error(f"Error processing query '{req.text}': {e}")
        raise HTTPException(status_code=500, detail=f"Failed to process query: {str(e)}")

@app.delete("/embed/{embedding_id}")
async def delete_embedding(embedding_id: str):
    try:
        if embedding_id in embeddings:
            del embeddings[embedding_id]
            rebuild_index()
            logger.info(f"Deleted embedding: {embedding_id}")
            return {"success": True, "message": f"Embedding {embedding_id} deleted"}
        else:
            raise HTTPException(status_code=404, detail=f"Embedding {embedding_id} not found")
    except Exception as e:
        logger.error(f"Error deleting embedding {embedding_id}: {e}")
        raise HTTPException(status_code=500, detail=f"Failed to delete embedding: {str(e)}")

@app.get("/embeddings")
def list_embeddings():
    return {
        "total_embeddings": len(embeddings),
        "embedding_ids": list(embeddings.keys())
    }

async def generate_openai_embedding(text: str) -> np.ndarray:
    """Generate embedding using OpenAI's text-embedding-ada-002 model"""
    try:
        response = client.embeddings.create(
            model="text-embedding-ada-002",
            input=text
        )
        return np.array(response.data[0].embedding)
    except Exception as e:
        logger.error(f"OpenAI embedding failed: {e}")
        # Fallback to our own embedding
        return generate_fallback_embedding(text)

def generate_fallback_embedding(text: str) -> np.ndarray:
    """Generate a more sophisticated fallback embedding based on text characteristics"""
    # Create a deterministic embedding based on text properties
    text_lower = text.lower()
    
    # Simple hash-based embedding
    hash_val = hash(text_lower)
    np.random.seed(hash_val)
    
    # Create embedding with some structure based on text length and content
    embedding = np.random.rand(1536)  # Same size as OpenAI's ada-002
    
    # Add some structure based on text characteristics
    length_factor = min(len(text) / 1000, 1.0)  # Normalize by expected max length
    embedding[:100] *= length_factor
    
    # Add cuisine-specific patterns (simple keyword matching)
    cuisine_keywords = {
        'italian': ['pasta', 'pizza', 'tomato', 'basil', 'olive'],
        'indian': ['curry', 'spice', 'rice', 'naan', 'masala'],
        'chinese': ['noodle', 'soy', 'ginger', 'wok', 'dumpling'],
        'mexican': ['taco', 'salsa', 'chili', 'corn', 'lime'],
        'japanese': ['sushi', 'miso', 'wasabi', 'nori', 'tempura']
    }
    
    for cuisine, keywords in cuisine_keywords.items():
        if any(keyword in text_lower for keyword in keywords):
            # Add cuisine-specific pattern to embedding
            cuisine_hash = hash(cuisine)
            np.random.seed(cuisine_hash)
            cuisine_pattern = np.random.rand(100)
            embedding[100:200] += cuisine_pattern * 0.1
    
    # Normalize the embedding
    embedding = embedding / np.linalg.norm(embedding)
    return embedding

def rebuild_index():
    """Rebuild the nearest neighbors index"""
    global index, ids
    if not embeddings:
        index = None
        ids = []
        return
    
    ids = list(embeddings.keys())
    X = np.stack(list(embeddings.values()))
    
    # Use cosine similarity for better semantic matching
    index = NearestNeighbors(
        n_neighbors=min(20, len(ids)), 
        metric="cosine",
        algorithm="brute"  # Better for cosine similarity
    ).fit(X)
    
    logger.info(f"Rebuilt index with {len(ids)} embeddings")

# Initialize index on startup
@app.on_event("startup")
async def startup_event():
    logger.info("Starting SpiceRoute Vector Service")
    rebuild_index()
