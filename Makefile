up_proxy:
	HTTPS_PROXY=$(HTTPS_PROXY)  docker-compose -f deploy/docker-compose/docker-compose.yaml up -d

up:
	docker-compose -f deploy/docker-compose/docker-compose.yaml up -d

down:
	docker-compose -f deploy/docker-compose/docker-compose.yaml down
