ORGANIZATION := alextanhongpin
NAME := scraper
REPOSITORY := ${ORGANIZATION}/${NAME}

VERSION := 1.0.0

docker:
	@docker build -t ${REPOSITORY} .

tag:
	@docker tag ${REPOSITORY}:latest ${REPOSITORY}:${VERSION}

run-local:
	@go run main.go -o repos.csv

run-docker:
	@docker run -it -v "$PWD:/data" alextanhongpin/scraper -o /data/repos.csv

push:
	@docker push ${REPOSITORY}:${VERSION}
	@echo pushed ${REPOSITORY}:${VERSION}