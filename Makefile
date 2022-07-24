# note: call scripts from /scripts

test-cover:
	go test .\... -coverprofile=coverage.txt -covermode count
	go tool cover -html=coverage