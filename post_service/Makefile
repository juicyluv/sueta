build:
	go build -o bin/post_ms app/cmd/main.go

run:
	go run app/cmd/main.go

test:
	(go test -v -race -timeout 1m -coverprofile cover.out ./app/internal/post/...; go tool cover -html=cover.out -o cover.html; rm cover.out)