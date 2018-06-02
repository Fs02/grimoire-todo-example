all: build start

dep:
	dep ensure
	go get github.com/Fs02/kamimai

migrate:
	export $$(cat .env | grep -v ^\# | xargs) && \
	kamimai --driver=mysql --dsn="mysql://$$MYSQL_USERNAME:$$MYSQL_PASSWORD@($$MYSQL_HOST:$$MYSQL_PORT)/$$MYSQL_DATABASE" --directory=./migrations sync

rollback:
	export $$(cat .env | grep -v ^\# | xargs) && \
	kamimai --driver=mysql --dsn="mysql://$$MYSQL_USERNAME:$$MYSQL_PASSWORD@($$MYSQL_HOST:$$MYSQL_PORT)/$$MYSQL_DATABASE" --directory=./migrations down

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

start:
	export $$(cat .env | grep -v ^\# | xargs) && \
	./grimoire-todo-example
