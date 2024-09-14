default: test

.PHONY: testall test lint docs

testall:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

docker-run:
	docker compose up -d --build

test:
	go test ./...

lint:
	golangci-lint run ./...

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.19.4
