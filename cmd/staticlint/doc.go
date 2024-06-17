/*
Staticlint - линтер, собранные с использованем go/analysis/multichecker и
состоящий из различных анализаторов.

# Использование

Без флагов (будут задействованы все анализаторы):

	$ go run cmd/staticlint/main.go ./...

С использованием флага (будет задействован заданный анализатор):

	$ go run cmd/staticlint/main.go -osexit ./...

С использованием только nilaway на пакете github.com/k0st1a/metrics
(nilaway_config нужен для задания проверяемого пакета)

	$ go run ./cmd/staticlint/main.go -nilaway_config.include-pkgs="github.com/k0st1a/metrics" ./...

С использованием всех анализаторов и nilaway на пакете github.com/k0st1a/metrics

	$ go run ./cmd/staticlint/main.go -nilaway_config.include-pkgs="github.com/k0st1a/metrics" ./...

Все флаги (предоставляются multichecker-ом) можно посмотреть:

	$ go run cmd/staticlint/main.go -h

# Используемые анализаторы

  - стандартные статическиt анализаторы пакета 'golang.org/x/tools/go/analysis/passes':

    -- golang.org/x/tools/go/analysis/passes/printf

    -- golang.org/x/tools/go/analysis/passes/shift

    -- golang.org/x/tools/go/analysis/passes/structtag

  - все анализаторы класса SA пакета staticcheck.io:

    -- Various misuses of the standard library (SA1 group) (https://staticcheck.io/docs/checks/#SA1)

    -- Concurrency issues (SA2 group) (https://staticcheck.io/docs/checks/#SA2)

    -- Testing issues (SA3 group) (https://staticcheck.io/docs/checks/#SA3)

    -- Code that isn't really doing anything (SA4 group) (https://staticcheck.io/docs/checks/#SA4)

    -- Correctness issues (SA5 group) (https://staticcheck.io/docs/checks/#SA5)

    -- Performance issues (SA6 groups) (https://staticcheck.io/docs/checks/#SA6)

    -- Dubious code constructs that have a high probability of being wrong (SA9 group)
    (https://staticcheck.io/docs/checks/#SA9)

  - анализаторы staticcheck пакета staticcheck.io:

    -- Code simplifications (S group) (https://staticcheck.io/docs/checks/#S)

    -- Stylistic issues (ST group) (https://staticcheck.io/docs/checks/#ST)

    -- Quickfixes (QF group) (https://staticcheck.io/docs/checks/#QF)

  - публичные анализаторы:

    -- nilaway (https://github.com/uber-go/nilaway)

    -- nilaway_config (https://github.com/uber-go/nilaway/blob/main/config/config.go)

    -- interfacebloat (https://github.com/sashamelentyev/interfacebloat)

  - собственный анализатор osexit, который запрещает прямой вызова os.Exit в функции main пакета main.
*/
package main
