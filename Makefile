DOCKER_COMPOSE=docker-compose
DC_RUN=$(DOCKER_COMPOSE) run --rm -T
DC_EXEC=$(DOCKER_COMPOSE) exec -T

install: docker-pull docker-build docker-up

pull docker-pull: # || true solves the fact that using --ignore-pull-failures does not always result in an exit code of 0
	$(DOCKER_COMPOSE) pull --quiet --ignore-pull-failures || true
build docker-build:
	$(DOCKER_COMPOSE) build --pull
up docker-up:
	$(DOCKER_COMPOSE) up -d --remove-orphans
down docker-down:
	$(DOCKER_COMPOSE) down --remove-orphans
stop docker-stop:
	$(DOCKER_COMPOSE) stop
logs docker-logs:
	$(DOCKER_COMPOSE) logs -f