build:
	go build -o tmp cmd/main.go

format:
	gofumpt -w .
