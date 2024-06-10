
build-docker:
	go mod vendor && cd .. && \
	docker build -t ai-fake-news .

run-docker:
	docker run -p 9012:9012 ai-fake-news

run:
	go run . --dev

test:
	go test -v --race ./...

lint:
	golangci-lint run

format:
	go fmt ./...