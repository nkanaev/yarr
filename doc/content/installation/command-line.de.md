---
title: Kommandozeile
description: Führen Sie yarr als Server aus und greifen Sie über einen beliebigen Browser darauf zu.
weight: 2
---

## Download

Laden Sie die Kommandozeilenversion von der [Release-Seite](https://github.com/nkanaev/yarr/releases/latest) herunter. Der Dateiname enthält kein `_gui`.

| Betriebssystem | Architektur | Download |
|---|---|---|
| macOS | Apple Silicon | [yarr_darwin_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_darwin_arm64.zip) |
| macOS | Intel | [yarr_darwin_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_darwin_amd64.zip) |
| Windows | ARM64 | [yarr_windows_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_windows_arm64.zip) |
| Windows | x86-64 | [yarr_windows_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_windows_amd64.zip) |
| Linux | ARM64 | [yarr_linux_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_arm64.zip) |
| Linux | ARMv7 | [yarr_linux_armv7.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_armv7.zip) |
| Linux | x86-64 | [yarr_linux_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_amd64.zip) |

## Ausführen

Starten Sie den Server:

```sh
./yarr
```

Standardmäßig lauscht der Server auf `127.0.0.1:7070`. yarr speichert die Daten automatisch im Konfigurationsordner des Benutzers.

## Beispiel

Führen Sie den Server auf allen Netzwerkschnittstellen mit Passwortschutz aus:

```sh
yarr -addr 0.0.0.0:7070 -auth alice:secret
```

Öffnen Sie `http://host:7070` in einem Browser und melden Sie sich mit `alice` / `secret` an.