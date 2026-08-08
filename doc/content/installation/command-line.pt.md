---
title: Linha de Comandos
description: Execute o yarr como servidor e aceda a partir de qualquer navegador.
weight: 2
---

## Download

Transfira a versão de linha de comandos a partir da [página de lançamentos](https://github.com/nkanaev/yarr/releases/latest). O nome do ficheiro não contém `_gui`.

| Sistema Operativo | Arquitetura | Transferência |
|---|---|---|
| macOS | Apple Silicon | [yarr_darwin_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_darwin_arm64.zip) |
| macOS | Intel | [yarr_darwin_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_darwin_amd64.zip) |
| Windows | ARM64 | [yarr_windows_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_windows_arm64.zip) |
| Windows | x86-64 | [yarr_windows_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_windows_amd64.zip) |
| Linux | ARM64 | [yarr_linux_arm64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_arm64.zip) |
| Linux | ARMv7 | [yarr_linux_armv7.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_armv7.zip) |
| Linux | x86-64 | [yarr_linux_amd64.zip](https://github.com/nkanaev/yarr/releases/latest/download/yarr_linux_amd64.zip) |

## Execução

Inicie o servidor:

```sh
./yarr
```

Por predefinição, o servidor escuta em `127.0.0.1:7070`. O yarr guarda os dados automaticamente na pasta de configuração do utilizador.

## Exemplo

Execute o servidor em todas as interfaces com proteção por palavra-passe:

```sh
yarr -addr 0.0.0.0:7070 -auth alice:secret
```

Abra `http://host:7070` num navegador e inicie sessão com `alice` / `secret`.