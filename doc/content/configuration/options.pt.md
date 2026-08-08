---
title: Opções
description: Opções de linha de comandos e as respetivas variáveis de ambiente equivalentes.
weight: 3
---

O servidor aceita opções como argumentos de linha de comandos e/ou variáveis de ambiente. A opção de linha de comandos tem precedência sobre a variável de ambiente correspondente.

| Opção | Variável de ambiente | Descrição |
| ------------ | -------------------- | ---------------------------------------------------------------------------- |
| `-addr` | `YARR_ADDR` | Endereço onde o servidor é executado (predefinição `127.0.0.1:7070`) |
| `-base` | `YARR_BASE` | Caminho base do URL do serviço |
| `-auth` | `YARR_AUTH` | Nome de utilizador e palavra-passe no formato `utilizador:palavra-passe` |
| `-auth-file` | `YARR_AUTHFILE` | Caminho para o ficheiro com `utilizador:palavra-passe`. Tem precedência sobre `-auth` |
| `-cert-file` | `YARR_CERTFILE` | Caminho para o ficheiro de certificado TLS |
| `-key-file` | `YARR_KEYFILE` | Caminho para o ficheiro de chave TLS |
| `-db` | `YARR_DB` | Caminho para o ficheiro de base de dados/armazenamento |
| `-log-file` | `YARR_LOGFILE` | Caminho para o ficheiro de registo (log) |
| `-open` | — | Abrir o servidor no navegador |

## HTTPS

Tanto `-cert-file` como `-key-file` são necessários para ativar o HTTPS.