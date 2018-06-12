JDB_MYSQL_DSN?=root:root@tcp($(shell docker-compose port mysql 3306))/testdb?parseTime=true
JDB_POSTGRES_DSN?=postgres://postgres:postgres@$(shell docker-compose port postgres 5432)/testdb?sslmode=disable
OS_TAG?=$(shell uname -s | tr '[:upper:]' '[:lower:]')
TAGS?=libsqlite3 $(OS_TAG) json1
TEST_ARGS?=-short -failfast
TEST_FULL_ARGS?=-count=1 -failfast

get:
	go get -v -t -d ./...

test:
	go test $(TEST_ARGS) \
		github.com/silas/jdb \
		github.com/silas/jdb/dialect/... \
		github.com/silas/jdb/internal/...

test_full:
	JDB_MYSQL_DSN="$(JDB_MYSQL_DSN)" JDB_POSTGRES_DSN="$(JDB_POSTGRES_DSN)" \
		go test -tags "$(TAGS)" $(TEST_FULL_ARGS) github.com/silas/jdb/...

update_json:
	./internal/json/update.sh

.PHONY: get test test_full update_json
