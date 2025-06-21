🛡️ Générateur d'Infographie Cybersécurité v1.1

    🚨 Outil open source pour créer des procédures d'urgence en cas de cyberattaque - Personnalisable, responsive et prêt à l'emploi

📋 Table des matières

    Aperçu
    Fonctionnalités
    Installation
    Utilisation
    Personnalisation
    Export et impression
    Licence
    Contribution
    Support

🎯 Aperçu

Le Générateur d'Infographie Cybersécurité est un outil web permettant de créer facilement des procédures d'urgence visuelles en cas de cyberattaque. Conçu pour les administrations publiques, collectivités territoriales et entreprises, il génère des infographies personnalisées et prêtes à être affichées.
Captures d'écran

Interface de création avec aperçu en temps réel

Exemple d'infographie générée
✨ Fonctionnalités
🎨 Personnalisation complète

    Thèmes prédéfinis : 4 palettes de couleurs professionnelles
    Couleurs personnalisables : En-tête, timeline, actions, arrière-plans
    Logo : Support des formats PNG, JPG, SVG (upload ou URL)
    Classifications : 4 niveaux de priorité personnalisables
    CSS avancé : Préfixes, isolation, styles personnalisés

📝 Gestion du contenu

    Timeline dynamique : Séquence d'étapes personnalisable
    Étapes illimitées : Ajout, suppression, duplication
    Actions formatées : Gras, italique, listes, sauts de ligne
    Images intégrées : Support d'images dans les actions
    70+ icônes : Émojis professionnels pré-sélectionnés

💾 Sauvegarde et export

    Sauvegarde automatique : Dans le navigateur (localStorage)
    Export projet : Format JSON pour partage/backup
    Export HTML : Page autonome avec tout intégré
    Export PDF : Multiples formats (A3, A4, optimisé)
    Export images : PNG 300 DPI, JPEG, WebP, SVG
    HTML impression : Optimisé pour Ctrl+P

🔒 Sécurité

    Validation des entrées : Protection XSS
    Sanitisation HTML : Balises autorisées uniquement
    Validation des URL : HTTPS requis pour les ressources externes
    Limites de sécurité : Taille des fichiers, nombre d'éléments
    CSP headers : Content Security Policy intégrée

🚀 Installation
Option 1 : Utilisation directe (recommandé)

    Téléchargez le fichier Generateur_Infographie-v1.3.html
    Ouvrez-le dans un navigateur moderne (Chrome, Firefox, Edge, Safari)
    C'est prêt ! Aucune installation requise

Option 2 : Hébergement web

bash

# Cloner le repository
git clone https://github.com/[votre-repo]/generateur-infographie-cyber.git

# Accéder au dossier
cd generateur-infographie-cyber

# Servir localement (exemple avec Python)
python -m http.server 8000

# Ouvrir dans le navigateur
# http://localhost:8000/Generateur_Infographie-v1.3.html

Option 3 : Intégration dans un site

html

<!-- Intégrer dans une iframe -->
<iframe src="Generateur_Infographie-v1.3.html" 
        width="100%" 
        height="800px" 
        frameborder="0">
</iframe>

📖 Utilisation
Démarrage rapide

    Personnalisez les couleurs : Utilisez les thèmes ou créez votre palette
    Ajoutez votre logo : Glissez-déposez ou indiquez une URL
    Modifiez le contenu : Titre, sous-titre, message d'urgence
    Configurez les étapes :
        Cliquez sur l'icône pour la changer
        Éditez le titre et la catégorie
        Ajoutez/modifiez les actions
    Prévisualisez : L'aperçu se met à jour en temps réel
    Exportez : HTML pour utilisation web ou PDF pour impression

Guide détaillé des fonctionnalités
🎨 Personnalisation des couleurs

javascript

// Structure des couleurs personnalisables
{
  header: "#dc3545",      // En-tête principal
  timeline: "#3498db",    // Éléments de timeline
  stepNumber: "#0056b3",  // Numéros des étapes
  footer: "#2c3e50",      // Pied de page
  emergency: "#ff6b6b",   // Message d'urgence
  bodyBg: "#f8f9fa",      // Arrière-plan général
  contentBg: "#ffffff",   // Contenu principal
  stepBg: "#ffffff"       // Cartes d'étapes
}

📝 Format des actions

html

<!-- Formatage supporté dans les actions -->
<strong>Texte en gras</strong>
<em>Texte en italique</em>
<u>Texte souligné</u>
<br> <!-- Saut de ligne -->
- Liste à puces

🖼️ Gestion des logos

    Formats supportés : PNG, JPG, SVG, GIF
    Taille max recommandée : 500 KB
    Positions disponibles : 6 emplacements (haut/bas × gauche/centre/droite)
    Redimensionnement : 50-300px

💾 Structure du projet JSON

json

{
  "version": "2.0",
  "timestamp": "2025-01-21T10:00:00.000Z",
  "projectName": "Mon infographie cyberattaque",
  "settings": {
    "mainTitle": "CYBERATTAQUE DÉTECTÉE",
    "mainSubtitle": "Procédure d'urgence",
    "colors": { ... }
  },
  "timeline": [ ... ],
  "steps": [ ... ],
  "logo": { ... },
  "classifications": { ... },
  "cssCustomization": { ... }
}

🎯 Personnalisation avancée
CSS personnalisé

Le générateur permet d'ajouter du CSS personnalisé avec :

    Préfixes de classes : Évite les conflits CSS
    Conteneur d'isolation : Limite la portée des styles
    Reset CSS optionnel : Normalise les styles de base

css

/* Exemple de CSS personnalisé */
.cyb-container {
  font-family: 'Ma Police Corporate', sans-serif;
}
.cyb-header {
  background-image: url('pattern.png');
}

Classifications personnalisées

Adaptez les niveaux de priorité à votre organisation :

javascript

// Exemples de personnalisation
P1/P2/P3/P4          // Numérique
Urgent/Normal/Info   // Textuel
Rouge/Orange/Vert    // Par couleur

📤 Export et impression
Formats d'export disponibles

Format	Usage	Caractéristiques
HTML	Web/Intranet	Autonome, responsive, interactif
PDF Optimisé	Impression A3/A4	Sans espaces blancs, 90%+ de la page
PNG 300 DPI	Impression pro	Haute résolution, qualité maximale
JPEG	Email/Web	Compression optimisée, taille réduite
WebP	Web moderne	Meilleure compression que JPEG
SVG	Vectoriel	Redimensionnable sans perte

Conseils d'impression

    Format recommandé : A3 portrait pour une lisibilité optimale
    Marges : Minimales (5-10mm)
    Couleurs : Activer l'impression des couleurs de fond
    Échelle : Ajuster à la page

⚖️ Licence

Ce projet est sous double licence :
🆓 Licence AGPL-3.0 (Gratuite)

    ✅ Usage personnel : Libre et gratuit
    ✅ Administrations publiques : Ministères, préfectures
    ✅ Collectivités : Mairies, départements, régions
    ✅ Établissements publics : Hôpitaux, écoles, universités
    ✅ Associations : À but non lucratif

Obligations : Publier les modifications, conserver les mentions de copyright
💼 Licence commerciale (Payante)

Requise pour :

    ❌ Entreprises privées (SARL, SAS, SA)
    ❌ Consultants/Freelances (usage client)
    ❌ Intégration dans produits commerciaux
    ❌ Services payants

Contact : [aaaa@bbb.ccc]
🤝 Contribution

Les contributions sont bienvenues ! Voici comment participer :

    Demander a être intégré au projet
    Fork le projet
    Créez votre branche (git checkout -b feature/AmazingFeature)
    Committez vos changements (git commit -m 'Add AmazingFeature')
    Push vers la branche (git push origin feature/AmazingFeature)
    Ouvrez une Pull Request

Guidelines de contribution

    Respecter la structure du code existant
    Commenter les nouvelles fonctionnalités
    Tester sur plusieurs navigateurs
    Maintenir la compatibilité avec les anciens projets

🐛 Bugs connus et limitations

    Safari : Certaines animations CSS peuvent être saccadées
    IE11 : Non supporté (utiliser Edge)
    Mobile : L'édition est optimisée pour desktop
    Taille max : 50 étapes, 20 actions par étape
    Images : Max 2MB par image d'action

📞 Support
Documentation

    Guide utilisateur complet
    FAQ
    Exemples de procédures

Obtenir de l'aide

    📧 Email : [aaaa@bbb.ccc]
    🐛 Issues : GitHub Issues
    💬 Discussions : GitHub Discussions

🙏 Remerciements

    Icônes : Emojis natifs Unicode
    Inspiration : ANSSI, Cybermalveillance.gouv.fr
    Bibliothèques : html2canvas, jsPDF

📈 Roadmap
Version 2.1 (prévue)

    Pas de roadmap précise pour le moment, à venir.

<div align="center">

Fait avec ❤️ pour la cybersécurité publique

"Chaque minute compte en cas de cyberattaque"

Site web • Documentation • Contact
</div>
