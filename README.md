# Library

The library project is part of a coding interview for Jordan Crabtree.

[project page](https://github.com/slcjordan/library)

## Table of contents
- Requirements
- Documentation
- Installation
- Configuration
- Getting Started

## Requirements

Standard unix tools, make and docker.

Nice to have: [Direnv](https://direnv.net/) can be useful for overriding environment variables or running the project outside of docker.

## Documentation

Library follows the convention of self-documenting code. Documentation gets pulled out and into markdown files viewable by markdown viewers or the internal project documentation site. To generate needed code, dockerfiles, test reports (this may take a while on the first run) and to view the internal documentation site:

```
make docs-wait
```

## Configuration

example .envrc file for local development:

```
export PGDATABASE=postgres
export PGHOST=librarydb
export PGUSER=libraryuser
export PGPASSWORD=*****
export PGPORT=5432
export LIBRARY_PG_CONNECTION_STRING="postgres://${PGUSER}:${PGPASSWORD}@localhost:${PGPORT}/${PGDATABASE}?sslmode=disable"
export LIBRARY_PG_CONNECT_TIMEOUT="3s"
export LIBRARY_HTTP_BASE_URL="/api/v1"
export LIBRARY_HTTP_LISTEN_ADDRESS="0.0.0.0:5082"
export LIBRARY_HTTP_MAX_LIST_SIZE="500"
```

## Getting started

List all make commands:
```
make help
```

Run the project locally on docker (published port 5082):
```
make run
```

Run project tests:
```
make test
```
