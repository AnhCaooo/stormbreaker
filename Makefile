APPLICATION_NAME ?= stormbreaker
TAG_VERSION ?= latest

docker: 
	go test ./... && docker build --tag ${APPLICATION_NAME}/${TAG_VERSION} .
