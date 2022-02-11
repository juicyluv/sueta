run:
	go run app/cmd/main.go

test:
	(go test -v -race -timeout 30s -coverprofile cover.out ./...; go tool cover -html=cover.out -o cover.html; rm cover.out)