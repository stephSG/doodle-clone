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

.PHONY: help build build-app build-image push-image deploy deploy-prd all

help:
	@echo "Doodle Clone - Build & Deployment"
	@echo ""
	@echo "Development:"
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
	@kubectl apply -f - <<EOF
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: doodle-postgres-pvc
  namespace: $(NAMESPACE)
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres-postgresql
  namespace: $(NAMESPACE)
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres-postgresql
  template:
    metadata:
      labels:
        app: postgres-postgresql
    spec:
      containers:
        - name: postgres-postgresql
          image: postgres:16-alpine
          env:
            - name: POSTGRES_USER
              value: doodle
            - name: POSTGRES_PASSWORD
              value: doodle123
            - name: POSTGRES_DB
              value: doodle_clone
          ports:
            - containerPort: 5432
          volumeMounts:
            - name: postgres-data
              mountPath: /var/lib/postgresql/data
      volumes:
        - name: postgres-data
          persistentVolumeClaim:
            claimName: doodle-postgres-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: postgres-postgresql
  namespace: $(NAMESPACE)
spec:
  selector:
    app: postgres-postgresql
  ports:
    - port: 5432
    targetPort: 5432
  type: ClusterIP
EOF
	@kubectl wait --for=condition=ready --timeout=120s pod -l app=postgres-postgresql -n $(NAMESPACE)

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
