# 🛡️ Générateur d'Infographie Cybersécurité

**Créez des fiches réflexe de réponse à incident professionnelles, en quelques minutes, 100 % en local.**

Outil destiné aux RSSI, équipes SSI et référents cybersécurité des administrations, collectivités, hôpitaux et associations : composez visuellement une affiche de procédure d'urgence (étapes, actions, timeline, logos, couleurs) et exportez-la en **PNG**, **PDF vectoriel** ou **HTML autonome**.

> **Version 2.2** — Application native (binaire unique, sans installation), éditeur WYSIWYG complet.
>
> Licence : **utilisation libre non commerciale** (particuliers, secteur public, hospitalier, associations). Voir [Licence](#-licence).

---

## 📋 Table des matières

- [Pourquoi cet outil](#-pourquoi-cet-outil)
- [Fonctionnalités](#-fonctionnalités)
- [Installation](#-installation)
- [Utilisation](#-utilisation)
- [Exports](#-exports)
- [Données et confidentialité](#-données-et-confidentialité)
- [Dépannage](#-dépannage)
- [Compatibilité des fichiers](#-compatibilité-des-fichiers)
- [Architecture (contributeurs)](#-architecture-contributeurs)
- [Licence](#-licence)
- [Historique des versions](#-historique-des-versions)

---

## 🎯 Pourquoi cet outil

En cas de cyberattaque, chaque minute compte : les bons réflexes doivent être affichés, lisibles, à jour. Les outils bureautiques classiques produisent des documents vite obsolètes et pénibles à maintenir ; les solutions en ligne posent un problème évident de confidentialité pour des procédures de sécurité.

Ce générateur répond aux deux exigences :

- **Rapidité** : une fiche réflexe complète et présentable en quelques minutes, modifiable en quelques secondes.
- **Souveraineté** : tout s'exécute sur votre poste. Aucun compte, aucun cloud, aucune donnée sortante — voir [Données et confidentialité](#-données-et-confidentialité).

## ✨ Fonctionnalités

### Édition WYSIWYG directe sur l'affiche

- **Cliquez un texte pour l'éditer en place** (titre, sous-titre, étapes, actions, timeline, messages) — barre de mise en forme flottante (gras, italique, souligné, barré), Entrée = saut de ligne, Échap = terminer.
- **Panneau contextuel** à la sélection de chaque élément : catégorie et couleurs d'une étape, type d'une action, symbole de timeline, couleur d'une zone… Le panneau se déplace à la souris.
- **Logos en placement libre** : glissez un logo n'importe où sur l'affiche (guides magnétiques sur les centres de page), redimensionnez-le à la poignée. Les six ancrages classiques (coins et centres d'en-tête/pied) restent disponibles.
- **Réordonnancement au glisser** : tirez la pastille numérotée d'une étape sur une autre pour permuter ; idem pour la timeline. Défilement automatique aux bords pour les affiches longues.
- **Zoom de prévisualisation** : Ajuster / 75 % / 100 %.
- **La prévisualisation est exacte** : la mise en page à l'écran est strictement identique, au pixel, à celle des exports.

### Personnalisation

- Thèmes de couleurs prédéfinis + couleur personnalisée par zone (en-tête, timeline, urgence, pied, fonds) et par étape.
- Icônes emoji par étape, timeline libre (éléments et séparateurs).
- Logos par fichier (intégrés au projet) ou par URL.
- Formats A4 / A3 / A2, portrait ou paysage.

### Fiabilité

- **Annuler / Rétablir** (Ctrl+Z / Ctrl+Y) couvrant toutes les modifications, y compris l'édition directe.
- **Sauvegarde automatique** continue ; projets multiples.
- Import / export **JSON** des projets ; ré-import possible depuis un export HTML autonome.
- Fermeture propre : l'application s'éteint d'elle-même ~15 s après la fermeture de la fenêtre.

## 📦 Installation

### Option A — Binaire (recommandé)

1. Téléchargez le binaire de la dernière version depuis l'onglet **[Releases](https://github.com/IiscsiI/GenerateurInfographie/releases)** du dépôt.
2. **Vérifiez l'intégrité** du fichier téléchargé — l'empreinte SHA-256 attendue est publiée avec chaque release :

   ```powershell
   Get-FileHash .\infographic-generator.exe -Algorithm SHA256
   ```

3. Lancez `infographic-generator.exe`. C'est tout : aucune installation, aucun droit administrateur requis.

Prérequis : Windows 10/11 avec **Microsoft Edge ou Google Chrome** installé (présent par défaut sur Windows). Linux et macOS : compiler depuis les sources (option B).

> ⚠️ Le binaire n'étant pas signé numériquement, Windows SmartScreen peut afficher un avertissement au premier lancement (« Informations complémentaires » → « Exécuter quand même »). La vérification SHA-256 ci-dessus est votre garantie d'intégrité.

### Option B — Compilation depuis les sources

Prérequis : [Go](https://go.dev/dl/) 1.23 ou supérieur.

```powershell
# Windows
git clone https://github.com/IiscsiI/GenerateurInfographie.git
cd GenerateurInfographie
.\build.ps1
```

```bash
# Linux / macOS
git clone https://github.com/IiscsiI/GenerateurInfographie.git
cd GenerateurInfographie
./build.sh
```

Le script télécharge les bibliothèques JavaScript embarquées (versions épinglées), compile le binaire et affiche son empreinte SHA-256. Options utiles : `-Console` (garder la console pour le débogage), `-ForceVendor` (re-télécharger les bibliothèques), `-Target linux` (compilation croisée).

## 🚀 Utilisation

Double-cliquez le binaire : l'application démarre un serveur strictement local (`127.0.0.1`, port aléatoire) et ouvre l'éditeur dans une fenêtre dédiée du navigateur. Fermez la fenêtre pour quitter — l'application s'arrête seule après ~15 secondes.

### Options de lancement

| Option | Effet |
|---|---|
| `-port 8080` | Port fixe au lieu d'un port aléatoire |
| `-data C:\chemin` | Répertoire des données (défaut : dossier `data` à côté du binaire) |
| `-browser "C:\...\chrome.exe"` | Impose le navigateur à utiliser (sinon détection automatique Edge/Chrome/Chromium/Brave) |
| `-no-kiosk` | Ouvre un onglet classique au lieu d'une fenêtre dédiée |
| `-no-browser` | Serveur seul, sans ouvrir de navigateur (l'application ne s'éteint alors pas automatiquement) |
| `-version` | Affiche la version |

### Prise en main express

1. **Cliquez n'importe quel texte de l'affiche** et tapez : c'est la façon la plus rapide de tout modifier.
2. Cliquez une étape, un logo ou une zone colorée : le **panneau contextuel** propose les réglages pertinents.
3. Glissez les logos, tirez les pastilles numérotées pour réordonner les étapes.
4. Le formulaire de gauche reste disponible pour une édition structurée (et le bouton « Ouvrir dans le formulaire » du panneau y mène directement).
5. **Export** : bouton `Export` → PNG, PDF ou HTML.

## 📤 Exports

| Format | Contenu | Usage type |
|---|---|---|
| **PNG** | Image haute résolution (72 / 150 / 300 DPI), format papier exact | Affichage écran, intégration documentaire |
| **PDF** | **Vectoriel** (texte net à toute échelle), une page au format choisi | Impression A4 → A2 |
| **HTML autonome** | Fichier unique, images intégrées en base64, lisible sur n'importe quel poste avec un simple navigateur, **ré-importable** dans le générateur (le projet complet y est embarqué) | Diffusion intranet, archivage, transfert entre postes |

Les trois exports partagent le même moteur de rendu que la prévisualisation : ce que vous voyez est ce que vous obtenez.

## 🔒 Données et confidentialité

Conçu pour des procédures de sécurité, l'outil applique ses propres règles :

- **Aucune donnée ne quitte le poste.** Pas de télémétrie, pas de compte, pas d'appel à des CDN : les bibliothèques web sont embarquées dans le binaire.
- Le serveur n'écoute que sur **127.0.0.1** et rejette les requêtes dont l'hôte ne correspond pas (protection contre le DNS rebinding).
- Les projets sont des fichiers **JSON lisibles** stockés dans `data/projects/<id>/` à côté du binaire (ou du répertoire `-data`) : sauvegarde et versionnage triviaux — copiez le dossier.
- Les images importées par fichier sont validées (types autorisés : PNG, JPEG, GIF, WebP ; contrôle des signatures binaires ; SVG refusé).

## 🔧 Dépannage

**« Erreur export : impossible de lancer le navigateur »** — L'export pilote une instance invisible d'Edge/Chrome. Trois stratégies de lancement sont tentées automatiquement et le message d'erreur détaille chaque tentative avec la sortie réelle du navigateur : lisez-le, la cause y figure. Solutions par ordre d'efficacité : fermer toutes les fenêtres du navigateur puis réessayer ; désigner un autre navigateur (`-browser "C:\Program Files\Google\Chrome\Application\chrome.exe"`) ; sous Edge, désactiver l'« Accélération du démarrage » (`edge://settings/system`) ; vérifier antivirus / AppLocker.

**Une fenêtre blanche apparaît brièvement pendant un export** — Comportement normal de la stratégie de secours « fenêtre hors écran » quand le mode invisible du navigateur est indisponible sur le poste.

**Écran tactile** — Le glisser tactile (logos, réordonnancement) est partiellement fonctionnel ; l'usage à la souris reste le mode supporté à ce stade.

**L'application semble encore active après fermeture** — Elle s'éteint ~15 s après la fermeture de la fenêtre. Au-delà, vérifiez le gestionnaire des tâches ; le journal (`-Console` lors du build) indique la cause.

**Le port est occupé** — Par défaut le port est aléatoire ; avec `-port`, choisissez-en un libre. Une seule instance tourne à la fois (un verrou l'assure) ; lancer le binaire une seconde fois ramène la fenêtre existante.

## 🔄 Compatibilité des fichiers

| Origine | Comportement |
|---|---|
| Projets **v4.x** (générateur Go v2.x) | Format natif (`_version: 4.1`) — les champs récents (position libre des logos) sont simplement ignorés par les versions plus anciennes |
| Projets **v1.5** (ancien générateur HTML) | **Migration automatique** à l'import |
| **Exports HTML autonomes** | Ré-importables directement (glisser le fichier dans Importer) |

## 🏗️ Architecture (contributeurs)

```
├── main.go                     # Démarrage, options, cycle de vie (heartbeat)
├── internal/
│   ├── browser/                # Détection navigateur, kiosk, exports chromedp
│   │   ├── detect.go           #   Détection Edge/Chrome/Chromium/Brave
│   │   ├── kiosk.go            #   Fenêtre applicative dédiée
│   │   └── export.go           #   PNG/PDF : lancement piloté, port DevTools
│   │                           #   explicite, 3 stratégies, fermeture propre
│   ├── project/                # Modèle v4.1, stockage, sanitisation, migration
│   └── server/                 # Routeur HTTP, API REST, sécurité (Host, nosniff)
└── web/                        # Interface embarquée dans le binaire (go:embed)
    ├── editor/                 # Éditeur Alpine.js (bibliothèques vendorées)
    └── render/                 # MOTEUR DE RENDU UNIQUE
```

Principes directeurs :

- **Un seul moteur de rendu** (`web/render/index.html`) sert la prévisualisation (iframe pilotée par postMessage), les exports PNG/PDF (chromedp) et l'export HTML autonome (auto-rendu). C'est ce qui garantit le WYSIWYG strict.
- La prévisualisation met en page à la **largeur exacte du viewport d'export** (1200 px CSS) ; l'affichage à l'échelle passe par le zoom.
- L'éditeur et la frame dialoguent par un protocole postMessage documenté dans les sources (`ig-render`, `ig-select`, `ig-edit`, `ig-move`…), avec silencieux anti-boucle pendant la frappe et les gestes.
- Toute entrée utilisateur est **sanitisée deux fois** (whitelist HTML côté éditeur et côté Go) ; les valeurs numériques sont bornées côté serveur.
- Tests : unitaires Go (`go test ./...`) ; le développement s'appuie sur des tests E2E chromedp pilotant l'éditeur réel.

Contributions bienvenues par issues et pull requests — en acceptant que le projet reste sous la licence non commerciale décrite ci-dessous.

## 📜 Licence

**Licence d'Utilisation Libre Non Commerciale (LULNC) v1.0** — texte complet dans [LICENSE](LICENSE).

En résumé (le fichier LICENSE fait foi) :

- ✅ **Utilisation libre et gratuite**, code source ouvert, modification et redistribution gratuite autorisées : particuliers, administrations d'État, collectivités territoriales, établissements publics de santé et hospitaliers, enseignement et recherche publics, associations à but non lucratif — y compris dans le cadre de leurs missions.
- ❌ **Sans licence commerciale** : vente ou monétisation du logiciel, prestations rémunérées réalisées avec (conseil, formation payante, livrables facturés), intégration à une offre commerciale ou SaaS, utilisation par ou pour une entité à but lucratif.
- 📄 Les documents que **vous** produisez avec l'outil vous appartiennent et se diffusent librement, tant que leur production n'est pas une prestation facturée.

Entreprises et prestataires : une licence commerciale peut être convenue — contact via le dépôt GitHub.

*Note : il s'agit d'une licence « source disponible, non commerciale », et non d'une licence open source au sens de l'OSI.*

## 📅 Historique des versions

- **2.3** — Modèles personnalisés : figez un projet (charte graphique, logos, tournures) comme modèle réutilisable depuis la bibliothèque ; bouton d'ajout d'étape directement dans l'affiche ; publication automatisée des binaires (GitHub Actions).
- **2.2** — Éditeur WYSIWYG complet : édition directe sur l'affiche, panneau contextuel déplaçable, placement libre des logos (glisser + guides + redimensionnement), réordonnancement des étapes et de la timeline au glisser, zoom de prévisualisation, alignement préviz/export garanti au pixel. Modèle de données v4.1.
- **2.1** — Exports fiabilisés sur tous postes (lancement navigateur piloté, port DevTools explicite, stratégies de secours, diagnostics détaillés) ; cycle de vie par battement de cœur (fermeture automatique fiable) ; export HTML autonome ré-importable ; durcissement sécurité (contrôle Host, validation des uploads, sanitisation idempotente).
- **2.0** — Réécriture en application native Go : binaire unique, interface embarquée, serveur strictement local, moteur de rendu unifié.
- **1.x** — Générateur HTML monopage d'origine (Alpine.js), migration automatique des projets assurée.

---

*Développé par Pierre-Eric Guillemin pour la communauté SSI du secteur public.*
