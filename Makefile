# for debug use make SHELL="sh -x"

# Получаем полный путь к клиенту докера
DOCKER := $(shell which docker)
DOCKERFILE := Dockerfile
PWD := $(shell pwd)

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
PG_DATABASE_DSN = "postgres://${PG_USER}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${PG_DB}?sslmode=disable"
PG_IMAGE = "postgres:13.13-bullseye"
PG_DOCKER_CONTEINER_NAME = "metrics-pg-13.3"

SERVER_PORT="8080"
SERVER_HOST="127.0.0.1"
SERVER_TRUSTED_SUBNET="127.0.0.1/32"
#SERVER_TRUSTED_SUBNET="192.168.1.1/32"
PPROF_SERVER_PORT="8086"
PPROF_SERVER_HOST="0.0.0.0"

METRICSTEST_ARGS = -test.v -source-path=.

# Используем := чтобы переменная содержала значение на на момент определения этой переменной, см
# https://ftp.gnu.org/old-gnu/Manuals/make-3.79.1/html_chapter/make_6.html#SEC59
# TODO: так нужно задать все переменные в Makefile
BUILD_VERSION := 0.0.1
BUILD_DATE := $(shell date -u +"%Y-%m-%d %H:%M:%S:%N %Z")
BUILD_COMMIT := $(shell git rev-parse HEAD)

.PHONY:proto-install
proto-install:
	# from https://grpc.io/docs/protoc-installation/#:~:text=Linux%2C%20using%20apt%20or%20apt%2Dget
	sudo apt install -y protobuf-compiler
	# from https://practicum.yandex.ru/learn/go-advanced/courses/65ce3d44-da98-4684-9499-465ff6cc6c64/sprints/226895/topics/30311053-9716-4af0-9a23-f4fa0725f918/lessons/fa184729-fbbd-4a1c-ae11-4e12f66b7f64/#:~:text=%D0%A3%D1%81%D1%82%D0%B0%D0%BD%D0%BE%D0%B2%D0%BA%D0%B0%20%D1%83%D1%82%D0%B8%D0%BB%D0%B8%D1%82%20gRPC
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	export PATH="${PATH}:$(go env GOPATH)/bin"

.PHONY:proto-generate
proto-generate:
	protoc --go_out=./internal/proto/ \
		   --go_opt=paths=source_relative \
		   --go-grpc_out=./internal/proto/ \
		   --go-grpc_opt=paths=source_relative \
		   ./internal/proto/model.proto

GOLANG_LDFLAGS := -ldflags "-X 'main.buildVersion=${BUILD_VERSION}' \
                            -X 'main.buildDate=${BUILD_DATE}' \
                            -X 'main.buildCommit=${BUILD_COMMIT}'"

.PHONY:build
build:
	go build -C ./cmd/agent/ -o agent -buildvcs=false ${GOLANG_LDFLAGS}
	go build -C ./cmd/server/ -o server -buildvcs=false ${GOLANG_LDFLAGS}

.PHONY:clean
clean:
	-rm -f ./cmd/agent/agent
	-rm -f ./cmd/server/server

.PHONY:godoc
godoc:
	go install golang.org/x/pkgsite/cmd/pkgsite@latest
	~/go/bin/pkgsite -open .

.PHONY:statictest
statictest:
	go vet -vettool=$$(which statictest) ./...

.PHONY:staticlint
staticlint:
	go run cmd/staticlint/main.go ./...

.PHONY:test
test: build statictest staticlint
	go test -v -race -count=1 ./...

.PHONY: cover
cover:
	mkdir -pv ./cover && \
	go test -v -coverpkg=./... -coverprofile=./cover/cover.profile.tmp ./... && \
	cat ./cover/cover.profile.tmp | grep -v "_easyjson.go" | grep -v "model.go" > ./cover/cover.profile && \
	rm ./cover/cover.profile.tmp && \
	go tool cover -func ./cover/cover.profile && \
	go tool cover -html ./cover/cover.profile -o ./cover/cover.html

.PHONY: cover-clean
cover-clean:
	rm -v -f ./cover/cover.profile ./cover/cover.html && \
	rm -v -r ./cover

.PHONY:test-analyzer
test-analyzer:
	go test -v -race -count=1 ./internal/pkg/analyzer/...

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
miter10: build statictest db-up
	SERVER_PORT=$$(random unused-port) ; \
	ADDRESS="localhost:$${SERVER_PORT}" ; \
	TEMP_FILE=$$(random tempfile) ; \
	METRICSTEST_ARGS="${METRICSTEST_ARGS} -test.run=TestIteration10[AB]" ; \
	metricstest $$METRICSTEST_ARGS \
				-binary-path=cmd/server/server \
				-agent-binary-path=cmd/agent/agent \
				-server-port=$$SERVER_PORT \
				-database-dsn=${PG_DATABASE_DSN} ;

.PHONY: miter11
miter11: build statictest db-up
	SERVER_PORT=$$(random unused-port) ; \
	ADDRESS="localhost:$${SERVER_PORT}" ; \
	TEMP_FILE=$$(random tempfile) ; \
	METRICSTEST_ARGS="${METRICSTEST_ARGS} -test.run=TestIteration11" ; \
	metricstest $$METRICSTEST_ARGS \
				-binary-path=cmd/server/server \
				-agent-binary-path=cmd/agent/agent \
				-server-port=$$SERVER_PORT \
				-database-dsn=${PG_DATABASE_DSN} ;

.PHONY: miter12
miter12: build statictest db-up
	#SERVER_PORT=$$(random unused-port) ;
	SERVER_PORT=8081 ; \
	ADDRESS="localhost:$${SERVER_PORT}" ; \
	TEMP_FILE=$$(random tempfile) ; \
	METRICSTEST_ARGS="${METRICSTEST_ARGS} -test.run=TestIteration12" ; \
	metricstest $$METRICSTEST_ARGS \
				-binary-path=cmd/server/server \
				-agent-binary-path=cmd/agent/agent \
				-server-port=$$SERVER_PORT \
				-database-dsn=${PG_DATABASE_DSN} ;

.PHONY: miter13
miter13: build statictest db-up
	SERVER_PORT=$$(random unused-port) ; \
	ADDRESS="localhost:$${SERVER_PORT}" ; \
	TEMP_FILE=$$(random tempfile) ; \
	METRICSTEST_ARGS="${METRICSTEST_ARGS} -test.run=TestIteration13" ; \
	metricstest $$METRICSTEST_ARGS \
				-binary-path=cmd/server/server \
				-agent-binary-path=cmd/agent/agent \
				-server-port=$$SERVER_PORT \
				-database-dsn=${PG_DATABASE_DSN} ;

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
						-database-dsn=${PG_DATABASE_DSN} ; \
		elif [ $$i -eq 14 ]; then \
			SERVER_PORT=$$(random unused-port) ; \
			ADDRESS="localhost:$${SERVER_PORT}" ; \
			TEMP_FILE=$$(random tempfile) ; \
			metricstest $$METRICSTEST_ARGS \
						-binary-path=cmd/server/server \
						-agent-binary-path=cmd/agent/agent \
						-server-port=$$SERVER_PORT \
						-database-dsn=${PG_DATABASE_DSN} ; \
						-key="$$TEMP_FILE" ; \
			go test -v -race ./... ; \
		fi ; \
		if [ $$? -eq 1 ]; then \
			break ; \
		fi ; \
    done

.PHONY: miter14
miter14: build statictest db-up
	SERVER_PORT=$$(random unused-port) ; \
	ADDRESS="localhost:$${SERVER_PORT}" ; \
	TEMP_FILE=$$(random tempfile) ; \
	metricstest -test.v -test.run=^TestIteration14$ \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-database-dsn=${PG_DATABASE_DSN} \
		-server-port="$$SERVER_PORT" \
		-key=$${TEMP_FILE} \
		-source-path=. ; \
	go test -v -race ./... ;

CRYPTO_DIR := ./crypto-key
CRYPTO_PRIVATE := ${CRYPTO_DIR}/private.pem
CRYPTO_PUBLIC := ${CRYPTO_DIR}/public.pem

.PHONY: crypto-key-generate
crypto-key-generate:
	mkdir -pv ${CRYPTO_DIR} && \
	openssl genrsa -out ${CRYPTO_PRIVATE} 4096 && \
	openssl rsa -in ${CRYPTO_PRIVATE} -outform PEM -pubout -out ${CRYPTO_PUBLIC}

.PHONY: crypto-key-clean
crypto-key-clean:
	rm -v -f  ${CRYPTO_PRIVATE} ${CRYPTO_PUBLIC} && \
	rm -d -f  ${CRYPTO_DIR}

HASH_KEY := "hash key"

.PHONY: server-run-with-args
server-run-with-args: build statictest db-up
	chmod +x ./cmd/server/server && \
		./cmd/server/server \
			-a ${SERVER_HOST}:${SERVER_PORT} \
			-d ${PG_DATABASE_DSN} \
			-p ${PPROF_SERVER_HOST}:${PPROF_SERVER_PORT} \
			-k ${HASH_KEY} \
			-t ${SERVER_TRUSTED_SUBNET} \
			-crypto-key ${CRYPTO_PRIVATE}

.PHONY: agent-run-with-args
agent-run-with-args: build statictest db-up
	chmod +x ./cmd/agent/agent && \
		./cmd/agent/agent \
			-a ${SERVER_HOST}:${SERVER_PORT} \
			-k ${HASH_KEY} \
			-crypto-key ${CRYPTO_PUBLIC}

.PHONY: pprof-mem-http
pprof-mem-http:
	go tool pprof -http=":9090" -seconds=30 http://${PPROF_SERVER_HOST}:${PPROF_SERVER_PORT}/debug/pprof/heap

.PHONY: pprof-mem-console
pprof-mem-console:
	go tool pprof -seconds=30 http://${PPROF_SERVER_HOST}:${PPROF_SERVER_PORT}/debug/pprof/heap

.PHONY: pprof-cpu-http
pprof-cpu-http:
	go tool pprof -http=":9090" -seconds=30 http://${PPROF_SERVER_HOST}:${PPROF_SERVER_PORT}/debug/pprof/profile

.PHONY: pprof-cpu-console
pprof-cpu-console:
	go tool pprof -seconds=30 http://${PPROF_SERVER_HOST}:${PPROF_SERVER_PORT}/debug/pprof/profile

.PHONY: pprof-mem-save
pprof-mem-save:
	curl --location http://${PPROF_SERVER_HOST}:${PPROF_SERVER_PORT}/debug/pprof/heap > ./profiles/result.pprof

.PHONY: pprofcompare
pprofcompare:
	go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof


.PHONY: db-up
db-up:
	PG_USER=${PG_USER} \
	PG_PASSWORD=${PG_PASSWORD} \
	PG_DB=${PG_DB} \
	PG_HOST=${PG_HOST} \
	PG_PORT=${PG_PORT} \
	PG_DATABASE_DSN=${PG_DATABASE_DSN} \
	PG_IMAGE=${PG_IMAGE} \
	PG_DOCKER_CONTEINER_NAME=${PG_DOCKER_CONTEINER_NAME} \
	docker compose -f ./docker-compose.yml up -d postgres

.PHONY: db-down
db-down:
	PG_USER=${PG_USER} \
	PG_PASSWORD=${PG_PASSWORD} \
	PG_DB=${PG_DB} \
	PG_HOST=${PG_HOST} \
	PG_PORT=${PG_PORT} \
	PG_DATABASE_DSN=${PG_DATABASE_DSN} \
	PG_IMAGE=${PG_IMAGE} \
	PG_DOCKER_CONTEINER_NAME=${PG_DOCKER_CONTEINER_NAME} \
	docker compose -f ./docker-compose.yml down postgres


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

# For run metrics in docker
UUID = $(shell cat /proc/sys/kernel/random/uuid)
METRICS_IMAGE := "metrics"
DOCKER_USER := ${USER}
UID := $(shell id -u)
GID := $(shell id -g)

DOCKER_PARAMS = \
    --volume ${PWD}:/home/${DOCKER_USER}/project \
	--volume ~/.vimrc:/home/${DOCKER_USER}/.vimrc:ro \
	--volume ~/.vim:/home/${DOCKER_USER}/.vim:ro \
	--volume ~/.gitignore:/home/${DOCKER_USER}/.gitignore:ro \
	--volume ~/git:/home/${DOCKER_USER}/git:ro \
    --tmpfs /tmp:exec,size=2G \
    --env UID=${UID} \
    --env GID=${GID} \
	--user ${UID}:${GID} \
    --volume="/etc/group:/etc/group:ro" \
    --volume="/etc/passwd:/etc/passwd:ro" \
    --volume="/etc/shadow:/etc/shadow:ro" \
    --name ${METRICS_IMAGE}-${UUID} \
    --privileged \
    --rm \
	-ti


.PHONY:cli
cli:
	@${DOCKER} run ${DOCKER_PARAMS} ${METRICS_IMAGE}:${BUILD_VERSION} bash

.PHONY:toolchain
toolchain:
	@${DOCKER} build \
		--build-arg DOCKER_USER=${DOCKER_USER} \
		--tag=${METRICS_IMAGE}:${BUILD_VERSION} \
		--pull ${NO_CACHE} \
        --rm -f ${DOCKERFILE} \
        ./
	@${DOCKER} tag ${METRICS_IMAGE}:${BUILD_VERSION} ${METRICS_IMAGE}:latest

