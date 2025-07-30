# SpiceRoute 🍽️

A microservices-based meal planning and recipe management application that helps users create personalized meal plans, discover recipes, and automate grocery ordering.

## 🏗️ Architecture

SpiceRoute is built using a microservices architecture with the following components:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Gateway       │    │   Profile       │
│                 │◄──►│   Service       │◄──►│   Service       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │                       │
                              ▼                       ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   Planner       │    │   Recipes       │
                       │   Service       │    │   Service       │
                       └─────────────────┘    └─────────────────┘
                              │                       │
                              ▼                       ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   Orderer       │    │   Vector        │
                       │   Service       │    │   Service       │
                       └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   Feedback      │
                       │   Service       │
                       └─────────────────┘
```

## 🚀 Services

### 1. **Gateway Service** (Go)

- **Port**: 8080
- **Purpose**: HTTP API gateway that routes requests to appropriate microservices
- **Endpoints**:
  - `POST /onboarding` - User preference setup
  - `POST /plan/generate` - Generate meal plans
- **Technology**: Go, Chi router, gRPC client

### 2. **Profile Service** (Go)

- **Port**: 50051
- **Purpose**: Manages user preferences and dietary restrictions
- **Features**:
  - Store user cuisine preferences
  - Track allergies and dietary restrictions
  - Manage budget constraints
  - Spice tolerance settings
- **Technology**: Go, gRPC, PostgreSQL

### 3. **Planner Service** (Python)

- **Port**: 50052
- **Purpose**: Generates optimized meal plans using constraint programming
- **Features**:
  - Calorie optimization
  - Budget constraints
  - Minimize cooking sessions
  - Balance nutrition and variety
- **Technology**: Python, FastAPI, OR-Tools (Google's optimization library)

### 4. **Recipes Service** (Go)

- **Port**: 50053
- **Purpose**: Manages recipe database and retrieval
- **Features**:
  - CRUD operations for recipes
  - Recipe search and filtering
  - Ingredient management
  - Nutritional information
- **Technology**: Go, gRPC, PostgreSQL

### 5. **Vector Service** (Python)

- **Port**: 50054
- **Purpose**: Semantic search for recipes using vector embeddings
- **Features**:
  - Recipe embedding generation
  - Similarity search
  - Content-based recommendations
- **Technology**: Python, FastAPI, NumPy, scikit-learn

### 6. **Orderer Service** (Python)

- **Port**: 50055
- **Purpose**: Automates grocery ordering through web automation
- **Features**:
  - DoorDash integration
  - Automated cart filling
  - Order confirmation
- **Technology**: Python, FastAPI, Playwright

### 7. **Feedback Service** (Go)

- **Port**: 50056
- **Purpose**: Collects and stores user feedback on meals
- **Features**:
  - Rating system
  - Skip tracking
  - Substitution logging
  - Comments and reviews
- **Technology**: Go, gRPC, PostgreSQL

## 🛠️ Technology Stack

### Backend

- **Languages**: Go, Python
- **Frameworks**: FastAPI, gRPC
- **Databases**: PostgreSQL
- **Optimization**: Google OR-Tools
- **ML/AI**: NumPy, scikit-learn
- **Web Automation**: Playwright

### Infrastructure

- **Containerization**: Docker
- **Orchestration**: Kubernetes
- **Cloud**: Google Cloud Platform
- **Infrastructure as Code**: Terraform
- **CI/CD**: GitHub Actions

## 📋 Prerequisites

- Go 1.21+
- Python 3.11+
- Docker
- PostgreSQL
- Terraform
- Google Cloud SDK (for deployment)

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd spiceroute
```

### 2. Set Up Python Environment

```bash
# Create virtual environment using uv
uv venv

# Activate virtual environment
source .venv/bin/activate

# Install Python dependencies
uv pip install -r requirements.txt

# Install Playwright browsers
playwright install
```

### 3. Set Up Go Environment

```bash
# Initialize Go module
go mod init spiceroute

# Install dependencies
go mod tidy

# Generate protobuf code
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/spiceroute.proto
```

### 4. Set Up Database

```bash
# Create PostgreSQL database
createdb spiceroute

# Set environment variable
export DB_DSN="postgresql://username:password@localhost:5432/spiceroute?sslmode=disable"
```

### 5. Run Services Locally

#### Start Profile Service

```bash
cd services/profile
go run main.go
```

#### Start Planner Service

```bash
cd services/planner
uvicorn app:app --reload --port 50052
```

#### Start Recipes Service

```bash
cd services/recipes
go run main.go
```

#### Start Vector Service

```bash
cd services/vector
uvicorn app:app --reload --port 50054
```

#### Start Orderer Service

```bash
cd services/orderer
uvicorn app:app --reload --port 50055
```

#### Start Feedback Service

```bash
cd services/feedback
go run main.go
```

#### Start Gateway Service

```bash
cd services/gateway
go run main.go
```

## 🐳 Docker Deployment

### Build Images

```bash
make build
```

### Run with Docker Compose

```bash
docker-compose up -d
```

## ☁️ Cloud Deployment

### 1. Set Up Google Cloud

```bash
# Authenticate with GCP
gcloud auth login

# Set project
gcloud config set project YOUR_PROJECT_ID
```

### 2. Deploy Infrastructure

```bash
cd infra/terraform
terraform init
terraform plan
terraform apply
```

### 3. Deploy to Kubernetes

```bash
# Apply Kubernetes manifests
kubectl apply -f infra/k8s/

# Or use the Makefile
make k8s-deploy
```

## 🔧 Configuration

### Environment Variables

| Service | Variable              | Description                   |
| ------- | --------------------- | ----------------------------- |
| All     | `DB_DSN`              | PostgreSQL connection string  |
| Gateway | `PROFILE_SERVICE_URL` | Profile service gRPC endpoint |
| Gateway | `PLANNER_SERVICE_URL` | Planner service gRPC endpoint |

### Database Schema

The application uses PostgreSQL with the following main tables:

- `preferences` - User dietary preferences
- `recipes` - Recipe database
- `feedback` - User feedback and ratings

## 🧪 Testing

### Run Go Tests

```bash
go test ./...
```

### Run Python Tests

```bash
pytest
```

## 📊 API Documentation

Once services are running, you can access:

- **Gateway API**: http://localhost:8080
- **Planner API**: http://localhost:50052/docs
- **Vector API**: http://localhost:50054/docs
- **Orderer API**: http://localhost:50055/docs

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For issues and questions:

1. Check the documentation
2. Search existing issues
3. Create a new issue with detailed information

---

**SpiceRoute** - Making meal planning deliciously simple! 🍳✨
