export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

include .env

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help


help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


docker-build:  ### build  image
	docker build -t crypto-scanner -f ./build/Dockerfile .
.PHONY: docker-build

compose-up: ### Run docker-compose for the .env
	docker-compose -f ./build/docker-compose.yml --env-file ./config/dev.env -p crypto_diploma up -d clickhouse grafana   
.PHONY: compose-up

run: ## start dev compose container
	 @$(MAKE) docker-build
	 @$(MAKE) compose-up
.PHONY: run 


deploy: ## deploy to prod
	 @$(MAKE) transfer path=apps/crypto-scanner 
	 @$(MAKE) docker-upload-remote path=apps/crypto-scanner/deployed_container.tar 
	 ssh ${HOST} 'cd apps/crypto-scanner && make compose-up'
	 ssh ${HOST} 'cd apps/crypto-scanner && rm Makefile deployed_container.tar && rm -rf config build'
	 echo Finished deploy
.PHONY: deploy 

cd: ## build and deploy to prod 
	 @$(MAKE) prepare
	 @$(MAKE) deploy 
.PHONY: cd 

prepare:
	 @$(MAKE) docker-build 
	 @$(MAKE) docker-save container=crypto-scanner:latest 
	 echo Prepare ended!
.PHONY: prepare

docker-upload-remote:
	ssh ${HOST} 'docker load --input $(path)'
.PHONY: docker-upload-remote

build-exe:## build exe file for program
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./crypto_building ./cmd/
.PHONY: build-exe

docker-save: ###save with container name
	docker save $(container) > deployed_container.tar
.PHONY: docker-save

transfer: ### transfer dev image and files for app
	 ssh ${HOST} 'mkdir -p $(path)' && rsync -av crypto_building conf.yaml Makefile deployed_container.tar .env config build migrations ${HOST}:$(path)   
.PHONY: transfer

remove-systemd-service:
	 ssh ${HOST} "rm -f /etc/systemd/system/crypto-parser.service"
.PHONY:remove-systemmd-service 

add-systemd-service:
	rsync -av ./deployment/crypto-parser.service  ${HOST}:/etc/systemd/system
.PHONY:add-systemmd-service 

start-daemon:
	ssh ${HOST} "systemctl daemon-reload && systemctl restart crypto-parser.service"
.PHONY: start-daemon

