ARTIFACT_PATH ?= "src/cmd"
lint: ## Lint the files
	@./scripts/golangci-lint.sh src

build:
	@cd src/ && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o exporter .

run:
	@cd src/ && go run *.go
