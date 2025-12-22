# Système de Connexions Réseau

## Vue d'ensemble

Le système de connexions permet de modéliser et visualiser les liens entre vos sous-réseaux. Vous pouvez définir différents types de connexions (VPN, Peering, NAT Gateway, etc.) et les voir représentées graphiquement dans le diagramme de mapping.

## Types de Connexions Supportés

### 1. VPN Site-à-Site
- **Couleur**: Violet (#8B5CF6)
- **Usage**: Connexions VPN entre sites distants
- **Exemple**: Liaison entre datacenter principal et site distant

### 2. Client OpenVPN
- **Couleur**: Cyan (#06B6D4)
- **Usage**: Connexions VPN client-serveur
- **Exemple**: Accès distant pour employés

### 3. NAT Gateway
- **Couleur**: Vert (#10B981)
- **Usage**: Passerelle de traduction d'adresses
- **Exemple**: Accès internet pour sous-réseaux privés

### 4. Internet Gateway
- **Couleur**: Orange (#F59E0B)
- **Usage**: Passerelle vers internet
- **Exemple**: Accès internet direct

### 5. Peering
- **Couleur**: Rose (#EC4899)
- **Usage**: Connexion directe entre VPCs/réseaux
- **Exemple**: Peering entre VPCs AWS

### 6. Transit Gateway
- **Couleur**: Indigo (#6366F1)
- **Usage**: Hub de connexion centralisé
- **Exemple**: AWS Transit Gateway

### 7. Direct Connect / ExpressRoute / Cloud Interconnect
- **Couleur**: Bleu (#3B82F6)
- **Usage**: Connexions dédiées vers le cloud
- **Exemple**: Ligne dédiée AWS Direct Connect

### 8. Load Balancer
- **Usage**: Répartition de charge
- **Exemple**: ALB/NLB distributing traffic

### 9. Firewall
- **Usage**: Filtrage et sécurité réseau
- **Exemple**: Pare-feu entre zones

### 10. Personnalisé
- **Usage**: Types de connexions spécifiques à votre infrastructure

## États des Connexions

### Actif (Active)
- **Indicateur**: ✓ vert
- **Opacité**: 100%
- **Description**: Connexion opérationnelle

### Inactif (Inactive)
- **Indicateur**: ✗ gris
- **Opacité**: 30%
- **Description**: Connexion désactivée

### En attente (Pending)
- **Indicateur**: ⚠ jaune
- **Opacité**: 60%
- **Style**: Ligne pointillée
- **Description**: Connexion en cours de configuration

### Erreur (Error)
- **Indicateur**: ⚠ rouge
- **Opacité**: 80%
- **Couleur**: Rouge
- **Description**: Connexion en erreur

## Gestion des Connexions

### Accès
- URL: `/subnets/connections`
- Navigation: Cliquez sur "Connexions" dans la barre de navigation

### Création d'une Connexion

1. **Informations de base**:
   - Sous-réseau source (obligatoire)
   - Sous-réseau cible (obligatoire)
   - Type de connexion (obligatoire)
   - Nom de la connexion (obligatoire)

2. **Informations optionnelles**:
   - Description détaillée
   - Bande passante (ex: "1Gbps", "100Mbps")
   - Latence en millisecondes
   - Coût mensuel en euros

3. **Métadonnées**:
   - Champ libre pour informations spécifiques

### Modification d'une Connexion

- Cliquez sur l'icône d'édition dans la carte de connexion
- Modifiez les informations souhaitées
- Sauvegardez les changements

### Suppression d'une Connexion

- Cliquez sur l'icône de suppression
- Confirmez la suppression dans la boîte de dialogue

## Visualisation dans le Mapping

### Représentation Graphique
- **Lignes colorées**: Chaque type de connexion a sa couleur
- **Épaisseur**: Connexions réseau plus épaisses (3px) que hiérarchie (2px)
- **Style**: Lignes pleines (actif) ou pointillées (en attente)
- **Opacité**: Varie selon l'état de la connexion

### Étiquettes
- Nom de la connexion affiché au centre de la ligne
- Visible uniquement pour les connexions réseau

### Légende
- Affichée automatiquement quand des connexions sont présentes
- Située en bas à droite du diagramme
- Liste les couleurs et types de connexions

## Cas d'Usage

### 1. Documentation d'Architecture
- Visualiser l'architecture réseau complète
- Documenter les flux de données
- Identifier les points de défaillance

### 2. Planification de Capacité
- Suivre la bande passante utilisée
- Identifier les goulots d'étranglement
- Planifier les montées en charge

### 3. Gestion des Coûts
- Suivre les coûts des connexions
- Optimiser les dépenses réseau
- Budgétiser les nouvelles connexions

### 4. Troubleshooting
- Identifier rapidement les connexions en erreur
- Visualiser les chemins réseau
- Diagnostiquer les problèmes de connectivité

### 5. Conformité et Audit
- Documenter les flux de données sensibles
- Vérifier la conformité aux politiques
- Préparer les audits de sécurité

## Bonnes Pratiques

### Nommage
- Utilisez des noms descriptifs et cohérents
- Incluez la direction du flux si pertinent
- Exemple: "VPN-Paris-Londres-Prod"

### Documentation
- Remplissez les descriptions pour les connexions critiques
- Documentez les configurations spéciales
- Maintenez les métadonnées à jour

### Monitoring
- Mettez à jour les états régulièrement
- Surveillez les connexions en erreur
- Planifiez la maintenance des connexions

### Organisation
- Groupez les connexions par fonction
- Utilisez des conventions de nommage
- Documentez les dépendances critiques

## Intégration Future

Le système est conçu pour être étendu avec :
- Import/export de configurations réseau
- Intégration avec des outils de monitoring
- Alertes automatiques sur les changements d'état
- Métriques de performance en temps réel
- Validation automatique de la connectivité
