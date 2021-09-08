prep:
	@go generate ./... && go fmt

tests:
	@make prep && go test -v -cover ./...

build:
	@make tests && go build cmd/tcsrdr/main.go