---
title: Opciones
description: Parámetros de la línea de comandos y sus variables de entorno equivalentes.
weight: 3
---

El servidor acepta opciones como argumentos de línea de comandos y/o variables de entorno. Un parámetro de línea de comandos tiene prioridad sobre su variable de entorno correspondiente.

| Opción | Variable de entorno | Descripción |
| ------------ | -------------------- | ---------------------------------------------------------------------------- |
| `-addr` | `YARR_ADDR` | Dirección en la que se ejecuta el servidor (por defecto `127.0.0.1:7070`) |
| `-base` | `YARR_BASE` | Ruta base de la URL del servicio |
| `-auth` | `YARR_AUTH` | Usuario y contraseña en formato `usuario:contraseña` |
| `-auth-file` | `YARR_AUTHFILE` | Ruta a un archivo con `usuario:contraseña`. Tiene prioridad sobre `-auth` |
| `-cert-file` | `YARR_CERTFILE` | Ruta al archivo del certificado TLS |
| `-key-file` | `YARR_KEYFILE` | Ruta al archivo de la clave privada TLS |
| `-db` | `YARR_DB` | Ruta al archivo de almacenamiento |
| `-log-file` | `YARR_LOGFILE` | Ruta al archivo de registro (log) |
| `-open` | — | Abrir el servidor en el navegador |

## HTTPS

Se requieren tanto `-cert-file` como `-key-file` para habilitar HTTPS.