# IPAM Data Seeding Scri
ssier contient des scripts pour peupler votre IPAM avec des données d'exemple.

## Prérequis

1. **k6 installé** - Outil de test de performance et d'API
   ```bash
   # macOS
   brew install k6
   
   # Ubuntu/Debian
   sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
   echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
   sudo apt-get update && sudo apt-get install k6
   
   # Autres systèmes
   # Voir: https://k6.io/docs/getting-started/installation/
   ```

2. **Backend IPAM en cours d'exécution**
   ```bash
   task dev:backend
   # ou
   cd backend && go run cmd/server/main.go
   ```

## Utilisation

### Option 1: Script bash (recommandé)
```bash
# Depuis la racine du projet
./scripts/seed-data.sh

# Ou avec une URL personnalisée
BASE_URL=http://localhost:8081 ./scripts/seed-data.sh
```

### Option 2: k6 directement
```bash
# Depuis la racine du projet
k6 run scripts/k6-seed-data.js

# Ou avec une URL personnalisée
BASE_URL=http://localhost:8081 k6 run scripts/k6-seed-data.js
```

## Données générées

Le script crée **20 subnets d'exemple** répartis sur :

### Cloud Providers
- **AWS** (3 subnets)
  - Production Web Tier (10.0.1.0/24)
  - Production App Tier (10.0.2.0/24)
  - Development VPC (10.1.0.0/16)

- **Azure** (3 subnets)
  - Production Frontend (10.2.1.0/24)
  - Production Backend (10.2.2.0/24)
  - Staging Environment (10.3.0.0/16)

- **GCP** (3 subnets)
  - Production Compute (10.4.1.0/24)
  - Production Database (10.4.2.0/24)
  - Development (10.5.0.0/16)

- **Scaleway** (2 subnets)
  - Production API (10.6.1.0/24)
  - Production Storage (10.6.2.0/24)

- **OVH** (1 subnet)
  - Production Web (10.7.1.0/24)

### On-Premise
- **Datacenters** (4 subnets)
  - Paris DC1 Management & Production
  - London DC1 Management & Production

- **Sites** (4 subnets)
  - New York Office
  - San Francisco Office
  - Tokyo Office
  - Berlin Office

## Personnalisation

Vous pouvez modifier le fichier `k6-seed-data.js` pour :
- Ajouter/supprimer des subnets
- Changer les plages IP
- Modifier les informations cloud
- Ajuster les descriptions

## Nettoyage

Pour supprimer toutes les données, vous pouvez :
1. Supprimer la base de données SQLite : `rm backend/data/ipam.db`
2. Redémarrer le backend pour recréer une base vide

## Dépannage

### "k6 command not found"
Installez k6 selon les instructions ci-dessus.

### "API is not available"
Vérifiez que le backend est en cours d'exécution sur le bon port (8081 par défaut).

### "Failed to create subnet"
Vérifiez les logs du backend pour voir les erreurs détaillées.
