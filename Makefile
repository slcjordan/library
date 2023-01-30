#
# I want to enable users to force re-running a target. PHONY targets are set up
# to run every time they are called. Sub-dependencies run only when their real
# targets are stale.
#

READMES=$(shell find . -type f -not -path '*node_modules*' -not -path './.cache*' -iname '*.md')
GOPKGS=$(shell find . -not -path './.cache*' -iname '*.go' | xargs dirname | sort | uniq)
GOFILES=$(shell find . -not -path './.cache*' -iname '*.go' -iname '*.go')
SQLC_QUERIES=$(shell find db/sqlc -not -path './.cache*' -iname '*.sql' | grep -v schema.sql)
MIGRATIONS=$(shell find db/migrate -not -path './.cache*' -iname '*.sql')
DOCS_APP_SRCS=$(shell find ./docs -not -path './.cache*' -not -path './docs/nginx.conf' -not -path './docs/node_modules/*' -type f)

# overrideable by environment variables
SQLC_VERSION?=1.16.0
NODE_VERSION?=18.7.0
POSTGRES_VERSION?=15.1
PGDATABASE?=postgres
PGHOST?=librarydb
PGUSER?=libraryuser
PGPASSWORD?=changeme
PGPORT?=5432
DB_CONN_STRING?=postgres://${PGUSER}:${PGPASSWORD}@${PGHOST}:${PGPORT}/${PGDATABASE}?sslmode=disable
PROMPT_MIGRATION?=$(shell bash -c 'read -p "Migration Identifier: " migration; echo $$migration')

# PHONY rules depend on the real targets of subtasks.

.PHONY: build
build: generate ## build docker image
	docker build --tag library-api --file dockerfiles/api .
	@touch .cache/make/build

.PHONY: test
test: go-test ## run all tests

.PHONY: run
run: generate docker-network postgres-wait ## run the api server on port 5082
	docker run \
		--interactive \
		--tty \
		--rm \
		--name library-local \
		--network library \
		--env LIBRARY_PG_CONNECT_TIMEOUT=3s \
		--env LIBRARY_PG_CONNECTION_STRING=${DB_CONN_STRING} \
		--env LIBRARY_HTTP_BASE_URL=/api/v1 \
		--env LIBRARY_HTTP_LISTEN_ADDRESS=0.0.0.0:5082 \
		--env LIBRARY_HTTP_MAX_LIST_SIZE=1000 \
		--volume ${PWD}:/go/src/github.com/slcjordan/library \
		--volume ${PWD}/.cache/pkg:/go/pkg \
		--workdir /go/src/github.com/slcjordan/library \
		--publish 5082:5082 \
		go-generate go run cmd/api/*.go

.PHONY: generate
generate: .cache/make/go-openapi .cache/make/go-sqlc .cache/make/go-generate ## run all code generation.

.PHONY: docs
docs: .cache/make/docs-manifest ## Build all docs

.PHONY: docs-clean
docs-clean: ## clean docs directory
	rm -r .cache/docs
	rm .cache/make/docs-readme
	rm .cache/make/docs-makefile

.PHONY: docs-app-build
docs-app-build: make-cache ## Build docs web app
	docker run \
		--interactive \
		--tty \
		--rm \
		--volume ${PWD}:${PWD} \
		--volume ${PWD}/.cache/node_modules:${PWD}/docs/node_modules \
		--workdir ${PWD}/docs \
		node:${NODE_VERSION} npm install && npm --prefix docs run build
	touch .cache/make/docs-app-build

.PHONY: docs-go
docs-go: make-cache docker-go-generate ## Generate documentation from godoc comments.
	ls ${GOPKGS} | xargs dirname | xargs -I {} mkdir -p ./.cache/docs/{}
	docker run \
		--interactive \
		--tty \
		--rm \
		--volume ${PWD}:/go/src/github.com/slcjordan/library \
		--volume ${PWD}/.cache/pkg:/go/pkg \
		--workdir /go/src/github.com/slcjordan/library \
		go-generate /bin/bash -c 'echo ${GOPKGS} | sed "s/ /\n/g" | xargs -I {} gomarkdoc --output ./.cache/docs/{}/godoc.md {}'
	@touch .cache/make/docs-go

.PHONY: docs-manifest
docs-manifest: make-cache .cache/make/docs-readme .cache/make/docs-makefile .cache/make/docs-go .cache/make/docs-openapi ## Generate documentation manifest file
	find .cache/docs -iname "*.md" -print | sed 's/^\.cache\/docs\///g'> .cache/docs/manifest.txt
	@touch .cache/make/docs-manifest

.PHONY: docs-makefile
docs-makefile: make-cache ## generate markdown for makefile help
	$(MAKE) help | grep -v '${PWD}' > .cache/docs/Makefile.md
	@touch .cache/make/docs-makefile

.PHONY: docs-openapi
docs-openapi: make-cache docker-openapi-codegen ## Generate documentation from openapi spec.
	ls ${GOPKGS} | xargs dirname | xargs -I {} mkdir -p ./.cache/docs/{}
	docker run \
		--interactive \
		--tty \
		--rm \
		--volume ${PWD}:${PWD} \
		--workdir ${PWD} \
		openapi-codegen bash -c "yq \
			--no-colors \
			-o=json \
			http/openapi.yaml > .cache/docs/openapi.json && \
		widdershins \
			--search false \
			--language_tabs 'go' 'javascript' \
			--summary .cache/docs/openapi.json \
			-o .cache/docs/api.md"
	@touch .cache/make/docs-openapi


.PHONY: docs-readme
docs-readme: make-cache ## stage readmes recursively for the documentation server
	ls ${READMES} | xargs dirname | xargs -I {} mkdir -p ./.cache/docs/{}
	ls ${READMES} | xargs -I {} cp {} ./.cache/docs/{}
	@touch .cache/make/docs-readme

# TODO implement live-reload
.PHONY: docs-start
docs-start: docker-network .cache/make/go-test .cache/make/docs-manifest .cache/make/docs-app-build ## start docs server at port 5081 if it isn't started
	docker container inspect --format='library-docs is {{.State.Status}}' library-docs || docker run \
		--detach \
		--interactive \
		--name library-docs \
		--rm \
		--tty \
		--volume ${PWD}/docs/nginx.conf:/etc/nginx/nginx.conf \
		--volume ${PWD}/.cache/docs-app:/usr/share/nginx/html/app \
		--volume ${PWD}/.cache/docs:/usr/share/nginx/html/docs/markdown \
		--volume ${PWD}/.cache/test:/usr/share/nginx/html/docs/test \
		--publish 5081:80 \
		--network 'library' \
		nginx

.PHONY: docs-stop
docs-stop:
	- docker stop library-docs

.PHONY: docs-wait ## start docs server if it isn't started (at port 5081) and wait for it to be ready
docs-wait: docs-start
	until docker run \
		--name libary-docs-wait \
		--rm \
		--network 'library' \
		curlimages/curl library-docs > /dev/null; \
	do \
			sleep 3; \
	done
	@echo "server is running."
	@echo "view docs at http://localhost:5081"



.PHONY: go-generate
go-generate: make-cache docker-go-generate ## Generate go code from go-generate comments.
	mkdir -p test
	@echo "(re)generating mocks"
	- rm -r test/mocks/*
	docker run \
		--interactive \
		--tty \
		--rm \
		--volume ${PWD}:/go/src/github.com/slcjordan/library \
		--volume ${PWD}/.cache/pkg:/go/pkg \
		--workdir /go/src/github.com/slcjordan/library \
		go-generate bash -c 'go generate -v ${GOPKGS}'
	@touch .cache/make/go-generate

.PHONY: go-openapi
go-openapi: make-cache docker-go-generate ## Generate go code from open-api.
	- docker run \
    --interactive \
    --tty \
    --rm \
    --volume `pwd`:`pwd` \
    --workdir `pwd` \
    --entrypoint oapi-codegen \
    go-generate \
      -config http/openapi-codegen.yaml \
      http/openapi.yaml > http/openapi.go && touch .cache/make/go-openapi || cat http/openapi.go
	@echo
	docker run \
		--interactive \
		--tty \
		--rm \
		--volume ${PWD}:/go/src/github.com/slcjordan/library \
		--volume ${PWD}/.cache/pkg:/go/pkg \
		--workdir /go/src/github.com/slcjordan/library \
		go-generate go fmt http/openapi.go

.PHONY: go-sqlc
go-sqlc: db/sqlc/schema.sql ## Generate go code from sqlc.
	- rm db/sqlc/*.sql.go
	docker run \
		--interactive \
		--tty \
		--rm \
		--volume ${PWD}:/repo \
		--workdir /repo/db/sqlc \
	kjconroy/sqlc:${SQLC_VERSION} generate
	@touch .cache/make/go-sqlc

.PHONY: go-test
go-test: make-cache generate postgres-wait ## Run go unit and integration tests
	docker run \
		--interactive \
		--tty \
		--rm \
		--network library \
		--volume ${PWD}:/go/src/github.com/slcjordan/library \
		--volume ${PWD}/.cache/pkg:/go/pkg \
		--workdir /go/src/github.com/slcjordan/library \
		--env LIBRARY_PG_CONNECTION_STRING=${DB_CONN_STRING} \
		--env LIBRARY_HTTP_LISTEN_ADDRESS='0.0.0.0' \
		--env LIBRARY_HTTP_BASE_URL=/api/v1 \
		go-generate go test \
			-v \
			-race \
			-covermode atomic \
			-coverpkg=$$(go list ./... | grep -v '\.cache' | paste -s -d, -) \
			-coverprofile ".cache/test/profile.cov" \
			$$(go list ./... | grep -v '\.cache') && go tool cover -html .cache/test/profile.cov -o .cache/test/cover.html
	@touch .cache/make/go-test

.PHONY: go-lint
go-lint: make-cache ## Run go lint checks
	docker run \
		--interactive \
		--tty \
		--rm \
		--network library \
		--volume ${PWD}:/go/src/github.com/slcjordan/library \
		--volume ${PWD}/.cache/lint:/root/.cache \
		--workdir /go/src/github.com/slcjordan/library \
		golangci/golangci-lint:v1.50.1 golangci-lint run -v

.DEFAULT_GOAL=help
.PHONY: help
help: ## Show this help.
	@echo "Make Commands:"
	@echo "---"
	@echo
	@cat $(MAKEFILE_LIST) | grep '^[a-z].*:.*##' | sed 's/\(.*\):.*##\(.*\)/* `make \1`:\2/'

.PHONY: postgres-start
postgres-start: make-cache docker-network ## Start postgres if it isn't started and return immediately.
	docker container inspect --format='postgres is {{.State.Status}}' ${PGHOST} || docker run \
		--detach \
		--name ${PGHOST} \
		--rm \
		--publish ${PGPORT}:${PGPORT} \
		--env POSTGRES_PASSWORD=${PGPASSWORD} \
		--env POSTGRES_USER=${PGUSER} \
		--env POSTGRES_DB=${PGDATABASE} \
		--env PGDATA=/var/lib/postgresql/data/pgdata \
		--network 'library' \
		--volume ${PWD}/.cache/data:/var/lib/postgresql/data \
		postgres:${POSTGRES_VERSION}

.PHONY: postgres-stop
postgres-stop:
	docker stop ${PGHOST}

.PHONY: postgres-wait
postgres-wait: postgres-start ## Start postgres if it isn't started and wait for it to be ready.
	until docker run \
		--name ${PGHOST}-wait \
		--rm \
		--network 'library' \
		postgres:${POSTGRES_VERSION} psql -d ${DB_CONN_STRING} -c 'SELECT 1'; \
	do \
			sleep 3; \
	done

.PHONY: psql
psql: ## psql client to the db
	$(MAKE) postgres-wait
	docker run \
		--interactive \
		--tty \
		--rm \
		--network library \
		--env PGDATABASE=${PGDATABASE} \
		--env PGHOST=${PGHOST} \
		--env PGPASSWORD=${PGPASSWORD} \
		--env PGPORT=${PGPORT} \
		--env PGUSER=${PGUSER} \
		--volume ${PWD}/db:/db \
		--workdir / \
		postgres:${POSTGRES_VERSION} psql

.PHONY: postgres-schema-dump
postgres-schema-dump: .cache/make/postgres-schema-migrate ## create postgres schema dump file under db/sqlc, which is necessary for sqlc
	$(MAKE) postgres-wait
	docker run \
		--interactive \
		--tty \
		--rm \
		--network library \
		--env PGDATABASE=${PGDATABASE} \
		--env PGHOST=${PGHOST} \
		--env PGPASSWORD=${PGPASSWORD} \
		--env PGPORT=${PGPORT} \
		--env PGUSER=${PGUSER} \
		--volume ${PWD}/db:/db \
		--workdir / \
		postgres:${POSTGRES_VERSION} pg_dump \
			--file db/sqlc/schema.sql \
			--schema-only

.PHONY: postgres-schema-migrate
postgres-schema-migrate: make-cache ## Perform unapplied database migrations.
	$(MAKE) postgres-wait
	docker run \
		--interactive \
		--tty \
		--rm \
		--network library \
		--volume ${PWD}:/go/src/github.com/slcjordan/library \
		--volume ${PWD}/.cache/pkg:/go/pkg \
		--workdir /go/src/github.com/slcjordan/library \
		go-generate goose -dir ./db/migrate postgres "${DB_CONN_STRING}" up
	@touch .cache/make/postgres-schema-migrate

.PHONY: postgres-schema-create-migration
postgres-schema-create-migration: make-cache ## Create a new empty schema migration.
	docker run \
		--interactive \
		--tty \
		--rm \
		--network library \
		--volume ${PWD}:/go/src/github.com/slcjordan/library \
		--volume ${PWD}/.cache/pkg:/go/pkg \
		--workdir /go/src/github.com/slcjordan/library \
		go-generate goose -dir ./db/migrate postgres "${DB_CONN_STRING}" create ${PROMPT_MIGRATION} sql
	@touch .cache/make/postgres-schema-migrate


# Docker

.PHONY: docker-go-generate
docker-go-generate:
	docker image inspect \
		--format 'docker image go-generate was created on {{.Created}}' \
		go-generate || docker build --tag go-generate --file dockerfiles/go-generate .

.PHONY: docker-network
docker-network:
	docker network inspect \
		--format 'docker network library was created on {{.Created}}' \
		library || docker network create library

.PHONY: docker-openapi-codegen
docker-openapi-codegen:
	docker image inspect \
		--format 'docker image openapi-codegen was created on {{.Created}}' \
		openapi-codegen || docker build --tag openapi-codegen --file dockerfiles/openapi-codegen .

.PHONY: make-cache
make-cache:
	@mkdir -p .cache/data .cache/docs .cache/docs-app .cache/make .cache/node_modules .cache/pkg .cache/test .cache/lint

# real rule targets depend on current step prerequisites.

db/sqlc/schema.sql: .cache/make/postgres-schema-migrate
	$(MAKE) postgres-schema-dump

.cache/make/docs-go: ${GOFILES}
	$(MAKE) docs-go

.cache/make/docs-makefile: $(MAKEFILE_LIST)
	$(MAKE) docs-makefile

.cache/make/docs-openapi: $(MAKEFILE_LIST)
	$(MAKE) docs-openapi

.cache/make/docs-readme: ${READMES}
	$(MAKE) docs-readme

.cache/make/go-generate: ${GOFILES}
	$(MAKE) go-generate

.cache/make/go-openapi: http/openapi-codegen.yaml http/openapi.yaml
	$(MAKE) go-openapi

.cache/make/go-sqlc: ${SQLC_QUERIES}
	$(MAKE) go-sqlc

.cache/make/postgres-schema-migrate: ${MIGRATIONS}
	$(MAKE) postgres-schema-migrate

.cache/make/docs-app-build: ${DOCS_APP_SRCS}
	$(MAKE) docs-app-build

.cache/make/docs-manifest: .cache/make/docs-go .cache/make/docs-makefile .cache/make/docs-readme
	$(MAKE) docs-manifest

.cache/make/go-test: ${GOFILES}
	$(MAKE) go-test

.cache/make/build: ${GOFILES}
	$(MAKE) build
