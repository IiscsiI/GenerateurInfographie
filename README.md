# 🛡️ Générateur d'Infographie Cybersécurité v1.5

> 🚨 Outil open source pour créer des procédures d'urgence en cas de cyberattaque - Personnalisable, responsive et prêt à l'emploi

## 📋 Table des matières

- [Aperçu](#-aperçu)
- [Fonctionnalités](#-fonctionnalités)
- [Installation](#-installation)
- [Utilisation](#-utilisation)
- [Personnalisation](#-personnalisation)
- [Export et impression](#-export-et-impression)
- [Licence](#-licence)
- [Contribution](#-contribution)
- [Support](#-support)

## 🎯 Aperçu

Le Générateur d'Infographie Cybersécurité est un outil web **tout-en-un** permettant de créer facilement des procédures d'urgence visuelles en cas de cyberattaque. Conçu pour les administrations publiques, collectivités territoriales et entreprises, il génère des infographies personnalisées et prêtes à être affichées.

### Captures d'écran

![Interface de création](screenshots/interface.png)
*Interface de création avec aperçu en temps réel*

![Exemple d'infographie](screenshots/example.png)
*Exemple d'infographie générée*

## ✨ Fonctionnalités

### 🎨 Personnalisation complète

- **Thèmes prédéfinis** : 4 palettes de couleurs professionnelles (Bleu Corporate, Rouge Urgence, Vert Sécurité, Sombre Professionnel)
- **Couleurs personnalisables** : 9 zones de couleurs (en-tête, timeline, étapes, arrière-plans, etc.)
- **Logo** : Support PNG, JPG, SVG, GIF (upload local ou URL)
- **Classifications** : 4 niveaux de priorité avec noms et couleurs personnalisables
- **CSS avancé** : Préfixes de classes, conteneur d'isolation, CSS personnalisé

### 📝 Gestion du contenu

- **Timeline dynamique** : Séquence d'étapes avec **150+ séparateurs** disponibles
- **Étapes flexibles** : Ajout, suppression, duplication (max 50 étapes)
- **Actions formatées** : Gras, italique, souligné, listes, sauts de ligne
- **Images dans les actions** : Support d'images avec positionnement et redimensionnement
- **400+ icônes** : Collection exhaustive d'émojis professionnels
- **Sélecteur d'icônes** : Interface visuelle avec catégories

### 💾 Sauvegarde et export

- **Sauvegarde automatique** : Dans le navigateur (localStorage)
- **Export projet** : Format JSON pour partage/backup
- **Export HTML** : Page autonome avec CSS et images intégrées
- **Export PDF** : Multiple formats (A3, A4, optimisé sans marges)
- **Export images** : PNG 300 DPI, JPEG, WebP, SVG
- **HTML impression** : Version optimisée pour impression navigateur

### 🔒 Sécurité

- **Validation des entrées** : Protection contre XSS
- **Sanitisation HTML** : Seules les balises sûres sont autorisées
- **Validation des URL** : Vérification des URLs externes
- **Limites de sécurité** : Taille des fichiers (5MB logos, 2MB images), nombre d'éléments
- **CSP headers** : Content Security Policy intégrée dans l'export

## 🚀 Installation

### Option 1 : Utilisation directe (recommandé)

1. Téléchargez le fichier `Generateur_Infographie-v1.5.html`
2. Ouvrez-le dans un navigateur moderne (Chrome, Firefox, Edge, Safari)
3. C'est prêt ! Aucune installation requise

### Option 2 : Hébergement web

```bash
# Cloner le repository
git clone https://github.com/[votre-repo]/generateur-infographie-cyber.git

# Accéder au dossier
cd generateur-infographie-cyber

# Servir localement (exemple avec Python)
python -m http.server 8000

# Ouvrir dans le navigateur
# http://localhost:8000/Generateur_Infographie-v1.5.html
