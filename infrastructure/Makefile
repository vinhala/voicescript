# For local work: Define a variable to check if .env file exists
ENV_FILE := .env

# Define the default stack name
STACK_NAME := stack-name
GITHUB_REPOSITORY_OWNER := github-repository-owner

# Check if the .env file exists, and include it if it does
ifneq ("$(wildcard $(ENV_FILE))","")
include $(ENV_FILE)
.EXPORT_ALL_VARIABLES:
endif

export $(shell sed 's/=.*//' $(ENV_FILE))

# these will speed up builds, for docker-compose >= 1.25
export COMPOSE_DOCKER_CLI_BUILD=1
export DOCKER_BUILDKIT=1

init:
	docker swarm init

# build local images
build-local:
	docker build -t ghcr.io/$(GITHUB_REPOSITORY_OWNER)/image-name:latest ../service-name

# deploy services without databases
deploy-local:
	docker stack deploy -c docker-compose.network-local.yml ${STACK_NAME} --detach=false --resolve-image=always

# deploy all services with databases and production urls and env
deploy-prod:
	docker stack deploy -c docker-compose.network-prod.yml ${STACK_NAME} --detach=false --resolve-image=always

# down all services
down:
	docker stack rm ${STACK_NAME}
