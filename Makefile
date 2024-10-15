DOCKER_USERNAME = anhcaoo
IMAGE_NAME = stormbreaker
TAGGED_VERSION = 1.0.0
DOCKER_IMAGE = ${IMAGE_NAME}:${TAGGED_VERSION} 

.PHONY: build tag push test docker

build: 
	docker build --tag ${DOCKER_IMAGE} .

tag: 
	docker tag ${DOCKER_IMAGE} ${DOCKER_USERNAME}/${DOCKER_IMAGE}

push: 
	docker push ${DOCKER_USERNAME}/${DOCKER_IMAGE}

test: 
	go test ./...

docker: test build