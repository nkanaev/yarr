---
title: Fever API
description: Подключение сторонних клиентов с поддержкой Fever к yarr.
weight: 5
---

yarr поддерживает Fever API по адресу `/fever`. Это позволяет
сторонним клиентам с поддержкой Fever (например, Reeder, Unread) подключаться к вашему
экземпляру yarr и читать ваши ленты.

## Включение Fever API

Fever API требует авторизации. Задайте имя пользователя и пароль на вашем
сервере yarr с помощью флага `-auth` (или переменной окружения `YARR_AUTH`):

```sh
yarr -auth username:password
```

## Настройка клиента с поддержкой Fever

1. Убедитесь, что ваш сервер yarr доступен для клиента.
2. В клиенте укажите URL вашего сервера yarr плюс `/fever`.
3. Введите имя пользователя и пароль, которые вы задали на сервере yarr.

## Примечания

Спецификация Fever API неточна, поэтому между сервером
и клиентом могут возникать проблемы совместимости.

Следующие приложения были протестированы с yarr:

> Разные приложения принимают разные форматы URL. Обратите внимание, включает ли URL
> схему `http://` и завершающий `/`.

| Приложение | Платформы | URL сервера в настройках |
| :-- | :-- | :-- |
| [Reeder](https://reederapp.com/) | macOS, iOS | `127.0.0.1:7070/fever` или `http://127.0.0.1:7070/fever` |
| [ReadKit](https://readkit.app/) | macOS, iOS | `http://127.0.0.1:7070/fever` |
| [Fluent Reader](https://github.com/yang991178/fluent-reader) | macOS, Windows | `http://127.0.0.1:7070/fever/` |
| [Unread](https://www.goldenhillsoftware.com/unread/) | iOS | `http://127.0.0.1:7070/fever` |
| [Fiery Feeds](https://voidstern.net/fiery-feeds) | macOS, iOS | `http://127.0.0.1:7070/fever` |