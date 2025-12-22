# Subnet Network Mapping

## Overview

La page de mapping des sous-réseaux fournit une visualisation interactive de votre infrastructure réseau. Elle permet de voir les relations entre les sous-réseaux, leur hiérarchie, et leur organisation par cloud provider.

## Accès

Accédez à la page de mapping via :
- URL: `/subnets/mapping`
- Navigation: Cliquez sur "Mapping" dans la barre de navigation principale

## Modes de Visualisation

### 1. Hierarchy Mode (Hiérarchie)
- Affiche les sous-réseaux dans une structure hiérarchique parent-enfant
- Les réseaux plus larges (préfixe plus petit) sont affichés en haut
- Les sous-réseaux contenus sont affichés en dessous avec des connexions visuelles
- Idéal pour comprendre la structure d'imbrication des réseaux

### 2. Network Mode (Réseau)
- Organise les sous-réseaux par ordre d'adresse IP
- Disposition en grille basée sur les plages IP
- Utile pour voir la distribution des adresses IP

### 3. Cloud Mode (Cloud)
- Groupe les sous-réseaux par fournisseur cloud
- Chaque provider a sa propre colonne
- Facilite la vue d'ensemble de l'infrastructure multi-cloud

## Fonctionnalités

### Filtres
- **Location**: Filtrer par emplacement géographique
- **Cloud Provider**: Filtrer par fournisseur cloud (AWS, Azure, GCP, Scaleway, OVH)
- Les filtres peuvent être combinés pour une recherche plus précise

### Interactions
- **Clic sur un nœud**: Affiche un panneau de détails avec toutes les informations du sous-réseau
- **Survol**: Affiche un tooltip rapide avec les informations essentielles
- **Molette de souris**: Zoom avant/arrière vers la position du curseur
- **Clic et glisser**: Déplacer le diagramme (pan)
- **Fullscreen**: Mode plein écran pour une meilleure visualisation
- **Export**: Exporter le diagramme (fonctionnalité à venir)

### Contrôles de Navigation

#### Boutons de contrôle
- **Zoom +**: Agrandir le diagramme
- **Zoom -**: Réduire le diagramme  
- **Reset**: Remettre le zoom à 100% et centrer
- **Ajuster**: Ajuster automatiquement le diagramme à la taille de l'écran

#### Raccourcis clavier
- **+ / =**: Zoom avant
- **-**: Zoom arrière
- **0**: Reset du zoom
- **F**: Ajuster à l'écran
- **Esc**: Fermer le panneau de détails

#### Navigation à la souris
- **Molette**: Zoom vers la position du curseur
- **Clic + glisser**: Déplacer le diagramme
- **Clic sur nœud**: Sélectionner/désélectionner
- **Survol**: Afficher tooltip

### Visualisation des Nœuds

Chaque nœud de sous-réseau affiche :
- **CIDR**: L'adresse réseau en notation CIDR
- **Nom**: Le nom du sous-réseau
- **Location**: L'emplacement géographique
- **Barre d'utilisation**: Indicateur visuel du taux d'utilisation
  - Vert: < 40%
  - Orange: 40-60%
  - Jaune: 60-80%
  - Rouge: ≥ 80%
- **Icône du provider**: Logo du fournisseur cloud (si applicable)
- **Pourcentage d'utilisation**: Affiché dans le coin inférieur droit

### Couleurs des Nœuds

Les bordures des nœuds sont colorées selon le cloud provider :
- **AWS**: Orange (#FF9900)
- **Azure**: Bleu (#0078D4)
- **GCP**: Bleu Google (#4285F4)
- **Scaleway**: Violet (#4F0599)
- **OVH**: Bleu foncé (#123F6D)
- **On-Premise**: Gris (#6B7280)

## Panneau de Détails

Lorsqu'un nœud est sélectionné, un panneau latéral affiche :
- Informations de base (CIDR, nom, location, type)
- Informations cloud (provider, région, compte)
- Statistiques d'utilisation (IPs allouées/totales)
- Détails réseau (network, broadcast, plage d'hôtes)

## Cas d'Usage

1. **Audit d'infrastructure**: Vue d'ensemble rapide de tous vos sous-réseaux
2. **Planification de capacité**: Identifier les réseaux saturés
3. **Documentation**: Visualisation claire pour la documentation technique
4. **Troubleshooting**: Comprendre rapidement les relations entre réseaux
5. **Multi-cloud management**: Vue unifiée de l'infrastructure multi-cloud

## Technologies Utilisées

- **React**: Framework UI
- **SVG**: Rendu des diagrammes
- **TypeScript**: Type safety
- **CSS Variables**: Thème clair/sombre

## Améliorations Futures

- Export du diagramme en PNG/SVG
- Zoom et pan interactifs
- Recherche de sous-réseaux dans le diagramme
- Affichage des connexions réseau entre sous-réseaux
- Mode de comparaison côte à côte
- Animations de transition entre les modes
