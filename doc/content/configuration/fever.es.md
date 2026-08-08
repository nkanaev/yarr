---
title: API Fever
description: Conecta clientes compatibles con Fever a yarr.
weight: 5
---

yarr admite la API Fever en el punto de acceso (endpoint) `/fever`. Esto permite que clientes de terceros compatibles con Fever (p. ej. Reeder, Unread) se conecten a tu instancia de yarr y lean tus feeds.

## Habilitar la API Fever

La API Fever requiere autenticación. Configura un usuario y contraseña en tu servidor yarr con el parámetro `-auth` (o la variable de entorno `YARR_AUTH`):

```sh
yarr -auth usuario:contraseña
```

## Configurar un cliente compatible con Fever

1. Asegúrate de que tu servidor yarr sea accesible desde el cliente.
2. En el cliente, indica la URL de tu servidor yarr seguida de `/fever`.
3. Introduce el usuario y la contraseña que configuraste en el servidor yarr.

## Notas

La especificación de la API Fever no es precisa, por lo que pueden surgir problemas de compatibilidad entre el servidor y el cliente.

Las siguientes aplicaciones han sido probadas con yarr:

> Distintas aplicaciones aceptan diferentes formatos de URL. Comprueba si la URL incluye el esquema `http://` y la barra final `/`.

| Aplicación | Plataformas | URL del servidor en la configuración |
| :-- | :-- | :-- |
| [Reeder](https://reederapp.com/) | macOS, iOS | `127.0.0.1:7070/fever` o `http://127.0.0.1:7070/fever` |
| [ReadKit](https://readkit.app/) | macOS, iOS | `http://127.0.0.1:7070/fever` |
| [Fluent Reader](https://github.com/yang991178/fluent-reader) | macOS, Windows | `http://127.0.0.1:7070/fever/` |
| [Unread](https://www.goldenhillsoftware.com/unread/) | iOS | `http://127.0.0.1:7070/fever` |
| [Fiery Feeds](https://voidstern.net/fiery-feeds) | macOS, iOS | `http://127.0.0.1:7070/fever` |