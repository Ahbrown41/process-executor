
.PHONY: build
build: ## - Building App
	@printf "\033[32m\xE2\x9c\x93 Building Windows App\n\033[0m"
	@env CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o process-orchestrator.exe
	@printf "\033[32m\xE2\x9c\x93 Building Mac App\n\033[0m"
	@env CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o process-orchestrator