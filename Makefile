# AnhCao 2024
DOCKER_USERNAME = anhcaoo
DOCKER_IMAGE = stormbreaker
TAGGED_VERSION = 1.0.0
DOCKER_CONTAINER = ${DOCKER_IMAGE}:${TAGGED_VERSION} 

.PHONY: build tag push test docker swagger

build: 
	docker build --tag ${DOCKER_CONTAINER} .

tag: 
	docker tag ${DOCKER_CONTAINER} ${DOCKER_USERNAME}/${DOCKER_CONTAINER}

push: 
	docker push ${DOCKER_USERNAME}/${DOCKER_CONTAINER}

test: 
	go test ./...

swagger: 
	swag init -g cmd/main.go

docker: test build