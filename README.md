# 📊 Doodle Clone - Système de Sondage de Dates

Application complète de type Doodle permettant de créer des événements, proposer des dates pour vote, et collecter les préférences des participants.

![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)
![Go](https://img.shields.io/badge/Go-1.24+-00ADD8E.svg)
![Vue](https://img.shields.io/badge/Vue-3.4+-4FC08D.svg)

## 🎯 Fonctionnalités

### 👤 Authentification
- **Connexion Google OAuth2** - Authentification en un clic
- **Email/Mot de passe** - Inscription traditionnelle avec hash bcrypt
- **Tokens JWT** - Access token (15min) + Refresh token (7 jours, httpOnly cookie)
- **Récupération de mot de passe** - Système de récupération par email

### 📋 Sondages
- **Création d'événements** - Titre, description, lieu, dates
- **Options de vote** : Oui, Non, Peut-être
- **Anonymat** - Possibilité de voter sans compte
- **Dates finales** - Fixer la date retenue
- **Privé** - Sondages accessibles uniquement via code d'accès unique

### 🗳️ Gestion des Votes
- **Vote multiple** - Permettre plusieurs sélections
- **Limite de votes** - Restreindre le nombre de votes par utilisateur
- **Votes anonymes** - Vote avec nom personnalisé
- **Mise à jour** - Modifier son vote à tout moment

### 🔔 Notifications
- **Rappel automatique** - X heures avant l'événement (configurable)
- **Notification date finale** - Quand la date est fixée
- **Paramétrable** - Activé/désactivé par l'admin

### 📤 Exports
- **PDF** - Export du sondage avec résultats
- **ICS** - Fichier calendrier (Google Calendar, Outlook)
- **CSV** - Données pour analyse

## 🏗️ Architecture

```
doodle-clone/
├── backend/                 # API Go (Gin + PostgreSQL)
│   ├── main.go              # Point d'entrée
│   ├── internal/
│   │   ├── config/          # Configuration variables d'environnement
│   │   ├── database/        # Connexion & migrations PostgreSQL
│   │   ├── models/          # Modèles de données
│   │   ├── handlers/        # API HTTP handlers
│   │   ├── middleware/      # Auth, CORS, Rate limiting
│   │   └── email/           # Envoi d'emails
│   └── .env                 # Variables d'environnement
│
└── frontend/                # Vue 3 SPA
    ├── src/
    │   ├── assets/          # Styles globaux, images
    │   ├── components/      # Composants réutilisables
    │   ├── router/          # Routes Vue Router
    │   ├── stores/          # Pinia state management
    │   ├── views/           # Pages de l'application
    │   ├── utils/           # Helpers (api, validators)
    │   └── main.js          # Point d'entrée Vue
    ├── index.html
    ├── package.json
    ├── vite.config.js
    └── tailwind.config.js   # Configuration Tailwind CSS
```

## 🚀 Installation

### Prérequis
- Go 1.24+
- Node.js 18+
- PostgreSQL 14+

### Backend

```bash
cd backend

# Créer le fichier .env
cp .env.example .env

# Modifier les variables d'environnement
nano .env
```

Variables requises :
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=doodle_clone
DB_USER=your_user
DB_PASSWORD=your_password

# JWT
JWT_SECRET=votre_clé_secrète_à_changer
REFRESH_SECRET=votre_autre_clé

# Frontend
FRONTEND_URL=http://localhost:5173

# Google OAuth (optionnel)
GOOGLE_CLIENT_ID=votre_client_id
GOOGLE_CLIENT_SECRET=votre_client_secret
GOOGLE_REDIRECT_URI=http://localhost:5173/auth/callback

# SMTP (pour les notifications)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=votre_email@gmail.com
SMTP_PASSWORD=votre_mot_de_passe_app
SMTP_FROM=Bot Doodle <noreply@example.com>
```

```bash
# Installer les dépendances
go mod download

# Lancer le serveur
make run
```

Le backend sera accessible sur `http://localhost:8080`

### Frontend

```bash
cd frontend

# Installer les dépendances
npm install

# Lancer le serveur de développement
npm run dev
```

Le frontend sera accessible sur `http://localhost:5173`

## 📚 API Documentation

### Authentification

#### Inscription
```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "Password123",
  "name": "John Doe"
}
```

#### Connexion
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "Password123"
}
```

#### Google OAuth
```http
GET /auth/google/login
```

### Sondages

#### Créer un sondage (authentifié)
```http
POST /api/polls
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Réunion d'équipe",
  "description": "Point sur l'avancement",
  "location": "Bureau A - 2ème étage",
  "dates": [
    {"start_time": "2026-03-01T10:00:00Z"},
    {"start_time": "2026-03-02T14:00:00Z"}
  ],
  "allow_maybe": true,
  "anonymous": true
}
```

#### Récupérer un sondage (par UUID ou code d'accès)
```http
GET /api/polls/{id_or_code}
```

#### Voter (anonyme ou authentifié)
```http
POST /api/polls/{id_or_code}/vote
Content-Type: application/json

{
  "votes": [
    {"date_option_id": "uuid-date-option", "response": "yes"}
  ],
  "user_name": "Marie Dupont"  // Requis si non authentifié
}
```

#### Fixer la date finale
```http
POST /api/polls/{id}/final
Authorization: Bearer <token>
Content-Type: application/json

{
  "date_option_id": "uuid-date-option"
}
```

### Routes

| Méthode | Route | Description | Auth |
|---------|-------|-------------|-----|
| GET | `/api/polls` | Liste des sondages publics | Non |
| GET | `/api/polls/:id` | Détails d'un sondage | Non |
| POST | `/api/polls/:id/vote` | Voter (anonyme ok) | Optionnel |
| POST | `/api/polls/:id/votes` | Voter (auth requis) | Oui |
| GET | `/api/polls/:id/export/pdf` | Export PDF | Non |
| GET | `/api/polls/:id/export/ics` | Export calendrier | Non |

## 🧪 Tests

```bash
cd backend
go test ./...

cd frontend
npm run test
```

## 📦 Déploiement

### Docker (recommandé)

```bash
docker-compose up -d
```

### Production

```bash
# Frontend
cd frontend
npm run build

# Backend
cd backend
go build -o bin/doodle-backend .
./bin/doodle-backend
```

## 🔧 Configuration Admin

Les notifications sont configurables via API :

```http
GET /api/notifications/settings
Authorization: Bearer <admin_token>
```

Paramètres :
- `reminder_enabled` : Activer les rappels (défaut: true)
- `reminder_hours` : Heures avant l'événement (défaut: 1)
- `final_date_enabled` : Notification date finale (défaut: true)
- `new_vote_enabled` : Notification nouveau vote (défaut: false)
- `new_comment_enabled` : Notification nouveau commentaire (défaut: false)

## 📄 Licence

MIT

## 👥 Contributeurs

Stéphane LE MINH NHUT

---

**Doodle Clone** - Une solution moderne de planification d'événements.
