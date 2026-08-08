---
title: API Fever
description: Connecter des clients compatibles Fever à yarr.
weight: 5
---

yarr prend en charge l'API Fever sur le point de terminaison `/fever`. Cela permet aux clients tiers compatibles Fever (ex. Reeder, Unread) de se connecter à votre instance yarr et de lire vos flux.

## Activer l'API Fever

L'API Fever nécessite une authentification. Définissez un nom d'utilisateur et un mot de passe sur votre serveur yarr avec l'option `-auth` (ou la variable d'environnement `YARR_AUTH`) :

```sh
yarr -auth nom_utilisateur:mot_de_passe
```

## Configurer un client compatible Fever

1. Assurez-vous que votre serveur yarr est accessible depuis le client.
2. Dans le client, renseignez l'URL de votre serveur yarr suivie de `/fever`.
3. Saisissez le nom d'utilisateur et le mot de passe configurés sur le serveur yarr.

## Remarques

La spécification de l'API Fever n'étant pas très précise, des problèmes de compatibilité peuvent survenir entre le serveur et le client.

Les applications suivantes ont été testées avec yarr :

> Différentes applications acceptent différents formats d'URL. Vérifiez si l'URL doit inclure le protocole `http://` et un slash final `/`.

| Application | Plateformes | URL du serveur de configuration |
| :-- | :-- | :-- |
| [Reeder](https://reederapp.com/) | macOS, iOS | `127.0.0.1:7070/fever` ou `http://127.0.0.1:7070/fever` |
| [ReadKit](https://readkit.app/) | macOS, iOS | `http://127.0.0.1:7070/fever` |
| [Fluent Reader](https://github.com/yang991178/fluent-reader) | macOS, Windows | `http://127.0.0.1:7070/fever/` |
| [Unread](https://www.goldenhillsoftware.com/unread/) | iOS | `http://127.0.0.1:7070/fever` |
| [Fiery Feeds](https://voidstern.net/fiery-feeds) | macOS, iOS | `http://127.0.0.1:7070/fever` |