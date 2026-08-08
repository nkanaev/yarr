---
title: API Fever
description: Ligar clientes compatíveis com Fever ao yarr.
weight: 5
---

O yarr suporta a API Fever no ponto de extremidade (endpoint) `/fever`. Isto permite que clientes de terceiros compatíveis com Fever (ex.: Reeder, Unread) se liguem à sua instância do yarr e leiam os seus feeds.

## Ativar a API Fever

A API Fever requer autenticação. Defina um nome de utilizador e palavra-passe no seu servidor yarr com a opção `-auth` (ou a variável de ambiente `YARR_AUTH`):

```sh
yarr -auth utilizador:palavrapasse
```

## Configurar um cliente compatível com Fever

1. Certifique-se de que o seu servidor yarr está acessível a partir do cliente.
2. No cliente, introduza o URL do seu servidor yarr seguido de `/fever`.
3. Introduza o nome de utilizador e a palavra-passe que configurou no servidor yarr.

## Notas

A especificação da API Fever não é precisa, pelo que podem ocorrer problemas de compatibilidade entre o servidor e o cliente.

As seguintes aplicações foram testadas com o yarr:

> Aplicações diferentes aceitam formatos de URL diferentes. Note se o URL inclui o esquema `http://` e a barra final `/`.

| Aplicação | Plataformas | URL do Servidor na Configuração |
| :-- | :-- | :-- |
| [Reeder](https://reederapp.com/) | macOS, iOS | `127.0.0.1:7070/fever` ou `http://127.0.0.1:7070/fever` |
| [ReadKit](https://readkit.app/) | macOS, iOS | `http://127.0.0.1:7070/fever` |
| [Fluent Reader](https://github.com/yang991178/fluent-reader) | macOS, Windows | `http://127.0.0.1:7070/fever/` |
| [Unread](https://www.goldenhillsoftware.com/unread/) | iOS | `http://127.0.0.1:7070/fever` |
| [Fiery Feeds](https://voidstern.net/fiery-feeds) | macOS, iOS | `http://127.0.0.1:7070/fever` |