# Yg-Scovery


Un crawler web rapide et efficace écrit en Go pour la découverte automatique de liens et l'exploration de sites web.


## Installation

### Installation rapide

```bash
go install github.com/ygp4ph/yg-scovery@latest
```

### Compilation depuis les sources

```bash
# Cloner le repository
git clone https://github.com/ygp4ph/yg-scovery.git
cd yg-scovery

# Installer les dépendances
go mod download

# Compiler
go build -o yg-scovery

# Ou installer directement
go install
```

## Utilisation

### Syntaxe de base

```bash
./yg-scovery -u <URL> [options]
```

### Options disponibles

| Flag | Alias | Description | Défaut |
|------|-------|-------------|--------|
| `-u` | `--url` | URL cible à crawler (requis) | - |
| `-d` | `--depth` | Profondeur maximale de récursion | 3 |
| `-e` | `--ext` | Afficher uniquement les liens externes | false |
| `-i` | `--int` | Afficher uniquement les liens internes | false |
| `-o` | `--output` | Sauvegarder les résultats en JSON | - |
| `-v` | `--verbose` | Afficher les erreurs détaillées | false |
| `-h` | `--help` | Afficher l'aide | - |

### Exemple

```text
~/CTF/HTB/en_cours $ ~/Projets/yg-scovery/yg-scovery -u ygp4ph.me -i

   __  ______ _      ______________ _   _____  _______  __
  / / / / __ `/_____/ ___/ ___/ __ \ | / / _ \/ ___/ / / /
 / /_/ / /_/ /_____(__  ) /__/ /_/ / |/ /  __/ /  / /_/ / 
 \__, /\__, /     /____/\___/\____/|___/\___/_/   \__, /  
/____//____/                                     /____/   v2.0.1
 
[INF] Scanning https://ygp4ph.me (Depth: 3)
[INF] Filter: Internal links only
[INT] https://ygp4ph.me/assets/pdp_anime.mp4
[INT] https://ygp4ph.me/assets/pdp.png
[INT] https://ygp4ph.me/
[INT] https://ygp4ph.me/script.js
[INT] https://ygp4ph.me/assets/rooftop.jpg
[INT] https://ygp4ph.me/styles.css
[INT] https://ygp4ph.me#links
[INT] https://ygp4ph.me/assets/favi.png
[INT] https://ygp4ph.me/Portfolio/
[INT] https://ygp4ph.me/writeups/
[INT] https://ygp4ph.me/writeups/writeups.css
[INT] https://ygp4ph.me/writeups/chemistry/
[INT] https://ygp4ph.me/writeups/trickster/
[INT] https://ygp4ph.me/assets/labalsa.jpg
[INT] https://ygp4ph.me/Portfolio/gallery.js
[INT] https://ygp4ph.me/#links
[INT] https://ygp4ph.me/writeups/trickster/image.png
```
