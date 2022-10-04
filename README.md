# Gas Price Test Task

Тестовое задание для [AppMagic](http://appmagic.rocks)

## Задание

Задача: реализовать сервис обработки истории цены gas в сети ethereum.

Вводные: существует исторический массив данных по цене -
https://github.com/CryptoRStar/GasPriceTestTask/blob/main/gas_price.json

Необходимо посчитать:
1) Сколько было потрачено gas помесячно.
2) Среднюю цену gas за день.
3) Частотное распределение цены по часам(за весь период).
4) Сколько заплатили за весь период (gas price * value).

Требования к сервису:
1) Данные должны загружаться удаленно.
2) Сервис должен вернуть все значения в виде json файла.
3) Данные должны быть посчитаны максимально быстро.

## Запуск

#### Вариант 1: Запустить готовый Docker образ

1. `docker pull ghcr.io/qulaz/gas-price-test:latest`
2. `docker run --rm -p 8000:8000 ghcr.io/qulaz/gas-price-test:latest`

#### Вариант 2: Собрать Docker образ самому

1. `docker build . -t gas-price-test:lastest`
2. `docker run --rm -p 8000:8000 gas-price-test:lastest`

#### Вариант 3: Собрать бинарник самому

**Внимание**: требуется go 1.19+

1. `go build -o server -a ./cmd/app/main.go`
2. Запустит собранный бинарник

### Проверка работоспособности

Сервер запущен и по-умолчанию доступен на порту `8000`.
Swagger UI доступен по адресу: http://localhost:8000/swagger/index.html

## Конфигурация

Конфигурация осуществляется через `.env` файл или переменные среды.
Переменные и их описание доступны в [.env.dist](.env.dist) файле.
По-умолчанию приложение запускается с оптимальными стандартными параметрами,
для запуска отдельно ничего настраивать не нужно.

В случае конфигурирования через `.env` файл, то он должен находиться в одной
директории с запускаемым бинарником приложения.
