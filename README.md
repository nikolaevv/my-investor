Мой инвестор
Сервис использует Tinkoff Invest Open API v2 через систему удалённого вызова gRPC для получения данных об акциях и их покупки.
Реализован как веб-сервер.

## Запросы
Созданы 4 запроса:
- Регистрация
- Аутентификация
- Получение информации об акции по тикеру
- Покупка акции по тикеру (с занесением покупки в БД)

## Используемые библиотеки
- gin (http веб-фреймворк)
- gorm (ORM для postgresql)
- grpc, protobuf (для связи с Tinkoff Invest API)
- viper (для конфигов)
- logrus (логирование)
- gomock (для юнит-тестирования)

## База данных
Используется реляционная база данных PostrgreSQL.
В ней хранятся:
- user (пользователь)
    - login
    - password_hash
- share (покупка акций пользователем)
    - ticker (тикер)
    - class_code (идентификатор биржи, например для Мосбиржи: TQBR, для СПб - SPBMX)
    - user_id (id пользователя-владельца)
    - quantity (кол-во лотов)

## Прочее
- Настроен docker compose (образы: go server + postgresql)
- Написаны юнит-тесты к хандлерам

Сервис был написан приблизительно за 15-20 часов

## Локальный деплой
`docker compose up`
