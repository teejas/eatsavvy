.PHONY: build-% build-rabbitmq run-% run-rabbitmq stop-rabbitmq deploy-frontend start-cf-tunnel clean build run \
	docker-build-api docker-build-worker docker-build docker-run-api docker-run-worker \
	docker-push-api docker-push-worker docker-push docker-network

DOCKER_NETWORK = eatsavvy-network

build-%:
	cd backend && go build -o bin/$* cmd/$*/main.go

build-rabbitmq:
	docker build -t rmq-delayed-exchange -f infra/DockerfileRabbitMQ .

build-frontend:
	cd frontend && npm run build

run-%:
	cd backend && go run cmd/$*/main.go

run-rabbitmq: docker-network
	docker run --detach --name rabbitmq --network ${DOCKER_NETWORK} -p 5672:5672 -p 15672:15672 rmq-delayed-exchange

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
	docker stop rabbitmq
	docker rm rabbitmq

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

docker-build: docker-build-api docker-build-worker

# Docker run targets
docker-run-api: docker-network
	docker run --name eatsavvy-api --network ${DOCKER_NETWORK} -p 8080:8080 eatsavvy-api

docker-run-worker: docker-network
	docker run --name eatsavvy-worker --network ${DOCKER_NETWORK} eatsavvy-worker

# Docker network
docker-network:
	@docker network inspect ${DOCKER_NETWORK} >/dev/null 2>&1 || docker network create ${DOCKER_NETWORK}

# OCIR push targets
# Required env vars: OCIR_REGION (e.g., iad.ocir.io), OCIR_NAMESPACE (tenancy namespace)
# Optional: OCIR_REPO (defaults to eatsavvy)
OCIR_REPO ?= eatsavvy

docker-push-api: docker-build-api
	docker tag eatsavvy-api ${OCIR_REGION}/${OCIR_NAMESPACE}/${OCIR_REPO}/api:latest
	docker push ${OCIR_REGION}/${OCIR_NAMESPACE}/${OCIR_REPO}/api:latest

docker-push-worker: docker-build-worker
	docker tag eatsavvy-worker ${OCIR_REGION}/${OCIR_NAMESPACE}/${OCIR_REPO}/worker:latest
	docker push ${OCIR_REGION}/${OCIR_NAMESPACE}/${OCIR_REPO}/worker:latest

docker-push: docker-push-api docker-push-worker