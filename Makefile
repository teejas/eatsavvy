.PHONY: build-% run-% run-rabbitmq stop-rabbitmq clean build run

build-%:
	go build -o bin/$* cmd/$*/main.go

run-%:
	go run cmd/$*/main.go

run-rabbitmq:
	docker run --detach --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:4-management

stop-rabbitmq:
	docker stop rabbitmq
	docker rm rabbitmq

build: build-api build-worker

run: run-api run-worker run-rabbitmq

clean:
	rm -rf bin/*