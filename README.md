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
| `-h` | `--help` | Afficher l'aide | - |

### Exemple

```text
~/Projets/Yg-scovery $ ./yg-scovery -u https://ygp4ph.me -i

_____.___.                  _________                                      
\__  |   | ____            /   _____/ ____  _______  __ ___________ ___.__.
 /   |   |/ ___\   ______  \_____  \_/ ___\/  _ \  \/ // __ \_  __ <   |  |
 \____   / /_/  > /_____/  /        \  \__(  <_> )   /\  ___/|  | \/\___  |
 / ______\___  /          /_______  /\___  >____/ \_/  \___  >__|   / ____|
 \/     /_____/                   \/     \/                \/       \/      v1.0.0
 
[INF] Scanning https://ygp4ph.me (Depth: 3)
[INF] Filter: Internal links only
[INT] https://ygp4ph.me/assets/pdp_anime.mp4
[INT] https://ygp4ph.me/T
[INT] https://ygp4ph.me/%EE
[INT] https://ygp4ph.me/%CD%BB%EB
[INT] https://ygp4ph.me/%A6%BD
[INT] https://ygp4ph.me/%9BB%B1%C4%C8%DF
[INT] https://ygp4ph.me/%92%EDV%E6
[INT] https://ygp4ph.me/assets/pdp.png
[INT] https://ygp4ph.me/%DF%9E%E8%E8&%EE
[INT] https://ygp4ph.me/styles.css
[INT] https://ygp4ph.me/assets/bg1.jpg
[INT] https://ygp4ph.me/assets/favi.png
[INT] https://ygp4ph.me/B
[INT] https://ygp4ph.me/index.html
[INT] https://ygp4ph.me/portfolio.html
[INT] https://ygp4ph.me/gallery.js
[INT] https://ygp4ph.me/index.html#links
[INT] https://ygp4ph.me#links
```
