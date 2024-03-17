LAST = 14
NUMBERS := $(shell seq 1 ${LAST})
ITERS := $(addprefix iter,${NUMBERS})

.DEFAULT_GOAL := build
.SILENT: ${ITERS}

PG_USER = "metrics-user"
PG_PASSWORD = "metrics-password"
PG_DB = "metrics-db"
PG_HOST = "localhost"
PG_PORT = "5432"
DATABASE_DSN = "postgres://${PG_USER}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${PG_DB}?sslmode=disable"
PG_IMAGE = "postgres:13.13-bullseye"
DOCKER_CONTEINER_NAME = "metrics-pg-13.3"

METRICSTEST_ARGS = -test.v -source-path=.

.PHONY:build
build:
	go build -C ./cmd/agent/ -o agent
	go build -C ./cmd/server/ -o server

.PHONY:clean
clean:
	-rm -f ./cmd/agent/agent
	-rm -f ./cmd/server/server

.PHONY:statictest
statictest:
	go vet -vettool=$$(which statictest) ./...


test: build statictest
	go test -v ./...

.PHONY: miter7
miter7: build statictest
	METRICSTEST_ARGS="${METRICSTEST_ARGS} -test.run=TestIteration7" ; \
	SERVER_PORT=$$(random unused-port) ; \
	ADDRESS="localhost:$${SERVER_PORT}" ; \
	TEMP_FILE=$$(random tempfile) ; \
	metricstest $$METRICSTEST_ARGS \
				-binary-path=cmd/server/server \
				-agent-binary-path=cmd/agent/agent \
				-server-port=$$SERVER_PORT ;

.PHONY: miter8
miter8: build statictest
	METRICSTEST_ARGS="${METRICSTEST_ARGS} -test.run=TestIteration8" ; \
	ADDRESS="localhost:8080" ; \
	TEMP_FILE=$$(random tempfile) ; \
	metricstest $$METRICSTEST_ARGS \
				-binary-path=cmd/server/server \
				-agent-binary-path=cmd/agent/agent \
				-server-port=8080 \
				-file-storage-path=$$TEMP_FILE ;

.PHONY: miter9
miter9: build statictest
	METRICSTEST_ARGS="${METRICSTEST_ARGS} -test.run=TestIteration9" ; \
			ADDRESS="localhost:8080" ;\
			TEMP_FILE="/tmp/metrics-db.json" ; \
			metricstest $$METRICSTEST_ARGS \
						-binary-path=cmd/server/server \
						-agent-binary-path=cmd/agent/agent \
						-server-port=8080 \
						-file-storage-path=$$TEMP_FILE ;

.PHONY: miter10
miter10: build statictest delay-after-db-run
	SERVER_PORT=$$(random unused-port) ; \
	ADDRESS="localhost:$${SERVER_PORT}" ; \
	TEMP_FILE=$$(random tempfile) ; \
	METRICSTEST_ARGS="${METRICSTEST_ARGS} -test.run=TestIteration10[AB]" ; \
	metricstest $$METRICSTEST_ARGS \
				-binary-path=cmd/server/server \
				-agent-binary-path=cmd/agent/agent \
				-server-port=$$SERVER_PORT \
				-database-dsn=${DATABASE_DSN} ;

.PHONY: miter11
miter11: build statictest delay-after-db-run
	SERVER_PORT=$$(random unused-port) ; \
	ADDRESS="localhost:$${SERVER_PORT}" ; \
	TEMP_FILE=$$(random tempfile) ; \
	METRICSTEST_ARGS="${METRICSTEST_ARGS} -test.run=TestIteration11" ; \
	metricstest $$METRICSTEST_ARGS \
				-binary-path=cmd/server/server \
				-agent-binary-path=cmd/agent/agent \
				-server-port=$$SERVER_PORT \
				-database-dsn=${DATABASE_DSN} ;

.PHONY: miter12
miter12: build statictest delay-after-db-run
	#SERVER_PORT=$$(random unused-port) ;
	SERVER_PORT=8081 ; \
	ADDRESS="localhost:$${SERVER_PORT}" ; \
	TEMP_FILE=$$(random tempfile) ; \
	METRICSTEST_ARGS="${METRICSTEST_ARGS} -test.run=TestIteration12" ; \
	metricstest $$METRICSTEST_ARGS \
				-binary-path=cmd/server/server \
				-agent-binary-path=cmd/agent/agent \
				-server-port=$$SERVER_PORT \
				-database-dsn=${DATABASE_DSN} ;

.PHONY: miter13
miter13: build statictest delay-after-db-run
	SERVER_PORT=$$(random unused-port) ; \
	ADDRESS="localhost:$${SERVER_PORT}" ; \
	TEMP_FILE=$$(random tempfile) ; \
	METRICSTEST_ARGS="${METRICSTEST_ARGS} -test.run=TestIteration13" ; \
	metricstest $$METRICSTEST_ARGS \
				-binary-path=cmd/server/server \
				-agent-binary-path=cmd/agent/agent \
				-server-port=$$SERVER_PORT \
				-database-dsn=${DATABASE_DSN} ;

.PHONY: ${ITERS}
${ITERS}: iter%: build statictest db-run;
	for i in $(shell seq 1 $*) ; do \
		METRICSTEST_ARGS="${METRICSTEST_ARGS} -test.run=TestIteration$$i[AB]?$$" ; \
		if [ $$i -eq 1 ]; then \
			metricstest $$METRICSTEST_ARGS \
						-binary-path=cmd/server/server ; \
		elif [ $$i -eq 2 ]; then \
			metricstest $$METRICSTEST_ARGS  \
						-agent-binary-path=cmd/agent/agent ; \
		elif [ $$i -eq 3 ]; then \
			metricstest $$METRICSTEST_ARGS \
						-binary-path=cmd/server/server \
						-agent-binary-path=cmd/agent/agent ; \
		elif [ $$i -ge 4 ] && [ $$i -le 7 ]; then \
			SERVER_PORT=$$(random unused-port) ; \
			ADDRESS="localhost:$${SERVER_PORT}" ; \
			TEMP_FILE=$$(random tempfile) ; \
			metricstest $$METRICSTEST_ARGS \
						-binary-path=cmd/server/server \
						-agent-binary-path=cmd/agent/agent \
						-server-port=$$SERVER_PORT ; \
		elif [ $$i -ge 8 ] && [ $$i -le 9 ]; then \
			SERVER_PORT=$$(random unused-port) ; \
			ADDRESS="localhost:$${SERVER_PORT}" ; \
			TEMP_FILE=$$(random tempfile) ; \
			metricstest $$METRICSTEST_ARGS \
						-binary-path=cmd/server/server \
						-agent-binary-path=cmd/agent/agent \
						-server-port=$$SERVER_PORT \
						-file-storage-path=$$TEMP_FILE ; \
		elif [ $$i -ge 10 ] && [ $$i -le 13 ]; then \
			SERVER_PORT=$$(random unused-port) ; \
			ADDRESS="localhost:$${SERVER_PORT}" ; \
			TEMP_FILE=$$(random tempfile) ; \
			metricstest $$METRICSTEST_ARGS \
						-binary-path=cmd/server/server \
						-agent-binary-path=cmd/agent/agent \
						-server-port=$$SERVER_PORT \
						-database-dsn=${DATABASE_DSN} ; \
		elif [ $$i -eq 14 ]; then \
			SERVER_PORT=$$(random unused-port) ; \
			ADDRESS="localhost:$${SERVER_PORT}" ; \
			TEMP_FILE=$$(random tempfile) ; \
			metricstest $$METRICSTEST_ARGS \
						-binary-path=cmd/server/server \
						-agent-binary-path=cmd/agent/agent \
						-server-port=$$SERVER_PORT \
						-database-dsn=${DATABASE_DSN} ; \
						-key="$$TEMP_FILE" ; \
			go test -v -race ./... ; \
		fi ; \
		if [ $$? -eq 1 ]; then \
			break ; \
		fi ; \
    done

.PHONY: delay-after-db-run
delay-after-db-run: db-run
	sleep 5s

.PHONY: db-run
db-run: db-image-pull db-stop
	-docker run \
	--name ${DOCKER_CONTEINER_NAME} \
	--rm -ti \
	-p ${PG_PORT}:5432 \
	-e POSTGRES_USER=${PG_USER} \
	-e POSTGRES_PASSWORD=${PG_PASSWORD} \
	-e POSTGRES_DB=${PG_DB} \
	-d ${PG_IMAGE}

.PHONY: db-stop
db-stop:
	-docker stop ${DOCKER_CONTEINER_NAME}

.PHONY: db-image-pull
db-image-pull:
	-docker image pull ${PG_IMAGE}


GOLANGCI_LINT_CACHE?=/tmp/praktikum-golangci-lint-cache

.PHONY: golangci-lint-run
golangci-lint-run: _golangci-lint-rm-unformatted-report

.PHONY: _golangci-lint-reports-mkdir
_golangci-lint-reports-mkdir:
	mkdir -p ./golangci-lint

.PHONY: _golangci-lint-run
_golangci-lint-run: _golangci-lint-reports-mkdir
	-docker run --rm \
    -v $(shell pwd):/app \
    -v $(GOLANGCI_LINT_CACHE):/root/.cache \
    -w /app \
    golangci/golangci-lint:v1.55.2 \
        golangci-lint run \
            -c .golangci.yml \
	> ./golangci-lint/report-unformatted.json

.PHONY: _golangci-lint-format-report
_golangci-lint-format-report: _golangci-lint-run
	cat ./golangci-lint/report-unformatted.json | jq > ./golangci-lint/report.json

.PHONY: _golangci-lint-rm-unformatted-report
_golangci-lint-rm-unformatted-report: _golangci-lint-format-report
	rm ./golangci-lint/report-unformatted.json

.PHONY: golangci-lint-clean
golangci-lint-clean:
	rm -rf ./golangci-lint
