export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help


help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

sqlc: ##Run orm generating
	docker run --rm -v "${CURDIR}:/src" -w /src sqlc/sqlc generate --no-remote 
.PHONY: sqlc

docker-build:  ### build  image
	docker build -t crypto-scanner -f ./build/Dockerfile .
.PHONY: docker-build

compose-up: ### Run docker-compose for the dev env
	docker compose -f ./build/docker-compose.yml --env-file ./config/dev.env -p crypto_diploma up --build -d crypto && docker compose logs -f
.PHONY: compose-up

run: ## start dev compose container
	 @$(MAKE) docker-build
	 @$(MAKE) compose-up
.PHONY: run 


