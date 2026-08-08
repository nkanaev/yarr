---
title: Ligne de commande
description: Exécuter yarr en tant que serveur et y accéder depuis n'importe quel navigateur.
weight: 2
---

## Téléchargement

Téléchargez la version en ligne de commande depuis la [page des versions (releases)](https://github.com/nkanaev/yarr/releases/latest). Le nom du fichier ne contient pas `_gui`.

| Système d'exploitation | Architecture | Téléchargement |
|---|---|---|
| macOS | Apple Silicon | [yarr_darwin_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_darwin_arm64.zip) |
| macOS | Intel | [yarr_darwin_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_darwin_amd64.zip) |
| Windows | ARM64 | [yarr_windows_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_windows_arm64.zip) |
| Windows | x86-64 | [yarr_windows_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_windows_amd64.zip) |
| Linux | ARM64 | [yarr_linux_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_arm64.zip) |
| Linux | ARMv7 | [yarr_linux_armv7.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_armv7.zip) |
| Linux | x86-64 | [yarr_linux_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_amd64.zip) |

## Exécution

Démarrez le serveur :

```sh
./yarr
```

Par défaut, le serveur écoute sur `127.0.0.1:7070`. yarr enregistre automatiquement les données dans le dossier de configuration de l'utilisateur.

## Exemple

Exécuter le serveur sur toutes les interfaces avec une protection par mot de passe :

```sh
yarr -addr 0.0.0.0:7070 -auth alice:secret
```

Ouvrez `http://host:7070` dans un navigateur et connectez-vous avec `alice` / `secret`.