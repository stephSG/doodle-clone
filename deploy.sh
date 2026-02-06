#!/bin/bash
set -e

# Configuration
SERVER="ubuntu@51.254.139.110"
SSH_KEY="~/.ssh/id_ed25519_2"
REMOTE_DIR="/home/ubuntu/kapsule"
PROJECTNAME="doodle"

# GitHub Container Registry Configuration
# Pour plus de sécurité, utiliser des variables d'environnement:
# export GITHUB_TOKEN="ghp_xxx"
GITHUB_USER="${GITHUB_USER:-stephSG}"
GITHUB_TOKEN="${GITHUB_TOKEN:-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}"

# Usage check
if [ -z "$1" ]; then
    echo "Usage: $0 <domain>"
    echo "Example: $0 doodle.kapsule.cloud"
    exit 1
fi

DOMAIN=$1

# Extract subdomain for namespace and image name (e.g. doodle.kapsule.cloud -> doodle)
if [[ "$DOMAIN" == *"."* ]]; then
    NAMESPACE=$(echo "$DOMAIN" | cut -d. -f1)
else
    NAMESPACE=$DOMAIN
fi

# Docker image (GitHub Container Registry)
DOCKER_REGISTRY="ghcr.io"
DOCKER_IMAGE_NAME="${DOCKER_REGISTRY}/$(echo $GITHUB_USER | tr '[:upper:]' '[:lower:]')/${NAMESPACE}"
GIT_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "latest")
IMAGE_NAME="${DOCKER_IMAGE_NAME}:${GIT_TAG}"

echo "=== Doodle Clone K3S Deployment to GHCR ==="
echo "Target Domain: $DOMAIN"
echo "Target Namespace: $NAMESPACE"
echo "Target Image: $IMAGE_NAME"
echo "Target Server: $SERVER"

# 1. Check GitHub credentials (token is set with default value)
echo "GitHub User: $GITHUB_USER"

# 2. Sync code and build Docker image on server
echo ""
echo "[1/5] Syncing code to server..."
rsync -avz --delete \
    --exclude '.git' \
    --exclude 'node_modules' \
    --exclude 'frontend/node_modules' \
    --exclude 'backend/bin' \
    --exclude 'backend/build' \
    --exclude 'frontend/dist' \
    --exclude '.env' \
    -e "ssh -i $SSH_KEY" \
    . $SERVER:$REMOTE_DIR/$PROJECTNAME/

echo "[2/5] Building Docker image on server..."
ssh -i $SSH_KEY $SERVER "cd $REMOTE_DIR/$PROJECTNAME && docker build -t $IMAGE_NAME ."

echo "[3/5] Pushing image to GitHub Container Registry..."
ssh -i $SSH_KEY $SERVER "docker push $IMAGE_NAME"

# 4. Deploy to K3S
echo "[4/5] Deploying to K3S..."

ssh -i $SSH_KEY $SERVER "
    export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
    export GH_USER='$(echo $GITHUB_USER | tr '[:upper:]' '[:lower:]')'
    export GH_TOKEN='$GITHUB_TOKEN'

    # Create namespace if not exists
    sudo kubectl create namespace $NAMESPACE --dry-run=client -o yaml | sudo kubectl apply -f -

    # Create ImagePullSecret for GHCR
    sudo kubectl create secret docker-registry ghcr-secret \
        --docker-server=ghcr.io \
        --docker-username=\$GH_USER \
        --docker-password=\$GH_TOKEN \
        --docker-email=\$GH_USER@users.noreply.github.com \
        -n $NAMESPACE \
        --dry-run=client -o yaml | sudo kubectl apply -f -

    # Create secret with env variables (if doesn't exist)
    sudo kubectl create secret generic doodle-env -n $NAMESPACE \\
        --from-literal=DB_NAME=doodle_clone \\
        --from-literal=DB_USER=doodle \\
        --from-literal=DB_PASSWORD=doodle123 \\
        --from-literal=JWT_SECRET=change-me-in-production \\
        --from-literal=REFRESH_SECRET=change-me-too \\
        --from-literal=GOOGLE_CLIENT_ID=your-client-id \\
        --from-literal=GOOGLE_CLIENT_SECRET=your-client-secret \\
        --from-literal=SMTP_HOST=smtp.gmail.com \\
        --from-literal=SMTP_PORT=587 \\
        --from-literal=SMTP_USER=your-email@gmail.com \\
        --from-literal=SMTP_PASSWORD=your-app-password \\
        --from-literal=SMTP_FROM=noreply@$DOMAIN \\
        --dry-run=client -o yaml | sudo kubectl apply -f - 2>/dev/null || true

    # Update and apply Kubernetes manifests
    cd $REMOTE_DIR/$PROJECTNAME

    # Update deployment with dynamic image and domain
    sed -i \"s|ghcr.io/stephsg/doodle:.*|$IMAGE_NAME|g\" kube/02_deployment.yaml
    sed -i \"s|\${DOMAIN}|$DOMAIN|g\" kube/02_deployment.yaml

    # Update ingress with dynamic domain
    sed -i \"s|\${DOMAIN}|$DOMAIN|g\" kube/03_ingress.yaml

    # Update namespace
    sed -i \"s|doodle-prd|$NAMESPACE|g\" kube/01_namespace.yaml
    sed -i \"s|doodle-prd|$NAMESPACE|g\" kube/02_deployment.yaml
    sed -i \"s|doodle-prd|$NAMESPACE|g\" kube/03_ingress.yaml
    sed -i \"s|doodle-prd|$NAMESPACE|g\" kube/04_pvc.yaml

    # Apply manifests
    sudo kubectl apply -f kube/01_namespace.yaml
    sudo kubectl apply -f kube/04_pvc.yaml
    sudo kubectl apply -f kube/02_deployment.yaml
    sudo kubectl apply -f kube/03_ingress.yaml

    # Restart deployment
    echo 'Restarting deployment...'
    sudo kubectl rollout restart deployment/doodle -n $NAMESPACE

    echo ''
    echo 'Waiting for deployment...'
    sudo kubectl wait --for=condition=available --timeout=180s deployment/doodle -n $NAMESPACE || true

    echo ''
    echo 'Pod status:'
    sudo kubectl get pods -n $NAMESPACE
"

echo ""
echo "=== Deployment complete! ==="
echo "App URL: https://$DOMAIN"
echo "Namespace: $NAMESPACE"
echo "Image: $IMAGE_NAME"
echo ""
echo "📝 To update secrets:"
echo "   ssh -i $SSH_KEY $SERVER"
echo "   sudo kubectl create secret generic doodle-env -n $NAMESPACE --from-literal=KEY=VALUE --dry-run=client -o yaml | sudo kubectl apply -f -"
echo ""
echo "🗄️  To access the database:"
echo "   ssh -i $SSH_KEY $SERVER"
echo "   sudo kubectl port-forward -n $NAMESPACE svc/postgres-postgresql 5432:5432"
