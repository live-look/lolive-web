# Lolive

Lolive is videostreaming service.

TODO: translate to English

# Документация по запуску (deprecated)

Используется make.

- `make build` - выполняет сборку приложения
- `make up` - запускает сервер, доступный по адресу www.camforchat.docker
- `make migrate` - выполняет миграции бд
- `make down` - полная остановка приложения. Все запущенные контейнеры будут остановлены.
- `make sh` - зайти в оболочку linux контейнера приложения.
- `make psql` - зайти в psql базы

Перед тем, как запустить локально приложение следует выполнить шаги:

1) Установить `dnsmasq`, и прописать в конфиг зону docker для разработки. Прописать в системе адрес резолвера 127.0.0.1.
2) Требуется сгенерировать самоподписные ssl-сертификаты. Без них стриминг не будет работать. Подробности в configs/cert,
там лежит скрипт генерации самоподписных сертификатов.
3) Собрать приложение - `make build`
4) Выполнить миграции - `make migrate`
5) Запустить - `make up`, открыть в браузере `https://www.camforchat.docker`
