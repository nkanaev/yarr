---
title: Options
description: Options en ligne de commande et leurs variables d'environnement équivalentes.
weight: 3
---

Le serveur accepte des options sous forme d'arguments en ligne de commande et/ou de variables d'environnement. Une option en ligne de commande a la priorité sur sa variable d'environnement correspondante.

| Option | Variable d'environnement | Description |
| ------------ | -------------------- | ---------------------------------------------------------------------------- |
| `-addr` | `YARR_ADDR` | Adresse sur laquelle le serveur s'exécute (par défaut `127.0.0.1:7070`) |
| `-base` | `YARR_BASE` | Chemin de base de l'URL du service |
| `-auth` | `YARR_AUTH` | Nom d'utilisateur et mot de passe au format `nom_utilisateur:mot_de_passe` |
| `-auth-file` | `YARR_AUTHFILE` | Chemin vers un fichier contenant `nom_utilisateur:mot_de_passe`. Prioritaire sur `-auth` |
| `-cert-file` | `YARR_CERTFILE` | Chemin vers le fichier de certificat TLS |
| `-key-file` | `YARR_KEYFILE` | Chemin vers le fichier de clé TLS |
| `-db` | `YARR_DB` | Chemin du fichier de stockage |
| `-log-file` | `YARR_LOGFILE` | Chemin du fichier de journalisation (log) |
| `-open` | — | Ouvrir le serveur dans le navigateur |

## HTTPS

`-cert-file` et `-key-file` sont tous deux requis pour activer le HTTPS.