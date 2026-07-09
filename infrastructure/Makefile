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

# these will speed up builds, for docker-compose >= 1.25
export COMPOSE_DOCKER_CLI_BUILD=1
export DOCKER_BUILDKIT=1

init:
	docker swarm init

# build local images
build-local:
	docker build -t ghcr.io/$(GITHUB_REPOSITORY_OWNER)/backend:latest ../api

# deploy services
deploy:
	docker stack deploy -c docker-compose.network.yml -c docker-compose.twenty.yml ${STACK_NAME} --detach=false --resolve-image=always
# down all services
down:
	docker stack rm ${STACK_NAME}
