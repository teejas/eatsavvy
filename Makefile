.PHONY: build-% build-rabbitmq run-% run-rabbitmq stop-rabbitmq deploy-frontend start-cf-tunnel clean build run \
	docker-build-api docker-build-worker docker-build docker-run-api docker-run-worker docker-run-rabbitmq \
	docker-stop-api docker-stop-worker docker-stop-rabbitmq docker-stop docker-network \
	docker-push-api docker-push-worker docker-push docker-network

DOCKER_NETWORK = eatsavvy-network

build-%:
	cd backend && go build -o bin/$* cmd/$*/main.go

build-rabbitmq:
	docker build -t rmq-delayed-exchange -f infra/docker/DockerfileRabbitMQ .

build-frontend:
	cd frontend && npm run build

run-%:
	cd backend && go run cmd/$*/main.go

run-rabbitmq: docker-network
	docker run --detach --name eatsavvy-rabbitmq --network ${DOCKER_NETWORK} -p 5672:5672 -p 15672:15672 rmq-delayed-exchange

run-frontend:
	cd frontend && npm run dev

deploy-frontend: build-frontend
	@echo "Deploying UI to Cloudflare Pages..."
	@if ! command -v wrangler &> /dev/null; then \
			echo "Installing Wrangler CLI..."; \
			npm install -g wrangler; \
	fi
	@cd frontend && wrangler pages deploy dist/ --project-name=eatsavvy-frontend
	@echo "Frontend deployed to Cloudflare Pages"

stop-rabbitmq:
	docker stop eatsavvy-rabbitmq
	docker rm eatsavvy-rabbitmq

build: build-api build-worker build-rabbitmq build-frontend

run: run-api run-worker run-rabbitmq run-frontend

start-cf-tunnel:
	cloudflared tunnel run --token ${CLOUDFLARED_TOKEN}

clean:
	rm -rf backend/bin/* frontend/dist/*

# Docker build targets
docker-build-api:
	docker build -t eatsavvy-api -f infra/docker/DockerfileAPI .

docker-build-worker:
	docker build -t eatsavvy-worker -f infra/docker/DockerfileWorker .

docker-build-rabbitmq: build-rabbitmq

docker-build: docker-build-api docker-build-worker docker-build-rabbitmq

# Docker run targets
docker-run-api: docker-network
	docker run --detach --name eatsavvy-api --network ${DOCKER_NETWORK} -p 8080:8080 eatsavvy-api

docker-run-worker: docker-network
	docker run --detach --name eatsavvy-worker --network ${DOCKER_NETWORK} eatsavvy-worker

docker-run-rabbitmq: run-rabbitmq

docker-stop-rabbitmq: stop-rabbitmq

docker-stop-api:
	docker stop eatsavvy-api
	docker rm eatsavvy-api

docker-stop-worker:
	docker stop eatsavvy-worker
	docker rm eatsavvy-worker

docker-run: docker-run-api docker-run-worker docker-run-rabbitmq

docker-stop: docker-stop-api docker-stop-worker docker-stop-rabbitmq

# Docker network
docker-network:
	@docker network inspect ${DOCKER_NETWORK} >/dev/null 2>&1 || docker network create ${DOCKER_NETWORK}

# OCIR push targets
# Required env vars: OCIR_REGION (e.g., iad.ocir.io), OCIR_NAMESPACE (tenancy namespace)
# Optional: OCIR_REPO (defaults to eatsavvy)
OCIR_REPO ?= eatsavvy

docker-push-api: docker-build-api
	docker tag eatsavvy-api ${OCIR_REGION}/${OCIR_NAMESPACE}/${OCIR_REPO}/api:v1.0.0
	docker push ${OCIR_REGION}/${OCIR_NAMESPACE}/${OCIR_REPO}/api:v1.0.0

docker-push-worker: docker-build-worker
	docker tag eatsavvy-worker ${OCIR_REGION}/${OCIR_NAMESPACE}/${OCIR_REPO}/worker:latest
	docker push ${OCIR_REGION}/${OCIR_NAMESPACE}/${OCIR_REPO}/worker:latest

docker-push-rabbitmq: docker-build-rabbitmq
	docker tag rmq-delayed-exchange ${OCIR_REGION}/${OCIR_NAMESPACE}/${OCIR_REPO}/rabbitmq:latest
	docker push ${OCIR_REGION}/${OCIR_NAMESPACE}/${OCIR_REPO}/rabbitmq:latest

docker-push: docker-push-api docker-push-worker docker-push-rabbitmq