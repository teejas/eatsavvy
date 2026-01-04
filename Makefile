.PHONY: build-% build-rabbitmq run-% run-rabbitmq stop-rabbitmq deploy-frontend start-cf-tunnel clean build run \
	docker-build-api docker-build-worker docker-build docker-run-api docker-run-worker docker-run-rabbitmq \
	docker-stop-api docker-stop-worker docker-stop-rabbitmq docker-stop docker-network \
	docker-push-api docker-push-worker docker-push docker-network

DOCKER_NETWORK = eatsavvy-network
PLATFORMS = linux/amd64,linux/arm64

# --------------------------------BUILD TARGETS--------------------------------
build-%:
	cd backend && go build -o bin/$* cmd/$*/main.go

build-rabbitmq:
	docker buildx build --platform ${PLATFORMS} -t rmq-delayed-exchange -f infra/docker/DockerfileRabbitMQ .

build-frontend:
	cd frontend && VITE_EATSAVVY_API_URL="https://api.eatsavvy.org" npm run build

build: build-api build-worker build-rabbitmq build-frontend

clean:
	rm -rf backend/bin/* frontend/dist/*

# Docker build targets
docker-build-api:
	docker buildx build --platform ${PLATFORMS} -t eatsavvy-api -f infra/docker/DockerfileAPI .

docker-build-worker:
	docker buildx build --platform ${PLATFORMS} -t eatsavvy-worker -f infra/docker/DockerfileWorker .

docker-build-rabbitmq: build-rabbitmq

docker-build: docker-build-api docker-build-worker docker-build-rabbitmq	

# --------------------------------LOCAL RUN TARGETS--------------------------------
run-%:
	cd backend && go run cmd/$*/main.go

run-rabbitmq: docker-network
	docker run --detach --name eatsavvy-rabbitmq --network ${DOCKER_NETWORK} -p 5672:5672 -p 15672:15672 rmq-delayed-exchange

run-frontend:
	cd frontend && VITE_EATSAVVY_API_URL="http://localhost:8080" npm run dev

stop-rabbitmq:
	docker stop eatsavvy-rabbitmq
	docker rm eatsavvy-rabbitmq

run: run-api run-worker run-rabbitmq run-frontend

start-cf-tunnel:
	cloudflared tunnel run --token ${CLOUDFLARED_TOKEN}

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

# Docker network
docker-network:
	@docker network inspect ${DOCKER_NETWORK} >/dev/null 2>&1 || docker network create ${DOCKER_NETWORK}

docker-remove-network:
	docker network rm ${DOCKER_NETWORK}

docker-run:
	$(MAKE) docker-run-rabbitmq
	@echo "Waiting 10s for RabbitMQ to start..."
	sleep 10
	$(MAKE) docker-run-api
	$(MAKE) docker-run-worker

docker-stop: docker-stop-api docker-stop-worker docker-stop-rabbitmq docker-remove-network


# --------------------------------DEPLOYMENT TARGETS--------------------------------
# OCIR push targets
# Required env vars: OCIR_REGION (e.g., iad.ocir.io), OCIR_NAMESPACE (tenancy namespace)
# Optional: OCIR_REPO (defaults to eatsavvy)
OCIR_REPO ?= eatsavvy
OCIR_REGION ?= sjc.ocir.io
OCIR_NAMESPACE ?= axifwvgpnbqk
TIMESTAMP := $(shell date +%Y%m%d-%H%M%S)

docker-push-api:
	docker buildx build --platform ${PLATFORMS} \
		-t ${OCIR_REGION}/${OCIR_NAMESPACE}/${OCIR_REPO}/api:latest \
		-t ${OCIR_REGION}/${OCIR_NAMESPACE}/${OCIR_REPO}/api:${TIMESTAMP} \
		-f infra/docker/DockerfileAPI --push .

docker-push-worker:
	docker buildx build --platform ${PLATFORMS} \
		-t ${OCIR_REGION}/${OCIR_NAMESPACE}/${OCIR_REPO}/worker:latest \
		-t ${OCIR_REGION}/${OCIR_NAMESPACE}/${OCIR_REPO}/worker:${TIMESTAMP} \
		-f infra/docker/DockerfileWorker --push .

docker-push-rabbitmq:
	docker buildx build --platform ${PLATFORMS} \
		-t ${OCIR_REGION}/${OCIR_NAMESPACE}/${OCIR_REPO}/rabbitmq:latest \
		-t ${OCIR_REGION}/${OCIR_NAMESPACE}/${OCIR_REPO}/rabbitmq:${TIMESTAMP} \
		-f infra/docker/DockerfileRabbitMQ --push .

docker-push: docker-push-api docker-push-worker docker-push-rabbitmq

deploy-rabbitmq: docker-push-rabbitmq
	cd infra/k8s && ./deploy.sh apply -f 02-rabbitmq.yaml
	kubectl rollout restart deployment rabbitmq -n eatsavvy

deploy-api: docker-push-api
	cd infra/k8s && ./deploy.sh apply -f 03-api.yaml
	kubectl rollout restart deployment api -n eatsavvy

deploy-worker: docker-push-worker
	cd infra/k8s && ./deploy.sh apply -f 04-worker.yaml
	kubectl rollout restart deployment worker -n eatsavvy

deploy-backend: deploy-rabbitmq deploy-api deploy-worker

deploy-frontend: build-frontend
	@echo "Deploying UI to Cloudflare Pages..."
	@if ! command -v wrangler &> /dev/null; then \
			echo "Installing Wrangler CLI..."; \
			npm install -g wrangler; \
	fi
	@cd frontend && wrangler pages deploy dist/ --project-name=eatsavvy-frontend
	@echo "Frontend deployed to Cloudflare Pages"