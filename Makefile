prep:
	@go fmt

tests:
	@make prep && go test ./...

build:
	@make tests && go build cmd/tcsrdr/main.go