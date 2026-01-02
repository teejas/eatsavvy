.PHONY: build-% build-rabbitmq run-% run-rabbitmq stop-rabbitmq clean build run

build-%:
	go build -o bin/$* cmd/$*/main.go

build-rabbitmq:
	docker build -t rmq-delayed-exchange -f DockerfileRabbitMQ .

run-%:
	go run cmd/$*/main.go

run-rabbitmq:
	docker run --detach --name rabbitmq -p 5672:5672 -p 15672:15672 rmq-delayed-exchange

stop-rabbitmq:
	docker stop rabbitmq
	docker rm rabbitmq

build: build-api build-worker build-rabbitmq

run: run-api run-worker run-rabbitmq

start-cf-tunnel:
	cloudflared tunnel run --token ${CLOUDFLARED_TOKEN}

clean:
	rm -rf bin/*