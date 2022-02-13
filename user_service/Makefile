build:
	go build -o bin/user_ms app/cmd/main.go

run:
	go run app/cmd/main.go

test:
	(go test -v -race -timeout 1m -coverprofile cover.out ./app/internal/user/...; go tool cover -html=cover.out -o cover.html; rm cover.out)