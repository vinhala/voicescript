# docker-stack-base
Base repo for setting up a VM with docker swarm stack

## Required GitHub secrets
- VM_HOST
- VM_USER_NAME
- VM_SSH_KEY

## Required GitHub variables
- STACK_NAME

## Usage
Use as template for creating new docker swarm stack deployments. As a first step run the Setup VM workflow to install dependencies, configure the user and create the caddy network. Will automatically skip steps that are not necessary anymore. The Setup VM workflow can thus be run without problem multiple times.