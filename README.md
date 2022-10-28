# TakingSurvey

## Описание

Скрипт, который успешно проходит опрос, расположенный по адресу: http://185.204.3.165

## Стек
- Go
- Конфигурация приложения [Viper](https://github.com/spf13/viper)
- Парсинг HTML [goquery](https://github.com/PuerkitoBio/goquery)
- RPS [rate](https://pkg.go.dev/golang.org/x/time/rate)

## Запуск
1. Установить стандартные значения в конфигурационном файлу `config/config.yaml`:
```
url: "http://185.204.3.165"
link: "/question/"
timeout: 10
cnt_workers: 5
rps: 10
```
2. Открыть терминал и набрать:
```
$ make
```
По стандарту запускается цель `build`

## Учебный материал
[Парсинг с конкуренцией в Golang](https://uproger.com/veb-skrejping-s-konkurencziej-v-golang/)
