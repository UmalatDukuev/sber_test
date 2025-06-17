DOCKER_IMAGE = sber_test
DOCKER_CONTAINER = sber_test_container

test:
	go test -v -cover ./...

lint:
	golangci-lint run

build:
	docker build -t $(DOCKER_IMAGE) .

run:
	docker run -p 8080:8080 --name $(DOCKER_CONTAINER) $(DOCKER_IMAGE)

stop:
	docker stop $(DOCKER_CONTAINER) || true
	docker rm $(DOCKER_CONTAINER) || true

clean:
	docker rmi $(DOCKER_IMAGE) || true
