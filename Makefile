# AnhCao 2024
TAGGED_VERSION = 1.0.0

.PHONY: build tag push test docker

build: 
	docker build --tag ${{ secrets.DOCKER_IMAGE }} .

tag: 
	docker tag ${{ secrets.DOCKER_IMAGE }} ${{ secrets.DOCKER_USERNAME }}/${{ secrets.DOCKER_IMAGE }}:${TAGGED_VERSION} 

push: 
	docker push ${{ secrets.DOCKER_USERNAME }}/${{ secrets.DOCKER_IMAGE }}:${TAGGED_VERSION} 

test: 
	go test ./...

docker: test build