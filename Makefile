IMAGE_NAME = iamragsraghu/user-service
TAG ?= latest

DB_HOST ?= 127.0.0.1
DB_PORT ?= 3306
DB_USER ?= root
DB_PASS ?= Wsxokn@123
DB_NAME ?= users_db

build:
	go build -v -o user-service ./cmd/user-service

run: build
	DB_HOST=$(DB_HOST) DB_PORT=$(DB_PORT) DB_USER=$(DB_USER) DB_PASS=$(DB_PASS) DB_NAME=$(DB_NAME) ./user-service

docker-build:
	docker build -t $(IMAGE_NAME):$(TAG) .

docker-run:
	docker run --rm -p 8080:8080 \
	-e DB_HOST=host.docker.internal \
	-e DB_PORT=3306 \
	-e DB_USER=$(DB_USER) \
	-e DB_PASS=$(DB_PASS) \
	-e DB_NAME=$(DB_NAME) \
	$(IMAGE_NAME):$(TAG)

up:
	docker-compose up -d

down:
	docker-compose down -v

restart: down up

logs:
	docker-compose logs -f user-service

k8s-apply:
	kubectl apply -f k8s/

k8s-delete:
	kubectl delete -f k8s/
