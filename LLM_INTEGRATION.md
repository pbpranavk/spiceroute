# SpiceRoute LLM Integration

This document describes the LLM (Large Language Model) enhancements added to the SpiceRoute microservices architecture.

## 🚀 Enhanced Services

### 1. Vector Service (`services/vector/app.py`)

**Enhanced Features:**

- **Real OpenAI Embeddings**: Uses OpenAI's `text-embedding-ada-002` model for high-quality semantic embeddings
- **Sophisticated Fallback**: Intelligent fallback embedding system when OpenAI is unavailable
- **Advanced Search**: Cosine similarity search with configurable thresholds
- **Metadata Support**: Store and query recipe metadata alongside embeddings
- **Health Monitoring**: Service health checks and embedding statistics

**Key Endpoints:**

- `POST /embed` - Create embeddings with metadata
- `POST /query` - Semantic search with similarity scores
- `DELETE /embed/{id}` - Remove embeddings
- `GET /embeddings` - List all stored embeddings
- `GET /health` - Service health check

**Example Usage:**

```python
# Create embedding
response = await client.post("/embed", json={
    "id": "recipe_001",
    "text": "Spicy Chicken Tikka Masala with aromatic spices",
    "metadata": {"cuisine": "indian", "spice_level": "medium"}
})

# Semantic search
response = await client.post("/query", json={
    "text": "I want something spicy and flavorful",
    "k": 5,
    "threshold": 0.7
})
```

### 2. Planner Service (`services/planner/app.py`)

**Enhanced Features:**

- **AI-Powered Recommendations**: GPT-3.5-turbo generates personalized meal planning advice
- **Intelligent Meal Suggestions**: Context-aware meal ideas based on preferences and restrictions
- **Nutrition Analysis**: Detailed nutrition summaries and goals
- **Cost Optimization**: Savings calculations and budget recommendations
- **Dietary Filtering**: Smart filtering based on dietary restrictions

**Key Endpoints:**

- `POST /plan` - Generate optimized meal plans with AI recommendations
- `POST /suggest` - Get AI-powered meal suggestions
- `GET /health` - Service health check

**Example Usage:**

```python
# Get meal suggestions
response = await client.post("/suggest", json={
    "user_id": "user_123",
    "preferences": {"cuisines": ["italian", "indian"], "spice_tolerance": "medium"},
    "dietary_restrictions": ["vegetarian"],
    "budget": 150.0,
    "days": 7
})

# Generate meal plan
response = await client.post("/plan", json={
    "user_id": "user_123",
    "days": 3,
    "dishes": [...],
    "daily_calories": 2000.0,
    "budget_week": 100.0,
    "preferences": {...},
    "dietary_restrictions": []
})
```

## 🔧 Configuration

### Environment Variables

Add to your Kubernetes secrets (`infra/k8s/secret.yaml`):

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: spiceroute-secret
  namespace: spiceroute
type: Opaque
stringData:
  DB_DSN: "<your-cloudsql-dsn>"
  OPENAI_API_KEY: "<your-openai-api-key>"
```

### Dependencies

Updated `requirements.txt`:

```
# OpenAI for embeddings and LLM integration
openai==1.3.7
```

## 🧪 Testing

Run the test suite to verify LLM integration:

```bash
# Set your OpenAI API key
export OPENAI_API_KEY="your-api-key-here"

# Run the test script
python test_llm_services.py
```

The test script demonstrates:

- Embedding creation and semantic search
- AI-powered meal suggestions
- Enhanced meal planning with recommendations
- Fallback behavior when OpenAI is unavailable

## 🎯 Key Benefits

### Vector Service Enhancements

1. **Better Semantic Understanding**: OpenAI embeddings capture nuanced meaning
2. **Improved Search Quality**: More accurate recipe recommendations
3. **Robust Fallback**: Works even without API access
4. **Scalable Architecture**: Ready for production vector databases

### Planner Service Enhancements

1. **Personalized Recommendations**: AI understands user preferences
2. **Nutritional Intelligence**: Smart nutrition goals and analysis
3. **Cost Optimization**: Savings calculations and budget advice
4. **Dietary Compliance**: Automatic filtering for restrictions

## 🔄 Fallback Behavior

Both services gracefully handle OpenAI API unavailability:

- **Vector Service**: Uses sophisticated hash-based embeddings with cuisine-specific patterns
- **Planner Service**: Provides curated fallback meal suggestions and basic recommendations

## 📊 Performance Considerations

- **Embedding Generation**: ~1-2 seconds per recipe with OpenAI
- **Search Queries**: Sub-second response times with in-memory index
- **AI Recommendations**: ~2-3 seconds for meal suggestions
- **Memory Usage**: ~1GB for 10,000 embeddings (in-memory)

## 🚀 Production Recommendations

1. **Vector Database**: Replace in-memory storage with Pinecone, Weaviate, or Qdrant
2. **Caching**: Implement Redis for frequently accessed embeddings
3. **Rate Limiting**: Add OpenAI API rate limiting
4. **Monitoring**: Add metrics for embedding generation and search performance
5. **Backup Embeddings**: Store fallback embeddings for critical recipes

## 🔐 Security

- API keys stored in Kubernetes secrets
- No sensitive data in embeddings
- Input validation and sanitization
- Error handling prevents information leakage

## 📈 Future Enhancements

1. **Multi-modal Embeddings**: Include recipe images
2. **Personalized Models**: Fine-tune on user feedback
3. **Recipe Generation**: AI-powered recipe creation
4. **Nutritional AI**: Advanced nutrition analysis
5. **Cost Prediction**: ML-based cost estimation
