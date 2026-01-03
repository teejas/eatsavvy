.PHONY: build-% build-rabbitmq run-% run-rabbitmq stop-rabbitmq clean build run

build-%:
	go build -o bin/$* cmd/$*/main.go

build-rabbitmq:
	docker build -t rmq-delayed-exchange -f DockerfileRabbitMQ .

build-ui:
	cd ui && npm run build

run-%:
	go run cmd/$*/main.go

run-rabbitmq:
	docker run --detach --name rabbitmq -p 5672:5672 -p 15672:15672 rmq-delayed-exchange

run-ui:
	cd ui && npm run dev

deploy-ui: build-ui
	@echo "Deploying UI to Cloudflare Pages..."
	@if ! command -v wrangler &> /dev/null; then \
			echo "Installing Wrangler CLI..."; \
			npm install -g wrangler; \
	fi
	@cd ui && wrangler pages deploy dist/ --project-name=eatsavvy-ui
	@echo "Frontend deployed to Cloudflare Pages"

stop-rabbitmq:
	docker stop rabbitmq
	docker rm rabbitmq

build: build-api build-worker build-rabbitmq build-ui

run: run-api run-worker run-rabbitmq run-ui

start-cf-tunnel:
	cloudflared tunnel run --token ${CLOUDFLARED_TOKEN}

clean:
	rm -rf bin/*