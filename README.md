# go-musthave-metrics-tpl

Шаблон репозитория для трека «Сервер сбора метрик и алертинга».

# Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

# Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-metrics-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

# Спецификация

Спецификация проекта находится в файле [SPECIFICATION.md](SPECIFICATION.md)

# Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

# Запуск тестов

Для запуска тестов можно использовать команду `make test`.

Для запуска тестов в docker контейнере используйте команду `make cli`, а затем команду `make test`.
Если докер контейнер не собран, то нужно предварительно вызвать команду `make toolchain`.

# Покрытие тестами

Для получения информации о покрытия тестами кодма надо использовать команду `make cover`.

Для для запуска команды `make cover` в docker контейнере предварительно используйте команду `make cli`.
Если докер контейнер не собран, то нужно предварительно вызвать команду `make toolchain`.

# TODO

TODO лист находится в файле [TODO.md](TODO.md)

# Проблема с godoc

При попытке использования godoc возникла проблема:
```
~/go/bin/godoc -http=:9090
using module mode; GOMOD=/home/konstantin/git/metrics/go.mod
2024/06/07 01:21:28 godoc: corpus fstree is nil
```

Поискал в интернете, нашел известный баг https://github.com/golang/go/issues/59431

Поэтому вместо godoc надо использовать pkgsite

# Использование pkgsite

Документацию по установке/использванию можно найти в самом pkgsite.

Для удобства напишу инстукцию по установке/запуску здесь:
* `go install golang.org/x/pkgsite/cmd/pkgsite@latest`
* `cd metrics` - заходим в папку и нашим репозиторием
* `~/go/bin/pkgsite -open .`:
```
2024/06/07 01:16:37 Info: go/packages.Load(["all"]) loaded 243 packages from . in 177.055605ms
2024/06/07 01:16:37 Info: go/packages.Load(std) loaded 289 packages from /usr/lib/go-1.22 in 328.622008ms
2024/06/07 01:16:37 Info: FetchDataSource: fetching github.com/k0st1a/metrics@v0.0.0
2024/06/07 01:16:37 Info: FetchDataSource: fetching std@latest
2024/06/07 01:16:37 Info: Listening on addr http://localhost:8080
2024/06/07 01:16:37 Info: FetchDataSource: fetched github.com/k0st1a/metrics@v0.0.0 using *fetch.goPackagesModuleGetter in 14.9159ms with error <nil>
2024/06/07 01:16:38 Info: FetchDataSource: fetched std@latest using *fetch.goPackagesModuleGetter in 106.506402ms with error <nil>
2024/06/07 01:16:38 Info: Failed to open browser window. Please visit http://localhost:8080 in your browser.
```
