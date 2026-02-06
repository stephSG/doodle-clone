PROJECTNAME := doodle
VERSION := 1.0
BUILD := $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
BUILD_DATE=$(shell date +%Y-%m-%d_%H:%M)
BUILD_TIME=$(shell date +%FT%T%z)

# Docker image settings
IMAGE := $(PROJECTNAME)
IMAGE_REGISTRY := techlab21.azurecr.io
IMAGE_FULL_PATH=$(IMAGE_REGISTRY)/$(IMAGE)

# Namespace
NAMESPACE := doodle-prd

.PHONY: help dev run run-dev build build-app build-image push-image deploy deploy-prd all

help:
	@echo "Doodle Clone - Build & Deployment"
	@echo ""
	@echo "Development:"
	@echo "  make dev           - Start backend + frontend"
	@echo "  make run           - Run backend locally"
	@echo "  make run-dev       - Run with hot reload"
	@echo ""
	@echo "Building:"
	@echo "  make build         - Build backend binary"
	@echo "  make build-app     - Build for Linux (production)"
	@echo "  make build-image   - Build Docker image"
	@echo "  make push-image    - Push Docker image to registry"
	@echo ""
	@echo "Deployment:"
	@echo "  make deploy        - Deploy to Kubernetes (dev)"
	@echo "  make deploy-prd     - Deploy to Kubernetes (prd)"
	@echo "  make all           - build-image push-image deploy"
	@echo ""
	@echo "Kubernetes:"
	@echo "  make kube-apply     - Apply all Kubernetes manifests"
	@echo "  make kube-delete    - Delete all Kubernetes resources"
	@echo ""
	@echo "Database:"
	@echo "  make db-create     - Create PostgreSQL in Kubernetes"

# =============================================================================
# Development
# =============================================================================

run:
	@echo "Starting backend..."
	@cd backend && go run main.go

run-dev:
	@echo "Starting backend with air..."
	@cd backend && air

dev:
	@echo "Starting development environment..."
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:5173"
	@-pkill -f "go run main.go" 2>/dev/null || true
	@-pkill -f "vite" 2>/dev/null || true
	@make -j2 run-backend-dev run-frontend-dev

run-backend-dev:
	@cd backend && go run main.go

run-frontend-dev:
	@cd frontend && npm run dev

# =============================================================================
# Building
# =============================================================================

build: build-app

build-app:
	@echo "Building backend for Linux..."
	@cd backend && mkdir -p build
	@cd backend && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o build/doodle-backend main.go
	@echo "Binary created: backend/build/doodle-backend"

build-front:
	@echo "Building frontend..."
	@cd frontend && npm run build

# =============================================================================
# Docker
# =============================================================================

build-image: build-app
	@echo "Building Docker image..."
	@docker build -t $(IMAGE_FULL_PATH):$(VERSION) -f Dockerfile .
	@docker tag $(IMAGE_FULL_PATH):$(VERSION) $(IMAGE_FULL_PATH):latest

push-image: build-image
	@echo "Pushing Docker image to registry..."
	@docker push $(IMAGE_FULL_PATH):$(VERSION)
	@docker push $(IMAGE_FULL_PATH):latest

# =============================================================================
# Kubernetes Deployment
# =============================================================================

kube-apply:
	@echo "Applying Kubernetes manifests..."
	@kubectl apply -f kube/01_namespace.yaml
	@echo "Waiting for namespace to be ready..."
	@sleep 2
	-kubectl apply -f kube/04_pvc.yaml
	@echo "Creating secrets..."
	@echo "  (Create secrets manually or use: kubectl create secret generic doodle-env --from-literal=KEY=VALUE)"
	@kubectl apply -f kube/02_deployment.yaml
	@kubectl apply -f kube/03_ingress.yaml
	@echo "Waiting for deployment to be ready..."
	@kubectl wait --for=condition=available --timeout=120s deployment/doodle -n $(NAMESPACE)
	@echo "Deployment complete!"

kube-delete:
	@echo "Deleting Kubernetes resources..."
	@kubectl delete -f kube/03_ingress.yaml --ignore-not-found=true
	@kubectl delete -f kube/02_deployment.yaml --ignore-not-found=true
	@kubectl delete -f kube/04_pvc.yaml --ignore-not-found=true
	@kubectl delete -f kube/01_namespace.yaml --ignore-not-found=true

deploy: kube-apply
	@echo "Deployed to $(NAMESPACE)"

deploy-prd: push-image
	@echo "Deploying to production..."
	@kubectl -n $(NAMESPACE) scale deployment doodle --replicas=0
	@kubectl -n $(NAMESPACE) scale deployment doodle --replicas=1

# =============================================================================
# Database
# =============================================================================

db-create:
	@echo "Creating PostgreSQL deployment..."
	@kubectl apply -f kube/04_pvc.yaml -n $(NAMESPACE)
	@echo "PostgreSQL PVC created"
	@echo "Note: Using shared PostgreSQL cluster at postgres.postgres.svc.cluster.local"

# =============================================================================
# Full pipeline
# =============================================================================

all: build-image push-image deploy-prd
	@echo "Full deployment completed!"

# =============================================================================
# Utility
# =============================================================================

logs:
	@kubectl logs -f -n $(NAMESPACE) deployment/doodle

status:
	@kubectl get all -n $(NAMESPACE)
	@kubectl get ingress -n $(NAMESPACE)

shell:
	@kubectl exec -it -n $(NAMESPACE) deployment/doodle -- /bin/sh
