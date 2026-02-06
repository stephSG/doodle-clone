# 📦 Deployment sur 3Ks Kubernetes

Ce guide explique comment déployer l'application Doodle Clone sur le cluster Kubernetes 3Ks.

## Prérequis

- `kubectl` configuré pour le cluster 3Ks
- `docker` installé
- Accès au registry `techlab21.azurecr.io`

## Structure

```
kube/
├── 01_namespace.yaml    # Namespace doodle-prd
├── 02_deployment.yaml   # Deployment + Service
├── 03_ingress.yaml      # Ingress avec TLS (doodle.kapsule.cloud)
└── 04_pvc.yaml          # PersistentVolumeClaim pour PostgreSQL
```

## Déploiement rapide

### 1. Cloner et se placer dans le projet

```bash
cd /path/to/doodle-clone
```

### 2. Lancer le script de déploiement

```bash
./deploy.sh doodle.kapsule.cloud
```

Le script va :
1. ✅ Builder l'image Docker
2. ✅ Pusher l'image vers le registry Azure
3. ✅ Appliquer les manifests Kubernetes
4. ✅ Créer les secrets (à mettre à jour après le déploiement)
5. ✅ Déployer l'application

### 3. Mettre à jour les secrets

Après le premier déploi, mettez à jour les secrets avec vos vraies valeurs :

```bash
kubectl create secret generic doodle-env -n doodle-prd \
  --from-literal=JWT_SECRET=votre-clé-secrète \
  --from-literal=REFRESH_SECRET=votre-autre-clé \
  --from-literal=GOOGLE_CLIENT_ID=votre-client-id \
  --from-literal=GOOGLE_CLIENT_SECRET=votre-client-secret \
  --from-literal=SMTP_HOST=smtp.gmail.com \
  --from-literal=SMTP_PORT=587 \
  --from-literal=SMTP_USER=votre-email@gmail.com \
  --from-literal=SMTP_PASSWORD=votre-mot-de-passe-app \
  --from-literal=SMTP_FROM=noreply@doodle.kapsule.cloud \
  --dry-run=client -o yaml | kubectl apply -f -
```

### 4. Créer la base de données PostgreSQL

```bash
make db-create
```

Ou manuellement :

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: doodle-postgres-pvc
  namespace: doodle-prd
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
  namespace: doodle-prd
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
  namespace: doodle-prd
spec:
  selector:
    app: postgres-postgresql
  ports:
    - port: 5432
    targetPort: 5432
  type: ClusterIP
EOF
```

## Commandes utiles

```bash
# Voir les logs
make logs
kubectl logs -f -n doodle-prd deployment/doodle

# Voir le statut
make status
kubectl get all -n doodle-prd

# Shell dans le conteneur
make shell
kubectl exec -it -n doodle-prd deployment/doodle -- /bin/sh

# Supprimer tout
make kube-delete

# Redéploiement (rolling update)
make deploy-prd
```

## Configuration Google OAuth

Pour configurer Google OAuth :

1. Allez sur [Google Cloud Console](https://console.cloud.google.com/)
2. Créez un projet OAuth 2.0
3. Ajoutez `https://doodle.kapsule.cloud/auth/google/callback` aux URI de redirection autorisées
4. Récupérez le Client ID et Secret
5. Mettez à jour le secret Kubernetes :

```bash
kubectl create secret generic doodle-env -n doodle-prd \
  --from-literal=GOOGLE_CLIENT_ID=votre-id \
  --from-literal=GOOGLE_CLIENT_SECRET=votre-secret \
  --dry-run=client -o yaml | kubectl apply -f -
```

## Accès

- **Application** : https://doodle.kapsule.cloud
- **API Swagger** : https://doodle.kapsule.cloud/swagger/index.html

## Dépannage

```bash
# Vérifier les pods
kubectl get pods -n doodle-prd

# Vérifier les logs
kubectl logs -n doodle-prd -l app=doodle

# Décrire le pod
kubectl describe pod -n doodle-prd <pod-name>

# Port-forward pour tests locaux
kubectl port-forward -n doodle-prd svc/doodle-backend 8080:8080
```
