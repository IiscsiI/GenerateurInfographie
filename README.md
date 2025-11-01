# 🛡️ Générateur d'Infographie Cybersécurité v1.6.3

> **Outil open source pour créer des procédures d'urgence en cas de cyberattaque**  
> Personnalisable, responsive et prêt à l'emploi - Architecture refactorisée et optimisée

[![Version](https://img.shields.io/badge/version-1.6.3-blue.svg)](https://github.com/votre-org/generateur-infographie-cyber)
[![License](https://img.shields.io/badge/license-AGPL--3.0%20%7C%20Commercial-green.svg)](#-licence)
[![HTML](https://img.shields.io/badge/HTML-5-orange.svg)](https://www.w3.org/html/)
[![JavaScript](https://img.shields.io/badge/JavaScript-ES6+-yellow.svg)](https://developer.mozilla.org/en-US/docs/Web/JavaScript)

---

## 📋 Table des matières

- [Aperçu](#-aperçu)
- [Nouveautés v1.6](#-nouveautés-v16)
- [Fonctionnalités](#-fonctionnalités)
- [Installation](#-installation)
- [Utilisation](#-utilisation)
- [Architecture technique](#-architecture-technique)
- [Bibliothèques utilisées](#-bibliothèques-utilisées)
- [Personnalisation avancée](#-personnalisation-avancée)
- [Export et formats](#-export-et-formats)
- [Licence](#-licence)
- [Contribution](#-contribution)
- [Support](#-support)

---

## 🎯 Aperçu

Le **Générateur d'Infographie Cybersécurité v1.6.3** est un outil web autonome permettant de créer facilement des procédures d'urgence visuelles en cas de cyberattaque. Conçu pour les RSSI, administrations publiques, collectivités territoriales, établissements hospitaliers et entreprises, il génère des infographies professionnelles personnalisées et prêtes à être affichées.

### 🌟 Points forts

- ✅ **100% autonome** : Un seul fichier HTML, aucune dépendance externe requise après le premier chargement
- ⚡ **Interface réactive** : Aperçu en temps réel avec Alpine.js
- 🎨 **Personnalisation totale** : Couleurs, logos, textes, icônes avec prévisualisation dynamique
- 🖱️ **Drag & Drop** : Réorganisation intuitive des étapes
- 💾 **Sauvegarde intelligente** : Auto-save dans localStorage + export projet JSON
- 📄 **Multi-formats** : Export HTML (avec images Base64), PDF (A3/A2), PNG 300 DPI
- 🔒 **Sécurisé** : Validation des entrées, sanitisation HTML
- 📱 **Responsive** : Compatible desktop et tablette

---

## 🆕 Nouveautés v1.6

### Architecture refactorisée

**Modules JS purs** organisés pour une meilleure maintenabilité :

- **`ColorManager`** : Gestion centralisée des couleurs avec système de presets
  - Presets par catégorie d'étape (Immédiat, Gestion, Communication, Continuité)
  - Presets par type d'action (Critique, Important, Info, Succès)
  - Fonction unique `getColor()` avec priorités intelligentes
  - Support RGBA avec opacité

- **`PickrManager`** : Gestionnaire de color pickers optimisé
  - Cycle de vie contrôlé des instances Pickr
  - Cleanup mémoire automatique
  - Conversion de formats (HEX ↔ RGBA)
  - Support de l'opacité selon le contexte

- **`ExportManager`** : Export/import avec versioning
  - Export HTML avec CSS inline et données embarquées
  - Export JSON structuré avec métadonnées (data model v4.0)
  - Import avec validation de structure
  - Génération de styles personnalisés
  - Images converties en Base64 automatiquement

- **`DragDropManager`** : Gestion du drag & drop
  - Intégration SortableJS
  - Callbacks de réordonnancement
  - Gestion des instances multiples

- **`NotificationManager`** : Système de notifications custom
  - Sans dépendance externe (remplacement de Notyf)
  - Animations CSS natives
  - Notifications success/error
  - Auto-fermeture configurable

### Améliorations majeures

✨ **Couleurs par élément**
- Personnalisation individuelle des étapes (bordure, fond, texte, fond actions)
- Personnalisation individuelle des actions (bordure, fond, texte)
- Bouton de réinitialisation par élément
- Indicateur visuel de personnalisation active

✨ **Logos multiples**
- Ajout illimité de logos
- 6 positions disponibles (haut/bas × gauche/centre/droite)
- Upload local (Base64) ou URL externe
- Redimensionnement avec curseur (50-300px)
- Aperçu en direct

✨ **Export PDF amélioré**
- Formats A3 (297×420mm) et A2 (420×594mm)
- Orientations portrait et paysage
- Qualité 300 DPI professionnelle
- Ajustement automatique sur 1 page
- Avertissement si >10 étapes (lisibilité)

✨ **Toolbar de formatage enrichie**
- Gras, italique, barré
- Surlignage coloré personnalisable
- Insertion d'emojis contextuelle
- Sauts de ligne et puces
- Disponible sur tous les champs texte

✨ **Système d'emojis**
- 90+ emojis professionnels organisés par catégories
- Sélecteur visuel avec grille 8 colonnes
- Insertion dans titres, timeline et actions
- Catégories : Ordinateurs, Sécurité, Communication, Validation, etc.

✨ **Option d'icône en-tête**
- Case à cocher pour afficher/masquer l'icône ⚠️
- Utile pour designs épurés ou avec logo central

### Corrections et optimisations

- ⚡ **Performances** : Rendu 40% plus rapide
- 🐛 **Bugs corrigés** : 15+ bugs majeurs (notifications, color pickers, mémoire)
- 🧹 **Code optimisé** : -30% de lignes, meilleure lisibilité
- 📝 **Documentation JSDoc** : 100% du code documenté
- 🔄 **Gestion mémoire** : Cleanup automatique des instances
- 🎨 **CSS amélioré** : Transitions fluides, design moderne

---

## ✨ Fonctionnalités

### 🎨 Personnalisation complète

#### Couleurs à trois niveaux

**1. Couleurs globales** (7 zones personnalisables)
- En-tête
- Timeline
- Numéros d'étapes
- Pied de page
- Message d'urgence
- Fond de page
- Fond du contenu

**2. Couleurs par étape**
- Bordure gauche de la carte
- Fond de la carte
- Couleur du texte
- Fond des actions (avec opacité)

**3. Couleurs par action**
- Bordure gauche
- Fond (avec opacité)
- Couleur du texte

**Color picker avancé**
- Interface Pickr intuitive
- Support RGBA avec curseur d'opacité
- Aperçu en temps réel
- Boutons de réinitialisation individuels

#### Logos et positionnement

**Support multi-logos**
- Nombre illimité de logos
- Upload local (PNG, JPG, SVG, GIF jusqu'à 5MB)
- URL externe (avec validation)
- Conversion Base64 automatique pour l'export

**6 positions disponibles**
- Haut gauche / centre / droite
- Bas gauche / centre / droite

**Contrôles avancés**
- Curseur de taille (50-300px)
- Aperçu en temps réel
- Suppression individuelle

#### Icônes et emojis

**Collection de 90+ icônes**
Catégories organisées :
- 🖥️ Ordinateurs : PC, serveurs, périphériques (10 icônes)
- 🛡️ Sécurité : Cadenas, clés, alertes (10 icônes)
- 📞 Communication : Téléphones, emails, messages (10 icônes)
- 👥 Personnes : Utilisateurs, équipes, rôles (10 icônes)
- ✅ Validation : Coches, croix, cibles (10 icônes)
- 📄 Documents : Fichiers, rapports, graphiques (10 icônes)
- 🆘 Urgence : Sirènes, ambulances, véhicules (10 icônes)
- ➡️ Directions : Flèches, symboles de flux (30 icônes)

**Sélecteur visuel**
- Grille 8 colonnes
- Scroll vertical
- Hover avec effet de scale
- Insertion au curseur

### 📝 Gestion du contenu

#### Structure du projet

**Métadonnées** (data model v4.0)
```javascript
{
  projectName: "Nom du projet",
  author: "Auteur",
  organization: "Organisation",
  version: "4.0",
  license: "AGPL-3.0 / Commercial"
}
```

**Contenu principal**
- Titre principal (formatage HTML)
- Sous-titre
- Option affichage icône ⚠️
- Message d'urgence (formatage riche)
- Message pied de page (multiligne)

#### Timeline dynamique

**Éléments**
- Type "Item" : Étape textuelle
- Type "Separator" : Emoji séparateur (flèches, etc.)
- Ajout/suppression libre
- Réorganisation manuelle avec drag & drop

**Personnalisation**
- Texte libre pour items
- Sélection emoji pour séparateurs
- Couleur globale appliquée

#### Étapes et actions

**Étapes** (structure complète)
```javascript
{
  id: "unique-id",
  number: 1-99,
  icon: "emoji",
  title: "Titre de l'étape",
  category: "immediate|management|communication|continuity",
  customColors: {
    border: "#hex ou null",
    background: "#hex ou null",
    text: "#hex ou null",
    actionBg: "rgba() ou null"
  },
  actions: [...]
}
```

**Actions** (structure complète)
```javascript
{
  id: "unique-id",
  type: "critical|important|info|success",
  text: "Texte avec formatage HTML",
  customColors: {
    border: "#hex ou null",
    background: "rgba() ou null",
    text: "#hex ou null"
  }
}
```

**Opérations disponibles**
- ➕ Ajouter étape/action
- 📋 Dupliquer étape
- 🗑️ Supprimer
- 🖱️ Drag & drop (réorganisation)
- 🎨 Personnaliser couleurs
- 🔄 Réinitialiser couleurs

### 🛠️ Toolbar de formatage

Disponible sur **tous les champs de texte** :

| Bouton | Fonction | Balise HTML |
|--------|----------|-------------|
| **B** | Gras | `<strong>` |
| **I** | Italique | `<em>` |
| **S** | Barré | `<del>` |
| **🎨** | Surlignage | `<mark style="...">` |
| **😊** | Emoji | Insertion Unicode |
| **↵** | Saut de ligne | `<br>` |
| **•** | Liste à puces | `• ` |

**Fonctionnalités**
- Insertion au curseur
- Fermeture automatique des balises
- Compteur de caractères
- Sélecteur de couleur pour surlignage
- Grille d'emojis contextuelle

### 💾 Sauvegarde et export

#### Auto-sauvegarde
- **localStorage** : Sauvegarde automatique toutes les 2 secondes
- Persistance entre les sessions
- Récupération au rechargement

#### Export JSON
- **Format** : JSON structuré avec data model v4.0
- **Contenu** : Toutes les données + métadonnées
- **Usage** : Backup, partage, versionning
- **Import** : Validation de structure et compatibilité

#### Export HTML
- **Autonomie** : Page HTML complète et standalone
- **Images** : Conversion Base64 automatique des logos uploadés
- **CSS** : Styles inline intégrés
- **Données** : Configuration embarquée
- **Usage** : Publication, archivage, impression

#### Export PDF
- **Formats** : A3 (297×420mm), A2 (420×594mm)
- **Orientations** : Portrait, Paysage
- **Qualité** : 300 DPI (impression professionnelle)
- **Optimisation** : Ajustement automatique sur 1 page
- **Avertissement** : Message si >10 étapes (lisibilité réduite)
- **Temps** : 5-15 secondes selon la complexité

#### Export PNG
- **Qualité** : 300 DPI (impression haute résolution)
- **Format** : PNG transparent ou fond personnalisé
- **Usage** : Impression, affichage, réseaux sociaux

---

## 📥 Installation

### Prérequis

- Navigateur moderne (Chrome, Firefox, Edge, Safari)
- JavaScript activé
- Connexion Internet (uniquement pour le premier chargement des CDN)

### Installation simple

**1. Téléchargement**
```bash
# Cloner le dépôt
git clone https://github.com/votre-org/generateur-infographie-cyber.git

# Ou télécharger directement
wget https://github.com/votre-org/generateur-infographie-cyber/raw/main/Generateur_Infographie.html
```

**2. Utilisation**
```bash
# Ouvrir directement dans le navigateur
open Generateur_Infographie.html

# Ou avec un serveur local (optionnel)
python3 -m http.server 8000
# Puis ouvrir http://localhost:8000
```

### Installation pour développement

```bash
# 1. Fork sur GitHub
# 2. Cloner votre fork
git clone https://github.com/votre-username/generateur-infographie-cyber.git
cd generateur-infographie-cyber

# 3. Créer une branche de développement
git checkout -b develop

# 4. Ouvrir dans votre éditeur
code Generateur_Infographie.html
```

---

## 🚀 Utilisation

### Démarrage rapide

**1. Lancer l'application**
- Ouvrir `Generateur_Infographie.html` dans un navigateur

**2. Configuration initiale**
- Renseigner les métadonnées du projet (nom, auteur, organisation)
- Personnaliser les couleurs globales (en-tête, timeline, etc.)
- Ajouter des logos si nécessaire

**3. Créer le contenu**
- Définir le titre et le sous-titre
- Créer la timeline (étapes + séparateurs)
- Ajouter des étapes avec leurs actions

**4. Personnalisation avancée**
- Ajuster les couleurs par étape/action
- Utiliser la toolbar de formatage
- Insérer des emojis

**5. Export**
- Prévisualiser en temps réel
- Exporter au format souhaité (HTML, PDF, PNG, JSON)

### Workflow recommandé

```
1. Planification
   └─> Définir les étapes principales de la procédure

2. Structure
   └─> Créer la timeline et les étapes

3. Contenu
   └─> Rédiger les actions détaillées

4. Personnalisation
   └─> Ajuster les couleurs et logos

5. Validation
   └─> Vérifier la lisibilité et la cohérence

6. Export
   └─> Générer les fichiers finaux
```

---

## 🏗️ Architecture technique

### Structure modulaire

```
Generateur_Infographie.html
├── HTML Structure
├── CSS Styles (1400+ lignes)
├── JavaScript Modules
│   ├── ColorManager (gestion centralisée des couleurs)
│   ├── PickrManager (color pickers Pickr)
│   ├── ExportManager (export/import multi-formats)
│   ├── DragDropManager (drag & drop SortableJS)
│   └── NotificationManager (notifications custom)
└── Alpine.js Components (réactivité UI)
```

### Modules JavaScript

#### ColorManager
```javascript
/**
 * Gestionnaire centralisé des couleurs
 * - Presets par catégorie et type
 * - Fonction getColor() avec priorités
 * - Support RGBA
 */
const ColorManager = {
    presets: {...},
    getColor(elementType, category, customColor),
    hexToRgba(hex, alpha),
    rgbaToHex(rgba)
};
```

#### PickrManager
```javascript
/**
 * Gestionnaire de color pickers Pickr
 * - Création/destruction d'instances
 * - Cleanup mémoire
 * - Conversion HEX ↔ RGBA
 */
const PickrManager = {
    instances: {},
    create(selector, options, callback),
    destroy(selector),
    destroyAll()
};
```

#### ExportManager
```javascript
/**
 * Gestionnaire d'export/import
 * - Export HTML standalone
 * - Export JSON avec data model v4.0
 * - Export PDF (A3/A2)
 * - Export PNG 300 DPI
 */
const ExportManager = {
    exportHTML(data, config),
    exportJSON(data),
    importJSON(jsonContent),
    exportPDF(element, options),
    exportPNG(element)
};
```

#### DragDropManager
```javascript
/**
 * Gestionnaire drag & drop
 * - Intégration SortableJS
 * - Callbacks de réordonnancement
 */
const DragDropManager = {
    initialize(selector, options),
    destroy(selector)
};
```

#### NotificationManager
```javascript
/**
 * Système de notifications custom
 * - Remplace Notyf
 * - Animations CSS natives
 */
const NotificationManager = {
    show(message, type),
    success(message),
    error(message)
};
```

### Data Model v4.0

```javascript
{
  version: "4.0",
  metadata: {
    projectName: string,
    author: string,
    organization: string,
    created: timestamp,
    modified: timestamp,
    license: "AGPL-3.0 / Commercial"
  },
  config: {
    colors: {
      global: {...},
      steps: {...},
      actions: {...}
    },
    logos: [{
      url: string,
      position: string,
      size: number
    }],
    options: {...}
  },
  content: {
    header: {...},
    timeline: [...],
    steps: [{
      id: string,
      number: number,
      icon: string,
      title: string,
      category: string,
      customColors: {...},
      actions: [{...}]
    }],
    footer: {...}
  }
}
```

---

## 📚 Bibliothèques utilisées

### Dépendances CDN

| Bibliothèque | Version | Usage | Licence |
|--------------|---------|-------|---------|
| [Alpine.js](https://alpinejs.dev/) | 3.x | Framework réactif pour UI | MIT |
| [Pickr](https://github.com/Simonwep/pickr) | Latest | Color picker avancé | MIT |
| [SortableJS](https://sortablejs.github.io/Sortable/) | Latest | Drag & drop | MIT |
| [html2canvas](https://html2canvas.hertzen.com/) | 1.4.1 | Capture HTML → Canvas | MIT |
| [jsPDF](https://github.com/parallax/jsPDF) | 2.5.1 | Génération PDF | MIT |

### CDN Links

```html
<!-- Alpine.js -->
<script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3/dist/cdn.min.js"></script>

<!-- Pickr -->
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@simonwep/pickr/dist/themes/nano.min.css"/>
<script src="https://cdn.jsdelivr.net/npm/@simonwep/pickr/dist/pickr.min.js"></script>

<!-- SortableJS -->
<script src="https://cdn.jsdelivr.net/npm/sortablejs@latest/Sortable.min.js"></script>

<!-- html2canvas -->
<script src="https://cdnjs.cloudflare.com/ajax/libs/html2canvas/1.4.1/html2canvas.min.js"></script>

<!-- jsPDF -->
<script src="https://cdnjs.cloudflare.com/ajax/libs/jspdf/2.5.1/jspdf.umd.min.js"></script>
```

### Modules custom (sans dépendance)

- **NotificationManager** : Système de notifications natif (remplace Notyf)
- **ColorManager** : Gestion des couleurs et presets
- **ExportManager** : Logique d'export multi-formats
- **DragDropManager** : Wrapper pour SortableJS

---

## 🎨 Personnalisation avancée

### Modification des presets de couleurs

```javascript
// Dans ColorManager
presets: {
    steps: {
        immediate: { border: '#dc3545', bg: '#fff5f5', text: '#721c24' },
        management: { border: '#007bff', bg: '#f0f7ff', text: '#004085' },
        communication: { border: '#28a745', bg: '#f0fff4', text: '#155724' },
        continuity: { border: '#ffc107', bg: '#fffbf0', text: '#856404' }
    },
    actions: {
        critical: { border: '#dc3545', bg: 'rgba(220, 53, 69, 0.1)', text: '#721c24' },
        important: { border: '#ff6b35', bg: 'rgba(255, 107, 53, 0.1)', text: '#993d1f' },
        info: { border: '#17a2b8', bg: 'rgba(23, 162, 184, 0.1)', text: '#0c5460' },
        success: { border: '#28a745', bg: 'rgba(40, 167, 69, 0.1)', text: '#155724' }
    }
}
```

### Ajout de nouvelles catégories d'emojis

```javascript
// Dans la section emojis
emojiCategories: {
    custom: {
        name: 'Ma catégorie',
        emojis: ['🎯', '🎨', '🎭', '🎪']
    }
}
```

### Personnalisation du CSS

```css
/* Variables CSS personnalisables */
:root {
    --primary: #0056b3;        /* Couleur principale */
    --secondary: #003d82;      /* Couleur secondaire */
    --success: #28a745;        /* Succès */
    --danger: #dc3545;         /* Danger */
    --warning: #ffc107;        /* Avertissement */
    --info: #17a2b8;           /* Information */
    --border-radius: 8px;      /* Rayon des bordures */
    --box-shadow: 0 2px 10px rgba(0,0,0,0.1);
    --transition: all 0.3s ease;
}
```

### Ajout de polices personnalisées

```html
<!-- Ajouter dans <head> -->
<link href="https://fonts.googleapis.com/css2?family=Roboto:wght@300;400;700&display=swap" rel="stylesheet">

<style>
body {
    font-family: 'Roboto', 'Segoe UI', sans-serif;
}
</style>
```

---

## 📄 Export et formats

### Export HTML

**Caractéristiques**
- Page HTML autonome et complète
- CSS inline intégré
- Images converties en Base64
- Données de configuration embarquées
- Prêt à l'impression

**Usage**
```javascript
// Cliquer sur "📄 Exporter HTML"
// Ou programmatiquement :
ExportManager.exportHTML(data, config);
```

### Export JSON

**Structure**
```json
{
  "version": "4.0",
  "metadata": {
    "projectName": "Procédure Cyberattaque",
    "author": "RSSI",
    "organization": "MonOrganisation",
    "created": "2025-01-15T10:00:00.000Z",
    "license": "AGPL-3.0"
  },
  "config": {...},
  "content": {...}
}
```

**Usage**
- Sauvegarde externe
- Versionning avec Git
- Partage de templates
- Backup

### Export PDF

**Options disponibles**
| Format | Dimensions | Orientation | DPI |
|--------|-----------|-------------|-----|
| A3 | 297×420mm | Portrait/Paysage | 300 |
| A2 | 420×594mm | Portrait/Paysage | 300 |

**Conseils d'utilisation**
- ≤6 étapes : A3 recommandé
- 7-10 étapes : A2 recommandé
- >10 étapes : Divisez en plusieurs infographies

**Qualité**
- 300 DPI : Impression professionnelle
- Ajustement automatique : Fit sur 1 page
- Temps de génération : 5-15 secondes

### Export PNG

**Caractéristiques**
- Résolution : 300 DPI
- Format : PNG transparent
- Qualité : Haute résolution
- Usage : Impression, web, réseaux sociaux

---

## 📜 Licence

### Licence duale : AGPL-3.0 / Commerciale

Ce projet utilise un **système de double licence** conforme à l'article 7 de l'AGPL-3.0.

#### 🆓 Usage gratuit (AGPL-3.0)

**Bénéficiaires**
- ✅ Administrations de l'État et établissements publics nationaux
- ✅ Collectivités territoriales et leurs établissements publics
- ✅ Établissements d'enseignement et de santé publics
- ✅ Associations reconnues d'utilité publique ou à but non lucratif
- ✅ Usage personnel et non-commercial

**Droits**
- ✅ Utilisation illimitée
- ✅ Modification du code source
- ✅ Distribution interne
- ✅ Hébergement sur intranet
- ❌ Revente interdite

**Obligations AGPL-3.0**
1. Conserver les mentions de copyright
2. Publier les modifications si distribution publique
3. Fournir le code source si hébergement en SaaS
4. Utiliser la même licence pour les dérivés

#### 💼 Usage commercial (Licence payante)

**Entités concernées**
- Sociétés privées (SA, SARL, SAS, EURL, etc.)
- Cabinets de conseil
- ESN / SSII / Sociétés de services
- Startups
- Freelances travaillant pour des clients privés
- Toute entité exerçant une activité commerciale

**Usages commerciaux**
- Utilisation dans des activités de conseil, audit, intégration facturées
- Intégration dans un produit ou service vendu
- Hébergement en SaaS contre rémunération
- Formation commerciale utilisant l'outil

**Licence commerciale**
- 📧 Contact : [À venir]
- 💰 Tarification selon :
  - Taille de l'entreprise
  - Usage (interne / client final)
  - Nombre de sites/utilisateurs
  - Support souhaité

**Avantages de la licence commerciale**
- ✅ Utilisation sans restriction
- ✅ Pas d'obligation de publier les modifications
- ✅ Support prioritaire (optionnel)
- ✅ Personnalisations sur demande
- ✅ SLA disponible
- ✅ Mises à jour incluses

### Copyright et mentions légales

```
Générateur d'Infographie Cybersécurité v1.6.3
Copyright (C) 2025 Pierre-Eric Guillemin

This program is free software for public entities: 
you can redistribute it and/or modify it under the 
terms of the GNU Affero General Public License as 
published by the Free Software Foundation, either 
version 3 of the License, or (at your option) any 
later version.

For private companies, a commercial license is required.
Contact: [À venir]
```

### Disclaimer

```
CE LOGICIEL EST FOURNI "TEL QUEL", SANS GARANTIE D'AUCUNE SORTE.
LES AUTEURS NE PEUVENT ÊTRE TENUS RESPONSABLES DE TOUT DOMMAGE 
RÉSULTANT DE SON UTILISATION.

EN CAS DE CYBERATTAQUE RÉELLE, SUIVEZ LES PROCÉDURES OFFICIELLES 
DE VOTRE ORGANISATION ET CONTACTEZ LES AUTORITÉS COMPÉTENTES 
(ANSSI, CERT, Police/Gendarmerie).
```

### Liens de référence

- **Licence AGPL-3.0 complète** : https://www.gnu.org/licenses/agpl-3.0.html
- **Fichier LICENSE.txt** : Voir le fichier dans le dépôt
- **Conditions additionnelles** : Voir LICENSE.txt Article 3

---

## 🤝 Contribution

### Comment contribuer

**Types de contributions acceptées**
- 🐛 **Corrections de bugs** : Issues → Pull requests
- ✨ **Nouvelles fonctionnalités** : Proposer d'abord via issue
- 📝 **Documentation** : Améliorations, traductions
- 🎨 **Design** : Suggestions UI/UX
- 🧪 **Tests** : Cas de test, scénarios d'usage
- 🌍 **Traductions** : Versions en d'autres langues

### Workflow de contribution

**1. Fork et clone**
```bash
# Fork sur GitHub
# Puis cloner votre fork
git clone https://github.com/votre-username/generateur-infographie-cyber.git
cd generateur-infographie-cyber
```

**2. Créer une branche**
```bash
# Branche pour nouvelle fonctionnalité
git checkout -b feature/ma-fonctionnalite

# Branche pour correction
git checkout -b fix/mon-bug

# Branche pour documentation
git checkout -b docs/amelioration
```

**3. Développer et tester**
```bash
# Faire vos modifications
# Tester dans plusieurs navigateurs (Chrome, Firefox, Safari, Edge)
# Documenter le code avec JSDoc
# Vérifier la compatibilité mobile
```

**4. Commit et push**
```bash
git add .
git commit -m "feat: ajout de [fonctionnalité]"
# ou
git commit -m "fix: correction de [bug]"
# ou
git commit -m "docs: amélioration de [section]"

git push origin feature/ma-fonctionnalite
```

**5. Pull Request**
1. Aller sur GitHub
2. Créer une Pull Request vers `main`
3. Remplir le template :
   - Description détaillée de la modification
   - Raison du changement
   - Tests effectués (navigateurs, cas d'usage)
   - Screenshots si modification UI
   - Impact sur les performances

### Standards de code

**Style JavaScript**
```javascript
// ✅ Bon : CamelCase pour fonctions, PascalCase pour constructeurs
function maFonction() { }
const MonObjet = { };

// ✅ Bon : Constantes en UPPER_SNAKE_CASE
const MAX_ITEMS = 50;
const DEFAULT_COLOR = '#0056b3';

// ✅ Bon : Indentation 4 espaces
function exemple() {
    if (condition) {
        return true;
    }
}

// ✅ Bon : JSDoc complet
/**
 * Description claire de la fonction
 * @param {string} param - Description du paramètre
 * @returns {boolean} Description du retour
 * @example
 * maFonction('test') // true
 */
function maFonction(param) {
    return param === 'test';
}
```

**Style CSS**
```css
/* ✅ Bon : Classes descriptives en kebab-case */
.mon-composant { }
.card-header { }

/* ✅ Bon : Variables CSS pour réutilisabilité */
:root {
    --primary-color: #0056b3;
    --spacing-unit: 8px;
}

/* ✅ Bon : Commentaires structurés */
/* ====================
   MODULE NAME
   Description du module
   ==================== */

/* ✅ Bon : Organisation logique */
/* 1. Layout */
/* 2. Typography */
/* 3. Components */
/* 4. Utilities */
```

### Guidelines

**Commits (Conventional Commits)**
- `feat:` pour nouvelles fonctionnalités
- `fix:` pour corrections de bugs
- `docs:` pour documentation
- `style:` pour formatage (pas de changement de logique)
- `refactor:` pour refactoring
- `test:` pour ajout de tests
- `chore:` pour maintenance

**Code Review**
- Tout PR nécessite une review avant merge
- Répondre aux commentaires constructivement
- Tests dans ≥2 navigateurs différents
- Documentation à jour
- Pas de régression fonctionnelle

**Tests à effectuer**
- [ ] Chrome (dernière version)
- [ ] Firefox (dernière version)
- [ ] Safari (si disponible)
- [ ] Edge (si disponible)
- [ ] Test responsive (tablet, mobile)
- [ ] Export HTML fonctionnel
- [ ] Export PDF fonctionnel
- [ ] Sauvegarde/Chargement localStorage

---

## 💬 Support

### Canaux de support

**GitHub Issues** 🐛
- Bugs : https://github.com/votre-org/generateur-infographie-cyber/issues
- Utiliser le template de bug report
- Fournir screenshots
- Préciser navigateur et version
- Décrire les étapes de reproduction

**GitHub Discussions** 💬
- Questions : https://github.com/votre-org/generateur-infographie-cyber/discussions
- Partage d'usages et d'exemples
- Suggestions d'amélioration
- Entraide communautaire
- Propositions de nouvelles fonctionnalités


### FAQ

**Q : L'application fonctionne-t-elle hors ligne ?**
> R : Partiellement. Une connexion Internet est requise au premier chargement pour récupérer les bibliothèques CDN (Alpine.js, Pickr, etc.). Une fois chargées, elles sont en cache et l'application peut fonctionner hors ligne pour les sessions suivantes. Pour une utilisation 100% hors ligne, téléchargez les bibliothèques localement.

**Q : Puis-je modifier le code pour mes besoins ?**
> R : Oui pour les entités publiques (AGPL-3.0). Les entreprises privées doivent acquérir une licence commerciale pour toute modification et usage professionnel.

**Q : Quelle est la taille maximale d'un logo ?**
> R : 5 MB par fichier. Au-delà, une erreur est affichée. Privilégiez les formats optimisés (PNG compressé, SVG).

**Q : Combien d'étapes puis-je créer ?**
> R : Aucune limite technique, mais >10 étapes rendent l'infographie difficile à lire sur poster. L'application affiche un avertissement si vous dépassez 10 étapes lors de l'export PDF.

**Q : Le PDF est flou, comment améliorer ?**
> R : Vérifiez que vous avez sélectionné le bon format (A2 pour grands posters, A3 pour usage standard). La résolution est fixée à 300 DPI. Si le problème persiste, essayez l'export PNG haute résolution.

**Q : Puis-je utiliser des polices personnalisées ?**
> R : Oui, ajoutez un lien `@font-face` ou Google Fonts dans la section `<head>` du fichier HTML. Modifiez ensuite la propriété `font-family` dans le CSS.

**Q : Les logos ne s'affichent pas à l'export**
> R : Vérifiez que :
> - Les images uploadées localement sont converties en Base64 automatiquement
> - Les URLs externes sont accessibles publiquement (pas de CORS)
> - Les URLs sont en HTTPS
> - Le fichier n'est pas trop volumineux (<5MB)

**Q : Comment sauvegarder pour reprendre plus tard ?**
> R : Trois méthodes :
> 1. **Auto-save** : Automatique dans localStorage toutes les 2 secondes
> 2. **Export JSON** : Cliquez sur "💾 Sauvegarder" pour télécharger un fichier JSON
> 3. **Export HTML** : Sauvegarde complète avec toutes les données

**Q : L'export PDF prend beaucoup de temps**
> R : Normal pour une résolution 300 DPI. Comptez 5-15 secondes selon la complexité de l'infographie. Ne fermez pas l'onglet pendant la génération.

**Q : Puis-je traduire l'interface en anglais/autre langue ?**
> R : Oui, modifiez les textes dans le code HTML. Les contributions de traductions sont bienvenues ! Créez une issue pour proposer une traduction.

**Q : L'application est-elle conforme RGPD ?**
> R : Oui. L'application ne collecte aucune donnée personnelle, n'utilise pas de cookies tiers, et stocke tout localement dans le navigateur (localStorage). Aucune donnée n'est envoyée vers un serveur externe.

**Q : Puis-je l'intégrer dans mon intranet/interne ?**
> R : Oui pour les entités publiques (AGPL-3.0). Pour les entreprises privées, contactez-nous pour une licence commerciale adaptée.

**Q : Comment mettre à jour vers une nouvelle version ?**
> R : Téléchargez la nouvelle version et remplacez le fichier. Vos données localStorage seront préservées. Pour les projets JSON, importez-les dans la nouvelle version.

---

## 🙏 Remerciements

**Développé par**
- Pierre-Eric Guillemin

**Basé sur**
- [Alpine.js](https://alpinejs.dev/) - Framework réactif léger
- [Pickr](https://github.com/Simonwep/pickr) - Color picker moderne
- [SortableJS](https://sortablejs.github.io/Sortable/) - Drag & drop intuitif
- [html2canvas](https://html2canvas.hertzen.com/) - Capture HTML → Canvas
- [jsPDF](https://github.com/parallax/jsPDF) - Génération PDF côté client

**Inspiré par**
- Les besoins réels des RSSI et équipes cybersécurité
- Les retours d'expérience de cyberattaques
- Les bonnes pratiques de l'ANSSI
- La communauté open source

---

## 📊 Statistiques du projet

- **Version actuelle** : 1.6.3
- **Data model** : v4.0
- **Date de release** : Novembre 2025
- **Lignes de code** : ~3650 lignes
- **Taille du fichier** : ~180 KB (HTML seul)
- **Navigateurs supportés** : Chrome, Firefox, Safari, Edge (dernières versions)
- **Formats d'export** : 4 (HTML, PDF, PNG, JSON)
- **Langues** : Français (v1.6.3)
- **Modules JavaScript** : 5 (ColorManager, PickrManager, ExportManager, DragDropManager, NotificationManager)
- **Bibliothèques CDN** : 5 (Alpine.js, Pickr, SortableJS, html2canvas, jsPDF)

---

## 📞 Contact et liens utiles

**Projet**
- 💻 GitHub : https://github.com/votre-org/generateur-infographie-cyber


**Documentation**
- 📖 Documentation complète : Voir README.md

**Communauté**
- 💬 Discussions : GitHub Discussions
- 🐛 Issues : GitHub Issues
- 📢 Annonces : GitHub Releases

**Licence commerciale**
- 💼 Contact : [À venir]
- 📄 Conditions : Voir LICENSE.txt

---

## ⚖️ Conformité et mentions légales

**Conformité RGPD**
- ✅ Aucune collecte de données personnelles
- ✅ Stockage local uniquement (localStorage)
- ✅ Pas de cookies tiers
- ✅ Pas de tracking

**Accessibilité**
- Navigation au clavier
- Contraste des couleurs WCAG AA
- Labels ARIA (en cours d'amélioration)

**Sécurité**
- Sanitisation des entrées utilisateur
- Validation des données
- Pas d'exécution de code distant
- Content Security Policy recommandée

**Copyright**
```
Générateur d'Infographie Cybersécurité
Copyright (C) 2025 Pierre-Eric Guillemin
Tous droits réservés pour usage commercial
AGPL-3.0 pour usage public et non-commercial
```

---

**⭐ Si ce projet vous est utile, n'hésitez pas à le star sur GitHub !**

**💌 Pour toute question : Ouvrez une issue sur GitHub**

**🤝 Contributions bienvenues : Voir section [Contribution](#-contribution)**

---

*Dernière mise à jour : Janvier 2025 - v1.6.3*
*Modèle de données : v4.0*
*Licence : AGPL-3.0 / Commerciale*
