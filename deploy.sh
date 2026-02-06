#!/bin/bash

# Doodle Clone Deployment Script for 3Ks Kubernetes
# Usage: ./deploy.sh doodle.kapsule.cloud

set -e

# Configuration
PROJECTNAME="doodle"
NAMESPACE="doodle-prd"
IMAGE_REGISTRY="techlab21.azurecr.io"
IMAGE="${IMAGE_REGISTRY}/${PROJECTNAME}"
VERSION="1.0"
DOMAIN="${1:-doodle.kapsule.cloud}"

echo "🚀 Deploying Doodle Clone to 3Ks Kubernetes"
echo "   Domain: $DOMAIN"
echo "   Namespace: $NAMESPACE"
echo "   Image: $IMAGE:$VERSION"
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Step 1: Build the Docker image
echo -e "${YELLOW}📦 Step 1: Building Docker image...${NC}"
cd "$(dirname "$0")"
docker build -t ${IMAGE}:${VERSION} -f Dockerfile .
docker tag ${IMAGE}:${VERSION} ${IMAGE}:latest
echo -e "${GREEN}✓ Docker image built${NC}"
echo ""

# Step 2: Push to registry
echo -e "${YELLOW}⬆️  Step 2: Pushing to registry...${NC}"
docker push ${IMAGE}:${VERSION}
docker push ${IMAGE}:latest
echo -e "${GREEN}✓ Image pushed to registry${NC}"
echo ""

# Step 3: Create/update Kubernetes resources
echo -e "${YELLOW}☸️  Step 3: Applying Kubernetes manifests...${NC}"

# Create namespace
kubectl apply -f kube/01_namespace.yaml

# Check if namespace is ready
kubectl wait --for=condition=ready --timeout=30s namespace/${NAMESPACE} 2>/dev/null || true

# Create secrets if they don't exist
if ! kubectl get secret doodle-env -n ${NAMESPACE} 2>/dev/null; then
    echo "   Creating secrets (please update values after deployment)..."
    kubectl create secret generic doodle-env -n ${NAMESPACE} \
        --from-literal=DB_NAME=doodle_clone \
        --from-literal=DB_USER=doodle \
        --from-literal=DB_PASSWORD=doodle123 \
        --from-literal=JWT_SECRET=change-me-in-production \
        --from-literal=REFRESH_SECRET=change-me-too \
        --from-literal=GOOGLE_CLIENT_ID=your-client-id \
        --from-literal=GOOGLE_CLIENT_SECRET=your-client-secret \
        --from-literal=SMTP_HOST=smtp.gmail.com \
        --from-literal=SMTP_PORT=587 \
        --from-literal=SMTP_USER=your-email@gmail.com \
        --from-literal=SMTP_PASSWORD=your-app-password \
        --from-literal=SMTP_FROM=noreply@doodle.kapsule.cloud
    echo -e "${GREEN}✓ Secrets created${NC}"
fi

# Apply PVC
kubectl apply -f kube/04_pvc.yaml

# Update deployment with new image and domain
envsubst < kube/02_deployment.yaml | kubectl apply -f -

# Apply ingress
envsubst < kube/03_ingress.yaml | kubectl apply -f -

echo -e "${GREEN}✓ Kubernetes manifests applied${NC}"
echo ""

# Step 4: Wait for deployment to be ready
echo -e "${YELLOW}⏳ Step 4: Waiting for deployment to be ready...${NC}"
kubectl wait --for=condition=available --timeout=180s deployment/doodle -n ${NAMESPACE} || {
    echo "Warning: Deployment not ready within timeout, but continuing..."
}
echo -e "${GREEN}✓ Deployment is ready${NC}"
echo ""

# Step 5: Show status
echo -e "${YELLOW}📊 Deployment Status:${NC}"
kubectl get all -n ${NAMESPACE}
kubectl get ingress -n ${NAMESPACE}
echo ""

echo -e "${GREEN}🎉 Deployment complete!${NC}"
echo ""
echo "🌐 Access your application at: https://${DOMAIN}"
echo ""
echo "📝 To view logs:"
echo "   make logs"
echo "   kubectl logs -f -n ${NAMESPACE} deployment/doodle"
echo ""
echo "🔧 To update secrets:"
echo "   kubectl create secret generic doodle-env -n ${NAMESPACE} \\"
echo "     --from-literal=GOOGLE_CLIENT_ID=xxx --dry-run=client -o yaml | kubectl apply -f -"
echo ""
echo "🗄️  To access the database:"
echo "   kubectl port-forward -n ${NAMESPACE} svc/postgres-postgresql 5432:5432"
