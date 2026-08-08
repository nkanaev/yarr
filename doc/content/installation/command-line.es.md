---
title: Línea de comandos
description: Ejecuta yarr como servidor y accede desde cualquier navegador.
weight: 2
---

## Descarga

Descarga la versión para línea de comandos desde la [página de lanzamientos](https://github.com/nkanaev/yarr/releases/latest). El nombre del archivo no contiene `_gui`.

| Sistema operativo | Arquitectura | Descarga |
|---|---|---|
| macOS | Apple Silicon | [yarr_darwin_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_darwin_arm64.zip) |
| macOS | Intel | [yarr_darwin_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_darwin_amd64.zip) |
| Windows | ARM64 | [yarr_windows_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_windows_arm64.zip) |
| Windows | x86-64 | [yarr_windows_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_windows_amd64.zip) |
| Linux | ARM64 | [yarr_linux_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_arm64.zip) |
| Linux | ARMv7 | [yarr_linux_armv7.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_armv7.zip) |
| Linux | x86-64 | [yarr_linux_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_amd64.zip) |

## Ejecución

Inicia el servidor:

```sh
./yarr
```

Por defecto, el servidor escucha en `127.0.0.1:7070`. yarr guarda los datos automáticamente en la carpeta de configuración del usuario.

## Ejemplo

Ejecutar el servidor en todas las interfaces de red con protección por contraseña:

```sh
yarr -addr 0.0.0.0:7070 -auth alice:secret
```

Abre `http://host:7070` en un navegador e inicia sesión con `alice` / `secret`.