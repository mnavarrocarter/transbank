fmt:
	gofmt -s -w .

test:
	go test -v ./...