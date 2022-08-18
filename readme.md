# gitlab-ci-bot

Предназначен для оповещений разработчиков о процессах GitLab CI, завершенных с ошибкой, через Telegram.

## Настройка и запуск

1. [Создание Telegram-бота](https://core.telegram.org/bots#3-how-do-i-create-a-bot)
1. Добавить бота к чату, в который хочется получать уведомления
1. [Узнать ID чата](https://sean-bradley.medium.com/get-telegram-chat-id-80b575520659)
1. [Настроить в GitLab CI отправку WebHook с событями Pipelines Events](https://docs.gitlab.com/ee/user/project/integrations/webhooks.html#configure-a-webhook-in-gitlab) и указать URL `http://host:port/webhook`
1. Прописать токен бота и ID чата в конфигурационном файле:
```toml
[TelegramBot]
Token = "<BOT_TOKEN>"
Chat = <CHAT_ID>

[Server]
Endpoint = "0.0.0.0:8080"
```

Запуск:

```bash
./gitlab-ci-bot -config <path_to_config>
```

## Запуск в Docker

```bash
docker-build -t gitlab-ci-bot .
docker run -e TELEGRAM_BOT_TOKEN=<token> -e TELEGRAM_BOT_CHAT=<chatID> -p 0.0.0.0:8080:8080/tcp gitlab-ci-bot
```