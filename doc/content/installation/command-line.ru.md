---
title: Командная строка
description: Запуск yarr в качестве сервера с доступом из любого браузера.
weight: 2
---

## Скачивание

Скачайте версию для командной строки со [страницы релизов](https://github.com/nkanaev/yarr/releases/latest). В имени файла нет `_gui`.

| ОС | Архитектура | Скачать |
|---|---|---|
| macOS | Apple Silicon | [yarr_darwin_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_darwin_arm64.zip) |
| macOS | Intel | [yarr_darwin_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_darwin_amd64.zip) |
| Windows | ARM64 | [yarr_windows_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_windows_arm64.zip) |
| Windows | x86-64 | [yarr_windows_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_windows_amd64.zip) |
| Linux | ARM64 | [yarr_linux_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_arm64.zip) |
| Linux | ARMv7 | [yarr_linux_armv7.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_armv7.zip) |
| Linux | x86-64 | [yarr_linux_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_amd64.zip) |

## Запуск

Запустите сервер:

```sh
./yarr
```

По умолчанию сервер слушает `127.0.0.1:7070`. yarr автоматически сохраняет данные в пользовательской папке конфигурации.

## Пример

Запустите сервер на всех интерфейсах с защитой паролем:

```sh
yarr -addr 0.0.0.0:7070 -auth alice:secret
```

Откройте `http://host:7070` в браузере и авторизуйтесь с логином `alice` и паролем `secret`.