fmt:
	gofmt -s -w .

test:
	go test -coverprofile cover.txt -v ./...