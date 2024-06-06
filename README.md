## Проблема с godoc

При попытке использования godoc возникла проблема:
```
~/go/bin/godoc -http=:9090
using module mode; GOMOD=/home/konstantin/git/metrics/go.mod
2024/06/07 01:21:28 godoc: corpus fstree is nil
```

Поискал в интернете, нашел известный баг https://github.com/golang/go/issues/59431

Поэтому вместо godoc надо использовать pkgsite

## Использование pkgsite

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
