default: run build format

set shell := ["zsh", "-c"]

run:
	source setenv.sh && air 

templ:
	source setenv.sh && templ generate -watch -v

build:
	go build -o tmp cmd/main.go

format:
	gofumpt -w .
