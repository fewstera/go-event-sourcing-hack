export GO111MODULE = on

eventsourcing-hack: install-deps
	go build -o eventsourcing-hack ./cmd/eventsourcing-hack/main.go

.PHONY: lint
lint:
	golint ./...

.PHONY: fmt
fmt:
	gofmt -s -w .

.PHONY: test
test: fmt lint
	go test -v -coverprofile=coverage/coverage.out ./...

.PHONY: view-coverage
view-coverage:
	cp ./coverage/coverage.out ./coverage.out
	go tool cover -html=./coverage.out

.PHONY: install-deps
install-deps:
	go get

.PHONY: provision-database
provision-database:
	go run cmd/create-dev-tables/main.go

.PHONY: start-dev
start-dev: install-deps provision-database
	./scripts/live-reload.sh

.PHONY: start
start: install-deps build provision-database
	./eventsourcing-hack

.PHONY: clean
clean:
	rm -rf ./eventsourcing-hack ./coverage/

# docker-rebuild will only need to be called when you have made changes
# to the Dockerfile or add new go dependencies. Not calling this
# on every start should speed up the dev start process.
.PHONY: docker-rebuild
docker-rebuild:
	docker-compose build

# docker-start-dev should be called to start the dev server. The docker-compose
# file overrides the /app dir with a volume to the root of this project. So
# changes to the go files should be picked up without rebuilding the image.
.PHONY: docker-start-dev
docker-start-dev:
	docker-compose up app

# docker-run-tests runs the tests via docker-compose.
.PHONY: docker-run-tests
docker-run-tests:
	docker-compose run app make test

# docker-kill should be used to kill and remove the created docker images. This
# can be used to get a fresh database or when things have gone funky.
.PHONY: docker-kill
docker-kill:
	docker-compose kill
	docker-compose rm -f
