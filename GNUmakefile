default: test

.PHONY: testall

testall:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

test:
	go test ./...

lint:
	golangci-lint run ./...

docs:
	tfplugindocs
